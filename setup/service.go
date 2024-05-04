package setup

import (
	"golang-note/configuration"
	"golang-note/service"
	"golang-note/util"

	"github.com/go-playground/validator/v10"
)

type ServiceSetup struct {
	UserService       service.UserService
	ContextService    service.ContextService
	PdfService        service.PdfService
	StreamJsonService service.StreamJsonService
}

func Service(mysql util.MysqlUtil, redis util.RedisUtil, validate *validator.Validate, config *configuration.Configuration, repository RepositorySetup) ServiceSetup {
	return ServiceSetup{
		UserService:       service.NewUserService(mysql, redis, validate, config.JwtKey, config.JwtAccessTokenExpireTime, config.JwtRefreshTokenExpireTime, repository.UserRepository, repository.PermissionRepository, repository.UserPermissionRepository, repository.RedisRepository),
		ContextService:    service.NewContextService(mysql.GetDb(), repository.ContextRepository),
		PdfService:        service.NewPdfService(),
		StreamJsonService: service.NewStreamJsonService(),
	}
}
