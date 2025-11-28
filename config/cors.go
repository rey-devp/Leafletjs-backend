package config

var allowedOrigins = []string{
	"http://localhost:3000",
}

func GetAllowedOrigins() []string {
	return allowedOrigins
}