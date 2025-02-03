package configs

type IConfig interface {
	GetDbUsername() string
	GetDbPassword() string
	GetDbHost() string
	GetDbPort() int
	GetDbDefaultName() string
}
