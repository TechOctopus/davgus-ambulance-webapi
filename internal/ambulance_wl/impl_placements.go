package ambulance_wl

import (
	"net/http"

	"github.com/TechOctopus/davgus-ambulance-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implPlacementsAPI struct {
}

func NewPlacementsApi() PlacementsAPI {
	return &implPlacementsAPI{}
}

func (o implPlacementsAPI) CreatePlacement(c *gin.Context) {
	value, exists := c.Get("db_service_placements")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db not found", "error": "db not found"})
		return
	}

	db, ok := value.(db_service.DbService[Placement])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db context is not of required type", "error": "cannot cast db context"})
		return
	}

	placement := Placement{}
	err := c.BindJSON(&placement)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid request body", "error": err.Error()})
		return
	}

	if placement.Id == "" {
		placement.Id = uuid.New().String()
	}

	err = db.CreateDocument(c, placement.Id, &placement)

	switch err {
	case nil:
		c.JSON(http.StatusCreated, placement)
	case db_service.ErrConflict:
		c.JSON(http.StatusConflict, gin.H{"status": "Conflict", "message": "Placement already exists", "error": err.Error()})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"status": "Bad Gateway", "message": "Failed to create placement in database", "error": err.Error()})
	}
}

func (o implPlacementsAPI) DeletePlacement(c *gin.Context) {
	value, exists := c.Get("db_service_placements")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service not found"})
		return
	}

	db, ok := value.(db_service.DbService[Placement])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service context is not of type DbService"})
		return
	}

	placementId := c.Param("placementId")
	err := db.DeleteDocument(c, placementId)

	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": "Placement not found", "error": err.Error()})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"status": "Bad Gateway", "message": "Failed to delete from database", "error": err.Error()})
	}
}

func (o implPlacementsAPI) GetPlacement(c *gin.Context) {
	updatePlacementFunc(c, func(c *gin.Context, placement *Placement) (*Placement, interface{}, int) {
		return nil, placement, http.StatusOK
	})
}

// GetPlacementForPatient - Custom method for fetching placement based on patient 
func (o implPlacementsAPI) GetPlacementForPatient(c *gin.Context) {
	// TODO: implement custom fetch or use FindAllDocuments and filter. 
	// As the assignment didn't explicitly implement complex queries, we use FindAll for now
	value, exists := c.Get("db_service_placements")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service_placements not found"})
		return
	}
	db, ok := value.(db_service.DbService[Placement])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service_placements context is not of required type"})
		return
	}
	placements, err := db.FindAllDocuments(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "Bad Gateway", "message": "Failed to load placements from database"})
		return
	}
	
	patientId := c.Param("patientId")
	for _, p := range placements {
		if p.PatientId == patientId {
			c.JSON(http.StatusOK, p)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": "Placement not found for patient"})
}

func (o implPlacementsAPI) GetPlacements(c *gin.Context) {
	value, exists := c.Get("db_service_placements")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service_placements not found"})
		return
	}

	db, ok := value.(db_service.DbService[Placement])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service_placements context is not of required type"})
		return
	}

	placements, err := db.FindAllDocuments(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "Bad Gateway", "message": "Failed to load placements from database"})
		return
	}

	c.JSON(http.StatusOK, placements)
}

func (o implPlacementsAPI) UpdatePlacement(c *gin.Context) {
	updatePlacementFunc(c, func(c *gin.Context, placement *Placement) (*Placement, interface{}, int) {
		var parsedPlacement Placement
		if err := c.ShouldBindJSON(&parsedPlacement); err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}

		parsedPlacement.Id = placement.Id

		return &parsedPlacement, parsedPlacement, http.StatusOK
	})
}
