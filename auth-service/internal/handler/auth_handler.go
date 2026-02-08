package handler

import (
	"context"
	"net/http"
	"time"

	"auth-service/internal/model"
	"auth-service/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	out, err := h.svc.Register(ctx, req)
	if err != nil {
		status := http.StatusBadRequest
		if err == service.ErrHashFailed || err == service.ErrIDGenFailed || err == service.ErrTokenFailed {
			status = http.StatusInternalServerError
		}
		if err == service.ErrInvalidCreds {
			status = http.StatusUnauthorized
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, out)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	out, err := h.svc.Login(ctx, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (h *AuthHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, model.MeResponse{
		ID:       c.MustGet("id").(int64),
		Username: c.MustGet("username").(string),
		Role:     c.MustGet("role").(string),
	})
}
