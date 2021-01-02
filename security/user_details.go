package security

import (
	"context"
)

// 敏感字符串，如：密码
type SensitiveString string

func (ss SensitiveString) String() string {
	return "***"
}

func (ss SensitiveString) Raw() string {
	return string(ss)
}

func (ss SensitiveString) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ss.String() + `"`), nil
}

type UserDetails interface {
	Principal
	GetPassword() string
	GetUsername() string
	// GetNameKind() NameKind
	IsExpired() bool
	IsLocked() bool
	IsEnabled() bool
}

type UserDetailsService interface {
	LoadUserById(ctx context.Context, userID string) (UserDetails, error)
	LoadUser(ctx context.Context, clientID string, username string) (UserDetails, error)
}
