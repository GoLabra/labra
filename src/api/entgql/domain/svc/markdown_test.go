package svc_test

import (
	"testing"

	"github.com/GoLabra/labra/src/api/entgql/entity"
	. "github.com/onsi/gomega"
)

func TestMarkdownFieldType(t *testing.T) {
	RegisterTestingT(t)
	Expect(string(entity.FieldTypeMarkdown)).To(Equal("Markdown"))
}
