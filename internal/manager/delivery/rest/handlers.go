package rest

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/netping-manager/internal/manager"
)

type handler struct {
	service manager.Service
}

func NewHandler(service manager.Service) (*handler, error) {
	if service == nil {
		return nil, errors.New("manager: delivery: service can't be nil")
	}
	return &handler{
		service: service,
	}, nil
}

func (h *handler) Register(router *echo.Group) error {
	router.POST("/config/location", h.addLocation())
	router.DELETE("/config/location", h.removeLocation())
	router.POST("/config/socket", h.addSocket())
	router.DELETE("/config/socket", h.removeSocket())

	router.POST("/control", h.setValue())
	router.GET("/control", h.getAll())
	return nil
}

func (h *handler) setValue() echo.HandlerFunc {
	type Request struct {
		Location int `json:"locationID"`
		Socket   int `json:"socketID"`
		Value    int `json:"newValue"`
	}
	return func(c echo.Context) error {
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		log.Println(req)
		v, err := h.service.ToggleSocket(c.Request().Context(), req.Socket, req.Location, req.Value)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func (h *handler) getAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		v, err := h.service.CheckAll(c.Request().Context())
		log.Println(err)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func (h *handler) addLocation() echo.HandlerFunc {
	type Request struct {
		Name         string `json:"name"`
		RealLocation string `json:"realLocation"`

		SNMPaddress   string `json:"snmpAddress"`
		SNMPport      int    `json:"snmpPort"`
		SNMPcommunity string `json:"snmpCommunity"`
	}
	return func(c echo.Context) error {
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		l := toServiceLocation(req.Name, req.RealLocation, req.SNMPaddress, req.SNMPcommunity, req.SNMPport)
		v, err := h.service.AddLocation(c.Request().Context(), l)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func (h *handler) removeLocation() echo.HandlerFunc {
	return func(c echo.Context) error {
		location := c.QueryParam("location")
		if location == "" {
			return c.JSON(http.StatusBadRequest, "There is no query param for 'name'")
		}
		id, err := strconv.Atoi(location)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		v, err := h.service.RemoveLocation(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func (h *handler) addSocket() echo.HandlerFunc {
	type Request struct {
		LocationID int    `json:"locationID"`
		SocketName string `json:"socketName"`
		SocketMIB  string `json:"socketMIB"`
	}
	return func(c echo.Context) error {
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		s := toServiceSocket(req.SocketName, req.SocketMIB)
		v, err := h.service.AddSocket(c.Request().Context(), s, req.LocationID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func (h *handler) removeSocket() echo.HandlerFunc {
	return func(c echo.Context) error {
		location := c.QueryParam("location")
		socket := c.QueryParam("socket")
		if location == "" || socket == "" {
			return c.JSON(http.StatusBadRequest, "This endpoint needs locationName and socketName in query params")
		}
		locationId, err := strconv.Atoi(location)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		socketId, err := strconv.Atoi(location)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		v, err := h.service.RemoveSocket(c.Request().Context(), socketId, locationId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}
