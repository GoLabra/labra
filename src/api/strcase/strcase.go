package strcase

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/iancoleman/strcase"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnake(str string) string {
	return strcase.ToSnake(str)
}

func NodeNameToGraphqlName(s string) string {
	return strings.ToLower(string(s[0])) + s[1:]
}

func ToLowerCamel(s string) string {
	words := strings.FieldsFunc(s, isSeparator)
	if len(words) == 1 {
		return strings.ToLower(words[0])
	}
	return strings.ToLower(words[0]) + pascalWords(words[1:])
}

func ToCamel(s string) string {
	words := strings.FieldsFunc(s, isSeparator)
	if len(words) == 1 {
		return words[0]
	}
	return words[0] + pascalWords(words[1:])
}

func ToPascal(s string) string {
	words := strings.FieldsFunc(s, isSeparator)
	if len(words) == 1 {
		return strings.ToTitle(words[0][:1]) + words[0][1:]
	}
	return pascalWords(words)
}

func ToTitle(val string) string {
	return strings.ToTitle(val[:1]) + val[1:]
}

func LowerFirstLetter(val string) string {
	return strings.ToLower(val[:1]) + val[1:]
}

func pascalWords(words []string) string {
	for i, w := range words {
		upper := strings.ToUpper(w)
		if _, ok := acronyms[upper]; ok {
			words[i] = upper
		} else {
			words[i] = strings.ToTitle(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, "")
}

func isSeparator(r rune) bool {
	return r == '_' || r == '-' || unicode.IsSpace(r)
}

var acronyms = map[string]struct{}{
	"ACL": struct{}{}, "API": struct{}{}, "ASCII": struct{}{}, "AWS": struct{}{}, "CPU": struct{}{}, "CSS": struct{}{}, "DNS": struct{}{}, "EOF": struct{}{}, "GB": struct{}{}, "GUID": struct{}{},
	"HCL": struct{}{}, "HTML": struct{}{}, "HTTP": struct{}{}, "HTTPS": struct{}{}, "ID": struct{}{}, "IP": struct{}{}, "JSON": struct{}{}, "KB": struct{}{}, "LHS": struct{}{}, "MAC": struct{}{},
	"MB": struct{}{}, "QPS": struct{}{}, "RAM": struct{}{}, "RHS": struct{}{}, "RPC": struct{}{}, "SLA": struct{}{}, "SMTP": struct{}{}, "SQL": struct{}{}, "SSH": struct{}{}, "SSO": struct{}{},
	"TCP": struct{}{}, "TLS": struct{}{}, "TTL": struct{}{}, "UDP": struct{}{}, "UI": struct{}{}, "UID": struct{}{}, "URI": struct{}{}, "URL": struct{}{}, "UTF8": struct{}{}, "UUID": struct{}{},
	"VM": struct{}{}, "XML": struct{}{}, "XMPP": struct{}{}, "XSRF": struct{}{}, "XSS": struct{}{},
}
