//go:build wireinject
// +build wireinject

package inject

import (
	"github.com/google/wire"

	"github.com/KScaesar/IsCoolLab2024/pkg/adapters"
	"github.com/KScaesar/IsCoolLab2024/pkg/adapters/cli"
	"github.com/KScaesar/IsCoolLab2024/pkg/adapters/database"
	"github.com/KScaesar/IsCoolLab2024/pkg/app"
)

//go:generate wire gen

func NewInfra(conf *database.GormConfing) (*adapters.Infra, error) {
	panic(wire.Build(
		wire.Struct(new(adapters.Infra), "*"),
		database.NewGrom,
	))
}

func NewAppService(infra *adapters.Infra) *app.Service {
	panic(wire.Build(
		// https://github.com/google/wire/blob/main/docs/guide.md#use-fields-of-a-struct-as-providers
		wire.FieldsOf(new(*adapters.Infra), "Database"),
		wire.Struct(new(app.Service), "*"),

		database.NewUserRepository,
		wire.Bind(new(app.UserRepository), new(*database.UserRepository)),

		database.NewFileSystemRepository,
		wire.Bind(new(app.FileSystemRepository), new(*database.FileSystemRepository)),

		app.NewUserUseCase,
		wire.Bind(new(app.UserService), new(*app.UserUseCase)),
	))
}

func NewRootCommand(infra *adapters.Infra) *cli.Command {
	panic(wire.Build(
		NewAppService,
		cli.NewRootCommand,
	))
}
