package rest

import (
	"database/sql"
	"gomicroservices/internal/user/model"
	"gomicroservices/internal/user/service"
	"gomicroservices/internal/util"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// var service usersvc.Service

// func init() {
// 	service = usersvc.NewService()
// }

type UserHandler struct {
	svc service.Service
}

func New(db *sql.DB) *UserHandler {
	return &UserHandler{
		svc: service.New(db),
	}
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	c.Logger().Infof("Reading user from service. id=%d", id)

	// ctx := context.Background()
	ctx := util.NewContextWithLogger(c)
	user, err := h.svc.GetUser(ctx, uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := util.NewContextWithLogger(c)

	var user model.User
	if err := c.Bind(&user); err != nil {
		return err
	}
	if err := h.svc.CreateUser(ctx, &user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, nil)
}
