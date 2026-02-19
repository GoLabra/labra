package entity

import (
	"fmt"
	"io"
	"net/mail"
	"strings"
)

type Email string

func (Email) ImplementsGraphQLType(name string) bool { return name == "Email" }

func (e *Email) UnmarshalGraphQL(input any) error {
	s, ok := input.(string)
	if !ok {
		return fmt.Errorf("email must be a string")
	}

	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return fmt.Errorf("email is required")
	}
	if len(s) > 254 {
		return fmt.Errorf("email too long")
	}

	addr, err := mail.ParseAddress(s)
	if err != nil || addr.Address != s || !strings.Contains(s, "@") {
		return fmt.Errorf("invalid email format")
	}

	*e = Email(s)
	return nil
}

func (e Email) MarshalGraphQL(w io.Writer) {
	fmt.Fprintf(w, "%q", string(e))
}
