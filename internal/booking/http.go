package booking

import (
	"fmt"
	"net/http"

	"github.com/slaengkast/shipping-api/internal/errors"

	"github.com/gin-gonic/gin"
)

type bookShippingRequest struct {
	Origin      string  `json:"origin" binding:"required"`
	Destination string  `json:"destination" binding:"required"`
	Weight      float32 `json:"weight" binding:"required"`
}

type bookShippingResponse struct {
	Id string `json:"id" binding:"required"`
}

type handler struct {
	bookingService Service
}

func NewHandler(bookingService Service) *handler {
	return &handler{bookingService: bookingService}
}

func (h handler) BookShipping(c *gin.Context) {
	var req bookShippingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.bookingService.BookShipping(c, req.Origin, req.Destination, req.Weight)

	if err != nil {
		handleError(c, err)
		return
	}

	res := bookShippingResponse{
		Id: string(id),
	}

  c.Header("Location", fmt.Sprintf("%s/%s", c.Request.URL.Path, id))
	c.JSON(http.StatusCreated, res)
}

type getBookingResponse struct {
	Id          string  `json:"id" binding:"required"`
	Origin      string  `json:"origin" binding:"required"`
	Destination string  `json:"destination" binding:"required"`
	Weight      float32 `json:"weight" binding:"required"`
	Price       float32 `json:"price" binding:"required"`
	Currency    string  `json:"currency" binding:"required"`
}

func (h handler) GetBooking(c *gin.Context) {
	id := c.Param("id")

	sh, err := h.bookingService.GetBooking(c, id)
	if err != nil {
		handleError(c, err)
		return
	}
	response := getBookingResponse{
		Id:          sh.Id(),
		Origin:      sh.Origin(),
		Destination: sh.Destination(),
		Weight:      sh.Weight(),
		Price:       sh.Price(),
		Currency:    "SEK",
	}

	c.JSON(http.StatusOK, response)
}

func handleError(c *gin.Context, err error) {
	apiError, ok := err.(errors.APIError)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	switch apiError.GetType() {
	case errors.ErrorNotFound:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.ErrorInput:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.ErrorConflict:
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.ErrorInternal:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
