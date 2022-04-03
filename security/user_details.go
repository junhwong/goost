package security

// 用户详情
type UserDetails interface {
	GetIdentity() string       // 返回此主体的唯一标识。
	GetTenantIdentity() string // 返回用户所属租户唯一标识。
	GetUsername() string       // 返回用户名
	IsExpired() bool           // 是否过期
	IsLocked() bool            // 是否锁定
	IsDisabled() bool          // 是否禁用
}

// type UserDetailsService interface {
// 	LoadUserById(ctx context.Context, id string) (UserDetails, error)
// 	LoadUser(ctx context.Context, username string) (UserDetails, error)
// 	LoadUserWithPassword(ctx context.Context, username string, passwd SensitiveString) (UserDetails, error)
// }
