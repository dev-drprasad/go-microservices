package rest

import (
	"gomicroservices/internal/order/model"
	"gomicroservices/internal/order/service"
	"gomicroservices/internal/util"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	service service.Service
}

func New(db *pgxpool.Pool) *OrderHandler {
	return &OrderHandler{
		service: service.New(db),
	}
}

type orderProductBody struct {
	OrderID   uint `json:"orderId"`
	ProductID uint `json:"productId"`
	Quantity  uint `json:"quantity"`
}

type orderBody struct {
	CustomerID uint `json:"customerId"`

	Products []orderProductBody `json:"products"`
}

func (o orderBody) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.CustomerID, validation.Required, validation.Min(uint(1))),
	)
}

func (h *OrderHandler) PlaceOrder(c echo.Context) error {
	ctx := util.NewContextWithLogger(c)

	var r orderBody
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, util.NewAPIError(http.StatusText(http.StatusBadRequest)))
	}

	if err := validation.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, util.NewAPIError(err.Error()))
	}

	products := []*model.OrderProduct{}
	for _, op := range r.Products {
		products = append(products, &model.OrderProduct{ProductID: op.ProductID, Quantity: op.Quantity})
	}
	order := model.Order{CustomerID: r.CustomerID, Products: products}
	if err := h.service.PlaceOrder(ctx, order); err != nil {
		c.Logger().Errorf("Failed to place an order: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, util.NewAPIError(http.StatusText(http.StatusInternalServerError)))
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *OrderHandler) GetOrder(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	branch, err := h.service.GetOrder(util.NewContextWithLogger(c), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.NewAPIError(http.StatusText(http.StatusInternalServerError)))
	}
	return c.JSON(http.StatusOK, branch)
}

func (h *OrderHandler) GetOrders(c echo.Context) error {
	orders, err := h.service.GetOrders(util.NewContextWithLogger(c))
	if err != nil {
		c.Logger().Errorf("Failed to get orders. %v", err)
		return c.JSON(http.StatusInternalServerError, util.NewAPIError(http.StatusText(http.StatusInternalServerError)))
	}
	return c.JSON(http.StatusOK, orders)
}
