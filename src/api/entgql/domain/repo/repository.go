package repo

import (
	"context"
	"reflect"
	"regexp"
	"strings"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/dialect/sql"
	"github.com/GoLabra/labra/src/api/entgql/ent"
	"github.com/GoLabra/labra/src/api/entgql/interfaces/repo"
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
	Tx         Txer
	Permission repo.Permission
	Role       repo.Role
	User       repo.User
}

func New(client *ent.Client) *Repository {
	return &Repository{
		// Node:                  NewNode(client),
		Tx:         NewTx(client),
		Permission: NewPermission(client),
		Role:       NewRole(client),
		User:       NewUser(client),
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
