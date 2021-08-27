package api

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.POST("/deposit", h.deposit)
	router.POST("/withdraw", h.withdraw)

	router.GET("/balance", h.balance)
	router.GET("/currency", h.currency)
	return router
}

type amountInput struct {
	Amount float64 `json:"amount" binding:"numeric"`
}

func (h *Handler) deposit(c *gin.Context) {

}

func (h *Handler) withdraw(c *gin.Context) {

}

type currencyInput struct {
	Currency string `form:"currency,omitempty" binding:"omitempty,oneof=SBP RUB"`
}

func (h *Handler) balance(c *gin.Context) {

}

func (h *Handler) currency(c *gin.Context) {

}
