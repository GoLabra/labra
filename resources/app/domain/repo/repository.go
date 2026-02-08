package repo

import (
	"app/ent"
	"app/interfaces/repo"
	"context"
	"reflect"
	"regexp"
	"strings"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/dialect/sql"
	adminRepo "github.com/GoLabra/labra/entgql/domain/repo"
	adminEnt "github.com/GoLabra/labra/entgql/ent"
	adminInterfaces "github.com/GoLabra/labra/entgql/interfaces/repo"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
)

var (
	ErrRepositoryNotSetInContext = "repository is not set in context"
)

type Txer interface {
	Create(ctx context.Context) (*ent.Tx, error)
	Commit(ctx context.Context, tx *ent.Tx) error
}

type Repository struct {
	// Node                  repo.Node
	Tx   Txer
	AdminUser              adminInterfaces.AdminUser
	Cycle              repo.Cycle
	File              adminInterfaces.File
	ForPermission              repo.ForPermission
	LifeCycleNot              repo.LifeCycleNot
	Miau              repo.Miau
	Role              adminInterfaces.Role
	User              repo.User
}

func New(client *ent.Client, adminClient *adminEnt.Client) *Repository {
	return &Repository{
		// Node:                  NewNode(client),
		Tx:   NewTx(client),
        AdminUser:              adminRepo.NewAdminUser(adminClient),
        Cycle:              NewCycle(client),
        File:              adminRepo.NewFile(adminClient),
        ForPermission:              NewForPermission(client),
        LifeCycleNot:              NewLifeCycleNot(client),
        Miau:              NewMiau(client),
        Role:              adminRepo.NewRole(adminClient),
        User:              NewUser(client),
	}
}

// TODO @David find better location for these (in api?)
func OrderFunc(o ent.OrderDirection, field string) func(*sql.Selector) {
	field = strcase.ToSnake(field)

	if o == entgql.OrderDirectionDesc {
		return ent.Desc(field)
	}

	return ent.Asc(field)
}

func CompareUniqueInput(a interface{}, b interface{}) bool {
	err := mapstructure.Decode(a, &b)
	if err != nil {
		return false
	}
	if reflect.ValueOf(b).IsZero() {
		return false
	}
	return true
}

func SnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([AZ][az]+)")
	var matchAllCap = regexp.MustCompile("([az09])([AZ])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
