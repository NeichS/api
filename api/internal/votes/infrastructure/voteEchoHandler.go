package infrastructure

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	sv "suffgo/internal/shared/domain/valueObjects"
	v "suffgo/internal/votes/application/useCases"
	d "suffgo/internal/votes/domain"

	se "suffgo/internal/shared/domain/errors"

	verrors "suffgo/internal/votes/domain/errors"

	"github.com/labstack/echo/v4"
)

type VoteEchoHandler struct {
	CreateVoteUsecase  *v.CreateUsecase
	DeleteVoteUsecase  *v.DeleteUsecase
	GetAllVoteUsecase  *v.GetAllUsecase
	GetVoteByIDUsecase *v.GetByIDUsecase
}

func NewVoteEchoHandler(
	createUC *v.CreateUsecase,
	deleteUC *v.DeleteUsecase,
	getAllUC *v.GetAllUsecase,
	getByIDUC *v.GetByIDUsecase,
) *VoteEchoHandler {
	return &VoteEchoHandler{
		CreateVoteUsecase:  createUC,
		DeleteVoteUsecase:  deleteUC,
		GetAllVoteUsecase:  getAllUC,
		GetVoteByIDUsecase: getByIDUC,
	}
}

func (h *VoteEchoHandler) CreateVote(c echo.Context) error {
	var req d.VoteCreateRequest
	// bindea el body del request (json) al dto
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userID, err := sv.NewID(req.UserID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	optionID, err := sv.NewID(req.OptionID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	vote := d.NeweVote(
		nil,
		userID,
		optionID,
	)

	//falta recepcionar al voto creado
	err = h.CreateVoteUsecase.Execute(*vote)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	//mapear dto

	response := map[string]interface{}{
		"success": "voto creado exitosamente.",
		"vote": "", 
	}

	return c.JSON(http.StatusCreated, response)
}

func (h *VoteEchoHandler) DeleteVote(c echo.Context) error {
	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	err = h.DeleteVoteUsecase.Execute(*id)

	if err != nil {
		if errors.Is(err, verrors.ErrVoteNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"success": "Vote deleted succesfully"})
}

func (h *VoteEchoHandler) GetAllVotes(c echo.Context) error {
	votes, err := h.GetAllVoteUsecase.Execute()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var votesDTO []d.VoteDTO
	for _, vote := range votes {
		voteDTO := &d.VoteDTO{
			ID:       vote.ID().Id,
			UserID:   vote.UserID().Id,
			OptionID: vote.OptionID().Id,
		}
		votesDTO = append(votesDTO, *voteDTO)
	}

	if votesDTO == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": verrors.ErrVoteNotFound.Error()})
	}

	response := map[string]interface{}{
		"success": "votos obtennidos correctamente.",
		"votes":   votesDTO,
	}

	return c.JSON(http.StatusOK, response)
}

func (h *VoteEchoHandler) GetVoteByID(c echo.Context) error {

	idParam := c.Param("id")
	idInput, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": se.ErrInvalidID.Error()})
	}

	id, _ := sv.NewID(uint(idInput))
	vote, err := h.GetVoteByIDUsecase.Execute(*id)

	if err != nil {
		if errors.Is(err, verrors.ErrVoteNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	voteDTO := &d.VoteDTO{
		ID:       vote.ID().Id,
		UserID:   vote.UserID().Id,
		OptionID: vote.OptionID().Id,
	}

	msg := fmt.Sprintf("voto con id %d obtenido exitosamente.", id.Id)
	response := map[string]interface{}{
		"success": msg,
		"votes":   voteDTO,
	}

	return c.JSON(http.StatusOK, response)
}
