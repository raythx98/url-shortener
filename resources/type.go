package resources

import (
	"github.com/raythx98/url-shortener/controller"
	"github.com/raythx98/url-shortener/repositories"
	"github.com/raythx98/url-shortener/service"
	"github.com/raythx98/url-shortener/tools/aws"
	"github.com/raythx98/url-shortener/tools/crypto"
	"github.com/raythx98/url-shortener/tools/postgres"
	"github.com/raythx98/url-shortener/tools/qrcode"
	"github.com/raythx98/url-shortener/tools/random"

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
	QrCode    qrcode.IQrCode
}

type Repositories struct {
	Repo repositories.IRepository
}

type Clients struct {
	S3 aws.IS3
}

type Services struct {
	auth      service.IAuth
	redirects service.IRedirects
	urls      service.IUrls
	users     service.IUsers
}

type Controllers struct {
	Auth      controller.IAuth
	Redirects controller.IRedirects
	Urls      controller.IUrls
	Users     controller.IUsers
}
