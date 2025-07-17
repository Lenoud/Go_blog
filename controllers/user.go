package controllers

import (
	"blog/models"
	"blog/storage/mysql"
	"blog/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Warn("注册请求参数无效", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := mysql.DB.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		utils.Warn("用户名或邮箱已存在",
			zap.String("username", req.Username),
			zap.String("email", req.Email))
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或邮箱已存在"})
		return
	}

	// 创建新用户
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error("密码加密失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	user := models.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err := mysql.DB.Create(&user).Error; err != nil {
		utils.Error("创建用户失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	utils.Info("用户注册成功",
		zap.String("username", user.Username),
		zap.String("email", user.Email))

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Warn("登录请求参数无效", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求参数"})
		return
	}

	// 查找用户
	var user models.User
	if err := mysql.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		utils.Warn("用户不存在", zap.String("username", req.Username))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.Warn("密码错误", zap.String("username", req.Username))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.Error("生成token失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "登录失败"})
		return
	}

	utils.Info("用户登录成功",
		zap.String("username", user.Username),
		zap.Uint("user_id", user.ID))

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
