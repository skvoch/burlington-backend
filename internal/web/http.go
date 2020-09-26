package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/skvoch/burlington-backend/tree/master/internal/models"
	"github.com/skvoch/burlington-backend/tree/master/internal/qr"
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
	//api.Handle(http.MethodGet, "/entity", t.setArea)
	//api.Handle(http.MethodPost, "/entity", t.setArea)
	api.Handle(http.MethodGet, "/qr", t.getQr)

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
	id := c.Query("id")

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
//
//func (t *serviceHTTP) getEntity(c *gin.Context){
//	//id := c.Params.ByName("id")
//	id, _ := c.Params.Get("id")
//	fmt.Print(id)
//	entity, err := service.GetEntity(id)
//	if err != nil{
//		c.JSON(http.StatusNotFound, gin.H{})
//		t.logger.Err(err).Msgf("id is %v", id)
//	}else{
//		c.JSON(http.StatusOK, entity)
//	}
//}
//
//func (t * serviceHTTP) setEntity(c *gin.Context){
//	var json models.Entity
//	if err := c.BindJSON(json); err != nil{
//		t.logger.Err(err)
//		return
//	}
//	res, err := t.service.SetEntity(json)
//	if err != nil{
//		c.JSON(http.StatusBadRequest, gin.H{})
//		t.logger.Err(err)
//	}else{
//		c.JSON(http.StatusOK, res)
//	}
//}

func (t *serviceHTTP) getQr(c *gin.Context){
	id := c.Query("id")
	img, err := qr.Generate(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &response{
			Error: "internal service error",
		})
	}else{
		c.JSON(http.StatusOK, img)
	}
}