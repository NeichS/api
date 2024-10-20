package repositories

import (
	"fmt"
	"suffgo/database"
	"suffgo/internal/user/entities"

	"github.com/labstack/gommon/log"
)

type userPostgresRepository struct {
	db database.Database
}

func NewUserPostgresRepository(db database.Database) UserRepository {
	return &userPostgresRepository{db: db}
}

func (r *userPostgresRepository) InsertUserData(in *entities.UserDto) error {
	data := &entities.User{
		Dni:      in.Dni,
		Mail:     in.Mail,
		Password: in.Password,
		Username: in.Username,
	}

	result := r.db.GetDb().Create(data)

	if result.Error != nil {
		log.Errorf("InsertUserData: %v", result.Error)
		return result.Error
	}

	log.Debugf("InsertUserData: %v", result.RowsAffected)
	return nil
}

// esta mal esto deberia devolver el modelo original y en Usecase lo pasas a DTO
func (r *userPostgresRepository) GetUserByID(id int) (*entities.UserSafeDto, error) {
	var user entities.User

	result := r.db.GetDb().Preload("CreatedRooms").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}

	userData := &entities.UserSafeDto{
		ID:           user.ID,
		Dni:          user.Dni,
		Mail:         user.Mail,
		Username:     user.Username,
		CreatedRooms: user.CreatedRooms,
	}

	return userData, nil
}

func (r *userPostgresRepository) DeleteUser(id int) error {
	result := r.db.GetDb().Delete(&entities.User{}, id)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

func (r *userPostgresRepository) FetchAll() ([]entities.User, error) {
	var users []entities.User

	result := r.db.GetDb().Preload("CreatedRooms").Find(&users)

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("User table is null")
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
