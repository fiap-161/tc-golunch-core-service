package handler

import (
	"context"
	"net/http"

	"github.com/fiap-161/tc-golunch-order-service/internal/admin/controller"
	"github.com/fiap-161/tc-golunch-order-service/internal/admin/dto"
	apperror "github.com/fiap-161/tc-golunch-order-service/internal/shared/errors"
	"github.com/fiap-161/tc-golunch-order-service/internal/shared/helper"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	adminController *controller.Controller
}

func New(adminController *controller.Controller) *Handler {
	return &Handler{
		adminController: adminController,
	}
}

// Register godoc
// @Summary      Register Admin
// @Description  Register a new admin user
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AdminRequestDTO  true  "Admin registration details"
// @Success      201      {object}  map[string]string     "Success message"
// @Failure      400      {object}  errors.ErrorDTO
// @Failure      500      {object}  errors.ErrorDTO
// @Router       /admin/register [post]
func (h *Handler) Register(c *gin.Context) {
	ctx := context.Background()

	var adminRequest dto.AdminRequestDTO
	if err := c.ShouldBindJSON(&adminRequest); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	err := h.adminController.Register(ctx, adminRequest)

	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// Login godoc
// @Summary      Admin Login
// @Description  Authenticates an admin user and returns a JWT token
// @Tags         Admin Domain
// @Accept       json
// @Produce      json
// @Param        request  body      dto.AdminRequestDTO  true  "Admin login credentials"
// @Success      200      {object}  TokenDTO
// @Failure      400      {object}  errors.ErrorDTO
// @Failure      401      {object}  errors.ErrorDTO
// @Failure      500      {object}  errors.ErrorDTO
// @Router       /admin/login [post]
func (h *Handler) Login(c *gin.Context) {
	ctx := context.Background()

	var adminRequest dto.AdminRequestDTO
	if err := c.ShouldBindJSON(&adminRequest); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	token, err := h.adminController.Login(ctx, adminRequest)

	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, &TokenDTO{
		TokenString: token,
	})
}

type TokenDTO struct {
	TokenString string `json:"token"`
}
