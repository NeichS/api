package infrastructure

import (
	"net/http"
	"strconv"
	u "suffgo/internal/users/application/useCases"

	d "suffgo/internal/users/domain"
	v "suffgo/internal/users/domain/valueObjects"

	sv "suffgo/internal/shared/domain/valueObjects"

	se "suffgo/internal/shared/domain/errors"

	"github.com/labstack/echo/v4"
)

type UserEchoHandler struct {
	CreateUserUsecase  *u.CreateUsecase
	DeleteUserUsecase  *u.DeleteUsecase
	GetAllUsersUsecase *u.GetAllUsecase
	GetUserByIDUsecase *u.GetByIDUsecase
	LoginUsecase       *u.LoginUsecase
}

// Constructor for UserEchoHandler
func NewUserEchoHandler(
	createUC *u.CreateUsecase,
	deleteUC *u.DeleteUsecase,
	getAllUC *u.GetAllUsecase,
	getByIDUC *u.GetByIDUsecase,
	loginUC *u.LoginUsecase,
) *UserEchoHandler {
	return &UserEchoHandler{
		CreateUserUsecase:  createUC,
		DeleteUserUsecase:  deleteUC,
		GetAllUsersUsecase: getAllUC,
		GetUserByIDUsecase: getByIDUC,
		LoginUsecase:       loginUC,
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

	if err := createSession(user.ID(), c); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	userDTO := &d.UserSafeDTO{
		ID:       user.ID().Id,
		Name:     user.FullName().Name,
		Lastname: user.FullName().Lastname,
		Username: user.Username().Username,
		Dni:      user.Dni().Dni,
		Email:    user.Email().Email,
	}

	response := map[string]interface{}{
		"sucess": "autenticación exitosa",
		"user":   userDTO,
	}

	// Devuelvo el id del usuario logueado
	return c.JSON(http.StatusOK, response)
}

func (h *UserEchoHandler) CreateUser(c echo.Context) error {
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
	// Map DTO to domain entity

	user := d.NewUser(
		nil,
		*fullname,
		*username,
		*dni,
		*email,
		*password,
	)

	// Call the use case
	user, err = h.CreateUserUsecase.Execute(*user)
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
	}

	response := map[string]interface{}{
		"sucess": "usuario creado exitosamente",
		"user":   userDTO,
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *UserEchoHandler) DeleteUser(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": invalidErr.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.DeleteUserUsecase.Execute(*id)
	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"sucess": "user deleted succesfully"})
}

func (h *UserEchoHandler) GetAllUsers(c echo.Context) error {
	users, err := h.GetAllUsersUsecase.Execute()
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
		}
		usersDTO = append(usersDTO, *userDTO)
	}

	return c.JSON(http.StatusOK, usersDTO)
}

func (h *UserEchoHandler) GetUserByID(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		invalidErr := &se.InvalidIDError{ID: idParam}
		return c.JSON(http.StatusBadRequest, map[string]string{"error": invalidErr.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	user, err := h.GetUserByIDUsecase.Execute(*id)
	if err != nil {
		if err.Error() == "user not found" {
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
	}
	return c.JSON(http.StatusOK, userDTO)
}

func (h *UserEchoHandler) Logout(c echo.Context) error {

	err := logout(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"sucess": "sesion cerrada exitosamente"})
}
