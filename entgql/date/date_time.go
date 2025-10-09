package date

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

const DateTimeFormat = DateOnlyFormat + "T" + TimeOnlyFormat

func MarshalDateTime(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		timeValue := time.Time(t)
		val := fmt.Sprintf("\"%s\"", timeValue.Format(DateTimeFormat))
		w.Write([]byte(val))
	})
}

func UnmarshalDateTime(v interface{}) (time.Time, error) {

	stringValue, ok := v.(string)

	if !ok {
		return time.Time{}, fmt.Errorf("Date value must be a string with \"%s\" format", DateTimeFormat)
	}

	timeValue, err := time.Parse(DateTimeFormat, stringValue)

	if err != nil {
		return time.Time{}, err
	}

	return timeValue, nil
}
