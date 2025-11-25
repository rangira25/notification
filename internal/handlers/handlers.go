package handlers

import (
	
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/rangira25/notification/internal/models"
	"github.com/rangira25/notification/internal/tasks"
)

type Handler struct {
	DB     *gorm.DB
	Redis  *redis.Client
	Asynq  *asynq.Client
	Logger *zap.Logger
}

func NewHandler(db *gorm.DB, r *redis.Client, a *asynq.Client, l *zap.Logger) *Handler {
	return &Handler{DB: db, Redis: r, Asynq: a, Logger: l}
}

type CreateUserReq struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	ReqID    string `json:"req_id" binding:"required"`
}

// CreateUser creates a user and enqueues a welcome email (idempotent)
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Warn("bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	idemKey := "idem:req:" + req.ReqID
	ok, err := h.Redis.SetNX(ctx, idemKey, "1", 24*time.Hour).Result()
	if err != nil {
		h.Logger.Error("redis error setnx", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	if !ok {
		c.JSON(http.StatusConflict, gin.H{"error": "duplicate request"})
		return
	}

	user := models.User{
		FullName: req.FullName,
		Email:    req.Email,
	}
	if err := h.DB.Create(&user).Error; err != nil {
		h.Redis.Del(ctx, idemKey)
		h.Logger.Error("db create user failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}

	// enqueue welcome email
	payload := &tasks.WelcomeEmailPayload{
		UserID: user.ID,
		Email:  user.Email,
		ReqID:  req.ReqID,
	}
	if _, err := tasks.EnqueueWelcomeEmail(h.Asynq, payload); err != nil {
		h.Redis.Del(ctx, idemKey)
		h.Logger.Error("enqueue failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "enqueue failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": user.ID})
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
