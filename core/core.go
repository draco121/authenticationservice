package core

import (
	"context"
	"fmt"
	"time"

	"github.com/draco121/authenticationservice/repository"
	"github.com/draco121/common/clients"
	"github.com/draco121/common/jwt"
	"github.com/draco121/common/models"
	"github.com/draco121/common/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type IAuthenticationService interface {
	PasswordLogin(ctx context.Context, loginInput *models.LoginInput) (*models.LoginOutput, error)
	Authenticate(ctx context.Context, token string) (*models.JwtCustomClaims, error)
	RefreshLogin(ctx context.Context, refreshToken string) (*models.LoginOutput, error)
	Logout(ctx context.Context, token string) error
}

type authenticationService struct {
	IAuthenticationService
	repo                 repository.IAuthenticationRepository
	userServiceApiClient clients.IUserServiceApiClient
}

func NewAuthenticationService(repository repository.IAuthenticationRepository, userServiceApiClient clients.IUserServiceApiClient) IAuthenticationService {
	us := authenticationService{
		repo:                 repository,
		userServiceApiClient: userServiceApiClient,
	}
	return us
}

func (s authenticationService) PasswordLogin(ctx context.Context, loginInput *models.LoginInput) (*models.LoginOutput, error) {
	user, err := s.userServiceApiClient.GetUserByEmail(loginInput.Email)
	if err != nil {
		return nil, err
	} else {
		if utils.CheckPasswordHash(loginInput.Password, user.Password) {
			session := models.Session{
				UserId:    user.ID.Hex(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				ID:        primitive.NewObjectID(),
			}
			id, err := s.repo.InsertOne(ctx, &session)
			if err != nil {
				return nil, err
			} else {
				claims := models.JwtCustomClaims{
					Email:     user.Email,
					UserId:    user.ID.Hex(),
					Role:      user.Role,
					SessionId: id,
				}
				token, err := jwt.GenerateJWT(&claims)
				if err != nil {
					return nil, err
				} else {
					refreshToken, err := jwt.GenerateRefreshToken(id)
					if err != nil {
						return nil, err
					} else {
						return &models.LoginOutput{
							Token:        token,
							RefreshToken: refreshToken,
						}, nil
					}
				}
			}
		} else {
			return nil, fmt.Errorf("invalid credentials")
		}
	}
}

func (s authenticationService) Authenticate(ctx context.Context, token string) (*models.JwtCustomClaims, error) {
	claims, err := jwt.VerifyJwtToken(token)
	if err != nil {
		if claims != nil {
			return &claims.JwtCustomClaims, err
		}
		return nil, err
	} else {
		_, err := s.repo.FindOneById(ctx, claims.JwtCustomClaims.SessionId)
		if err != nil {
			return nil, err
		}
		return &claims.JwtCustomClaims, nil
	}
}

func (s authenticationService) RefreshLogin(ctx context.Context, refreshToken string) (*models.LoginOutput, error) {
	claims, err := jwt.VerfiyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	} else {
		session, err := s.repo.FindOneById(ctx, claims.SessionId)
		if err != nil {
			return nil, err
		} else {
			session, err = s.repo.UpdateOne(ctx, session)
			if err != nil {
				return nil, err
			} else {
				user, err := s.userServiceApiClient.GetUserById(session.UserId)
				if err != nil {
					return nil, err
				} else {
					claims := models.JwtCustomClaims{
						Email:     user.Email,
						Role:      user.Role,
						UserId:    user.ID.Hex(),
						SessionId: session.ID.Hex(),
					}
					newToken, err := jwt.GenerateJWT(&claims)
					if err != nil {
						return nil, err
					} else {
						return &models.LoginOutput{
							Token:        newToken,
							RefreshToken: refreshToken,
						}, nil
					}
				}
			}
		}
	}
}

func (s authenticationService) Logout(ctx context.Context, token string) error {
	claims, err := jwt.VerifyJwtToken(token)
	if err != nil {
		return err
	} else {
		_, err := s.repo.DeleteOneById(ctx, claims.JwtCustomClaims.SessionId)
		if err != nil {
			return err
		}
		return nil
	}
}
