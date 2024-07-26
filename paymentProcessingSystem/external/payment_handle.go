package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/models"
)

const paymentGatewayURL = "https://api.paymentgateway.com/payment"

func MakePayment(paymentReq models.PaymentRequest) (*models.PaymentResponse, error) {
	// 將支付請求轉換為JSON
	requestBody, err := json.Marshal(paymentReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payment request: %v", err)
	}

	// 創建HTTP請求
	req, err := http.NewRequest("POST", paymentGatewayURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 發送HTTP請求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 檢查HTTP響應狀態
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment gateway returned non-OK status: %v", resp.Status)
	}

	// 解析響應
	var paymentResp models.PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&paymentResp); err != nil {
		return nil, fmt.Errorf("failed to decode payment response: %v", err)
	}

	return &paymentResp, nil
}

func CheckPaymentStatus(transactionID string) (*models.PaymentResponse, error) {
	statusURL := fmt.Sprintf("%s/status/%s", paymentGatewayURL, transactionID)

	// 创建HTTP请求
	req, err := http.NewRequest("GET", statusURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}

	// 发送HTTP请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("payment service returned non-OK status: %v", resp.Status)
	}

	// 解析响应
	var statusResp models.PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		return nil, fmt.Errorf("failed to decode status response: %v", err)
	}

	return &statusResp, nil
}
