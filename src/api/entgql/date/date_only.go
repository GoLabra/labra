package date

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

const DateOnlyFormat = time.DateOnly

func MarshalDateOnly(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		timeValue := time.Time(t)
		val := fmt.Sprintf("\"%s\"", timeValue.Format(DateOnlyFormat))
		w.Write([]byte(val))
	})
}

func UnmarshalDateOnly(v interface{}) (time.Time, error) {

	stringValue, ok := v.(string)

	if !ok {
		return time.Time{}, fmt.Errorf("Date value must be a string with \"%s\" format", time.DateOnly)
	}

	timeValue, err := time.Parse(time.DateOnly, stringValue)

	if err != nil {
		return time.Time{}, err
	}

	return timeValue, nil
}
