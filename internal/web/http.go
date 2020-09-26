package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/skvoch/burlington-backend/tree/master/internal/models"
	"github.com/skvoch/burlington-backend/tree/master/internal/service"
	"github.com/skvoch/burlington-backend/tree/master/internal/qr"
	"golang.org/x/net/context"
	"net/http"
)

type serviceHTTP struct {
	bindURL      string
	router       *gin.Engine
	server       *http.Server

	logger zerolog.Logger
}

func New(bindUrl string, logger zerolog.Logger)*serviceHTTP{
	res := serviceHTTP{
		router: gin.Default(),
		bindURL: bindUrl,
		logger: logger,
	}
	res.registerHandlers()
	return &res
}

func (t *serviceHTTP) registerHandlers() {

	api := t.router.Group("/api")

	api.Handle(http.MethodGet, "/area", t.getArea)
	api.Handle(http.MethodPost, "/area", t.setArea)
	api.Handle(http.MethodGet, "/entity", t.setArea)
	api.Handle(http.MethodPost, "/entity", t.setArea)
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
	//id := c.Params.ByName("id")
	id, _ := c.Params.Get("id")
	fmt.Print(id)
	area, err := service.GetArea(id)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{})
		t.logger.Err(err).Msgf("id is %v", id)
	}else{
		c.JSON(http.StatusOK, area)
	}
}

func (t *serviceHTTP) setArea(c *gin.Context){
	var json models.Area
	if err := c.BindJSON(json); err != nil{
		t.logger.Err(err)
		return
	}
	res, err := service.SetAreaModel(json)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{})
		t.logger.Err(err)
	}else{
		c.JSON(http.StatusOK, res)
	}
}

func (t *serviceHTTP) getEntity(c *gin.Context){
	//id := c.Params.ByName("id")
	id, _ := c.Params.Get("id")
	fmt.Print(id)
	entity, err := service.GetEntity(id)
	if err != nil{
		c.JSON(http.StatusNotFound, gin.H{})
		t.logger.Err(err).Msgf("id is %v", id)
	}else{
		c.JSON(http.StatusOK, entity)
	}
}

func (t * serviceHTTP) setEntity(c *gin.Context){
	var json models.Entity
	if err := c.BindJSON(json); err != nil{
		t.logger.Err(err)
		return
	}
	res, err := service.SetEntity(json)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{})
		t.logger.Err(err)
	}else{
		c.JSON(http.StatusOK, res)
	}
}

func (t *serviceHTTP) getQr(c *gin.Context){
	id := c.Params.Get("id")
	img, err := qr.Generate(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
	}else{
		c.JSON(http.StatusOK, img)
	}
}