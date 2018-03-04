package api

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/andrexus/imposm-api/service"
	"github.com/andrexus/imposm-api/enums"
)

func (api *API) SearchNearbyTransportPoints(ctx echo.Context) error {
	r := new(service.NearbyPointsSearchRequest)
	if err := ctx.Bind(r); err != nil {
		response := &MessageResponse{Status: enums.Error, Message: err.Error()}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	items, err := api.imposm.FindNearbyTransportPoints(r)
	if err != nil {
		response := &MessageResponse{Status: enums.Error, Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response := &ListResponse{Page: 1, PageSize: len(items), Total: len(items), Items: items}
	return ctx.JSON(http.StatusOK, response)
}

func (api *API) SearchNearbyStreets(ctx echo.Context) error {
	r := new(service.NearbyPointsSearchRequest)
	if err := ctx.Bind(r); err != nil {
		response := &MessageResponse{Status: enums.Error, Message: err.Error()}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	items, err := api.imposm.FindNearbyStreets(r)
	if err != nil {
		response := &MessageResponse{Status: enums.Error, Message: err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response := &ListResponse{Page: 1, PageSize: len(items), Total: len(items), Items: items}
	return ctx.JSON(http.StatusOK, response)
}
