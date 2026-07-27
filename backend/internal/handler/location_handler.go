package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/internal/model"
	"backend/internal/repository"
)

// LocationGetter — methods expected from the repository by the handler
type LocationGetter interface {
	GetByID(ctx context.Context, id int64) (*model.Location, error)
	GetByMapillaryID(ctx context.Context, mid int64) (*model.Location, error)
	GetRandom(ctx context.Context) (*model.Location, error)
	GetRandomByCountryCode(ctx context.Context, cc string) (*model.Location, error)
	PatchByID(ctx context.Context, id int64, patch *model.Location) error
	PatchByMapillaryID(ctx context.Context, mid int64, patch *model.Location) error
	BulkCreate(ctx context.Context, locations []*model.Location) (int, error)
	DeleteByID(ctx context.Context, id int64) error
	DeleteByMapillaryID(ctx context.Context, mid int64) error
}

type LocationHandler struct {
	repo LocationGetter
}

func NewLocationHandler(repo LocationGetter) *LocationHandler {
	return &LocationHandler{repo: repo}
}

// validationStatus translates a model validation error into an HTTP status code and message
func validationStatus(err error) (int, string) {
	switch {
	case errors.Is(err, model.ErrMissingFields):
		return http.StatusBadRequest, "missing required fields"
	case errors.Is(err, model.ErrInvalidLatitude):
		return http.StatusBadRequest, "invalid latitude"
	case errors.Is(err, model.ErrInvalidLongitude):
		return http.StatusBadRequest, "invalid longitude"
	case errors.Is(err, model.ErrInvalidCountryCode):
		return http.StatusBadRequest, "invalid country code"
	default:
		return http.StatusBadRequest, "invalid payload"
	}
}

// GET /locations/:id
func (h *LocationHandler) GetLocationByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	loc, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
			return
		}
		log.Printf("Error getting location %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting location"})
		return
	}

	c.JSON(http.StatusOK, loc)
}

// GET /locations/mapillary/:mid
func (h *LocationHandler) GetLocationByMID(c *gin.Context) {
	ctx := c.Request.Context()

	mid, err := strconv.ParseInt(c.Param("mid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mapillary id"})
		return
	}

	loc, err := h.repo.GetByMapillaryID(ctx, mid)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
			return
		}
		log.Printf("Error getting location by mapillary_id %d: %v", mid, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting location"})
		return
	}

	c.JSON(http.StatusOK, loc)
}

// GET /locations/random
func (h *LocationHandler) GetRandomLocation(c *gin.Context) {
	ctx := c.Request.Context()

	loc, err := h.repo.GetRandom(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no locations available"})
			return
		}
		log.Printf("Error getting random location: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting random location"})
		return
	}

	c.JSON(http.StatusOK, loc)
}

// GET /locations/random/country/:cc
func (h *LocationHandler) GetRandomLocationByCC(c *gin.Context) {
	ctx := c.Request.Context()
	cc := c.Param("cc")

	if len(cc) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid country code"})
		return
	}

	loc, err := h.repo.GetRandomByCountryCode(ctx, cc)
	if err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "no locations for this country"})
			return
		}
		log.Printf("Error getting random location for country %s: %v", cc, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting random location"})
		return
	}

	c.JSON(http.StatusOK, loc)
}

// PATCH /locations/:id
func (h *LocationHandler) UpdateLocationByID(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var patch model.Location
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if err := patch.ValidateForPatch(); err != nil {
		status, msg := validationStatus(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	if err := h.repo.PatchByID(ctx, id, &patch); err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
			return
		}
		log.Printf("Error patching location %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating location"})
		return
	}

	c.Status(http.StatusNoContent)
}

// PATCH /locations/mapillary/:mid
func (h *LocationHandler) UpdateLocationByMID(c *gin.Context) {
	ctx := c.Request.Context()

	mid, err := strconv.ParseInt(c.Param("mid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mapillary id"})
		return
	}

	var patch model.Location
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if err := patch.ValidateForPatch(); err != nil {
		status, msg := validationStatus(err)
		c.JSON(status, gin.H{"error": msg})
		return
	}

	if err := h.repo.PatchByMapillaryID(ctx, mid, &patch); err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
			return
		}
		log.Printf("Error patching location by mapillary_id %d: %v", mid, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error updating location"})
		return
	}

	c.Status(http.StatusNoContent)
}

// POST /locations
func (h *LocationHandler) AddLocations(c *gin.Context) {
	ctx := c.Request.Context()

	var inputs []*model.Location
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if len(inputs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty location list"})
		return
	}

	// strict validation: any invalid entry makes the entire request fail
	for i, in := range inputs {
		if err := in.ValidateForCreate(); err != nil {
			status, msg := validationStatus(err)
			c.JSON(status, gin.H{"error": msg, "index": i})
			return
		}
	}

	inserted, err := h.repo.BulkCreate(ctx, inputs)
	if err != nil {
		log.Printf("Error creating locations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating locations"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"inserted": inserted,
		"total":    len(inputs),
	})
}

// DELETE /locations/:id
func (h *LocationHandler) DeleteLocation(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.DeleteByID(ctx, id); err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
			return
		}
		log.Printf("Error deleting location %d: %v", id, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting location"})
		return
	}

	c.Status(http.StatusNoContent)
}

// DELETE /locations/mapillary/:mid
func (h *LocationHandler) DeleteLocationByMID(c *gin.Context) {
	ctx := c.Request.Context()

	mid, err := strconv.ParseInt(c.Param("mid"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mapillary id"})
		return
	}

	if err := h.repo.DeleteByMapillaryID(ctx, mid); err != nil {
		if errors.Is(err, repository.ErrLocationNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "location not found"})
			return
		}
		log.Printf("Error deleting location by mapillary_id %d: %v", mid, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error deleting location"})
		return
	}

	c.Status(http.StatusNoContent)
}