package web

// type serverSettings struct {
// 	delegate                func(*serverSettings)
// 	GracefulShutdownTimeout time.Duration
// 	Addr                    string
// 	Mode                    string
// }

// func (settings *serverSettings) apply(target *serverSettings) {
// 	settings.delegate(target)
// }

// ServerOption 服务器配置项
// type ServerOption interface {
// 	apply(settings *serverSettings)
// }

// func WithGracefulShutdownTimeout(d time.Duration) ServerOption {
// 	return &serverSettings{delegate: func(settings *serverSettings) {
// 		settings.GracefulShutdownTimeout = d
// 	}}
// }
// func WithAddr(addr string) ServerOption {
// 	return &serverSettings{delegate: func(settings *serverSettings) {
// 		settings.Addr = addr
// 	}}
// }
