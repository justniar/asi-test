package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/your_project/models"
	"github.com/your_project/utils"
	"gorm.io/gorm"
)

type ClientHandler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func (h *ClientHandler) CreateClient(c *gin.Context) {
	var input models.Client

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("client_logo")
	if err == nil {
		filePath, err := utils.SaveFile(file, header.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
			return
		}
		input.ClientLogo = filePath
	} else {
		input.ClientLogo = "/uploads/no-image.jpg"
	}

	input.Slug = utils.GenerateSlug(input.Name)
	input.CreatedAt = time.Now()

	h.DB.Create(&input)

	clientJSON, _ := json.Marshal(input)
	h.Redis.Set(context.Background(), fmt.Sprintf("client:%s", input.Slug), clientJSON, time.Hour)

	c.JSON(http.StatusCreated, input)
}

func (h *ClientHandler) UpdateClient(c *gin.Context) {
	slug := c.Param("slug")
	var client models.Client

	if err := h.DB.Where("slug = ?", slug).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Client not found"})
		return
	}

	var input models.Client
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("client_logo")
	if err == nil {
		filePath, err := utils.SaveFile(file, header.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
			return
		}
		client.ClientLogo = filePath
	}

	client.Name = input.Name
	client.Address = input.Address
	client.PhoneNumber = input.PhoneNumber
	client.City = input.City
	client.UpdatedAt = time.Now()

	h.DB.Save(&client)

	h.Redis.Del(context.Background(), fmt.Sprintf("client:%s", slug))

	clientJSON, _ := json.Marshal(client)
	h.Redis.Set(context.Background(), fmt.Sprintf("client:%s", slug), clientJSON, time.Hour)

	c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) GetClient(c *gin.Context) {
	slug := c.Param("slug")
	cacheKey := fmt.Sprintf("client:%s", slug)

	cachedClient, err := h.Redis.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var client models.Client
		json.Unmarshal([]byte(cachedClient), &client)
		c.JSON(http.StatusOK, client)
		return
	}

	var client models.Client
	if err := h.DB.Where("slug = ? AND deleted_at IS NULL", slug).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Client not found"})
		return
	}

	clientJSON, _ := json.Marshal(client)
	h.Redis.Set(context.Background(), cacheKey, clientJSON, time.Hour)

	c.JSON(http.StatusOK, client)
}

func (h *ClientHandler) DeleteClient(c *gin.Context) {
	slug := c.Param("slug")
	var client models.Client

	if err := h.DB.Where("slug = ?", slug).First(&client).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Client not found"})
		return
	}

	now := time.Now()
	client.DeletedAt = &now
	h.DB.Save(&client)

	h.Redis.Del(context.Background(), fmt.Sprintf("client:%s", slug))

	c.JSON(http.StatusOK, gin.H{"message": "Client deleted"})
}
