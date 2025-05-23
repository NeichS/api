package mappers

import (
	"suffgo/internal/rooms/domain"
	v "suffgo/internal/rooms/domain/valueObjects"
	m "suffgo/internal/rooms/infrastructure/models"
	sv "suffgo/internal/shared/domain/valueObjects"
)

func DomainToModel(room *domain.Room) *m.Room {

	return &m.Room{
		ID:          room.ID().Id,
		IsFormal:    room.IsFormal().IsFormal,
		Name:        room.Name().Name,
		Code:        room.Code().Code,
		Description: room.Description().Description,
		AdminID:     room.AdminID().Id,
		State:       room.State().CurrentState,
		Image:       room.Image().Image,
	}
}

func ModelToDomain(roomModel *m.Room) (*domain.Room, error) {
	id, err := sv.NewID(roomModel.ID)
	if err != nil {
		return nil, err
	}
	isFormal, err := v.NewIsFormal(roomModel.IsFormal)
	if err != nil {
		return nil, err
	}
	name, err := v.NewName(roomModel.Name)
	if err != nil {
		return nil, err
	}
	adminID, err := sv.NewID(roomModel.AdminID)
	if err != nil {
		return nil, err
	}

	description, err := v.NewDescription(roomModel.Description)
	if err != nil {
		return nil, err
	}
	code, err := v.NewInviteCode(roomModel.Code)
	if err != nil {
		return nil, err
	}
	image, err := v.NewImage(roomModel.Image)
	if err != nil {
		return nil, err
	}

	state, err := v.NewState(roomModel.State)
	if err != nil {
		return nil, err
	}

	return domain.NewRoom(id, *isFormal, code,*name, adminID, *description, image, state), nil
}

