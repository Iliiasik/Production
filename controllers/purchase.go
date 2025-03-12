package controllers

import (
	"github.com/gin-gonic/gin"
)

func PurchaseForm(c *gin.Context) {
	c.HTML(200, "purchase.html", gin.H{})
}
