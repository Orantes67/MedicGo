package repositories

import "ApiMedicGO/src/feature/register/domain/entities"

type UserRepository interface {
	Save(user *entities.User) error
	FindByEmail(email string) (*entities.User, error)
}
