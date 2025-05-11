package date

import (
	"fmt"
	"io"
	"time"

	"github.com/99designs/gqlgen/graphql"
)

const TimeOnlyFormat = "15:04:05.000"

func MarshalTimeOnly(t time.Time) graphql.Marshaler {
	return graphql.WriterFunc(func(w io.Writer) {
		timeValue := time.Time(t)
		val := fmt.Sprintf("\"%s\"", timeValue.Format(TimeOnlyFormat))
		w.Write([]byte(val))
	})
}

func UnmarshalTimeOnly(v interface{}) (time.Time, error) {

	stringValue, ok := v.(string)

	if !ok {
		return time.Time{}, fmt.Errorf("Date value must be a string with \"%s\" format", TimeOnlyFormat)
	}

	timeValue, err := time.Parse(TimeOnlyFormat, stringValue)

	if err != nil {
		return time.Time{}, err
	}

	return timeValue, nil
}
