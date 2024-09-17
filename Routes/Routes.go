package routes

import (
	"github.com/gin-gonic/gin"
	controllers "github.com/idontknowtoobrother/stripe-go-lang/Controllers"
)

func SetupRoutes(ctrl controllers.Controller) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	v1.GET("/products", ctrl.GetProducts)
	v1.POST("/products", ctrl.CreateProduct)
	v1.GET("/config", ctrl.Config)
	v1.POST("/create-payment-intent", ctrl.HandleCreatePaymentIntent)

	return r
}
