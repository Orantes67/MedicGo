package entities

const (
	// Roles del sistema
	RoleNurse         = "enfermero"
	RoleDoctor        = "doctor"
	RoleJefeDoctor    = "jefe_doctor"
	RoleJefeEnfermera = "jefe_enfermera"

	// Áreas / Especialidades del hospital
	SpecialtyUrgencias       = "urgencias"       // Atención inmediata y triaje
	SpecialtyHospitalizacion = "hospitalizacion"  // Habitaciones para pacientes internados
	SpecialtyUCI             = "uci"              // Unidad de Cuidados Intensivos — pacientes críticos
)
