package resources

import (
	"github.com/raythx98/go-zap/controller"
	"github.com/raythx98/go-zap/repositories"
	"github.com/raythx98/go-zap/service"
	"github.com/raythx98/go-zap/tools/crypto"
	"github.com/raythx98/go-zap/tools/postgres"
	"github.com/raythx98/go-zap/tools/random"

	"github.com/raythx98/gohelpme/tool/basicauth"
	"github.com/raythx98/gohelpme/tool/jwthelper"
	"github.com/raythx98/gohelpme/tool/logger"
	"github.com/raythx98/gohelpme/tool/validator"
)

type Tools struct {
	Validator validator.IValidator
	Log       logger.ILogger
	Db        postgres.IPostgres
	Jwt       jwthelper.IJwt
	BasicAuth basicauth.IAuth
	Crypto    crypto.ICrypto
	Random    random.IRandom
}

type Repositories struct {
	Repo repositories.IRepository
}

type Clients struct {
}

type Services struct {
	auth      service.IAuth
	redirects service.IRedirects
	urls      service.IUrls
}

type Controllers struct {
	Auth      controller.IAuth
	Redirects controller.IRedirects
	Urls      controller.IUrls
}
