package api

import (
	"context"
	"errors"
	"golang-rest-api-template/pkg/auth"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/models"
	"golang-rest-api-template/pkg/response"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	LoginHandler(c *gin.Context)
	RegisterHandler(c *gin.Context)
}

// userRepository holds shared resources like database
type userRepository struct {
	DB  database.Database
	Ctx *context.Context
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db database.Database, ctx *context.Context) *userRepository {
	return &userRepository{
		DB:  db,
		Ctx: ctx,
	}
}

// @BasePath /api/v1

// LoginHandler godoc
// @Summary Authenticate a user
// @Schemes
// @Description Authenticates a user using username and password, returns a JWT token if successful
// @Tags user
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   user     body    models.LoginUser     true        "User login object"
// @Success 200 {string} string "JWT Token"
// @Failure 400 {string} string "Bad Request"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /login [post]
func (r *userRepository) LoginHandler(c *gin.Context) {
	var incomingUser models.User
	var dbUser models.User

	if err := c.ShouldBindJSON(&incomingUser); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := r.DB.Where("username = ?", incomingUser.Username).First(&dbUser).Error(); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Unauthorized(c, "Invalid username or password")
		} else {
			response.InternalServerError(c, "Database error")
		}
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(incomingUser.Password)); err != nil {
		response.Unauthorized(c, "Invalid username or password")
		return
	}

	token, err := auth.GenerateToken(dbUser.Username, dbUser.ID)
	if err != nil {
		response.InternalServerError(c, "Error generating token")
		return
	}

	response.OK(c, gin.H{"token": token})
}

// RegisterHandler godoc
// @Summary Register a new user
// @Schemes http
// @Description Registers a new user with the given username and password
// @Tags user
// @Security ApiKeyAuth
// @Accept  json
// @Produce  json
// @Param   user     body    models.LoginUser     true        "User registration object"
// @Success 201 {string} string	"Successfully registered"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /register [post]
func (r *userRepository) RegisterHandler(c *gin.Context) {
	var user models.LoginUser

	if err := c.ShouldBindJSON(&user); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		response.InternalServerError(c, "Could not hash password")
		return
	}

	newUser := models.User{Username: user.Username, Password: hashedPassword}

	if err := r.DB.Create(&newUser).Error; err != nil {
		response.InternalServerError(c, "Could not save user")
		return
	}

	response.Created(c, gin.H{"message": "Registration successful"})
}
