package id

import "github.com/oklog/ulid/v2"

func New() string {
	return ulid.Make().String()
}

func Parse(s string) (ulid.ULID, error) {
	return ulid.Parse(s)
}
