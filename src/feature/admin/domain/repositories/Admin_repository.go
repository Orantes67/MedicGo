package repositories

import "ApiMedicGO/src/feature/admin/domain/entities"

// AdminRepository define el contrato de acceso a datos para el feature admin.
type AdminRepository interface {
	// GetMetrics devuelve las métricas generales del hospital (tab Métricas).
	GetMetrics() (*entities.AdminMetrics, error)

	// GetUsuariosByRol devuelve la lista de usuarios filtrados por rol ("enfermero" o "doctor").
	GetUsuariosByRol(rol string) ([]*entities.UsuarioResumen, error)

	// FindByLicenseNumber busca un usuario por número de colegiado (para validar duplicados).
	FindByLicenseNumber(licencia string) (*entities.UsuarioResumen, error)

	// SaveUsuario persiste un nuevo usuario (enfermero o doctor) en la colección users.
	SaveUsuario(nombre, email, licencia, hashedPassword, rol, especialidad string) (*entities.UsuarioResumen, error)

	// AsignarEnfermeroADoctor vincula un enfermero a un doctor (actualiza campo doctor_id).
	AsignarEnfermeroADoctor(enfermeroID, doctorID string) error

	// GetDistribucionAreas devuelve la cantidad de pacientes por área hospitalaria.
	GetDistribucionAreas() ([]*entities.AreaDistribucion, error)
}
