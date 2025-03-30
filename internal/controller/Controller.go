package controller

import (
	"Gravitum/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (c *UserController) RegisterUserRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("", c.CreateUser)
		userGroup.GET("/:id", c.GetUserByID)
		userGroup.PUT("/:id", c.UpdateUser)
	}
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	var request struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.userService.CreateUser(
		request.FirstName,
		request.LastName,
		request.Email,
		request.Password,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, user)
}

func (c *UserController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	user, err := c.userService.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	id := ctx.Param("id")

	var request struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email" binding:"omitempty,email"`
		Password  string `json:"password" binding:"omitempty,min=8"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем текущего пользователя
	existingUser, err := c.userService.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}
	if existingUser == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if request.FirstName != "" {
		existingUser.FirstName = request.FirstName
	}
	if request.LastName != "" {
		existingUser.LastName = request.LastName
	}
	if request.Email != "" {
		existingUser.Email = request.Email
	}
	if request.Password != "" {
		existingUser.Password = request.Password
	}

	updatedUser, err := c.userService.UpdateUser(existingUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, updatedUser)
}
