package main

import (
	"log"
	"os"
	"strings"

	"github.com/TechOctopus/davgus-ambulance-webapi/api"
	"github.com/TechOctopus/davgus-ambulance-webapi/internal/ambulance_wl"
	"github.com/gin-gonic/gin"
)

func main() {
	log.Printf("Server started")
	port := os.Getenv("AMBULANCE_API_PORT")
	if port == "" {
		port = "8080"
	}
	environment := os.Getenv("AMBULANCE_API_ENVIRONMENT")
	if !strings.EqualFold(environment, "production") { // case insensitive comparison
		gin.SetMode(gin.DebugMode)
	}
	engine := gin.New()
	engine.Use(gin.Recovery())

	handleFunctions := &ambulance_wl.ApiHandleFunctions{
		DepartmentsAPI: ambulance_wl.NewDepartmentsApi(),
		PatientsAPI:    ambulance_wl.NewPatientsApi(),
		PlacementsAPI:  ambulance_wl.NewPlacementsApi(),
	}

	ambulance_wl.NewRouterWithGinEngine(engine, *handleFunctions)

	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
