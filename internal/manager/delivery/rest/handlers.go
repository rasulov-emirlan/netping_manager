package rest

import (
	"errors"
	"net/http"
	"strconv"
	"time"

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
	router.GET("/info", h.getSocketTypes())
	router.POST("/config/socket", h.addSocket())
	router.DELETE("/config/socket/:id", h.removeSocket())
	router.GET("/config/sockets", h.listSockets())

	router.POST("/control", h.setValue())
	router.GET("/control", h.getAll())
	return nil
}

func (h *handler) setValue() echo.HandlerFunc {
	type Request struct {
		Socket   int `json:"socketID"`
		TurnOn    bool `json:"turnON"`
	}
	return func(c echo.Context) error {
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		go func () {
			err := h.service.ToggleSocket(c.Request().Context(), req.Socket, req.TurnOn)
			if err != nil {
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		}()
		select {
		case <- time.After(5 * time.Second):
			return c.JSON(http.StatusInternalServerError, "looks like we could not connect to netping")
		case <- c.Request().Context().Done():
		}
		return c.NoContent(http.StatusOK)
	}
}

func (h *handler) getAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		locationAddress := c.QueryParam("location")
		if locationAddress == "" {
			return c.JSON(http.StatusBadRequest, "you need location in query params")
		}
		v, err := h.service.CheckAll(c.Request().Context(), locationAddress)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func (h *handler) addSocket() echo.HandlerFunc {
	type Request struct {
		LocationAddress string `json:"locationAddress"`
		SocketName 		string `json:"socketName"`
		SocketMIB  		string `json:"socketMIB"`
		SocketType 		int    `json:"socketType"`
	}
	return func(c echo.Context) error {
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if req.SocketType < 1 || req.SocketType > 3 {
			return c.JSON(http.StatusBadRequest, "Incorrect socket type")
		}
		s := toServiceSocket(req.SocketName, req.SocketMIB, req.SocketType)
		v, err := h.service.AddSocket(c.Request().Context(), s, req.LocationAddress)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func (h *handler) removeSocket() echo.HandlerFunc {
	return func(c echo.Context) error {
		socket := c.Param("id")
		if socket == "" {
			return c.JSON(http.StatusBadRequest, "This endpoint needs locationName and socketName in query params")
		}
		socketId, err := strconv.Atoi(socket)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err = h.service.RemoveSocket(c.Request().Context(), socketId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.NoContent(http.StatusOK)
	}
}

func (h *handler) listSockets() echo.HandlerFunc {
	return func(c echo.Context) error {
		v, err := h.service.ListAllSockets(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func (h *handler) getSocketTypes() echo.HandlerFunc {
	type Response struct {
		TypeID int `json:"typeID"`
		TypeName string `jons:"name"`
	}
	return func(c echo.Context) error {
		// TODO: change this code to an actual service call
		resp := []*Response{
			{TypeID: manager.TypeUnknown, TypeName: "unknown"},
			{TypeID: manager.TypeGenerator, TypeName: "generator"},
			{TypeID: manager.TypeAC, TypeName: "air conditioner"},
			{TypeID: manager.TypeHeater, TypeName: "heater"},
		}
		return c.JSON(http.StatusOK, resp)
	}
}