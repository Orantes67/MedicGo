package application

import (
	"errors"
	"os"
	"time"

	"ApiMedicGO/src/feature/login/domain/entities"
	"ApiMedicGO/src/feature/login/domain/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// LoginUseCase maneja la lógica de negocio del login.
type LoginUseCase struct {
	repo repositories.UserRepository
}

func NewLoginUseCase(repo repositories.UserRepository) *LoginUseCase {
	return &LoginUseCase{repo: repo}
}

// LoginRequest es el DTO de entrada.
type LoginRequest struct {
	LicenseNumber string `json:"license_number" binding:"required"`
	Password      string `json:"password" binding:"required"`
}

// LoginResponse es el DTO de salida.
type LoginResponse struct {
	Token string        `json:"token"`
	User  *entities.User `json:"user"`
}

// Claims define el payload del JWT.
type Claims struct {
	UserID        string `json:"user_id"`
	LicenseNumber string `json:"license_number"`
	Role          string `json:"role"`
	jwt.RegisteredClaims
}

// Execute valida las credenciales y retorna un JWT si son correctas.
func (uc *LoginUseCase) Execute(req LoginRequest) (*LoginResponse, error) {
	user, err := uc.repo.FindByLicenseNumber(req.LicenseNumber)
	if err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	token, err := generateJWT(user)
	if err != nil {
		return nil, errors.New("error al generar el token")
	}

	user.Password = ""
	return &LoginResponse{Token: token, User: user}, nil
}

func generateJWT(user *entities.User) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET no está configurado en las variables de entorno")
	}

	claims := Claims{
		UserID:        user.ID.Hex(),
		LicenseNumber: user.LicenseNumber,
		Role:          user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
