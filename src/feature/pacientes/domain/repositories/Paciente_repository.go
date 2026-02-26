package repositories

import "ApiMedicGO/src/feature/pacientes/domain/entities"

type PacienteRepository interface {
	Save(paciente *entities.Paciente) error
	FindAll() ([]*entities.Paciente, error)
	FindByID(id string) (*entities.Paciente, error)
	FindByDoctor(doctorID string) ([]*entities.Paciente, error)
	FindByEnfermero(enfermeroID string) ([]*entities.Paciente, error)
	Assign(pacienteID string, doctorID *string, enfermeroID *string, nombreEnfermero string) error
	Update(id string, paciente *entities.Paciente) error
	Delete(id string) error
}
