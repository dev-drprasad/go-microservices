package rest

import (
	"gomicroservices/internal/organization/model"
	"gomicroservices/internal/organization/service"
	"gomicroservices/internal/util"
	"net/http"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

type OrganizationHandler struct {
	service service.Service
}

func New(db *pgxpool.Pool) *OrganizationHandler {
	return &OrganizationHandler{
		service: service.New(db),
	}
}

type organizationRequest struct {
	Name string `json:"name"`
}

func (o organizationRequest) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Name, validation.Required, validation.Length(3, 60)),
	)
}

func (h *OrganizationHandler) CreateOrganization(c echo.Context) error {
	ctx := util.NewContextWithLogger(c)

	var r organizationRequest
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, util.BadRequest)
	}

	if err := validation.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	organization := model.Organization{Name: r.Name}
	if err := h.service.CreateOrganization(ctx, organization); err != nil {
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}

	return c.JSON(http.StatusOK, nil)
}

type branchRequest struct {
	Name           string `json:"name"`
	PhoneNumber    string `json:"phoneNumber"`
	Address        string `json:"address"`
	OrganizationID uint   `json:"organizationId"`
}

func (o branchRequest) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Name, validation.Required, validation.Length(3, 60)),
		validation.Field(&o.PhoneNumber, validation.Required, validation.Length(10, 12)),
	)
}

func (h *OrganizationHandler) CreateBranch(c echo.Context) error {
	ctx := util.NewContextWithLogger(c)

	var r branchRequest
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, util.BadRequest)
	}

	if err := validation.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	branch := model.Branch{Name: r.Name, Address: r.Address, PhoneNumber: r.PhoneNumber, OrganizationID: r.OrganizationID}
	if err := h.service.CreateBranch(ctx, branch); err != nil {
		c.Logger().Errorf("Failed to create branches. organizationId=%v : %v", r.OrganizationID, err)
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *OrganizationHandler) GetBranch(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	branch, err := h.service.GetBranch(util.NewContextWithLogger(c), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}
	return c.JSON(http.StatusOK, branch)
}

func (h *OrganizationHandler) GetBranches(c echo.Context) error {
	branches, err := h.service.GetBranches(util.NewContextWithLogger(c))
	if err != nil {
		c.Logger().Errorf("Failed to get branches of org. %v", err)
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}
	return c.JSON(http.StatusOK, branches)
}

func (h *OrganizationHandler) GetBranchesByOrganization(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	branches, err := h.service.GetBranchesByOrganization(util.NewContextWithLogger(c), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}
	return c.JSON(http.StatusOK, branches)
}

func (h *OrganizationHandler) GetOrganizations(c echo.Context) error {
	organizations, err := h.service.GetOrganizations(util.NewContextWithLogger(c))
	if err != nil {
		c.Logger().Errorf("Failed to get branches of org. %v", err)
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}
	return c.JSON(http.StatusOK, organizations)
}
