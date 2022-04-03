package web

// func RequiredAuthorized() func(Context) {
// 	return func(ctx Context) {
// 		auth := security.AuthenticationFromContext(ctx)
// 		if !auth.IsAuthenticated() {
// 			panic(ErrUnauthorized)
// 		}
// 		ctx.Next()
// 	}
// }

// //Authroize
// func Authority(perms ...security.Permission) RouteOption {
// 	if len(perms) == 0 {
// 		panic("perms cannot be nil")
// 	}
// 	handler := func(ctx Context) {
// 		auth := security.AuthenticationFromContext(ctx)
// 		for _, p := range perms {
// 			if !auth.IsGranted(p) {
// 				panic(&errors.AccessDeniedError{Any: false, Denied: []security.Permission{p}})
// 			}
// 		}
// 		ctx.Next()
// 	}
// 	return &routeOptions{do: func(mo *routeOptions) {
// 		mo.BeforHandlers = append(mo.BeforHandlers, handler)
// 	}}
// }
// func AuthorityAny(perms ...security.Permission) RouteOption {
// 	if len(perms) == 0 {
// 		panic("perms cannot be nil")
// 	}
// 	handler := func(ctx Context) {
// 		auth := security.AuthenticationFromContext(ctx)
// 		for _, p := range perms {
// 			if auth.IsGranted(p) {
// 				ctx.Next()
// 				return
// 			}
// 		}
// 		panic(&errors.AccessDeniedError{Any: true, Denied: perms})
// 	}
// 	return &routeOptions{do: func(mo *routeOptions) {
// 		mo.BeforHandlers = append(mo.BeforHandlers, handler)
// 	}}
// }
