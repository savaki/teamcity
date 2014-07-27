package teamcity

import "os"

const (
	TEAMCITY_URL      = "TEAMCITY_URL"
	TEAMCITY_USERNAME = "TEAMCITY_USERNAME"
	TEAMCITY_PASSWORD = "TEAMCITY_PASSWORD"
)

type Auth struct {
	Username string
	Password string
}

func NewFromEnv() *Auth {
	return New(os.Getenv(TEAMCITY_USERNAME), os.Getenv(TEAMCITY_PASSWORD))
}

func New(username, password string) *Auth {
	return &Auth{
		Username: username,
		Password: password,
	}
}
