package config

import (
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	Stage             string `default:"development"`
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
	AwsS3Bucket       string `required:"true"`
	AwsRegion         string `required:"true"`
	AwsAccessKey      string `required:"true"`
	AwsSecretKey      string `required:"true"`
}

func Load() *Specification {
	var s Specification
	envconfig.MustProcess("APP_URLSHORTENER", &s)
	return &s
}

func (s *Specification) IsDevelopment() bool {
	return strings.EqualFold(s.Stage, "development")
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

func (s *Specification) GetAwsS3Bucket() string {
	return s.AwsS3Bucket
}

func (s *Specification) GetAwsRegion() string {
	return s.AwsRegion
}

func (s *Specification) GetAwsAccessKey() string {
	return s.AwsAccessKey
}

func (s *Specification) GetAwsSecretKey() string {
	return s.AwsSecretKey
}
