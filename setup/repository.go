package setup

import "golang-note/repository"

type RepositorySetup struct {
	UserRepository           repository.UserRepository
	PermissionRepository     repository.PermissionRepository
	UserPermissionRepository repository.UserPermissionRepository
	RedisRepository          repository.RedisRepository
	ContextRepository        repository.ContextRepository
}

func Repository() RepositorySetup {
	return RepositorySetup{
		UserRepository:           repository.NewUserRepository(),
		PermissionRepository:     repository.NewPermissionRepository(),
		UserPermissionRepository: repository.NewUserPermissionRepository(),
		RedisRepository:          repository.NewRedisRepository(),
		ContextRepository:        repository.NewContextRepository(),
	}
}
