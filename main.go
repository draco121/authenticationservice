package main

import (
	"os"

	"github.com/draco121/authenticationservice/controllers"
	"github.com/draco121/authenticationservice/core"
	"github.com/draco121/authenticationservice/repository"
	"github.com/draco121/authenticationservice/routes"

	"github.com/draco121/common/clients"
	"github.com/draco121/common/database"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func RunApp() {
	db := database.NewMongoDatabaseDefaults()
	repo := repository.NewAuthenticationRepository(db)
	userServiceApiClient := clients.NewUserServiceApiClient(os.Getenv("USER_SERVICE_BASEURL"))
	service := core.NewAuthenticationService(repo, userServiceApiClient)
	controllers := controllers.NewControllers(service)
	router := gin.Default()
	routes.RegisterRoutes(controllers, router)
	router.Run()
}
func main() {
	godotenv.Load()
	RunApp()
}