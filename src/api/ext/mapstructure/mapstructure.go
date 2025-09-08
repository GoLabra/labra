package mapstructure

import (
	"fmt"
	"reflect"
	"time"
	"database/sql"
	"github.com/mitchellh/mapstructure"
)

// Decode is a drop-in replacement for mapstructure.Decode that fixes the time.Time bug
// where time.Time fields are decoded to zero values instead of proper time values.
// This addresses the issue described in: https://github.com/go-viper/mapstructure/issues/20
func Decode(input, output interface{}) error {
	config := &mapstructure.DecoderConfig{
		DecodeHook: timeHook(),
		Result:     output,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return fmt.Errorf("failed to create decoder: %w", err)
	}

	return decoder.Decode(input)
}

// timeHook returns a decoder hook that properly handles time.Time conversion
func timeHook() mapstructure.DecodeHookFunc {
		timePtrType := reflect.TypeOf(&time.Time{})
		timeType := reflect.TypeOf(time.Time{})
		strMapType := reflect.TypeOf(map[string]interface{}{})
		sqlNullTimeType := reflect.TypeOf(sql.NullTime{})
		sqlNullTimePtrType := reflect.TypeOf(&sql.NullTime{})
		return func(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
			//for struct => struct conversions mapstructure currently employs the tactic to always decode it to a map
			//and then convert the map to the new struct. For the time.Time struct this fails as the time struct has no
			//public fields.
			//To work around this we json encode/decode the time
			//We also support sql.NullTime and gorm.DeletedAt
			if t == strMapType {
				if f == timePtrType {
					return encodeTime(data.(*time.Time))
				} else if f == timeType {
					d := data.(time.Time)
					return encodeTime(&d)
				} else if f == sqlNullTimeType {
					d := data.(sql.NullTime)
					if d.Valid {
						return encodeTime(&d.Time)
					} else {
						return encodeTime(nil)
					}
				} else if f == sqlNullTimePtrType {
					d := data.(*sql.NullTime)
					if d.Valid {
						return encodeTime(&d.Time)
					} else {
						return encodeTime(nil)
					}
				} 
			}


			if f == strMapType && (t == sqlNullTimeType || t == sqlNullTimePtrType || t == timeType || t == timePtrType) {
				to, err := decodeTime(data.(map[string]interface{}))
				if err != nil {
					return nil, err
				}

				if t == timeType {
					if to == nil {
						//as we cannot return nil, return zero time instead
						return time.Time{}, nil
					}
					return *to, nil
				} else if t == timePtrType {
					if to == nil {
						//we can't just return to as to is a typed nil; currently mapstructure does not handle that well
						return nil, nil
					} else {
						return to, nil
					}
				} else if t == sqlNullTimeType {
					if to == nil {
						return sql.NullTime{Time: time.Time{}, Valid: false}, nil
					} else {
						return sql.NullTime{Time: *to, Valid: true}, nil
					}
				} else if t == sqlNullTimePtrType {
					if to == nil {
						return &sql.NullTime{Time: time.Time{}, Valid: false}, nil
					} else {
						return &sql.NullTime{Time: *to, Valid: true}, nil
					}
				}
			}

			return data, nil
		}
	}


func encodeTime(data *time.Time) (map[string]interface{}, error) {
	if data == nil {
		return map[string]interface{}{"_json": []uint8("null")}, nil
	}
	json, err := data.MarshalJSON()
	return map[string]interface{}{"_json": json}, err
}

func decodeTime(mapData map[string]interface{}) (*time.Time, error) {
	to := &time.Time{}
	json, ok := mapData["_json"]
	if !ok {
		return nil, nil //TODO: error
	}
	err := to.UnmarshalJSON(json.([]byte))
	if to.IsZero() {
		return nil, nil
	}
	return to, err
}