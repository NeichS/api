package infrastructure

import (
	"github.com/labstack/echo/v4"
	userInfr "suffgo/internal/users/infrastructure"
)

func InitializeRoomEchoRouter(e *echo.Echo, handler *RoomEchoHandler) {

	roomGroup := e.Group("/v1/rooms")
	roomGroup.GET("", handler.GetAllRooms)
	roomGroup.Use(userInfr.AuthMiddleware)
	roomGroup.POST("", handler.CreateRoom)
	roomGroup.DELETE("/:id", handler.DeleteRoom)
	roomGroup.GET("/:id", handler.GetRoomByID)
	roomGroup.GET("/myRooms", handler.GetRoomsByAdmin)
}
