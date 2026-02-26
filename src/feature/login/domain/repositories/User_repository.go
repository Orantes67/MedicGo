package repositories

import "ApiMedicGO/src/feature/login/domain/entities"

type UserRepository interface {
	FindByLicenseNumber(licenseNumber string) (*entities.User, error)
	CreateUser(user *entities.User) error
}
