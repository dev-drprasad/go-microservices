package rest

import (
	"gomicroservices/internal/user/model"
	"gomicroservices/internal/user/service"
	"gomicroservices/internal/util"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/labstack/echo/v4"
)

// var service usersvc.Service

// func init() {
// 	service = usersvc.NewService()
// }

type UserHandler struct {
	svc service.Service
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func New(db *pgxpool.Pool) *UserHandler {
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

func (h *UserHandler) GetUsers(c echo.Context) error {
	ctx := util.NewContextWithLogger(c)
	users, err := h.svc.GetUsers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

type userRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	BranchID uint   `json:"branchId"`
	Role     string `json:"role"`
}

func (u userRequest) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required, validation.Length(5, 60)),
		validation.Field(&u.Username, validation.Required, validation.Length(3, 20)),
		validation.Field(&u.Password, validation.Required, validation.Length(5, 20)),
		validation.Field(&u.BranchID, validation.Required),
		validation.Field(&u.Role, validation.Required, validation.In("admin", "staff")),
	)
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	ctx := util.NewContextWithLogger(c)

	var r userRequest
	if err := c.Bind(&r); err != nil {
		c.Logger().Errorf("Bind of payload -> user failed. %v", err)
		return c.JSON(http.StatusBadRequest, util.NewAPIError(http.StatusText(http.StatusBadRequest)))
	}
	if err := validation.Validate(r); err != nil {
		c.Logger().Errorf("Payload validation failed. %v", err)
		return c.JSON(http.StatusBadRequest, util.NewAPIError(err.Error()))
	}

	user := model.User{
		Name:     r.Name,
		Username: r.Username,
		Password: r.Password,
		BranchID: r.BranchID,
		Role:     r.Role,
	}

	if err := h.svc.CreateUser(ctx, &user); err != nil {
		if err == service.ErrInvalidRequest {
			return c.JSON(http.StatusBadRequest, util.NewAPIError(http.StatusText(http.StatusBadRequest)))
		}
		c.Logger().Errorf("Error occured while creating user. %v", err)
		return c.JSON(http.StatusInternalServerError, util.NewAPIError(http.StatusText(http.StatusInternalServerError)))
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *UserHandler) Login(c echo.Context) error {
	ctx := util.NewContextWithLogger(c)
	var payload LoginPayload
	if err := c.Bind(&payload); err != nil {
		return err
	}

	token, err := h.svc.Login(ctx, payload.Username, payload.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			return c.JSON(http.StatusUnauthorized, util.NewAPIError(http.StatusText(http.StatusUnauthorized)))
		} else {
			c.Logger().Errorf("Error occured while login %v", err)
			return c.JSON(http.StatusInternalServerError, util.NewAPIError(http.StatusText(http.StatusInternalServerError)))
		}
	}
	return c.JSON(http.StatusOK, map[string]string{"token": token})
}
