// Package utils groups helper functions used across the API. This file is
// located in src/api/utils.
package utils

// P returns a pointer to the provided value.
func P[T any](val T) *T {
	return &val
}
