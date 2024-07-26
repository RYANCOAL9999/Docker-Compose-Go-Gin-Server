package external

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/models"
)

const transferServiceURL = "https://api.bank.com/transfer"

func MakeTransfer(transferReq models.TransferRequest) (*models.PaymentResponse, error) {
	// 将转账请求转换为JSON
	requestBody, err := json.Marshal(transferReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer request: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", transferServiceURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送HTTP请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 检查HTTP响应状态
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
