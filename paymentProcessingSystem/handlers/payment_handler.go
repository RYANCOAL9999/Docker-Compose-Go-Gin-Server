package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/databases"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/external"
	"github.com/RYANCOAL9999/SpinnrTechnologyInterview/paymentProcessingSystem/models"
	"github.com/gin-gonic/gin"
)

func MakeCreditCardPayment(payment models.Payment) *models.PaymentResponse {
	paymentReq := models.PaymentRequest{
		Amount:         payment.Amount,
		CardNumber:     payment.Describle.CardNumber,
		ExpirationDate: payment.Describle.ExpirationDate,
		CVV:            payment.Describle.Key,
		Currency:       payment.Describle.Currency,
	}

	paymentResp, err := external.MakePayment(paymentReq)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
		return nil
	}
	return paymentResp
}

func MakeBankTransfer(payment models.Payment) *models.PaymentResponse {
	paymentReq := models.TransferRequest{
		SenderAccount:   payment.Describle.Sender,
		ReceiverAccount: payment.Describle.Receiver,
		Amount:          payment.Amount,
		Currency:        payment.Describle.Currency,
		Description:     payment.Describle.Description,
	}

	paymentResp, err := external.MakeTransfer(paymentReq)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
		return nil
	}
	return paymentResp
}

func MakeThirdPartyPayment(payment models.Payment) *models.PaymentResponse {
	paymentReq := models.PaymentRequest{
		Amount:         payment.Amount,
		CardNumber:     payment.Describle.CardNumber,
		ExpirationDate: payment.Describle.ExpirationDate,
		CVV:            payment.Describle.Key,
		Currency:       payment.Describle.Currency,
		Description:    payment.Describle.Description,
	}

	// 创建支付请求
	paymentResp, err := external.MakePayment(paymentReq)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
		return nil
	}

	fmt.Printf("Payment successful: %+v\n", paymentResp)

	// 查询支付状态
	statusResp, err := external.CheckPaymentStatus(paymentResp.TransactionID)
	if err != nil {
		fmt.Printf("Check payment status failed: %v\n", err)
		return nil
	}
	return statusResp
}

func MakeBlockchainPayment(payment models.Payment) *models.PaymentResponse {

	paymentReq := models.BlockchainPaymentRequest{
		SenderAddress:   payment.Describle.Sender,
		ReceiverAddress: payment.Describle.Receiver,
		Amount:          payment.Amount,
		Currency:        payment.Describle.Currency,
		PrivateKey:      payment.Describle.Key,
	}

	// 创建支付请求
	paymentResp, err := external.MakeBlockchainPayment(paymentReq)
	if err != nil {
		fmt.Printf("Blockchain payment failed: %v\n", err)
		return nil
	}

	fmt.Printf("Blockchain payment successful: %+v\n", paymentResp)

	// 查询支付状态
	statusResp, err := external.CheckBlockchainPaymentStatus(paymentResp.TransactionID)
	if err != nil {
		fmt.Printf("Check blockchain payment status failed: %v\n", err)
		return nil
	}
	return statusResp

}

func ShowPayment(c *gin.Context, db *sql.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	payment, err := databases.GetPayment(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payment)
}

func CreatePayment(c *gin.Context, db *sql.DB) {
	var payment models.Payment
	if err := c.BindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// var item *models.PaymentResponse

	// switch method := payment.Method; method {
	// case "CreditCardPayment":
	// 	item = MakeCreditCardPayment(payment)
	// case "BankTransfer":
	// 	item = MakeCreditCardPayment(payment)
	// case "ThirdPartyPayment":
	// 	item = MakeThirdPartyPayment(payment)
	// case "BlockchainPayment":
	// 	item = MakeBlockchainPayment(payment)
	// }

	paymentID, err := databases.AddPayment(db, payment)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var result models.PaymentResult
	result.ID = paymentID
	result.Status = "Success"
	// result.Status = item.Status
	c.JSON(http.StatusCreated, result)
}
