package application

import (
	"errors"

	"ApiMedicGO/src/feature/register/domain/entities"
	"ApiMedicGO/src/feature/register/domain/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUseCase maneja la lógica de negocio del registro.
type RegisterUseCase struct {
	repo repositories.UserRepository
}

func NewRegisterUseCase(repo repositories.UserRepository) *RegisterUseCase {
	return &RegisterUseCase{repo: repo}
}

// RegisterRequest es el DTO de entrada.
type RegisterRequest struct {
	Name          string `json:"name"           binding:"required"`
	LicenseNumber string `json:"license_number" binding:"required"`
	Specialty     string `json:"specialty"      binding:"required,oneof=urgencias hospitalizacion uci cirugia pediatria"`
	Email         string `json:"email"          binding:"required,email"`
	Password      string `json:"password"       binding:"required,min=8"`
	Role          string `json:"role"           binding:"required,oneof=enfermero doctor administrador"`
}

// RegisterResponse es el DTO de salida.
type RegisterResponse struct {
	Message string        `json:"message"`
	User    *entities.User `json:"user"`
}

// Execute registra un nuevo usuario si el número de colegiado no existe.
func (uc *RegisterUseCase) Execute(req RegisterRequest) (*RegisterResponse, error) {
	existing, _ := uc.repo.FindByLicenseNumber(req.LicenseNumber)
	if existing != nil {
		return nil, errors.New("el número de colegiado ya está registrado")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error al procesar la contraseña")
	}

	user := &entities.User{
		ID:            primitive.NewObjectID(),
		Name:          req.Name,
		LicenseNumber: req.LicenseNumber,
		Email:         req.Email,
		Password:      string(hashedPassword),
		Role:          req.Role,
		Specialty:     req.Specialty,
	}

	if err := uc.repo.Save(user); err != nil {
		return nil, errors.New("error al guardar el usuario")
	}

	user.Password = ""
	return &RegisterResponse{
		Message: "Usuario registrado exitosamente",
		User:    user,
	}, nil
}
