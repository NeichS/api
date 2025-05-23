package infrastructure

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	u "suffgo/internal/users/application/useCases"

	d "suffgo/internal/users/domain"
	v "suffgo/internal/users/domain/valueObjects"

	sv "suffgo/internal/shared/domain/valueObjects"

	se "suffgo/internal/shared/domain/errors"

	uerr "suffgo/internal/users/domain/errors"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type UserEchoHandler struct {
	CreateUserUsecase     *u.CreateUsecase
	DeleteUserUsecase     *u.DeleteUsecase
	GetAllUsersUsecase    *u.GetAllUsecase
	GetUserByIDUsecase    *u.GetByIDUsecase
	GetUserByEmailUsecase *u.GetByEmailUsecase
	LoginUsecase          *u.LoginUsecase
	RestoreUsecase        *u.RestoreUsecase
	ChangePasswordUsecase *u.ChangePassword
	UpdateUsecase         *u.UpdateUsecase
	GetUsersByRoomUsecase *u.GetUsersByRoom
}

// Constructor for UserEchoHandler
func NewUserEchoHandler(
	createUC *u.CreateUsecase,
	deleteUC *u.DeleteUsecase,
	getAllUC *u.GetAllUsecase,
	getByIDUC *u.GetByIDUsecase,
	getByEmailUC *u.GetByEmailUsecase,
	loginUC *u.LoginUsecase,
	restoreUC *u.RestoreUsecase,
	changePassUC *u.ChangePassword,
	updateUC *u.UpdateUsecase,
	getByRoomUC *u.GetUsersByRoom,
) *UserEchoHandler {
	return &UserEchoHandler{
		CreateUserUsecase:     createUC,
		DeleteUserUsecase:     deleteUC,
		GetAllUsersUsecase:    getAllUC,
		GetUserByIDUsecase:    getByIDUC,
		GetUserByEmailUsecase: getByEmailUC,
		LoginUsecase:          loginUC,
		RestoreUsecase:        restoreUC,
		ChangePasswordUsecase: changePassUC,
		UpdateUsecase:         updateUC,
		GetUsersByRoomUsecase: getByRoomUC,
	}
}

func (u *UserEchoHandler) Login(c echo.Context) error {

	var req d.LoginRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	username, err := v.NewUserName(req.Username)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	pass, err := v.NewPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := u.LoginUsecase.Execute(*username, *pass)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	if err := createSession(user.ID(), user.FullName().Name, c); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userDTO := &d.UserSafeDTO{
		ID:       user.ID().Id,
		Name:     user.FullName().Name,
		Lastname: user.FullName().Lastname,
		Username: user.Username().Username,
		Dni:      user.Dni().Dni,
		Email:    user.Email().Email,
		Image:    user.Image().URL(),
	}

	response := map[string]interface{}{
		"success": "autenticación exitosa",
		"user":    userDTO,
	}

	// Devuelvo el id del usuario logueado
	return c.JSON(http.StatusOK, response)
}

func (u *UserEchoHandler) CreateUser(c echo.Context) error {
	var req d.UserCreateRequest
	// bindea el body del request (json) al dto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	fullname, err := v.NewFullName(req.Name, req.Lastname)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	username, err := v.NewUserName(req.Username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	dni, err := v.NewDni(req.Dni)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	email, err := v.NewEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	password, err := v.NewPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	hashed, err := v.HashPassword(password.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	image, err := v.NewImage(req.Image)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	// Map DTO to domain entity

	user := d.NewUser(
		nil,
		*fullname,
		*username,
		*dni,
		*email,
		*hashed,
		image,
	)

	// Call the use case
	user, err = u.CreateUserUsecase.Execute(*user)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	userDTO := &d.UserSafeDTO{
		ID:       user.ID().Id,
		Name:     user.FullName().Name,
		Lastname: user.FullName().Lastname,
		Username: user.Username().Username,
		Dni:      user.Dni().Dni,
		Email:    user.Email().Email,
		Image:    user.Image().URL(),
	}

	response := map[string]interface{}{
		"success": "usuario creado exitosamente",
		"user":    userDTO,
	}

	return c.JSON(http.StatusCreated, response)
}

func (u *UserEchoHandler) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	id, _ := sv.NewID(uint(idInput))

	currentUserID, err := GetAuthenticatedUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	err = u.DeleteUserUsecase.Execute(*id, *currentUserID)
	if err != nil {
		if errors.Is(err, uerr.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "user deleted succesfully"})
}

func (u *UserEchoHandler) GetAllUsers(c echo.Context) error {
	users, err := u.GetAllUsersUsecase.Execute()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var usersDTO []d.UserSafeDTO
	for _, user := range users {
		userDTO := &d.UserSafeDTO{
			ID:       user.ID().Id,
			Name:     user.FullName().Name,
			Lastname: user.FullName().Lastname,
			Username: user.Username().Username,
			Dni:      user.Dni().Dni,
			Email:    user.Email().Email,
			Image:    user.Image().URL(),
		}
		usersDTO = append(usersDTO, *userDTO)
	}

	return c.JSON(http.StatusOK, usersDTO)
}

func (u *UserEchoHandler) GetUserByID(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	user, err := u.GetUserByIDUsecase.Execute(*id)
	if err != nil {
		if errors.Is(err, uerr.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	userDTO := &d.UserSafeDTO{
		ID:       user.ID().Id,
		Name:     user.FullName().Name,
		Lastname: user.FullName().Lastname,
		Username: user.Username().Username,
		Dni:      user.Dni().Dni,
		Email:    user.Email().Email,
		Image:    user.Image().URL(),
	}
	return c.JSON(http.StatusOK, userDTO)
}

func (u *UserEchoHandler) GetUserByEmail(c echo.Context) error {

	var request struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.Bind(&request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format: " + err.Error(),
		})
	}

	// Validar el formato del email
	email, err := v.NewEmail(request.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid email format: " + err.Error(),
		})
	}

	// Buscar el usuario por email
	user, err := u.GetUserByEmailUsecase.Execute(*email)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Error retrieving user: " + err.Error(),
		})
	}

	if user == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "User not found with email: " + request.Email,
		})
	}

	userDTO := &d.UserSafeDTO{
		ID:       user.ID().Id,
		Name:     user.FullName().Name,
		Lastname: user.FullName().Lastname,
		Username: user.Username().Username,
		Dni:      user.Dni().Dni,
		Email:    user.Email().Email,
		Image:    user.Image().URL(),
	}

	// Devolver el usuario si fue encontrado
	return c.JSON(http.StatusOK, userDTO)
}

func (u *UserEchoHandler) Logout(c echo.Context) error {

	err := logout(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "sesion cerrada exitosamente"})
}

// handler para saber si esta autenticado
func (u *UserEchoHandler) CheckAuth(c echo.Context) error {
	_, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "usuario no autenticado"})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "usuario autenticado"})
}

func (u *UserEchoHandler) Restore(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	id, _ := sv.NewID(uint(idInput))
	err = u.RestoreUsecase.Execute(*id)

	if err != nil {
		if errors.Is(err, uerr.ErrUserNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, map[string]string{"succes": "user restored succesfully"})

}

func (u *UserEchoHandler) ChangePassword(c echo.Context) error {
	var req struct {
		NewPassword string `json:"new_password"`
	}
	id, err := GetAuthenticatedUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format: " + err.Error(),
		})
	}

	newPassword, err := v.NewPassword(req.NewPassword)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid password format: " + err.Error(),
		})
	}
	err = u.ChangePasswordUsecase.Execute(*id, *newPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to change password: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"success": "Password changed successfully",
	})
}

func (u *UserEchoHandler) GetUsersByRoom(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	roomid, _ := sv.NewID(idInput)

	users, err := u.GetUsersByRoomUsecase.Execute(*roomid)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	var usersDTO []d.UserSafeDTO
	for _, user := range users {
		userDTO := &d.UserSafeDTO{
			ID:       user.ID().Id,
			Name:     user.FullName().Name,
			Lastname: user.FullName().Lastname,
			Username: user.Username().Username,
			Dni:      user.Dni().Dni,
			Email:    user.Email().Email,
			Image:    user.Image().URL(),
		}
		usersDTO = append(usersDTO, *userDTO)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": "users fetched sucessfully",
		"users":   usersDTO,
	})
}

func (h *UserEchoHandler) Update(c echo.Context) error {
	log.Println("Update user handler called")
	id, err := GetAuthenticatedUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	// Bind del cuerpo de la solicitud
	var req d.UserSafeDTO
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Validar y crear los value objects
	fullName, err := v.NewFullName(req.Name, req.Lastname)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	username, err := v.NewUserName(req.Username)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	dni, err := v.NewDni(req.Dni)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	email, err := v.NewEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	currentUser, err := h.GetUserByIDUsecase.Execute(*id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	image, err := v.NewImage(req.Image)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	currentPassword := currentUser.Password()

	// Crear el objeto User con los datos actualizados
	user := d.NewUser(
		id,
		*fullName,
		*username,
		*dni,
		*email,
		currentPassword, // No actualizamos la contraseña
		image,
	)

	// Llamar al caso de uso para actualizar el usuario
	updatedUser, err := h.UpdateUsecase.Execute(user)
	if err != nil {
		if err.Error() == "unauthorized" {
			return c.JSON(http.StatusMethodNotAllowed, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	// Crear el DTO para la respuesta
	userDTO := &d.UserSafeDTO{
		ID:       updatedUser.ID().Id,
		Name:     updatedUser.FullName().Name,
		Lastname: updatedUser.FullName().Lastname,
		Username: updatedUser.Username().Username,
		Dni:      updatedUser.Dni().Dni,
		Email:    updatedUser.Email().Email,
		Image:    updatedUser.Image().URL(),
	}

	// Devolver la respuesta
	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": "user updated successfully",
		"user":    userDTO,
	})
}
