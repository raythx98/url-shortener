package configs

import "github.com/kelseyhightower/envconfig"

type Specification struct {
	Debug         bool   `default:"true"`
	ServerPort    int    `default:"5051"`
	DbUsername    string `required:"true"`
	DbPassword    string `required:"true"`
	DbHost        string `required:"true"`
	DbPort        int    `required:"true"`
	DbDefaultName string `required:"true"`
}

func Load() *Specification {
	var s Specification
	envconfig.MustProcess("APP_URLSHORTENER", &s)
	return &s
}

func (s *Specification) GetDbUsername() string {
	return s.DbUsername
}

func (s *Specification) GetDbPassword() string {
	return s.DbPassword
}

func (s *Specification) GetDbHost() string {
	return s.DbHost
}

func (s *Specification) GetDbPort() int {
	return s.DbPort
}

func (s *Specification) GetDbDefaultName() string {
	return s.DbDefaultName
}
