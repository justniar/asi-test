package handler

import (
	models "ASI/model"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type ClientHandler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func (h *ClientHandler) GetClients(c *gin.context) {
	var clients []models.Client
	h.DB.where("deleted_at IS NULL").Find(&clients)
	c.JSON(http.StatusOk, clients)
}

func (h *ClientHandler) GetClient(c *gin.context) {
	slug := c.Param("slug")
	cacheKey := fmt.Sprintf("client:%s", slug)

	cachedClient, err := h.Redis.Get(ctx, cacheKey).Result()
	if err == nil {
		var client models.Client
		json.Unmarshall([]byte(cachedClient), &client)
		c.JSON(http.StatusOK, client)
		return
	}

	var client models.Client
	if err := h.DB.where("slug = ? AND deleted_at is null", slug).first(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message: client not found"})
		return
	}

	clientJSON, _ := json.Marshal(client)
	h.Redis.Set(ctx, cacheKey, clientJSON, time.Hour)
	c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) CreateClient(c *gin.Context) {
	var input models.Client
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.Slug = utils.GeneratedSlug(input.Name)
	input.CreatedAt = time.Now()

	h.DB.Create(&input)

	clientJSON, _ := json.Marshal(input)
	h.Redis.Set(ctx, fmt.Sprintf("client:%s", input.Slug), clientJSON, time.Hour)

	c.JSON(http.StatusCreated, input)
}

func (h *ClientHandler) UpdateClient(c *gin.Context) {

}
