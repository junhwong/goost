package web

type Controller interface {
	// OnError(ctx Context, err error) error
	Init(ctx ServerContext, r ControllerRoutes) error
}
