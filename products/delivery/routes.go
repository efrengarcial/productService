package delivery

import (
	"context"
	"github.com/apps/productService/internal/mid"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"

	"github.com/apps/productService/products"
	"github.com/gin-gonic/gin"
)

const fiveSecondsTimeout = time.Second * 5

type delivery struct {
	usecase products.ProductService
}

func (d *delivery) Create(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	product := &products.Product{}
	if err := c.ShouldBindJSON(&product); err != nil {
		_ = c.Error(err)
		return
	}

	if err := d.usecase.Create(ctx, product); err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusCreated, "Ok")
}

func (d *delivery) Get(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	id := c.Param("id")

	product, err := d.usecase.Get(ctx, id)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, product)
}

func (d *delivery) GetAll(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), fiveSecondsTimeout)
	defer cancel()

	products, err := d.usecase.GetAll(ctx)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, products)
}


func rootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"text": "Welcome to gin lambda server.",
	})
}

// Routes -
func Routes(usecase products.ProductService, logger *logrus.Logger) (*gin.Engine , error) {
	// set server mode
	gin.SetMode(gin.DebugMode)

	r := gin.New()

	// Global middleware
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(mid.Error(logger))

	delivery := &delivery{usecase}

	r.POST("/products", delivery.Create)
	r.GET("/products", delivery.GetAll)
	r.GET("/products/:id", delivery.Get)
	r.GET("/", rootHandler)

	return r, nil
}
