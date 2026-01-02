package services

import (
	"context"
	"log"

	"github.com/fiap-161/tc-golunch-order-service/internal/order/entity/enum"
	"github.com/fiap-161/tc-golunch-order-service/internal/order/usecases"
	"github.com/fiap-161/tc-golunch-order-service/internal/shared/httpclient"
)

// OrderServiceHTTPGateway provides order operations via HTTP
type OrderServiceHTTPGateway struct {
	orderUseCase      *usecases.UseCases
	paymentClient     *httpclient.PaymentServiceClient
	productionClient  *httpclient.ProductionServiceClient
}

// OrderData represents order data for external services
type OrderData struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func NewOrderServiceHTTPGateway(orderUseCase *usecases.UseCases) *OrderServiceHTTPGateway {
	return &OrderServiceHTTPGateway{
		orderUseCase:     orderUseCase,
		paymentClient:    httpclient.NewPaymentServiceClient(),
		productionClient: httpclient.NewProductionServiceClient(),
	}
}

// FindByID returns order data for external services
func (g *OrderServiceHTTPGateway) FindByID(ctx context.Context, orderID string) (OrderData, error) {
	order, err := g.orderUseCase.FindByID(ctx, orderID)
	if err != nil {
		return OrderData{}, err
	}

	return OrderData{
		ID:     order.Entity.ID,
		Status: string(order.Status),
	}, nil
}

// Update updates order status and notifies other services
func (g *OrderServiceHTTPGateway) Update(ctx context.Context, orderData OrderData) (OrderData, error) {
	currentOrder, err := g.orderUseCase.FindByID(ctx, orderData.ID)
	if err != nil {
		return OrderData{}, err
	}

	oldStatus := string(currentOrder.Status)
	currentOrder.Status = enum.OrderStatus(orderData.Status)

	updatedOrder, updateErr := g.orderUseCase.Update(ctx, currentOrder)
	if updateErr != nil {
		return OrderData{}, updateErr
	}

	// Notify other services about status change
	g.notifyStatusChange(ctx, orderData.ID, oldStatus, orderData.Status)

	return OrderData{
		ID:     updatedOrder.Entity.ID,
		Status: string(updatedOrder.Status),
	}, nil
}

// CreatePaymentForOrder creates payment when order is created
func (g *OrderServiceHTTPGateway) CreatePaymentForOrder(ctx context.Context, orderID string) error {
	_, err := g.paymentClient.CreatePayment(ctx, orderID)
	if err != nil {
		log.Printf("Failed to create payment for order %s: %v", orderID, err)
		return err
	}
	return nil
}

// NotifyProductionService notifies production service about order changes
func (g *OrderServiceHTTPGateway) NotifyProductionService(ctx context.Context, orderID, status string) error {
	// Only notify production for specific statuses
	if status == "paid" || status == "preparing" || status == "ready" {
		err := g.productionClient.NotifyNewOrder(ctx, orderID, status)
		if err != nil {
			log.Printf("Failed to notify production service for order %s: %v", orderID, err)
			return err
		}
	}
	return nil
}

// notifyStatusChange handles notifications to other services when status changes
func (g *OrderServiceHTTPGateway) notifyStatusChange(ctx context.Context, orderID, oldStatus, newStatus string) {
	// Async notifications to avoid blocking the main flow
	go func() {
		// Notify production service
		if err := g.NotifyProductionService(context.Background(), orderID, newStatus); err != nil {
			log.Printf("Error notifying production service: %v", err)
		}
	}()
}