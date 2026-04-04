package ambulance_wl

import (
	"net/http"

	"github.com/TechOctopus/davgus-ambulance-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implPatientsAPI struct {
}

func NewPatientsApi() PatientsAPI {
	return &implPatientsAPI{}
}

// CreatePatient - Saves new patient definition
func (o *implPatientsAPI) CreatePatient(c *gin.Context) {
	value, exists := c.Get("db_service_patients")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db not found", "error": "db not found"})
		return
	}

	db, ok := value.(db_service.DbService[Patient])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db context is not of required type", "error": "cannot cast db context"})
		return
	}

	patient := Patient{}
	err := c.BindJSON(&patient)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid request body", "error": err.Error()})
		return
	}

	if patient.Id == "" {
		patient.Id = uuid.New().String()
	}

	err = db.CreateDocument(c, patient.Id, &patient)

	switch err {
	case nil:
		c.JSON(http.StatusCreated, patient)
	case db_service.ErrConflict:
		c.JSON(http.StatusConflict, gin.H{"status": "Conflict", "message": "Patient already exists", "error": err.Error()})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"status": "Bad Gateway", "message": "Failed to create patient in database", "error": err.Error()})
	}
}

// DeletePatient - Deletes specific patient
func (o *implPatientsAPI) DeletePatient(c *gin.Context) {
	value, exists := c.Get("db_service_patients")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service not found"})
		return
	}

	db, ok := value.(db_service.DbService[Patient])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service context is not of type DbService"})
		return
	}

	patientId := c.Param("patientId")
	err := db.DeleteDocument(c, patientId)

	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": "Patient not found", "error": err.Error()})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"status": "Bad Gateway", "message": "Failed to delete patient from database", "error": err.Error()})
	}
}

// GetPatient - Provides details about patient
func (o *implPatientsAPI) GetPatient(c *gin.Context) {
	updatePatientFunc(c, func(c *gin.Context, patient *Patient) (*Patient, interface{}, int) {
		return nil, patient, http.StatusOK
	})
}

// GetPatients - Provides the list of patients
func (o *implPatientsAPI) GetPatients(c *gin.Context) {
	value, exists := c.Get("db_service_patients")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service_patients not found"})
		return
	}

	db, ok := value.(db_service.DbService[Patient])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service_patients context is not of required type"})
		return
	}

	patients, err := db.FindAllDocuments(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "Bad Gateway", "message": "Failed to load patients from database"})
		return
	}

	c.JSON(http.StatusOK, patients)
}

// UpdatePatient - Updates specific patient
func (o *implPatientsAPI) UpdatePatient(c *gin.Context) {
	updatePatientFunc(c, func(c *gin.Context, patient *Patient) (*Patient, interface{}, int) {
		var parsedPatient Patient
		if err := c.ShouldBindJSON(&parsedPatient); err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}

		parsedPatient.Id = patient.Id

		return &parsedPatient, parsedPatient, http.StatusOK
	})
}
