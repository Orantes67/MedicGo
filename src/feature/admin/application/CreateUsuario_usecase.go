package application

import (
	"errors"

	"ApiMedicGO/src/feature/admin/domain/entities"
	"ApiMedicGO/src/feature/admin/domain/repositories"

	"golang.org/x/crypto/bcrypt"
)

// CreateUsuarioUseCase permite al administrador crear enfermeros o doctores.
type CreateUsuarioUseCase struct {
	repo repositories.AdminRepository
}

func NewCreateUsuarioUseCase(repo repositories.AdminRepository) *CreateUsuarioUseCase {
	return &CreateUsuarioUseCase{repo: repo}
}

// CreateUsuarioRequest es el DTO de entrada para crear un usuario desde el panel admin.
type CreateUsuarioRequest struct {
	Nombre        string `json:"nombre"         binding:"required"`
	LicenseNumber string `json:"license_number" binding:"required"`
	Email         string `json:"email"          binding:"required,email"`
	Password      string `json:"password"       binding:"required,min=8"`
	Rol           string `json:"rol"            binding:"required,oneof=enfermero doctor"`
	Especialidad  string `json:"especialidad"   binding:"required,oneof=urgencias hospitalizacion uci cirugia pediatria"`
}

func (uc *CreateUsuarioUseCase) Execute(req CreateUsuarioRequest) (*entities.UsuarioResumen, error) {
	// Verificar que el número de colegiado no exista ya
	existing, _ := uc.repo.FindByLicenseNumber(req.LicenseNumber)
	if existing != nil {
		return nil, errors.New("el número de colegiado ya está registrado")
	}

	// Hash de la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error al procesar la contraseña")
	}

	return uc.repo.SaveUsuario(
		req.Nombre,
		req.Email,
		req.LicenseNumber,
		string(hashedPassword),
		req.Rol,
		req.Especialidad,
	)
}
