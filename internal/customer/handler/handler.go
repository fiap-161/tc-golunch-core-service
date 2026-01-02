package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/fiap-161/tc-golunch-order-service/internal/customer/controller"
	"github.com/fiap-161/tc-golunch-order-service/internal/customer/dto"
	apperror "github.com/fiap-161/tc-golunch-order-service/internal/shared/errors"
	"github.com/fiap-161/tc-golunch-order-service/internal/shared/helper"
)

type Handler struct {
	customerController *controller.Controller
}

func New(customerController *controller.Controller) *Handler {
	return &Handler{
		customerController: customerController,
	}
}

// Create godoc
// @Summary      Creates a new customer
// @Description  Creates a customer based on the information sent in the request body
// @Tags         Customer Domain
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CustomerRequestDTO  true  "Customer data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  errors.ErrorDTO
// @Failure      500      {object}  errors.ErrorDTO
// @Router       /customer/register [post]
func (h *Handler) Create(c *gin.Context) {
	ctx := context.Background()

	var customerDTO dto.CustomerRequestDTO
	if err := c.ShouldBindJSON(&customerDTO); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	id, err := h.customerController.Create(ctx, customerDTO)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":      id,
		"message": "Customer created successfully",
	})
}

// Identify godoc
// @Summary      Identifies customer by CPF
// @Description  Returns a JWT token when identifying the customer by CPF
// @Tags         Customer Domain
// @Accept       json
// @Produce      json
// @Param        cpf   path      string     true  "Customer CPF"
// @Success      200   {object}  dto.TokenDTO
// @Failure      404   {object}  errors.ErrorDTO
// @Failure      500   {object}  errors.ErrorDTO
// @Router       /customer/identify/{cpf} [get]
func (h *Handler) Identify(c *gin.Context) {
	ctx := context.Background()
	cpf := c.Param("cpf")

	token, err := h.customerController.Identify(ctx, cpf)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.TokenDTO{
		TokenString: token,
	})
}

// Anonymous godoc
// @Summary      Generates anonymous customer
// @Description  Generates a JWT token for an anonymous customer (without CPF)
// @Tags         Customer Domain
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.TokenDTO
// @Failure      500  {object}  errors.ErrorDTO
// @Router       /customer/anonymous [get]
func (h *Handler) Anonymous(c *gin.Context) {
	ctx := context.Background()

	token, err := h.customerController.Identify(ctx, "")
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.TokenDTO{
		TokenString: token,
	})
}
