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

type PaymentServiceClient struct {
	baseURL    string
	httpClient *http.Client
}

type PaymentRequest struct {
	OrderID string `json:"order_id"`
}

type PaymentResponse struct {
	ID      string `json:"id"`
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
	QRCode  string `json:"qr_code,omitempty"`
}

type OrderUpdateRequest struct {
	Status string `json:"status"`
}

func NewPaymentServiceClient() *PaymentServiceClient {
	baseURL := os.Getenv("PAYMENT_SERVICE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8082"
	}

	return &PaymentServiceClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *PaymentServiceClient) CreatePayment(ctx context.Context, orderID string) (*PaymentResponse, error) {
	payload := PaymentRequest{OrderID: orderID}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/payment", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment service returned status %d", resp.StatusCode)
	}

	var paymentResp PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &paymentResp, nil
}

func (c *PaymentServiceClient) GetPaymentStatus(ctx context.Context, paymentID string) (*PaymentResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/payment/"+paymentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment service returned status %d", resp.StatusCode)
	}

	var paymentResp PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &paymentResp, nil
}
