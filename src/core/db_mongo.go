package core

import (
	"context"
	"log"
	"os"
	"time"

	"ApiMedicGO/src/feature/login/domain/entities"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var DB *mongo.Database

func ConnectMongoDB() {
	uri := os.Getenv("MONGO_CREDENTIALS")
	if uri == "" {
		log.Fatal("MONGO_CREDENTIALS environment variable not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	DB = client.Database("ApiMedicGO")
	log.Println("✅ Connected to MongoDB successfully")

	// Inicializar usuario administrador predefinido
	initializeAdminUser()
}

// initializeAdminUser crea el usuario administrador si no existe.
func initializeAdminUser() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	usersCollection := DB.Collection("users")

	// Verificar si el usuario admin ya existe
	result := usersCollection.FindOne(ctx, bson.M{"license_number": "ADMIN001"})
	if result.Err() == nil {
		log.Println("✅ Admin user already exists")
		return
	}

	// Hash de la contraseña: "admin123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("❌ Error hashing admin password: %v", err)
		return
	}

	// Crear usuario administrador
	adminUser := entities.User{
		Name:          "Administrador",
		LicenseNumber: "ADMIN001",
		Email:         "admin@medicgo.com",
		Password:      string(hashedPassword),
		Role:          entities.RoleAdmin,
		Specialty:     "",
	}

	_, err = usersCollection.InsertOne(ctx, adminUser)
	if err != nil {
		log.Printf("❌ Error creating admin user: %v", err)
		return
	}

	log.Println("✅ Admin user created successfully")
	log.Println("   📧 Usuario: ADMIN001")
	log.Println("   🔐 Contraseña: admin123")
}
