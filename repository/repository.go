package repository

import (
	"context"
	"errors"
	"github.com/draco121/common/clients"
	"time"

	"github.com/draco121/common/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type IAuthenticationRepository interface {
	InsertOne(ctx context.Context, session *models.Session) (string, error)
	UpdateOne(ctx context.Context, session *models.Session) (*models.Session, error)
	FindOneById(ctx context.Context, id string) (*models.Session, error)
	DeleteOneById(ctx context.Context, id string) (*models.Session, error)
	GetUserByEmail(userId string) (*models.User, error)
	GetUserById(userId string) (*models.User, error)
}

type authenticationRepository struct {
	IAuthenticationRepository
	db                   *mongo.Database
	userServiceApiClient clients.IUserServiceApiClient
}

func NewAuthenticationRepository(database *mongo.Database, userServiceApiClient clients.IUserServiceApiClient) IAuthenticationRepository {
	repo := authenticationRepository{
		db:                   database,
		userServiceApiClient: userServiceApiClient,
	}
	return &repo
}

func (ur authenticationRepository) InsertOne(ctx context.Context, session *models.Session) (string, error) {
	result, err := ur.db.Collection("sessions").InsertOne(ctx, session)
	if err != nil {
		return "", err
	} else {
		id := result.InsertedID.(primitive.ObjectID)
		return id.Hex(), nil
	}
}

func (ur authenticationRepository) UpdateOne(ctx context.Context, session *models.Session) (*models.Session, error) {
	filter := bson.M{"_id": session.ID}
	update := bson.M{"$set": bson.M{
		"updatedAt": time.Now(),
	}}
	result := models.Session{}
	err := ur.db.Collection("sessions").FindOneAndUpdate(ctx, filter, update).Decode(&result)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else {
		return &result, nil
	}
}

func (ur authenticationRepository) FindOneById(ctx context.Context, id string) (*models.Session, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	} else {
		filter := bson.D{{Key: "_id", Value: objectId}}
		result := models.Session{}
		err := ur.db.Collection("sessions").FindOne(ctx, filter).Decode(&result)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		} else {
			return &result, nil
		}
	}
}

func (ur authenticationRepository) DeleteOneById(ctx context.Context, id string) (*models.Session, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	} else {
		filter := bson.D{{Key: "_id", Value: objectId}}
		result := models.Session{}
		err := ur.db.Collection("users").FindOneAndDelete(ctx, filter).Decode(&result)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		} else {
			return &result, nil
		}
	}
}

func (ur authenticationRepository) GetUserById(userId string) (*models.User, error) {
	return ur.userServiceApiClient.GetUserById(userId)
}

func (ur authenticationRepository) GetUserByEmail(email string) (*models.User, error) {
	return ur.userServiceApiClient.GetUserByEmail(email)
}
