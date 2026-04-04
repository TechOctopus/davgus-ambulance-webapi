package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/TechOctopus/davgus-ambulance-webapi/api"
	"github.com/TechOctopus/davgus-ambulance-webapi/internal/ambulance_wl"
	"github.com/TechOctopus/davgus-ambulance-webapi/internal/db_service"
	"github.com/gin-contrib/cors"
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

	corsMiddleware := cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
	engine.Use(corsMiddleware)

	// setup context update middleware
	departmentDb := db_service.NewMongoService[ambulance_wl.Department](db_service.MongoServiceConfig{
		Collection: "departments",
	})
	patientDb := db_service.NewMongoService[ambulance_wl.Patient](db_service.MongoServiceConfig{
		Collection: "patients",
	})
	placementDb := db_service.NewMongoService[ambulance_wl.Placement](db_service.MongoServiceConfig{
		Collection: "placements",
	})
	defer departmentDb.Disconnect(context.Background())
	defer patientDb.Disconnect(context.Background())
	defer placementDb.Disconnect(context.Background())

	engine.Use(func(ctx *gin.Context) {
		ctx.Set("db_service_departments", departmentDb)
		ctx.Set("db_service_patients", patientDb)
		ctx.Set("db_service_placements", placementDb)
		ctx.Next()
	})

	handleFunctions := &ambulance_wl.ApiHandleFunctions{
		DepartmentsAPI: ambulance_wl.NewDepartmentsApi(),
		PatientsAPI:    ambulance_wl.NewPatientsApi(),
		PlacementsAPI:  ambulance_wl.NewPlacementsApi(),
	}

	ambulance_wl.NewRouterWithGinEngine(engine, *handleFunctions)

	engine.GET("/openapi", api.HandleOpenApi)
	engine.Run(":" + port)
}
