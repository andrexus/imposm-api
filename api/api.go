package api

import (
	"context"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/andrexus/imposm-api/conf"
	"github.com/andrexus/imposm-api/enums"
	"github.com/andrexus/imposm-api/service"
	"github.com/labstack/echo"
	"database/sql"
)

// API is the data holder for the API
type API struct {
	config *conf.Config
	log    *logrus.Entry
	db     *sql.DB
	echo   *echo.Echo

	// Services used by the API
	imposm service.ImposmService
}

type ListResponse struct {
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
	Total    int         `json:"total"`
	Items    interface{} `json:"items"`
}

type MessageResponse struct {
	Status  enums.APIResponseStatus `json:"status"`
	Message string                  `json:"message"`
	Errors  []ErrorResponseItem     `json:"errors,omitempty"`
}

type ErrorResponseItem struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

// Start will start the API on the specified port
func (api *API) Start() error {
	return api.echo.Start(fmt.Sprintf(":%d", api.config.API.Port))
}

// Stop will shutdown the engine internally
func (api *API) Stop() error {
	logrus.Info("Stopping API server")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return api.echo.Shutdown(ctx)
}

// NewAPI will create an api instance that is ready to start
func NewAPI(config *conf.Config, db *sql.DB) *API {
	api := &API{
		config: config,
		log:    logrus.WithField("component", "api"),
		db:     db,
	}

	api.imposm = service.NewImposmService(service.NewImposmRepository(db))

	// add the endpoints
	e := echo.New()
	e.HideBanner = true
	//e.Use(api.logRequest)

	g := e.Group("/api/v1")

	// transport points
	g.POST("/transport-points", api.SearchNearbyTransportPoints)
	g.POST("/streets", api.SearchNearbyStreets)

	api.echo = e

	return api
}

func (api *API) logRequest(f echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		req := ctx.Request()
		logger := api.log.WithFields(logrus.Fields{
			"method": req.Method,
			"path":   req.URL.Path,
		})
		ctx.Set(loggerKey, logger)

		logger.WithFields(logrus.Fields{
			"user_agent": req.UserAgent(),
			"ip_address": ctx.RealIP(),
		}).Info("Request")

		err := f(ctx)
		if err != nil {
			ctx.Error(err)
		}
		return err
	}
}
