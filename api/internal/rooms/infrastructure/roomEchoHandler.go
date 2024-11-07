package infrastructure

import (
	"net/http"
	"strconv"
	u "suffgo/internal/rooms/application/useCases"

	d "suffgo/internal/rooms/domain"
	v "suffgo/internal/rooms/domain/valueObjects"

	sv "suffgo/internal/shared/domain/valueObjects"

	se "suffgo/internal/shared/domain/errors"

	"github.com/labstack/echo/v4"
)

type RoomEchoHandler struct {
	CreateRoomUsecase  *u.CreateUsecase
	DeleteRoomUsecase  *u.DeleteUsecase
	GetAllUsecase      *u.GetAllUsecase
	GetRoomByIDUsecase *u.GetByIDUsecase
	GetByAdminUsecase  *u.GetByAdminUsecase
}

func NewRoomEchoHandler(
	creatUC *u.CreateUsecase,
	deleteUC *u.DeleteUsecase,
	getAllUC *u.GetAllUsecase,
	getByIDUC *u.GetByIDUsecase,
	getByAdminUC *u.GetByAdminUsecase,
) *RoomEchoHandler {
	return &RoomEchoHandler{
		CreateRoomUsecase:  creatUC,
		DeleteRoomUsecase:  deleteUC,
		GetAllUsecase:      getAllUC,
		GetRoomByIDUsecase: getByIDUC,
		GetByAdminUsecase:  getByAdminUC,
	}
}

func (h *RoomEchoHandler) CreateRoom(c echo.Context) error {
	var req d.RoomCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	linkInvite, err := v.NewLinkInvite(req.LinkInvite)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	isFormal, err := v.NewIsFormal(req.IsFormal)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	name, err := v.NewName(req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	adminID, err := sv.NewID(req.AdminID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	room := d.NewRoom(
		nil,
		*linkInvite,
		*isFormal,
		*name,
		adminID,
	)

	err = h.CreateRoomUsecase.Execute(*room)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusCreated, req)
}

func (h *RoomEchoHandler) DeleteRoom(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"message": invalidErr.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.DeleteRoomUsecase.Execute(*id)
	if err != nil {
		if err.Error() == "room not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]string{"success": "room deleted succesfully"})
}

func (h *RoomEchoHandler) GetAllRooms(c echo.Context) error {
	rooms, err := h.GetAllUsecase.Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	var roomsDTO []d.RoomDTO
	for _, room := range rooms {
		roomDTO := &d.RoomDTO{
			ID:         room.ID().Id,
			LinkInvite: room.LinkInvite().LinkInvite,
			IsFormal:   room.IsFormal().IsFormal,
			Name:       room.Name().Name,
			AdminID:    room.AdminID().Id,
		}
		roomsDTO = append(roomsDTO, *roomDTO)
	}
	return c.JSON(http.StatusOK, roomsDTO)
}

func (h *RoomEchoHandler) GetRoomByID(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"message": invalidErr.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	room, err := h.GetRoomByIDUsecase.Execute(*id)
	if err != nil {
		if err.Error() == "room not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"message": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	roomDTO := &d.RoomDTO{
		ID:         room.ID().Id,
		LinkInvite: room.LinkInvite().LinkInvite,
		IsFormal:   room.IsFormal().IsFormal,
		Name:       room.Name().Name,
		AdminID:    room.AdminID().Id,
	}
	return c.JSON(http.StatusOK, roomDTO)
}

func (h *RoomEchoHandler) GetRoomsByAdmin(c echo.Context) error {

	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"message": invalidErr.Error()})
	}
	id, _ := sv.NewID(uint(idInput))
	rooms, err := h.GetByAdminUsecase.Execute(*id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	var roomsDTO []d.RoomDTO
	for _, room := range rooms {
		roomDTO := &d.RoomDTO{
			ID:         room.ID().Id,
			LinkInvite: room.LinkInvite().LinkInvite,
			IsFormal:   room.IsFormal().IsFormal,
			Name:       room.Name().Name,
			AdminID:    room.AdminID().Id,
		}
		roomsDTO = append(roomsDTO, *roomDTO)
	}
	return c.JSON(http.StatusOK, roomsDTO)
}
