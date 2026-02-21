package repositories

import "ApiMedicGO/src/feature/login/domain/entities"

type UserRepository interface {
	FindByEmail(email string) (*entities.User, error)
}
