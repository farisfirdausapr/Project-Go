package handler

import (
	"fmt"
	"net/http"
	"project/auth"
	"project/model"
	u "project/model/user"
	reposerviceUser "project/reposervice/reposervice-user"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService reposerviceUser.Service
	authService auth.Service
}

func NewUserHandler(userService reposerviceUser.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input u.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := model.FormatValidationError(err)
		errorsMessage := gin.H{"error": errors}

		response := model.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	user, err := h.userService.RegistrasiUser(input)
	if err != nil {
		response := model.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		response := model.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := u.FormatUser(user, token)

	response := model.APIResponse("Account has been register", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input u.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := model.FormatValidationError(err)
		errorsMessage := gin.H{"error": errors}

		response := model.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	user, err := h.userService.Login(input)
	if err != nil {
		errorsMessage := gin.H{"error": err.Error()}

		response := model.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(user.ID)
	if err != nil {
		response := model.APIResponse("Login failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := u.FormatUser(user, token)

	response := model.APIResponse("Successfuly login", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input u.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := model.FormatValidationError(err)
		errorsMessage := gin.H{"error": errors}

		response := model.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsAvailableEmail(input)
	if err != nil {
		errorsMessage := gin.H{"error": "Server error"}

		response := model.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorsMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_abailable": isEmailAvailable,
	}

	var metaMessage string

	if isEmailAvailable {
		metaMessage = "Email has been register"
	} else {
		metaMessage = "Email is available"
	}

	response := model.APIResponse(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvataric(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := model.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(u.User)
	userID := currentUser.ID

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := model.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}

		response := model.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}

	response := model.APIResponse("Avatar successfuly uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
