package teamcity

import "os"

type Auth struct {
	Username string
	Password string
}

func NewFromEnv() *Auth {
	return New(os.Getenv("TEAMCITY_USERNAME"), os.Getenv("TEAMCITY_PASSWORD"))
}

func New(username, password string) *Auth {
	return &Auth{
		Username: username,
		Password: password,
	}
}
