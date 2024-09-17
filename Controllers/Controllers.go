package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	models "github.com/idontknowtoobrother/stripe-go-lang/Models"
	repository "github.com/idontknowtoobrother/stripe-go-lang/Repository"
	utils "github.com/idontknowtoobrother/stripe-go-lang/Utils"
	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/paymentintent"
)

type Controller interface {
	GetProducts(c *gin.Context)
	CreateProduct(c *gin.Context)
	Config(c *gin.Context)
	HandleCreatePaymentIntent(c *gin.Context)
}

type controller struct {
	repo repository.Repo
	ctx  context.Context
}

func NewController(ctx context.Context, repo repository.Repo) Controller {
	return &controller{
		repo: repo,
		ctx:  ctx,
	}
}

func (ctrl *controller) GetProducts(c *gin.Context) {
	products, err := ctrl.repo.GetAll()
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (ctrl *controller) CreateProduct(c *gin.Context) {
	var product models.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := ctrl.repo.Create(&product); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"product": product})
}

func (ctrl *controller) Config(c *gin.Context) {
	err := godotenv.Load()
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"publishableKey": os.Getenv("STRIPE_PUBLISHABLE_KEY"),
	})
}

func (ctrl *controller) HandleCreatePaymentIntent(c *gin.Context) {

	var product models.Product
	stripe.Key = utils.GetEnv("STRIPE_SECRET_KEY")

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(product.Uuid)

	data, err := ctrl.repo.GetByUuid(product.Uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(int64(data.Price))
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(data.Price)),
		Currency: stripe.String(string(stripe.CurrencyTHB)),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"clientSecret": pi.ClientSecret})
}
