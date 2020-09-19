package rest

import (
	"database/sql"
	"gomicroservices/internal/organization/service"
	"gomicroservices/internal/util"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// func GetUsers(c echo.Context) error {
// 	users := []model.User{{"admin2", "0000"}, {"admin", "0000"}}
// 	return c.JSON(http.StatusOK, users)
// }

type OrganizationHandler struct {
	service service.Service
}

func New(db *sql.DB) *OrganizationHandler {
	return &OrganizationHandler{
		service: service.New(db),
	}
}

func (h *OrganizationHandler) GetBranch(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user := h.service.GetBranch(util.NewContextWithLogger(c), uint64(id))
	return c.JSON(http.StatusOK, user)
}
