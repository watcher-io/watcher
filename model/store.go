package model

import "context"

type ObjectStore interface {
	Upload(context.Context, string, string, string) error
	Fetch(context.Context, string, string) (string, error)
}
