package resources

import (
	"context"
	"github.com/raythx98/url-shortener/tools/supabase"
	"time"

	"github.com/raythx98/url-shortener/controller"
	"github.com/raythx98/url-shortener/service"
	"github.com/raythx98/url-shortener/sqlc/db"
	"github.com/raythx98/url-shortener/tools/config"
	"github.com/raythx98/url-shortener/tools/zerologger"

	"github.com/raythx98/gohelpme/tool/basicauth"
	"github.com/raythx98/gohelpme/tool/crypto"
	"github.com/raythx98/gohelpme/tool/jwthelper"
	"github.com/raythx98/gohelpme/tool/postgres"
	"github.com/raythx98/gohelpme/tool/random"
	"github.com/raythx98/gohelpme/tool/validator"
)

func RegisterRepos(_ context.Context, tools Tools) *db.Queries {
	return db.New(tools.DbPool)
}

func RegisterControllers(_ context.Context, urlShortenerSvc Services, tools Tools) Controllers {
	return Controllers{
		Auth:      controller.NewAuth(urlShortenerSvc.auth, tools.Validator, tools.Log),
		Redirects: controller.NewRedirects(urlShortenerSvc.redirects, tools.Validator, tools.Log),
		Urls:      controller.NewUrls(urlShortenerSvc.urls, tools.Validator, tools.Log),
		Users:     controller.NewUsers(urlShortenerSvc.users, tools.Validator, tools.Log),
	}
}

func RegisterServices(_ context.Context, urlMappingRepo *db.Queries, tools Tools) Services {
	return Services{
		auth:      service.NewAuth(urlMappingRepo, tools.Log, tools.Jwt, tools.Crypto),
		redirects: service.NewRedirects(urlMappingRepo, tools.Log),
		urls:      service.NewUrls(urlMappingRepo, tools.Log, tools.Random),
		users:     service.NewUsers(urlMappingRepo, tools.Log, tools.Crypto),
	}
}

func CreateTools(ctx context.Context, config *config.Specification) Tools {

	// TODO: Simplified implementation, change to AWS S3 later.
	supabase.New(config.SupabaseKey)

	validate := validator.New()
	log := zerologger.New(config.Debug)
	dbPool := postgres.NewPool(ctx, config, log)
	jwtHelper := jwthelper.New(jwthelper.Config{
		Issuer:               "raythx98@gmail.com",
		Audiences:            []string{"raythx98@gmail.com"},
		AccessTokenValidity:  1 * time.Hour,
		RefreshTokenValidity: 24 * time.Hour,
	}, config)
	basicAuth := basicauth.New(config)
	cryptoTool := crypto.New(crypto.DefaultConfig())
	randomTool := random.New()

	return Tools{
		Validator: validate,
		Log:       log,
		DbPool:    dbPool,
		Jwt:       jwtHelper,
		BasicAuth: basicAuth,
		Crypto:    cryptoTool,
		Random:    randomTool,
	}
}
