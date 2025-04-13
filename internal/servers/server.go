package servers

import "context"

type BackendServer interface {
	Start() error
	ConfigureRoutes()
	Stop(ctx context.Context)
}
