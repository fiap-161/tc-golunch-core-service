package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/fiap-161/tc-golunch-core-service/internal/product/controller"
	"github.com/fiap-161/tc-golunch-core-service/internal/product/dto"
	apperror "github.com/fiap-161/tc-golunch-core-service/internal/shared/errors"
	"github.com/fiap-161/tc-golunch-core-service/internal/shared/helper"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	controller *controller.Controller
}

func New(controller *controller.Controller) *Handler {
	return &Handler{controller: controller}
}

// ListCategories List Categories godoc
// @Summary      List Categories
// @Description  List Categories
// @Tags         Product Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Success      200   {array}   string
// @Failure      401  {object}  errors.ErrorDTO
// @Router       /product/categories [get]
func (h *Handler) ListCategories(c *gin.Context) {
	ctx := context.Background()
	c.JSON(http.StatusOK, h.controller.ListCategories(ctx))
}

// UploadImage godoc
// @Summary      Upload product image
// @Description  Uploads a JPEG or PNG image (max 5MB) and returns its public URL
// @Tags         Product Domain
// @Accept       multipart/form-data
// @Produce      json
// @Param        image  formData  file  true  "Product image (JPEG or PNG, max 5MB)"
// @Success      201    {object}  dto.ImageURLDTO
// @Failure      401  {object}  errors.ErrorDTO
// @Failure      400    {object}  errors.ErrorDTO  "Image is missing, invalid, or too large"
// @Failure      500    {object}  errors.ErrorDTO  "Internal error while processing the image"
// @Router       /products/image [post]
func (h *Handler) UploadImage(c *gin.Context) {
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			MessageError: "Validation error",
			Message:      "Image is required or too large (max 5MB).",
		})
		return
	}

	ctx := context.Background()
	imageURL, err := h.controller.UploadImage(ctx, fileHeader)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, dto.ImageURLDTO{ImageURL: imageURL})
}

// Create Product godoc
// @Summary      Create Product
// @Description  Create a new product
// @Tags         Product Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.ProductRequestDTO true "Product to create. Note category is an integer number. See [GET] /product/categories to get a valid category_id"
// @Success      201  {object}  dto.ProductResponseDTO
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      401  {object}  errors.ErrorDTO
// @Router       /product/ [post]
func (h *Handler) Create(c *gin.Context) {
	var productDTO dto.ProductRequestDTO

	if err := c.ShouldBindJSON(&productDTO); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	created, err := h.controller.Create(context.Background(), productDTO)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetAll Get All Products by Category godoc
// @Summary      Get all products by category
// @Description  Returns all products. Optionally, filter by category using query param. Categories must match those returned from [GET] /product/categories.
// @Tags         Product Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        category query string false "Category name (e.g., 'drink', 'meal', 'side', 'dessert')"
// @Success      200  {object}  dto.ProductListResponseDTO
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      401  {object}  errors.ErrorDTO
// @Router       /product [get]
func (h *Handler) GetAllByCategory(c *gin.Context) {
	query := c.Query("category")
	query = strings.ToUpper(query)
	query = strings.ReplaceAll(query, " ", "")

	list, err := h.controller.GetAllByCategory(context.Background(), query)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, list)
}

// Update Product godoc
// @Summary      Update Product
// @Description  Update an existing product by ID
// @Tags         Product Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id       path      string                         true  "Product ID"
// @Param        request  body      dto.ProductRequestUpdateDTO true  "Product data to update"
// @Success      200      {object}  dto.ProductResponseDTO
// @Failure      400      {object}  errors.ErrorDTO
// @Failure      401  {object}  errors.ErrorDTO
// @Router       /product/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var productUpdateDTO dto.ProductRequestUpdateDTO

	if err := c.ShouldBindJSON(&productUpdateDTO); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "Invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	updated, err := h.controller.Update(context.Background(), id, productUpdateDTO)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

// Delete Product godoc
// @Summary      Delete Product
// @Description  Delete a product by ID
// @Tags         Product Domain
// @Security BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      204  "No Content"
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      401  {object}  errors.ErrorDTO
// @Router       /product/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.controller.Delete(context.Background(), id)

	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) ValidateIfProductExists(c *gin.Context) {
	id := c.Param("id")

	_, err := h.controller.FindByID(context.Background(), id)
	if err != nil {
		helper.HandleError(c, err)
		return
	}

	c.Next()
}
