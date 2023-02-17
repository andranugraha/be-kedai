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
	}
	Cache = CacheConfig{
		Host:     GetEnv("REDIS_HOST", "localhost"),
		Port:     GetEnv("REDIS_PORT", "6379"),
		Password: GetEnv("REDIS_PASS", ""),
	}
)
