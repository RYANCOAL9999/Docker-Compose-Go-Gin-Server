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

	var url string = "https://api.paymentgateway.com/payment"

	paymentReq := models.PaymentRequest{
		Amount:         payment.Amount,
		CardNumber:     payment.Describle.CardNumber,
		ExpirationDate: payment.Describle.ExpirationDate,
		CVV:            payment.Describle.Key,
		Currency:       payment.Describle.Currency,
	}

	// Create payment request
	paymentResp, err := external.MakePayment(paymentReq, url)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
		return nil
	}
	return paymentResp
}

func MakeBankTransfer(payment models.Payment) *models.PaymentResponse {

	const url string = "https://api.bank.com/transfer"

	paymentReq := models.TransferRequest{
		SenderAccount:   payment.Describle.Sender,
		ReceiverAccount: payment.Describle.Receiver,
		Amount:          payment.Amount,
		Currency:        payment.Describle.Currency,
		Description:     payment.Describle.Description,
	}

	// Create payment request
	paymentResp, err := external.MakePayment(paymentReq, url)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
		return nil
	}
	return paymentResp
}

func MakeThirdPartyPayment(payment models.Payment) *models.PaymentResponse {

	const url string = "https://api.thirdParty.com/payment"

	paymentReq := models.PaymentRequest{
		Amount:         payment.Amount,
		CardNumber:     payment.Describle.CardNumber,
		ExpirationDate: payment.Describle.ExpirationDate,
		CVV:            payment.Describle.Key,
		Currency:       payment.Describle.Currency,
		Description:    payment.Describle.Description,
	}

	// Create payment request
	paymentResp, err := external.MakePayment(paymentReq, url)
	if err != nil {
		fmt.Printf("Payment failed: %v\n", err)
		return nil
	}

	fmt.Printf("Payment successful: %+v\n", paymentResp)

	// Query payment status
	statusResp, err := external.CheckPaymentStatus(paymentResp.TransactionID, url)
	if err != nil {
		fmt.Printf("Check payment status failed: %v\n", err)
		return nil
	}
	return statusResp
}

func MakeBlockchainPayment(payment models.Payment) *models.PaymentResponse {

	var url string = "https://api.blockchainplatform.com/transaction"

	paymentReq := models.BlockchainPaymentRequest{
		SenderAddress:   payment.Describle.Sender,
		ReceiverAddress: payment.Describle.Receiver,
		Amount:          payment.Amount,
		Currency:        payment.Describle.Currency,
		PrivateKey:      payment.Describle.Key,
	}

	// Create payment request
	paymentResp, err := external.MakePayment(paymentReq, url)
	if err != nil {
		fmt.Printf("Blockchain payment failed: %v\n", err)
		return nil
	}

	fmt.Printf("Blockchain payment successful: %+v\n", paymentResp)

	// check payment status
	statusResp, err := external.CheckPaymentStatus(paymentResp.TransactionID, url)
	if err != nil {
		fmt.Printf("Check blockchain payment status failed: %v\n", err)
		return nil
	}
	return statusResp
}

// @Summary      Retrieve a payment by ID
// @Description  Get details of a specific payment identified by its ID from the database.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "Payment ID"
// @Success      200  {object}  models.Payment  "Payment details"
// @Failure      400  {object}  models.ErrorResponse  "Invalid ID supplied"
// @Failure      500  {object}  models.ErrorResponse  "Internal server error"
// @Router       /payments/{id} [get]
func ShowPayment(c *gin.Context, db *sql.DB) {
	id, _ := strconv.Atoi(c.Param("id"))
	payment, err := databases.GetPayment(db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, payment)
}

// @Summary      Create a new payment
// @Description  Create a new payment entry in the database using the provided payment details. The payment can be of various methods including credit card, bank transfer, third-party, or blockchain.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        payment  body  models.Payment  true  "Payment details to be created"
// @Success      201  {object}  models.PaymentResult  "Payment created successfully with the payment ID"
// @Failure      400  {object}  models.ErrorResponse  "Bad request due to invalid input"
// @Failure      500  {object}  models.ErrorResponse  "Internal server error"
// @Router       /payments [post]
func CreatePayment(c *gin.Context, db *sql.DB) {
	var payment models.Payment
	if err := c.BindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	var item *models.PaymentResponse

	switch method := payment.Method; method {
	case "CreditCardPayment":
		item = MakeCreditCardPayment(payment)
	case "BankTransfer":
		item = MakeBankTransfer(payment)
	case "ThirdPartyPayment":
		item = MakeThirdPartyPayment(payment)
	case "BlockchainPayment":
		item = MakeBlockchainPayment(payment)
	}

	paymentID, err := databases.AddPayment(db, payment)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	var result models.PaymentResult
	result.ID = paymentID
	result.Status = item.Status
	c.JSON(http.StatusCreated, result)
}
