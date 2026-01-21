package api

import (
	"context"
	"golang-rest-api-template/pkg/cache"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/middleware"
	"time"

	docs "golang-rest-api-template/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"golang.org/x/time/rate"
)

func ContextMiddleware(bookRepository BookRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("appCtx", bookRepository)
		c.Next()
	}
}

func NewRouter(logger *zap.Logger, mongoCollection *mongo.Collection, db database.Database, redisClient cache.Cache, ctx *context.Context) *gin.Engine {
	bookRepository := NewBookRepository(db, redisClient, ctx)
	userRepository := NewUserRepository(db, ctx)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.APM()) // APM middleware for tracing
	r.Use(ContextMiddleware(bookRepository))

	//r.Use(gin.Logger())
	r.Use(middleware.RequestResponseLogger(logger, mongoCollection))
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.Security())
		r.Use(middleware.Xss())
	}
	r.Use(middleware.Cors())
	r.Use(middleware.RateLimiter(rate.Every(1*time.Minute), 60)) // 60 requests per minute

	docs.SwaggerInfo.BasePath = "/api/v1"

	// Base API group with common middleware
	v1 := r.Group("/api/v1")
	v1.Use(middleware.APIKeyAuth())

	// Public routes (no auth required)
	v1.GET("/", bookRepository.Healthcheck)

	// Auth routes (login/register)
	auth := v1.Group("/")
	{
		auth.POST("/login", userRepository.LoginHandler)
		auth.POST("/register", userRepository.RegisterHandler)
	}

	// Protected routes (JWT required)
	protected := v1.Group("/books")
	protected.Use(middleware.JWTAuth())
	{
		protected.GET("", bookRepository.FindBooks)
		protected.POST("", bookRepository.CreateBook)
		protected.GET("/:id", bookRepository.FindBook)
		protected.PUT("/:id", bookRepository.UpdateBook)
		protected.DELETE("/:id", bookRepository.DeleteBook)
	}
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
