package entities

const (
	// Roles del sistema
	RoleNurse   = "enfermero"
	RoleDoctor  = "doctor"
	RoleAdmin   = "administrador"

	// Áreas / Especialidades del hospital
	SpecialtyUrgencias       = "urgencias"        // Atención inmediata y triaje
	SpecialtyHospitalizacion = "hospitalizacion"  // Habitaciones para pacientes internados
	SpecialtyUCI             = "uci"              // Unidad de Cuidados Intensivos — pacientes críticos
	SpecialtyCirugia         = "cirugia"          // Cirugía general y especialidades quirúrgicas
	SpecialtyPediatria       = "pediatria"        // Atención médica infantil
)
