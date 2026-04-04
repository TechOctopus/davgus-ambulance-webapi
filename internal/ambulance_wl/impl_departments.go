package ambulance_wl

import (
	"net/http"

	"github.com/TechOctopus/davgus-ambulance-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implDepartmentsAPI struct {
}

func NewDepartmentsApi() DepartmentsAPI {
	return &implDepartmentsAPI{}
}

// GetDepartment - Provides details about department
func (o *implDepartmentsAPI) GetDepartment(c *gin.Context) {
	updateDepartmentFunc(c, func(c *gin.Context, department *Department) (*Department, interface{}, int) {
		return nil, department, http.StatusOK
	})
}

// GetDepartments - Provides the list of departments
func (o *implDepartmentsAPI) GetDepartments(c *gin.Context) {
	value, exists := c.Get("db_service_departments")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service_departments not found",
				"error":   "db_service_departments not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Department])
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service_departments context is not of required type",
				"error":   "cannot cast db_service_departments context",
			})
		return
	}

	departments, err := db.FindAllDocuments(c)
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to load departments from database",
				"error":   err.Error(),
			})
		return
	}

	c.JSON(http.StatusOK, departments)
}

// UpdateDepartment - Updates specific department
func (o *implDepartmentsAPI) UpdateDepartment(c *gin.Context) {
	updateDepartmentFunc(c, func(c *gin.Context, department *Department) (*Department, interface{}, int) {
		var parsedDepartment Department
		if err := c.ShouldBindJSON(&parsedDepartment); err != nil {
			return nil, gin.H{
				"status":  http.StatusBadRequest,
				"message": "Invalid request body",
				"error":   err.Error(),
			}, http.StatusBadRequest
		}
		
		// Map the fields from parsed structure to the document retrieved
		// Assuming UUID shouldn't change
		parsedDepartment.Id = department.Id

		return &parsedDepartment, parsedDepartment, http.StatusOK
	})
}

func (o implDepartmentsAPI) CreateDepartment(c *gin.Context) {
	value, exists := c.Get("db_service_departments")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db not found",
				"error":   "db not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Department])
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db context is not of required type",
				"error":   "cannot cast db context to db_service.DbService",
			})
		return
	}

	department := Department{}
	err := c.BindJSON(&department)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{
				"status":  "Bad Request",
				"message": "Invalid request body",
				"error":   err.Error(),
			})
		return
	}

	if department.Id == "" {
		department.Id = uuid.New().String()
	}

	err = db.CreateDocument(c, department.Id, &department)

	switch err {
	case nil:
		c.JSON(
			http.StatusCreated,
			department,
		)
	case db_service.ErrConflict:
		c.JSON(
			http.StatusConflict,
			gin.H{
				"status":  "Conflict",
				"message": "Department already exists",
				"error":   err.Error(),
			},
		)
	default:
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to create department in database",
				"error":   err.Error(),
			},
		)
	}
}

func (o implDepartmentsAPI) DeleteDepartment(c *gin.Context) {
	value, exists := c.Get("db_service_departments")
	if !exists {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service not found",
				"error":   "db_service not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[Department])
	if !ok {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": "db_service context is not of type db_service.DbService",
				"error":   "cannot cast db_service context to db_service.DbService",
			})
		return
	}

	departmentId := c.Param("departmentId")
	err := db.DeleteDocument(c, departmentId)

	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Department not found",
				"error":   err.Error(),
			},
		)
	default:
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to delete department from database",
				"error":   err.Error(),
			})
	}
}
