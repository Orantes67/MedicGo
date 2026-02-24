package repositories

import "ApiMedicGO/src/feature/register/domain/entities"

type UserRepository interface {
	Save(user *entities.User) error
	FindByLicenseNumber(licenseNumber string) (*entities.User, error)
}
