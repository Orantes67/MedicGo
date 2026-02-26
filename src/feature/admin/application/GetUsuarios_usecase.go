package application

import (
	"ApiMedicGO/src/feature/admin/domain/entities"
	"ApiMedicGO/src/feature/admin/domain/repositories"
)

// GetUsuariosUseCase lista los usuarios del sistema filtrados por rol.
type GetUsuariosUseCase struct {
	repo repositories.AdminRepository
}

func NewGetUsuariosUseCase(repo repositories.AdminRepository) *GetUsuariosUseCase {
	return &GetUsuariosUseCase{repo: repo}
}

// GetUsuariosResponse agrupa enfermeros y doctores para el tab Usuarios.
type GetUsuariosResponse struct {
	Enfermeros []*entities.UsuarioResumen `json:"enfermeros"`
	Doctores   []*entities.UsuarioResumen `json:"doctores"`
}

func (uc *GetUsuariosUseCase) Execute() (*GetUsuariosResponse, error) {
	enfermeros, err := uc.repo.GetUsuariosByRol("enfermero")
	if err != nil {
		return nil, err
	}
	doctores, err := uc.repo.GetUsuariosByRol("doctor")
	if err != nil {
		return nil, err
	}
	return &GetUsuariosResponse{
		Enfermeros: enfermeros,
		Doctores:   doctores,
	}, nil
}
