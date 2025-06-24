package utils

import "github.com/samborkent/uuidv7"

func NewUUIDV7() string {
	return uuidv7.New().String()
}
