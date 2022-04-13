package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/netping-manager/watcher"
)

func setValue(w *watcher.Watcher) echo.HandlerFunc {
	type Request struct {
		Location string `json:"locationName"`
		Socket   string `json:"socketName"`
		Value    int    `json:"newValue"`
	}
	return func(c echo.Context) error {
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		v, err := w.ToggleSocket(req.Location, req.Socket, req.Value)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func getAll(w *watcher.Watcher) echo.HandlerFunc {
	return func(c echo.Context) error {
		v, err := w.Walk()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}
