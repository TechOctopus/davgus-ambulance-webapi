package ambulance_wl

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implPlacementsAPI struct {
}

func NewPlacementsApi() PlacementsAPI {
	return &implPlacementsAPI{}
}

func (o implPlacementsAPI) CreatePlacement(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPlacementsAPI) DeletePlacement(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPlacementsAPI) GetPlacement(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPlacementsAPI) GetPlacementForPatient(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPlacementsAPI) GetPlacements(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implPlacementsAPI) UpdatePlacement(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
