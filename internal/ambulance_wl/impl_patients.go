package ambulance_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implPatientsAPI struct {
}

func NewPatientsApi() PatientsAPI {
	return &implPatientsAPI{}
}

func (o implPatientsAPI) CreatePatient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPatientsAPI) DeletePatient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPatientsAPI) GetPatient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPatientsAPI) GetPatients(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPatientsAPI) UpdatePatient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
