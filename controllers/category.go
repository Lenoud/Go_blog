package controllers

import (
	"blog/models"
	"blog/storage/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=50"`
	Slug string `json:"slug" binding:"required,min=2,max=50"`
}

func CreateCategory(c *gin.Context) {
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 检查分类名是否已存在
	var existingCategory models.Category
	if err := mysql.DB.Where("name = ? OR slug = ?", req.Name, req.Slug).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类名或别名已存在"})
		return
	}

	category := models.Category{
		Name: req.Name,
		Slug: req.Slug,
	}

	if err := mysql.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建分类失败"})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := mysql.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取分类列表失败"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func GetCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category
	if err := mysql.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if err := mysql.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "分类不存在"})
		return
	}

	// 检查新的分类名或别名是否与其他分类冲突
	var existingCategory models.Category
	if err := mysql.DB.Where("(name = ? OR slug = ?) AND id != ?", req.Name, req.Slug, id).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "分类名或别名已存在"})
		return
	}

	category.Name = req.Name
	category.Slug = req.Slug

	if err := mysql.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新分类失败"})
		return
	}

	c.JSON(http.StatusOK, category)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := mysql.DB.Delete(&models.Category{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除分类失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}
