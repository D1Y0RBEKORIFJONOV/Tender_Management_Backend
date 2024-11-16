package app

import "context"

type Auth interface {
	CreateUser(ctx context.Context)
}
