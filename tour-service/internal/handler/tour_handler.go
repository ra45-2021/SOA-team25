package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"tour-service/internal/model"
	"tour-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type TourHandler struct {
	svc         *service.TourService
	minioClient *minio.Client
	bucketName  string
}

func NewTourHandler(svc *service.TourService, mc *minio.Client, bucket string) *TourHandler {
	return &TourHandler{svc: svc, minioClient: mc, bucketName: bucket}
}

func (h *TourHandler) CreateTour(c *gin.Context) {
	var tour model.Tour
	if err := c.ShouldBindJSON(&tour); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.CreateTour(c.Request.Context(), &tour); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tour)
}

func (h *TourHandler) AddCheckpoint(c *gin.Context) {
	tourID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}
	defer file.Close()

	objectName := fmt.Sprintf("tour-%d/%s", tourID, header.Filename)
	_, err = h.minioClient.PutObject(context.Background(), h.bucketName, objectName, file, header.Size, minio.PutObjectOptions{
		ContentType: header.Header.Get("Content-Type"),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "upload failed"})
		return
	}

	cp := model.Checkpoint{
		TourID:      uint(tourID),
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
		Latitude:    parseFloat(c.PostForm("latitude")),
		Longitude:   parseFloat(c.PostForm("longitude")),
		ImageURL:    fmt.Sprintf("http://localhost:9000/%s/%s", h.bucketName, objectName),
	}

	if err := h.svc.AddCheckpoint(&cp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cp)
}

func (h *TourHandler) GetMyTours(c *gin.Context) {
	authorIDStr := c.Query("authorId")
	authorID, err := strconv.ParseInt(authorIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid authorId"})
		return
	}

	tours, err := h.svc.GetAuthorTours(authorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tours)
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}