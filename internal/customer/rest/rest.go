package rest

import (
	"gomicroservices/internal/customer/model"
	"gomicroservices/internal/customer/service"
	"gomicroservices/internal/util"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type CustomerHandler struct {
	service service.Service
}

func New(db *pgxpool.Pool) *CustomerHandler {
	return &CustomerHandler{
		service: service.New(db),
	}
}

type customerRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	Zipcode     string `json:"zipcode"`
	PhoneNumber string `json:"phoneNumber"`
}

func (o customerRequest) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Name, validation.Required, validation.Length(2, 60)),
		validation.Field(&o.Address, validation.Required, validation.Length(2, 255)),
		validation.Field(&o.Zipcode, validation.Required, validation.Length(5, 10)),
		validation.Field(&o.PhoneNumber, validation.Required, validation.Length(10, 15)),
	)
}

func (h *CustomerHandler) AddCustomer(c echo.Context) error {
	ctx := util.NewContextWithLogger(c)

	var r customerRequest
	if err := c.Bind(&r); err != nil {
		c.Logger().Errorf("Failed to bind customer payload : %v", err)
		return c.JSON(http.StatusBadRequest, util.NewAPIError(http.StatusText(http.StatusBadRequest)))
	}

	if err := validation.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, util.NewAPIError(err.Error()))
	}

	customer := model.Customer{Name: r.Name, Address: r.Address, Zipcode: r.Zipcode, PhoneNumber: r.PhoneNumber}
	if err := h.service.AddCustomer(ctx, &customer); err != nil {
		c.Logger().Errorf("Failed to create customer : %v", err)
		return c.JSON(http.StatusInternalServerError, util.NewAPIError(http.StatusText(http.StatusInternalServerError)))
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *CustomerHandler) GetCustomer(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	customer, err := h.service.GetCustomer(util.NewContextWithLogger(c), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.NewAPIError(http.StatusText(http.StatusInternalServerError)))
	}
	return c.JSON(http.StatusOK, customer)
}

func (h *CustomerHandler) UpdateCustomer(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var r customerRequest
	if err := c.Bind(&r); err != nil {
		c.Logger().Errorf("Failed to bind customer payload : %v", err)
		return c.JSON(http.StatusBadRequest, util.NewAPIError(http.StatusText(http.StatusBadRequest)))
	}

	if err := validation.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, util.NewAPIError(err.Error()))
	}

	customer := model.Customer{Name: r.Name, Address: r.Address, Zipcode: r.Zipcode, PhoneNumber: r.PhoneNumber}
	err := h.service.UpdateCustomer(util.NewContextWithLogger(c), uint(id), customer)
	if err != nil {
		c.Logger().Errorf("Failed to update customer : %v", err)
		return c.JSON(http.StatusInternalServerError, util.NewAPIError(http.StatusText(http.StatusInternalServerError)))
	}
	return c.JSON(http.StatusOK, nil)
}

func (h *CustomerHandler) GetCustomers(c echo.Context) error {
	customers, err := h.service.GetCustomers(util.NewContextWithLogger(c))
	if err != nil {
		c.Logger().Errorf("Failed to get customers : %v", err)
		return c.JSON(http.StatusInternalServerError, util.NewAPIError(http.StatusText(http.StatusInternalServerError)))
	}
	return c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandler) NewCustomersCount(c echo.Context) error {
	counts, err := h.service.NewCustomersCount(util.NewContextWithLogger(c))
	if err != nil {
		c.Logger().Errorf("Failed to get new customers count : %v", err)
		return c.JSON(http.StatusInternalServerError, util.NewAPIError(http.StatusText(http.StatusInternalServerError)))
	}
	return c.JSON(http.StatusOK, counts)
}
