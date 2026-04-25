package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	Data    DataConfig    `yaml:"data"`
	Service ServiceConfig `yaml:"service"`
	Secrets SecretsConfig
}

// Data config
type DataConfig struct {
	User         UserConfig         `yaml:"user"`
	UserAuthData UserAuthDataConfig `yaml:"user_auth_data"`
	Token        TokenConfig        `yaml:"token"`
}

// Data/User config
type UserConfig struct {
	Cache    UserCacheConfig `yaml:"cache"`
	Username UsernameConfig  `yaml:"username"`
}

// Data/User/Cache config
type UserCacheConfig struct {
	Expires int `yaml:"expires" env:"DATA_USER_CACHE_EXPIRES" env-default:"30"`
}

// Data/User/Username config
type UsernameConfig struct {
	MinLength int `yaml:"min_length" env:"DATA_USER_NAME_MIN_LENGTH" env-default:"3"`
	MaxLength int `yaml:"max_length" env:"DATA_USER_NAME_MAX_LENGTH" env-default:"20"`
}

// Data/UserAuthData config
type UserAuthDataConfig struct {
	Cache UserAuthDataCacheConfig `yaml:"cache"`
}

// Data/UserAuthData/Cache config
type UserAuthDataCacheConfig struct {
	Expires int `yaml:"expires" env:"DATA_USER_AUTH_DATA_CACHE_EXPIRES"`
}

// Data/Token config
type TokenConfig struct {
	Refresh RefreshTokenConfig `yaml:"refresh"`
	Access  AccessTokenConfig  `yaml:"access"`
}

// Data/Token/Refresh config
type RefreshTokenConfig struct {
	MinTTL int `yaml:"min_ttl" env:"DATA_TOKEN_REFRESH_MIN_TTL" env-default:"5"`
	TTL    int `yaml:"ttl" env:"DATA_TOKEN_REFRESH_TTL" env-default:"60000"`
}

// Data/Token/Access config
type AccessTokenConfig struct {
	TTL int `yaml:"ttl" env:"DATA_TOKEN_ACCESS_TTL" env-default:"15"`
}

// Service config
type ServiceConfig struct {
	Auth AuthServiceConfig `yaml:"auth"`
}

// Service/Auth config
type AuthServiceConfig struct {
	MinPasswordLength int `yaml:"min_password_length" env:"SERVICE_AUTH_MIN_PASSWORD_LENGTH" env-default:"8"`
}
type SecretsConfig struct {
	JWTGenerationKey string `yaml:"-" env:"ACCESS_TOKEN_JWT_KEY" env-default:"mockkey"`
}

func MustLoad() *Config {
	var cfg Config

	err := cleanenv.ReadConfig("config/config.yaml", &cfg)

	if err != nil {
		panic(err)
	}

	return &cfg
}
