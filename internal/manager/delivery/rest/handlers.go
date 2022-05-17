package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rasulov-emirlan/netping-manager/internal/manager"

	usersHelpers "github.com/rasulov-emirlan/netping-manager/internal/users/delivery/rest"

	"go.uber.org/zap"
)

type handler struct {
	service manager.Service
	log     *zap.SugaredLogger
	jwtkey  []byte
}

func NewHandler(service manager.Service, jwkey []byte, logger *zap.SugaredLogger) (*handler, error) {
	if service == nil {
		return nil, errors.New("manager: delivery: service can't be nil")
	}
	return &handler{
		service: service,
		log:     logger,
		jwtkey:  jwkey,
	}, nil
}

func (h *handler) Register(router *echo.Group) error {
	router.GET("/config/socket/type", h.getSocketTypes())
	router.POST("/config/socket", h.addSocket(), usersHelpers.CheckRole(h.jwtkey, true))
	router.PATCH("/config/socket/:id", h.updateSocket(), usersHelpers.CheckRole(h.jwtkey, true))
	router.DELETE("/config/socket/:id", h.removeSocket(), usersHelpers.CheckRole(h.jwtkey, true))
	router.GET("/config/sockets", h.listSockets(), usersHelpers.CheckRole(h.jwtkey, false))

	router.POST("/control", h.setValue(), usersHelpers.CheckRole(h.jwtkey, false))
	router.GET("/control/:id", h.getAll(), usersHelpers.CheckRole(h.jwtkey, false))
	return nil
} 

func (h *handler) setValue() echo.HandlerFunc {
	type Request struct {
		Socket int  `json:"socketID"`
		TurnOn bool `json:"turnON"`
	}
	return func(c echo.Context) error {
		defer h.log.Sync()
		claims, ok := c.Get(usersHelpers.UserInfoFromContext).(*usersHelpers.Claims)
		if !ok {
			return c.NoContent(http.StatusBadRequest)
		}
		h.log.Infow("Trying to toggle a socket", zap.Int("userid", claims.UserID))
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := h.service.ToggleSocket(c.Request().Context(), req.Socket, req.TurnOn)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		h.log.Infow("Toggled a socket", zap.Int("userid", claims.UserID))
		return c.NoContent(http.StatusOK)
	}
}

func (h *handler) getAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		location := c.Param("id")
		if location == "" {
			return c.JSON(http.StatusBadRequest, "you need location in query params")
		}
		id, err := strconv.Atoi(location)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		v, err := h.service.CheckAll(c.Request().Context(), id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, v)
	}
}

func (h *handler) addSocket() echo.HandlerFunc {
	type Request struct {
		Location   int    `json:"locationID"`
		SocketName string `json:"socketName"`
		SocketMIB  string `json:"socketMIB"`
		SocketType int    `json:"socketType"`
	}
	return func(c echo.Context) error {
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if req.SocketType < 1 || req.SocketType > 4 {
			return c.JSON(http.StatusBadRequest, "Incorrect socket type")
		}
		s := toServiceSocket(req.SocketName, req.SocketMIB, req.SocketType)
		v, err := h.service.AddSocket(c.Request().Context(), s, req.Location)
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
		TypeID   int    `json:"typeID"`
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

func (h *handler) updateSocket() echo.HandlerFunc {
	type Request struct {
		Name          string `json:"name"`
		SNMPaddress   string `json:"snmpAddress"`
		SNMPcommunity string `json:"snmpCommunity"`
		SNMPport      int    `json:"snmpPort"`
		SNMPmib       string `json:"snmpMib"`
		ObjectType    int    `json:"objectType"`
	}
	return func(c echo.Context) error {
		tmp := c.Param("id")
		id, err := strconv.Atoi(tmp)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		req := &Request{}
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		sock := toServiceSocket(req.Name, req.SNMPmib, req.ObjectType)
		sock.SNMPaddress = req.SNMPaddress
		s, err := h.service.UpdateSocket(c.Request().Context(), sock, id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, s)
	}
}
