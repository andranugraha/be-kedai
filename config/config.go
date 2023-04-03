package config

var (
	AppName       = "Kedai"
	SecretKey     = GetEnv("SECRET_KEY", "secret_key")
	HashKey       = GetEnv("HASH_KEY", "secret_key")
	MerchantCode  = GetEnv("MERCHANT_CODE", "code")
	PlatformFee   = GetEnv("PLATFORM_FEE", "0")
	AdminEmail    = GetEnv("ADMIN_EMAIL", "")
	AdminPassword = GetEnv("ADMIN_PASSWORD", "")
	DB            = DBConfig{
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
	Mailgun = MailgunConfig{
		PRIVATE_API_KEY: GetEnv("MAILGUN_PRIVATE_API_KEY", ""),
		DOMAIN:          GetEnv("MAILGUN_DOMAIN", ""),
		API_BASE_URL:    GetEnv("MAILGUN_API_BASE_URL", ""),
		SENDER:          GetEnv("MAILGUN_SENDER", ""),
	}
	Origin                = GetArrayENV("ORIGIN", []string{"http://localhost:3000"})
	DefaultProfilePicture = GetEnv("DEFAULT_PROFILE_PICTURE", "")
)
