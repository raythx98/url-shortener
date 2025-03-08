package config

import "github.com/kelseyhightower/envconfig"

type Specification struct {
	Debug             bool   `default:"true"`
	ServerPort        int    `default:"5051"`
	DbUsername        string `required:"true"`
	DbPassword        string `required:"true"`
	DbHost            string `required:"true"`
	DbPort            int    `required:"true"`
	DbDefaultName     string `required:"true"`
	JwtSecret         string `required:"true"`
	BasicAuthUsername string `required:"true"`
	BasicAuthPassword string `required:"true"`
	SupabaseKey       string `required:"true"`
}

func Load() *Specification {
	var s Specification
	envconfig.MustProcess("APP_URLSHORTENER", &s)
	return &s
}

func (s *Specification) GetHmacSecret() []byte {
	return []byte(s.JwtSecret)
}

func (s *Specification) GetJwtSecret() string {
	return s.JwtSecret
}

func (s *Specification) GetBasicAuthUsername() []byte {
	return []byte(s.BasicAuthUsername)
}

func (s *Specification) GetBasicAuthPassword() []byte {
	return []byte(s.BasicAuthPassword)
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
