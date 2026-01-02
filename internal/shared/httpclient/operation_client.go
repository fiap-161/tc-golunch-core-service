package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type OperationServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

type OperationNotificationRequest struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func NewOperationServiceClient() *OperationServiceClient {
	baseURL := os.Getenv("OPERATION_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8083"
	}

	return &OperationServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *OperationServiceClient) NotifyNewOrder(ctx context.Context, orderID, status string) error {
	payload := OperationNotificationRequest{
		OrderID: orderID,
		Status:  status,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal notification request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/admin/orders/notify", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("operation service returned status %d", resp.StatusCode)
	}

	return nil
}

func (c *OperationServiceClient) UpdateOrderStatus(ctx context.Context, orderID, status string) error {
	payload := OperationNotificationRequest{
		OrderID: orderID,
		Status:  status,
	}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal update request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", c.baseURL+"/admin/orders/"+orderID, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("operation service returned status %d", resp.StatusCode)
	}

	return nil
}
