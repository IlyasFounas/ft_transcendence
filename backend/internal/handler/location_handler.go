package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"backend/internal/models"
	"backend/internal/repository"
)

// LocationGetter - interface that defines the methods for interacting with location data in the repository.
type LocationGetter interface {
	GetByID(ctx context.Context, id int64) (*models.Location, error)
	GetByMapillaryID(ctx context.Context, mid int64) (*models.Location, error)
	GetRandom(ctx context.Context) (*models.Location, error)
	GetRandomByCountryCode(ctx context.Context, cc string) (*models.Location, error)
	PatchByID(ctx context.Context, id int64, patch models.LocationPatch) error
	PatchByMapillaryID(ctx context.Context, mid int64, patch models.LocationPatch) error
	BulkCreate(ctx context.Context, locations []*models.Location) (int, error)
	DeleteByID(ctx context.Context, id int64) error
	DeleteByMapillaryID(ctx context.Context, mid int64) error
}

type LocationHandler struct {
	repo LocationGetter
}

func NewLocationHandler(repo LocationGetter) *LocationHandler {
	return &LocationHandler{repo: repo}
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

	var patch models.LocationPatch
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if patch.Latitude != nil && (*patch.Latitude < -90 || *patch.Latitude > 90) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid latitude"})
		return
	}
	if patch.Longitude != nil && (*patch.Longitude < -180 || *patch.Longitude > 180) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid longitude"})
		return
	}
	if patch.CountryCode != nil && len(*patch.CountryCode) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid country code"})
		return
	}

	if err := h.repo.PatchByID(ctx, id, patch); err != nil {
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

	var patch models.LocationPatch
	if err := c.ShouldBindJSON(&patch); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if patch.Latitude != nil && (*patch.Latitude < -90 || *patch.Latitude > 90) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid latitude"})
		return
	}
	if patch.Longitude != nil && (*patch.Longitude < -180 || *patch.Longitude > 180) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid longitude"})
		return
	}
	if patch.CountryCode != nil && len(*patch.CountryCode) != 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid country code"})
		return
	}

	if err := h.repo.PatchByMapillaryID(ctx, mid, patch); err != nil {
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

	var inputs []models.LocationInput
	if err := c.ShouldBindJSON(&inputs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	if len(inputs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty location list"})
		return
	}

	locations := make([]*models.Location, 0, len(inputs))
	var skipped int
	for _, in := range inputs {
		if in.Lat < -90 || in.Lat > 90 || in.Long < -180 || in.Long > 180 {
			skipped++
			continue
		}
		locations = append(locations, &models.Location{
			MID:  in.MID,
			Lat:  in.Lat,
			Long: in.Long,
			CC:   in.CC,
		})
	}

	if len(locations) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no valid entries"})
		return
	}

	inserted, err := h.repo.BulkCreate(ctx, locations)
	if err != nil {
		log.Printf("Error creating locations: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error creating locations"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"inserted": inserted,
		"skipped":  skipped,
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

	if err := h.repo.Delete(ctx, id); err != nil {
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