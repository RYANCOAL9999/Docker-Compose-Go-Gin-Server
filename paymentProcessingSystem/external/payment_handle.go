package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/models"
)

func MakePayment(paymentReq interface{}, url string) (*models.PaymentResponse, error) {

	// Convert the transfer request to JSON
	requestBody, err := json.Marshal(paymentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	//Send HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("transfer service returned non-OK status: %v", resp.Status)
	}

	// 解析响应
	var transferResp models.PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&transferResp); err != nil {
		return nil, fmt.Errorf("failed to decode transfer response: %v", err)
	}

	return &transferResp, nil
}

func CheckPaymentStatus(transactionID string, url string) (*models.PaymentResponse, error) {
	statusURL := fmt.Sprintf("%s/status/%s", url, transactionID)

	// Create HTTP request
	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}

	// Send HTTP request
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check HTTP response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment service returned non-OK status: %v", resp.Status)
	}

	// Parse response
	var statusResp models.PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %v", err)
	}

	return &statusResp, nil
}
