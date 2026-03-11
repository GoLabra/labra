package handler

import (
	"testing"

	"github.com/GoLabra/labra/entgql/ent"
	"github.com/GoLabra/labra/entgql/ent/enttest"
	_ "github.com/mattn/go-sqlite3"
)

func newTestAdminEntClient(t *testing.T) *ent.Client {
	t.Helper()
	return enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
}
