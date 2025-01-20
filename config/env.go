package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/go-sql-driver/mysql"
	"github.com/lpernett/godotenv"
)

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	JWTSecret              string
	JWTExpirationInSeconds int64
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	var defaultJWTExpiration int64 = 3600 * 24 * 7
	return Config{
		PublicHost:             getEnv("PUBLIC_HOST", "http://localhost"),
		Port:                   getEnv("PORT", "8080"),
		DBUser:                 getEnv("DB_USER", "root"),
		DBPassword:             getEnv("DB_PASSWORD", "Pa$$w0rd"),
		DBAddress:              fmt.Sprintf("%s:%s", getEnv("DB_HOST", "127.0.0.1"), getEnv("DB_PORT", "5432")),
		DBName:                 getEnv("DB_NAME", "ecom"),
		JWTSecret:              getEnv("JWT_SECRET", "67hytrrenjkoplhgrestnh90))898"),
		JWTExpirationInSeconds: getEnvAsInt("JWT_EXPIRATION_IN_SECONDS", defaultJWTExpiration),
	}
}

func GetDBConfig() mysql.Config {
	return mysql.Config{
		User: Envs.DBUser,
		Passwd: Envs.DBPassword,
		Addr: Envs.DBAddress,
		DBName: Envs.DBName,
		Net: "tcp",
		AllowNativePasswords: true,
		ParseTime: true,
	}
}

func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}

	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return i
		}
	}

	return fallback
}
