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
	Error string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func New(bindUrl string, logger zerolog.Logger, service *service.Service)*serviceHTTP{
	res := serviceHTTP{
		router: gin.Default(),
		bindURL: bindUrl,
		logger: logger,
		service: service,
	}
	res.registerHandlers()
	return &res
}

type serviceHTTP struct {
	bindURL      string
	router       *gin.Engine
	server       *http.Server
	service *service.Service

	logger zerolog.Logger
}

func (t *serviceHTTP) registerHandlers() {

	api := t.router.Group("/api")

	api.Handle(http.MethodGet, "/area", t.getArea)
	api.Handle(http.MethodPost, "/area", t.setArea)
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

func (t *serviceHTTP) getArea(c *gin.Context) {
	id, _ := c.Params.Get("id")

	area, err := t.service.GetArea(id)

	if err != nil{
		c.JSON(http.StatusNotFound, &response{
			Error: "area not found",
		})
		return
	}

	c.JSON(http.StatusOK, area)
}

func (t *serviceHTTP) setArea(c *gin.Context){
	var area models.Area

	if err := c.BindJSON(&area); err != nil{
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
