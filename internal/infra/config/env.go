package config

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

// Env holds all environment configuration variables
type Env struct {
	// Application settings
	AppName string `env:"APP_NAME,required"`
	AppPort int    `env:"APP_PORT,required"`

	// JWT authentication settings
	JWTSecret string `env:"JWT_SECRET,required"`

	// PostgreSQL database settings
	PostgresUsername string `env:"POSTGRES_USERNAME,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPort     int    `env:"POSTGRES_PORT,required"`
	PostgresDB       string `env:"POSTGRES_DB,required"`
	PostgresSSL      string `env:"POSTGRES_SSL,required"`

	// Supabase storage settings
	StorageURL    string `env:"SUPABASE_URL,required"`
	StorageToken  string `env:"SUPABASE_KEY,required"`
	StorageBucket string `env:"SUPABASE_BUCKET,required"`

	// Midtrans payment settings
	MidtransKey string `env:"MIDTRANS_SERVER_KEY,required"`

	// AI service integration settings
	VistaraAIURL string `env:"VISTARA_AI_URL" envDefault:"http://localhost:5000"`
	VistaraAIKey string `env:"VISTARA_AI_KEY" envDefault:"vistara-be-service-key"`
}

// LoadEnv loads and validates environment variables
func LoadEnv() (*Env, error) {
	config := new(Env)
	if err := env.Parse(config); err != nil {
		return nil, err
	}
	return config, nil
}
