package controllers

import (
	"net/http"

	"github.com/draco121/authenticationservice/core"
	"github.com/draco121/common/models"

	"github.com/gin-gonic/gin"
)

type Controllers struct {
	service core.IAuthenticationService
}

func NewControllers(service core.IAuthenticationService) Controllers {
	c := Controllers{
		service: service,
	}
	return c
}

func (s *Controllers) Login(c *gin.Context) {
	var loginInput models.LoginInput
	if c.ShouldBind(&loginInput) != nil {
		c.JSON(400, gin.H{
			"message": "data validation error",
		})
	} else {
		res, err := s.service.PasswordLogin(c.Request.Context(), &loginInput)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result": res,
			})
		}
	}
}

func (s *Controllers) Authenticate(c *gin.Context) {
	id := c.GetHeader("Authentication")
	if id != "" {
		result, err := s.service.Authenticate(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, result)
		}
	} else {
		c.Status(http.StatusForbidden)
	}

}

func (s *Controllers) RefreshLogin(c *gin.Context) {
	refreshToken := c.GetHeader("refreshToken")
	if refreshToken != "" {
		result, err := s.service.RefreshLogin(c.Request.Context(), refreshToken)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, result)
		}
	}
}

func (s *Controllers) Logout(c *gin.Context) {
	token := c.GetHeader("Authentication")
	err := s.service.Logout(c.Request.Context(), token)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	} else {
		c.Status(http.StatusNoContent)
	}
}
