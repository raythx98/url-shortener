package resources

import (
	"context"
	"time"

	"github.com/raythx98/url-shortener/controller"
	"github.com/raythx98/url-shortener/repositories"
	"github.com/raythx98/url-shortener/service"
	"github.com/raythx98/url-shortener/tools/config"
	"github.com/raythx98/url-shortener/tools/qrcode"
	"github.com/raythx98/url-shortener/tools/reqctx"
	"github.com/raythx98/url-shortener/tools/zerologger"

	"github.com/raythx98/gohelpme/tool/aws"
	"github.com/raythx98/gohelpme/tool/basicauth"
	"github.com/raythx98/gohelpme/tool/crypto"
	"github.com/raythx98/gohelpme/tool/jwthelper"
	"github.com/raythx98/gohelpme/tool/postgres"
	"github.com/raythx98/gohelpme/tool/random"
	"github.com/raythx98/gohelpme/tool/validator"
)

func RegisterRepos(_ context.Context, tools Tools) Repositories {
	return Repositories{Repo: repositories.NewRepository(tools.Db)}
}

func RegisterControllers(_ context.Context, urlShortenerSvc Services, tools Tools) Controllers {
	return Controllers{
		Auth:      controller.NewAuth(urlShortenerSvc.auth, tools.Validator, tools.Log),
		Redirects: controller.NewRedirects(urlShortenerSvc.redirects, tools.Validator, tools.Log),
		Urls:      controller.NewUrls(urlShortenerSvc.urls, tools.Validator, tools.Log),
		Users:     controller.NewUsers(urlShortenerSvc.users, tools.Validator, tools.Log),
	}
}

func RegisterClients(ctx context.Context, config *config.Specification) (Clients, error) {
	awsTool, err := aws.NewConfig(ctx, config)
	if err != nil {
		return Clients{}, err
	}

	return Clients{
		S3: aws.NewS3(awsTool),
	}, nil
}

func RegisterServices(_ context.Context, config *config.Specification, repo Repositories, clients Clients,
	tools Tools) Services {
	return Services{
		auth:      service.NewAuth(repo.Repo, tools.Log, tools.Jwt, tools.Crypto, tools.ReqCtx),
		redirects: service.NewRedirects(repo.Repo, tools.Log),
		urls:      service.NewUrls(config, repo.Repo, clients.S3, tools.Log, tools.Random, tools.QrCode, tools.ReqCtx),
		users:     service.NewUsers(repo.Repo, tools.Log, tools.Crypto, tools.ReqCtx),
	}
}

func CreateTools(ctx context.Context, config *config.Specification) Tools {
	validate := validator.New()
	log := zerologger.New(config.Debug)
	db := postgres.New(ctx, config, log)
	jwtHelper := jwthelper.New(jwthelper.Config{
		Issuer:               "raythx98@gmail.com",
		Audiences:            []string{"raythx98@gmail.com"},
		AccessTokenValidity:  1 * time.Hour,
		RefreshTokenValidity: 24 * time.Hour,
	}, config)
	basicAuth := basicauth.New(config)
	cryptoTool := crypto.New(crypto.DefaultConfig())
	randomTool := random.New()
	qrCode := qrcode.New()
	reqCtx := reqctx.New()

	return Tools{
		Validator: validate,
		Log:       log,
		Db:        db,
		Jwt:       jwtHelper,
		BasicAuth: basicAuth,
		Crypto:    cryptoTool,
		Random:    randomTool,
		QrCode:    qrCode,
		ReqCtx:    reqCtx,
	}
}
