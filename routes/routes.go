package routes

import (
	"blog/controllers"
	"blog/middleware"

	"github.com/gin-gonic/gin"
)

// Init 初始化路由
func Init() *gin.Engine {
	r := gin.New()

	// 使用中间件
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	// API路由组
	api := r.Group("/api")
	{
		// 用户认证
		api.POST("/register", controllers.Register)
		api.POST("/login", controllers.Login)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(middleware.Auth())
		{
			// 文章管理
			auth.POST("/posts", controllers.CreatePost)
			auth.PUT("/posts/:slug", controllers.UpdatePost)
			auth.DELETE("/posts/:slug", controllers.DeletePost)

			// 分类管理
			auth.POST("/categories", controllers.CreateCategory)
			auth.PUT("/categories/:id", controllers.UpdateCategory)
			auth.DELETE("/categories/:id", controllers.DeleteCategory)
		}

		// 公开路由
		api.GET("/posts", controllers.GetPosts)
		api.GET("/posts/:slug", controllers.GetPost)
		api.GET("/categories", controllers.GetCategories)
		api.GET("/categories/:id", controllers.GetCategory)
	}

	return r
}
