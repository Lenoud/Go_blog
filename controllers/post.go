package controllers

import (
	"blog/models"
	"blog/storage/mongodb"
	"blog/storage/mysql"
	"blog/storage/redis"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type PostRequest struct {
	Title      string `json:"title" binding:"required,min=2,max=200"`
	Slug       string `json:"slug" binding:"required,min=2,max=200"`
	Content    string `json:"content" binding:"required"`
	Summary    string `json:"summary" binding:"max=500"`
	CategoryID uint   `json:"category_id" binding:"required"`
	Status     string `json:"status" binding:"required,oneof=draft published"`
}

func CreatePost(c *gin.Context) {
	var req PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查分类是否存在
	var category models.Category
	if err := mysql.DB.First(&category, req.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类不存在"})
		return
	}

	// 检查别名是否已存在
	var existingPost models.Post
	if err := mysql.DB.Where("slug = ?", req.Slug).First(&existingPost).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章别名已存在"})
		return
	}

	userID := c.GetUint("userID")
	post := models.Post{
		Title:      req.Title,
		Slug:       req.Slug,
		Content:    req.Content,
		Summary:    req.Summary,
		AuthorID:   userID,
		CategoryID: req.CategoryID,
		Status:     req.Status,
	}

	if err := mysql.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	// 将文章存入MongoDB
	if _, err := mongodb.Database.Collection("posts").InsertOne(context.Background(), post); err != nil {
		// 记录插入MongoDB失败的日志，但不中断流程
	}

	// 清除相关缓存
	ctx := context.Background()
	redis.Client.Del(ctx, "posts:list")
	redis.Client.Del(ctx, fmt.Sprintf("category:%d:posts", req.CategoryID))

	c.JSON(http.StatusCreated, post)
}

func GetPosts(c *gin.Context) {
	ctx := context.Background()
	cacheKey := "posts:list"

	// 尝试从缓存获取
	cachedPosts, err := redis.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var posts []models.Post
		if err := json.Unmarshal([]byte(cachedPosts), &posts); err == nil {
			c.JSON(http.StatusOK, posts)
			return
		}
	}

	// 从数据库获取
	var posts []models.Post
	if err := mysql.DB.Preload("Author").Preload("Category").Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}

	// 存入缓存
	if postsJSON, err := json.Marshal(posts); err == nil {
		redis.Client.Set(ctx, cacheKey, postsJSON, 1*time.Hour)
	}

	c.JSON(http.StatusOK, posts)
}

func GetPost(c *gin.Context) {
	slug := c.Param("slug")
	ctx := context.Background()
	cacheKey := fmt.Sprintf("post:%s", slug)

	// 尝试从缓存获取
	cachedPost, err := redis.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		var post models.Post
		if err := json.Unmarshal([]byte(cachedPost), &post); err == nil {
			c.JSON(http.StatusOK, post)
			return
		}
	}

	// 从数据库获取
	var post models.Post
	if err := mysql.DB.Preload("Author").Preload("Category").Where("slug = ?", slug).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 存入缓存
	if postJSON, err := json.Marshal(post); err == nil {
		redis.Client.Set(ctx, cacheKey, postJSON, 1*time.Hour)
	}

	c.JSON(http.StatusOK, post)
}

func UpdatePost(c *gin.Context) {
	slug := c.Param("slug")
	var req PostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var post models.Post
	if err := mysql.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	userID := c.GetUint("userID")
	if post.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限修改此文章"})
		return
	}

	// 检查分类是否存在
	var category models.Category
	if err := mysql.DB.First(&category, req.CategoryID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类不存在"})
		return
	}

	// 检查新的别名是否与其他文章冲突
	if req.Slug != slug {
		var existingPost models.Post
		if err := mysql.DB.Where("slug = ? AND id != ?", req.Slug, post.ID).First(&existingPost).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "文章别名已存在"})
			return
		}
	}

	post.Title = req.Title
	post.Slug = req.Slug
	post.Content = req.Content
	post.Summary = req.Summary
	post.CategoryID = req.CategoryID
	post.Status = req.Status

	if err := mysql.DB.Save(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}

	// 更新MongoDB中的文章
	filter := bson.M{"slug": slug}
	update := bson.M{"$set": post}
	if _, err := mongodb.Database.Collection("posts").UpdateOne(context.Background(), filter, update); err != nil {
		// 记录更新MongoDB失败的日志
	}

	// 清除相关缓存
	ctx := context.Background()
	redis.Client.Del(ctx, "posts:list")
	redis.Client.Del(ctx, fmt.Sprintf("post:%s", slug))
	redis.Client.Del(ctx, fmt.Sprintf("category:%d:posts", post.CategoryID))

	c.JSON(http.StatusOK, post)
}

func DeletePost(c *gin.Context) {
	slug := c.Param("slug")
	var post models.Post
	if err := mysql.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	userID := c.GetUint("userID")
	if post.AuthorID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此文章"})
		return
	}

	if err := mysql.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}

	// 从MongoDB删除文章
	if _, err := mongodb.Database.Collection("posts").DeleteOne(context.Background(), bson.M{"slug": slug}); err != nil {
		// 记录从MongoDB删除失败的日志
	}

	// 清除相关缓存
	ctx := context.Background()
	redis.Client.Del(ctx, "posts:list")
	redis.Client.Del(ctx, fmt.Sprintf("post:%s", slug))
	redis.Client.Del(ctx, fmt.Sprintf("category:%d:posts", post.CategoryID))

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
