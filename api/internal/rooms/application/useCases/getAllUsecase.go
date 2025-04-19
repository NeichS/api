package usecases

import (
	"suffgo/internal/rooms/domain"
)

type (
	GetAllUsecase struct {
		getAllRepository domain.RoomRepository
	}
)

func NewGetAllUsecase(repository domain.RoomRepository) *GetAllUsecase {
	return &GetAllUsecase{
		getAllRepository: repository,
	}
}

func (s *GetAllUsecase) Execute() ([]domain.Room, error) {

	rooms, err := s.getAllRepository.GetAll()

	if err != nil {
		return nil, err
	}

	updatedRooms := make([]domain.Room, 0, len(rooms))

	updatedRooms = append(updatedRooms, rooms...)

	return updatedRooms, nil
}
