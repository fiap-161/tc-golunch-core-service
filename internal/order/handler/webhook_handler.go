package handler

import (
	"context"
	"net/http"

	"github.com/fiap-161/tc-golunch-core-service/internal/order/controller"
	"github.com/fiap-161/tc-golunch-core-service/internal/order/entity/enum"
	apperror "github.com/fiap-161/tc-golunch-core-service/internal/shared/errors"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	controller *controller.Controller
}

type PaymentWebhookRequest struct {
	OrderID string `json:"order_id" binding:"required"`
	Status  string `json:"status" binding:"required"`
}

func NewWebhookHandler(controller *controller.Controller) *WebhookHandler {
	return &WebhookHandler{controller: controller}
}

// PaymentWebhook godoc
// @Summary      Payment Status Webhook
// @Description  Receive payment status updates from Payment Service
// @Tags         Webhooks
// @Accept       json
// @Produce      json
// @Param        request body PaymentWebhookRequest true "Payment status update"
// @Success      200  {object}  map[string]any
// @Failure      400  {object}  errors.ErrorDTO
// @Failure      404  {object}  errors.ErrorDTO
// @Router       /webhook/payment [post]
func (h *WebhookHandler) PaymentWebhook(c *gin.Context) {
	var webhookReq PaymentWebhookRequest
	if err := c.ShouldBindJSON(&webhookReq); err != nil {
		c.JSON(http.StatusBadRequest, apperror.ErrorDTO{
			Message:      "invalid request body",
			MessageError: err.Error(),
		})
		return
	}

	// Find the order
	orderDAO, err := h.controller.FindByID(context.Background(), webhookReq.OrderID)
	if err != nil {
		c.JSON(http.StatusNotFound, apperror.ErrorDTO{
			Message:      "order not found",
			MessageError: err.Error(),
		})
		return
	}

	// Update order status based on payment status
	var newStatus enum.OrderStatus
	switch webhookReq.Status {
	case "paid", "approved":
		newStatus = enum.OrderStatusReceived // Move to received when payment is confirmed
	case "cancelled", "rejected":
		newStatus = enum.OrderStatusAwaitingPayment // Keep awaiting payment for failed payments
	default:
		// Unknown status, ignore
		c.JSON(http.StatusOK, gin.H{"message": "webhook received, status ignored"})
		return
	}

	// Update the order with new status
	orderDAO.Status = newStatus
	_, err = h.controller.Update(context.Background(), orderDAO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apperror.ErrorDTO{
			Message:      "failed to update order status",
			MessageError: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "payment status updated successfully",
		"order_id": webhookReq.OrderID,
		"status":   string(newStatus),
	})
}