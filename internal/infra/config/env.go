package config

import (
	"github.com/caarlos0/env/v11"
	_ "github.com/joho/godotenv/autoload"
)

type Env struct {
	AppName string `env:"APP_NAME,required"`
	AppPort int    `env:"APP_PORT,required"`

	JWTSecret string `env:"JWT_SECRET,required"`

	PostgresUsername string `env:"POSTGRES_USERNAME,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPort     int    `env:"POSTGRES_PORT,required"`
	PostgresDB       string `env:"POSTGRES_DB,required"`
	PostgresSSL      string `env:"POSTGRES_SSL,required"`

	StorageURL    string `env:"SUPABASE_URL,required"`
	StorageToken  string `env:"SUPABASE_KEY,required"`
	StorageBucket string `env:"SUPABASE_BUCKET,required"`

	MidtransKey string `env:"MIDTRANS_SERVER_KEY,required"`

	// AI Service Configuration
	VistaraAIURL string `env:"VISTARA_AI_URL" envDefault:"http://localhost:5000"`
}

func LoadEnv() (*Env, error) {

	_env := new(Env)

	if err := env.Parse(_env); err != nil {
		return nil, err
	}

	return _env, nil
}
