package ambulance_wl

import (
	"net/http"

	"github.com/TechOctopus/davgus-ambulance-webapi/internal/db_service"
	"github.com/gin-gonic/gin"
)

type documentUpdater[T any] func(ctx *gin.Context, document *T) (updatedDocument *T, responseContent interface{}, status int)

func updateDocumentFunc[T any](ctx *gin.Context, dbKey string, idParam string, updater documentUpdater[T]) {
	value, exists := ctx.Get(dbKey)
	if !exists {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": dbKey + " not found",
				"error":   dbKey + " not found",
			})
		return
	}

	db, ok := value.(db_service.DbService[T])
	if !ok {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{
				"status":  "Internal Server Error",
				"message": dbKey + " context is not of required type",
				"error":   "cannot cast context to db_service.DbService",
			})
		return
	}

	documentId := ctx.Param(idParam)

	document, err := db.FindDocument(ctx, documentId)

	switch err {
	case nil:
		// continue
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Document not found",
				"error":   err.Error(),
			},
		)
		return
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to load document from database",
				"error":   err.Error(),
			})
		return
	}

	updatedDocument, responseObject, status := updater(ctx, document)

	if updatedDocument != nil {
		err = db.UpdateDocument(ctx, documentId, updatedDocument)
	} else {
		err = nil // redundant but for clarity
	}

	switch err {
	case nil:
		if responseObject != nil {
			ctx.JSON(status, responseObject)
		} else {
			ctx.AbortWithStatus(status)
		}
	case db_service.ErrNotFound:
		ctx.JSON(
			http.StatusNotFound,
			gin.H{
				"status":  "Not Found",
				"message": "Document was deleted while processing the request",
				"error":   err.Error(),
			},
		)
	default:
		ctx.JSON(
			http.StatusBadGateway,
			gin.H{
				"status":  "Bad Gateway",
				"message": "Failed to update document in database",
				"error":   err.Error(),
			})
	}
}

func updateDepartmentFunc(ctx *gin.Context, updater documentUpdater[Department]) {
	updateDocumentFunc(ctx, "db_service_departments", "departmentId", updater)
}

func updatePatientFunc(ctx *gin.Context, updater documentUpdater[Patient]) {
	updateDocumentFunc(ctx, "db_service_patients", "patientId", updater)
}

func updatePlacementFunc(ctx *gin.Context, updater documentUpdater[Placement]) {
	updateDocumentFunc(ctx, "db_service_placements", "placementId", updater)
}
