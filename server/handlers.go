package server

import (
	"log"
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
		log.Println(req)
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

func addLocation(w *watcher.Watcher) echo.HandlerFunc {
	type Request struct {
		LocationName string `json:"locationName"`
		Address      string `json:"address"`
		Community    string `json:"community"`
		Port         int    `json:"port"`
	}
	return func(c echo.Context) error {
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		v, err := w.AddLocation(req.LocationName, req.Address, req.Community, req.Port)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func removeLocation(w *watcher.Watcher) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.QueryParam("name")
		if name == "" {
			return c.JSON(http.StatusBadRequest, "There is no query param for 'name'")
		}
		v, err := w.RemoveLocation(name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func addSocket(w *watcher.Watcher) echo.HandlerFunc {
	type Request struct {
		LocationName string `json:"locationName"`
		SocketName   string `json:"socketName"`
		SocketMIB    string `json:"socketMIB"`
	}
	return func(c echo.Context) error {
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		v, err := w.AddSocket(req.LocationName, req.SocketName, req.SocketMIB)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func removeSocket(w *watcher.Watcher) echo.HandlerFunc {
	return func(c echo.Context) error {
		locationName := c.QueryParam("locationName")
		socketName := c.QueryParam("socketName")
		if locationName == "" || socketName == "" {
			return c.JSON(http.StatusBadRequest, "This endpoint needs locationName and socketName in query params")
		}
		v, err := w.RemoveSocket(locationName, socketName)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}
