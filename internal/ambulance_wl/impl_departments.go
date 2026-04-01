package ambulance_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implDepartmentsAPI struct {
}

func NewDepartmentsApi() DepartmentsAPI {
	return &implDepartmentsAPI{}
}

func (o implDepartmentsAPI) GetDepartment(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implDepartmentsAPI) GetDepartments(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implDepartmentsAPI) UpdateDepartment(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
