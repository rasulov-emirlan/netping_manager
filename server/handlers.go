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
		if err := w.ToggleSocket(req.Location, req.Socket, req.Value); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, "everything is cool")
	}
}

func getAll(w *watcher.Watcher) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := w.Walk(); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.String(http.StatusOK, "everything is cool")
	}
}
