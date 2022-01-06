package api

import (
	"github.com/gin-gonic/gin"

	"github.com/BenjaminCallahan/my-bank-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Handler represent of entity handler
type Handler struct {
	service *service.Service
}

// NewHandler creates s new entity of Handler
func NewHandler(service *service.Service) *Handler {
	return &Handler{service}
}

// InitRoutes init all roues of web application
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.POST("/deposit", h.deposit)
	router.POST("/withdraw", h.withdraw)

	router.GET("/balance", h.balance)
	router.GET("/currency", h.currency)
	return router
}

// input body of amount
type amountInput struct {
	Amount float64 `json:"amount" binding:"numeric"`
}

// deposit allows to deposit amount the amount to a bank account
func (h *Handler) deposit(c *gin.Context) {
	var input amountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := "invalid input body"
		logrus.WithField("handler", "deposit").Errorf("%s: %s\n", errMsg, err.Error())
		newErrorResponse(c, http.StatusBadRequest, errMsg)
		return
	}

	if err := h.service.BankAccount.AddFunds(input.Amount); err != nil {
		errMsg := "failed to add amount"
		logrus.WithField("handler", "deposit").Errorf("%s: %s\n", errMsg, err.Error())
		newErrorResponse(c, http.StatusInternalServerError, errMsg)
		return
	}

	c.Status(http.StatusOK)
}

// withdraw allows to withdrawing  amount from a bank account
func (h *Handler) withdraw(c *gin.Context) {

}

type currencyInput struct {
	Currency string `form:"currency,omitempty" binding:"omitempty,oneof=SBP RUB"`
}

func (h *Handler) balance(c *gin.Context) {

}

func (h *Handler) currency(c *gin.Context) {

}
