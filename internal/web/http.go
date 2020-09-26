package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/skvoch/burlington-backend/tree/master/internal/models"
	"github.com/skvoch/burlington-backend/tree/master/internal/service"
	"golang.org/x/net/context"
	"net/http"
)

type response struct {
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func New(bindUrl string, logger zerolog.Logger, service *service.Service) *serviceHTTP {
	res := serviceHTTP{
		router:  gin.Default(),
		bindURL: bindUrl,
		logger:  logger,
		service: service,
	}
	res.registerHandlers()
	return &res
}

type serviceHTTP struct {
	bindURL string
	router  *gin.Engine
	server  *http.Server
	service *service.Service

	logger zerolog.Logger
}

func (t *serviceHTTP) registerHandlers() {

	api := t.router.Group("/api")

	api.Handle(http.MethodGet, "/area", t.getArea)
	api.Handle(http.MethodPost, "/area", t.setArea)
	api.Handle(http.MethodPost, "/find_path", t.findPath)
}

func (t *serviceHTTP) Run() {
	t.server = &http.Server{
		Addr:    t.bindURL,
		Handler: t.router,
	}

	go func() {
		if err := t.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			t.logger.Fatal().Err(err).Msg("failed to start listening http")
		}
	}()
}

func (t *serviceHTTP) Stop() error {
	t.logger.Info().Msg("shutdown http server...")

	if err := t.server.Shutdown(context.Background()); err != nil {
		return err
	}

	return nil
}

func (t *serviceHTTP) findPath(ctx *gin.Context) {
	type Request struct {
		AreaName string `json:"area_name"`
		FromX    int64  `json:"from_x"`
		FromY    int64  `json:"from_y"`
		FromZ    int64  `json:"from_z"`

		ToX int64 `json:"to_x"`
		ToY int64 `json:"to_y"`
		ToZ int64 `json:"to_z"`
	}

	var requestData Request

	if err := ctx.BindJSON(&requestData); err != nil {
		ctx.JSON(http.StatusNotFound, &response{
			Error: "failed to bind json",
		})
		return
	}
	result, err := t.service.FindPath(requestData.AreaName,
		models.XYZ{
			X: requestData.FromX,
			Y: requestData.FromY,
			Z: requestData.FromZ,
		},
		models.XYZ{
			X: requestData.ToX,
			Y: requestData.ToY,
			Z: requestData.ToZ,
		})
	if err != nil {
		ctx.JSON(http.StatusNotFound, &response{
			Error: err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (t *serviceHTTP) getArea(c *gin.Context) {
	id, _ := c.Params.Get("id")

	area, err := t.service.GetArea(id)

	if err != nil {
		c.JSON(http.StatusNotFound, &response{
			Error: "area not found",
		})
		return
	}

	c.JSON(http.StatusOK, area)
}

func (t *serviceHTTP) setArea(c *gin.Context) {
	var area models.Area

	if err := c.BindJSON(&area); err != nil {
		t.logger.Error().Err(err).Msg("binding json to object")
		c.JSON(http.StatusBadRequest, &response{
			Error: "incorrect input json",
		})
		return
	}

	if err := t.service.SetArea(area); err != nil {
		c.JSON(http.StatusInternalServerError, &response{
			Error: "internal service error",
		})
		return
	}
	c.JSON(http.StatusOK, response{
		Message: "ok",
	})
}
