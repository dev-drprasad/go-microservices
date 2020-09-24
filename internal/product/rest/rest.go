package rest

import (
	"crypto/rand"
	"fmt"
	"gomicroservices/internal/product/model"
	"gomicroservices/internal/product/service"
	"gomicroservices/internal/util"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
)

const maxUploadSize = 1 * 1024 * 1024
const uploadPath = "./uploads"

type ProductHandler struct {
	service service.Service
}

func New(db *pgxpool.Pool) *ProductHandler {
	return &ProductHandler{
		service: service.New(db),
	}
}

type brandRequest struct {
	Name string `json:"name"`
}

func (o brandRequest) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Name, validation.Required, validation.Length(2, 20)),
	)
}

func (h *ProductHandler) CreateBrand(c echo.Context) error {
	ctx := util.NewContextWithLogger(c)

	var r brandRequest
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, util.BadRequest)
	}

	if err := validation.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	branch := model.Brand{Name: r.Name}
	if err := h.service.CreateBrand(ctx, branch); err != nil {
		c.Logger().Errorf("Failed to create brand : %v", err)
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *ProductHandler) GetBrand(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	branch, err := h.service.GetBrand(util.NewContextWithLogger(c), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}
	return c.JSON(http.StatusOK, branch)
}

func (h *ProductHandler) GetBrands(c echo.Context) error {
	branches, err := h.service.GetBrands(util.NewContextWithLogger(c))
	if err != nil {
		c.Logger().Errorf("Failed to get brands : %v", err)
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}
	return c.JSON(http.StatusOK, branches)
}

type categoryRequest struct {
	Name string `json:"name"`
}

func (o categoryRequest) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Name, validation.Required, validation.Length(2, 60)),
	)
}

func (h *ProductHandler) CreateCategory(c echo.Context) error {
	ctx := util.NewContextWithLogger(c)

	var r categoryRequest
	if err := c.Bind(&r); err != nil {
		return c.JSON(http.StatusBadRequest, util.BadRequest)
	}

	if err := validation.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	category := model.Category{Name: r.Name}
	if err := h.service.CreateCategory(ctx, category); err != nil {
		c.Logger().Errorf("Failed to create category : %v", err)
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *ProductHandler) GetCategory(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	category, err := h.service.GetCategory(util.NewContextWithLogger(c), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}
	return c.JSON(http.StatusOK, category)
}

func (h *ProductHandler) GetCategories(c echo.Context) error {
	categories, err := h.service.GetCategories(util.NewContextWithLogger(c))
	if err != nil {
		c.Logger().Errorf("Failed to get categories : %v", err)
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}
	return c.JSON(http.StatusOK, categories)
}

type productRequest struct {
	Name       string       `json:"name"`
	Cost       util.Float64 `json:"cost"`
	SellPrice  util.Float64 `json:"sellPrice"`
	BrandID    uint         `json:"brandId"`
	CategoryID uint         `json:"categoryId"`
	ImageURLs  []string     `json:"imageUrls"`
	Stock      uint         `json:"stock"`
}

func (o productRequest) Validate() error {
	return validation.ValidateStruct(&o,
		validation.Field(&o.Name, validation.Required, validation.Length(2, 60)),
		validation.Field(&o.BrandID, validation.Required, validation.Min(uint(1))),
		validation.Field(&o.CategoryID, validation.Required, validation.Min(uint(1))),
		validation.Field(&o.Cost, validation.Required, validation.Min(0.1), validation.Max(999999.99)),
		validation.Field(&o.SellPrice, validation.Required, validation.Min(0.1), validation.Max(999999.99)),
		validation.Field(&o.Stock, validation.Required, validation.Min(uint(1))),
	)
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	ctx := util.NewContextWithLogger(c)

	var r productRequest
	if err := c.Bind(&r); err != nil {
		c.Logger().Errorf("Failed to bind product payload : %v", err)
		return c.JSON(http.StatusBadRequest, util.BadRequest)
	}

	if err := validation.Validate(r); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	product := model.Product{Name: r.Name, Cost: float64(r.Cost), SellPrice: float64(r.SellPrice), BrandID: r.BrandID, CategoryID: r.CategoryID, ImageURLs: []string{}, Stock: uint(r.Stock)}
	if err := h.service.CreateProduct(ctx, product); err != nil {
		c.Logger().Errorf("Failed to create product : %v", err)
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}

	return c.JSON(http.StatusOK, nil)
}

func (h *ProductHandler) GetProduct(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	category, err := h.service.GetProduct(util.NewContextWithLogger(c), uint(id))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}
	return c.JSON(http.StatusOK, category)
}

func (h *ProductHandler) GetProducts(c echo.Context) error {
	categories, err := h.service.GetProducts(util.NewContextWithLogger(c))
	if err != nil {
		c.Logger().Errorf("Failed to get products : %v", err)
		return c.JSON(http.StatusInternalServerError, util.Internal)
	}
	return c.JSON(http.StatusOK, categories)
}

func (h *ProductHandler) UploadImages(c echo.Context) error {
	if err := c.Request().ParseMultipartForm(maxUploadSize); err != nil {

		c.Logger().Errorf("Could not parse multipart form: %v", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "could not parse form"})
	}

	files := c.Request().MultipartForm.File["files"]

	filenames := []string{}
	for _, file := range files {
		if file.Size > maxUploadSize {
			c.Logger().Errorf("File is too big")
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "file is too big"})
		}

		src, err := file.Open()
		if err != nil {
			c.Logger().Errorf("Failed to open file : %v", err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		defer src.Close()

		fileBytes, err := ioutil.ReadAll(src)
		if err != nil {
			c.Logger().Errorf("Failed to read file : %v", err)
			return c.JSON(http.StatusInternalServerError, nil)
		}

		filetype := http.DetectContentType(fileBytes)
		if filetype != "image/jpeg" && filetype != "image/jpg" {
			c.Logger().Errorf("Expected jpeg/jpg, but got %v", filetype)
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "only jpeg/jpg is allowed"})
		}

		b := make([]byte, 10)
		if _, err = rand.Read(b); err != nil {
			c.Logger().Errorf("Failed to generate random name : %v", err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		filename := fmt.Sprintf("%x%d%s", b, time.Now().UnixNano(), path.Ext(file.Filename))

		dst, err := os.Create(path.Join(uploadPath, filename))
		if err != nil {
			c.Logger().Errorf("Failed to create dest file : %v", err)
			return c.JSON(http.StatusInternalServerError, nil)
		}
		defer dst.Close()

		_, err = dst.Write(fileBytes)
		if err != nil {
			c.Logger().Errorf("Failed to write to dest file : %v", err)
			return c.JSON(http.StatusInternalServerError, nil)
		}

		filenames = append(filenames, filename)
	}

	return c.JSON(http.StatusOK, filenames)
}
