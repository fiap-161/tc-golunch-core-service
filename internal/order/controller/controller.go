package controller

import (
	"context"
	"log"

	"github.com/fiap-161/tc-golunch-order-service/internal/order/dto"
	"github.com/fiap-161/tc-golunch-order-service/internal/order/gateway"
	"github.com/fiap-161/tc-golunch-order-service/internal/order/gateway/services"
	"github.com/fiap-161/tc-golunch-order-service/internal/order/interfaces"
	"github.com/fiap-161/tc-golunch-order-service/internal/order/presenter"
	"github.com/fiap-161/tc-golunch-order-service/internal/order/usecases"
)

type Controller struct {
	orderUseCase *usecases.UseCases
	httpGateway  *services.OrderServiceHTTPGateway
}

func Build(orderGateway *gateway.Gateway, productService interfaces.ProductService, productOrderService interfaces.ProductOrderService) *Controller {
	orderUseCase := usecases.Build(orderGateway, productService, productOrderService)
	httpGateway := services.NewOrderServiceHTTPGateway(orderUseCase)

	return &Controller{
		orderUseCase: orderUseCase,
		httpGateway:  httpGateway,
	}
}

func (c *Controller) Create(ctx context.Context, orderDTO dto.CreateOrderDTO) (dto.OrderDAO, error) {
	order, err := c.orderUseCase.CreateCompleteOrder(ctx, orderDTO)
	if err != nil {
		return dto.OrderDAO{}, err
	}

	// After order is created, trigger payment creation asynchronously
	go func() {
		if err := c.httpGateway.CreatePaymentForOrder(context.Background(), order.Entity.ID); err != nil {
			log.Printf("Failed to create payment for order %s: %v", order.Entity.ID, err)
		}
	}()

	// Notify production service about new order
	go func() {
		if err := c.httpGateway.NotifyProductionService(context.Background(), order.Entity.ID, string(order.Status)); err != nil {
			log.Printf("Failed to notify production service for order %s: %v", order.Entity.ID, err)
		}
	}()

	presenter := presenter.Build()
	return presenter.FromEntityToDAO(order), nil
}

func (c *Controller) GetAll(ctx context.Context, id string) ([]dto.OrderDAO, error) {
	presenter := presenter.Build()

	orders, err := c.orderUseCase.GetAllOrById(ctx, id)
	if err != nil {
		return nil, err
	}

	return presenter.FromEntityListToDAOList(orders), nil
}

func (c *Controller) GetPanel(ctx context.Context) ([]dto.OrderDAO, error) {
	presenter := presenter.Build()

	orders, err := c.orderUseCase.GetPanel(ctx)
	if err != nil {
		return nil, err
	}

	return presenter.FromEntityListToDAOList(orders), nil
}

func (c *Controller) FindByID(ctx context.Context, id string) (dto.OrderDAO, error) {
	presenter := presenter.Build()

	order, err := c.orderUseCase.FindByID(ctx, id)
	if err != nil {
		return dto.OrderDAO{}, err
	}

	return presenter.FromEntityToDAO(order), nil
}

func (c *Controller) Update(ctx context.Context, orderDTO dto.OrderDAO) (dto.OrderDAO, error) {
	presenter := presenter.Build()

	order := dto.FromOrderDAO(orderDTO)
	oldStatus := string(order.Status)

	updated, err := c.orderUseCase.Update(ctx, order)
	if err != nil {
		return dto.OrderDAO{}, err
	}

	// Notify other services about status change if status changed
	newStatus := string(updated.Status)
	if oldStatus != newStatus {
		go func() {
			if err := c.httpGateway.NotifyProductionService(context.Background(), updated.Entity.ID, newStatus); err != nil {
				log.Printf("Failed to notify production service for order %s status change: %v", updated.Entity.ID, err)
			}
		}()
	}

	return presenter.FromEntityToDAO(updated), nil
}
