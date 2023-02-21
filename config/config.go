package config

var (
	AppName   = "Kedai"
	SecretKey = GetEnv("SECRET_KEY", "secret_key")
	DB        = DBConfig{
		Host:     GetEnv("DB_HOST", "localhost"),
		Port:     GetEnv("DB_PORT", "5432"),
		Username: GetEnv("DB_USER", ""),
		Password: GetEnv("DB_PASS", ""),
		DbName:   GetEnv("DB_NAME", ""),
		SslMode:  GetEnv("DB_SSLMODE", "disable"),
	}
	Cache = CacheConfig{
		Host:     GetEnv("REDIS_HOST", "localhost"),
		Port:     GetEnv("REDIS_PORT", "6379"),
		Username: GetEnv("REDIS_USER", ""),
		Password: GetEnv("REDIS_PASS", ""),
	}
	Origin = GetArrayENV("ORIGIN", []string{"http://localhost:3000"})
)
