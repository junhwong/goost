package security

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
)

// Principal 表示主体的抽象概念，它可以用来表示任何实体，例如，个人、组织。
type Principal interface {
	// GetIdentity 返回此主体的唯一标识。
	GetIdentity() string
	// GetTenantIdentity 返回此主体的所属租户唯一标识。
	GetTenantIdentity() string
}

type Authentication interface {
	// 凭据
	GetCredentials() Credentials // token,sessionid
	// 主体
	GetPrincipal() Principal
	// 用户明细
	GetDetails() UserDetails
	// 是否认证
	IsAuthenticated() bool
	// IsGranted 确定当前用户是否拥有指定的权限
	IsGranted(perm Permission) bool
}

// 认证管理器
type AuthenticationManager interface {
	Authenticate(context.Context, Credentials) (Authentication, error)
}

const (
	AuthenticationContextKey = "security.AuthenticationContextKey"
)

// AuthenticationFromContext 从上下文中获取授权对象，如果未找到则返回匿名对象
func AuthenticationFromContext(ctx context.Context) Authentication {
	if auth, ok := ctx.Value(AuthenticationContextKey).(Authentication); ok && auth != nil {
		return auth
	}
	a := &anonymousAuthentication{}
	return a
}

// //https://docs.microsoft.com/zh-cn/dotnet/api/system.security.claims.claimtypes?view=netcore-3.1
// //https://docs.microsoft.com/zh-cn/dotnet/api/system.security.claims.claim?view=netcore-3.1
// type Claim struct {
// 	Type  string
// 	Name  string
// 	Value string
// }

//https://docs.microsoft.com/zh-cn/dotnet/api/system.net.networkcredential?view=netcore-3.1
//https://docs.microsoft.com/zh-cn/dotnet/api/system.net.credentialcache?view=netcore-3.1
type Credentials interface {
	GetKind() string
	// GetClaim(name string) *Claim
	Get(name string) interface{}
	GetString(name string) string
}

//

type NetworkCredentials struct {
	kind       string
	credential string // base64 tid:uid:pwd
	// Method   string
	// ClientIP string //RemoteAddress
	// Host     string //domain
	// Referer  string

	claims map[string]interface{} // 域
	// Audience  string // aud(Audience)：代表这个JWT的接收对象
	// ExpiresAt int64  // exp(Expiration time)：是一个时间戳，代表这个JWT的过期时间
	// Id        string // jti(JWT ID)：是JWT的唯一标识
	// IssuedAt  int64  // iat(Issued at)：是一个时间戳，代表这个JWT的签发时间
	// Issuer    string // iss(Issuser)：代表这个JWT的签发主体,颁发者身份标识
	// NotBefore int64  // nbf(Not Before)：是一个时间戳，代表这个JWT生效的开始时间，意味着在这个时间之前验证JWT是会失败的
	// Subject   string // sub(Subject)：代表这个JWT的主体，即它的所有人
}

func NewNetworkCredentials(kind string) *NetworkCredentials {
	return &NetworkCredentials{
		kind:   kind,
		claims: make(map[string]interface{}),
	}
}

// 由3部分组成，以:分割。tid?:uid:pwd
func ParseBasicCredential(s string) (tenantID, user, passwd string, err error) {
	var data []byte
	data, err = base64.StdEncoding.DecodeString(s)
	if err != nil {
		err = fmt.Errorf("invalid credential: %v", err)
		return
	}
	arr := strings.Split(string(data), ":")
	switch len(arr) {
	case 3:
		tenantID = arr[0]
		user = arr[1]
		passwd = arr[2]
	case 2:
		user = arr[0]
		passwd = arr[1]
	}
	return
}

func (token *NetworkCredentials) SetClaim(key string, val interface{}) error {
	switch key {
	case "credential":
		if s, ok := val.(string); ok {
			token.credential = s
			if strings.EqualFold(token.kind, "Basic") {
				// var err error
				// var tid string
				// tid, token.Name, token.Password, err = ParseBasicCredential(s)
				// if err != nil {
				// 	return err
				// }
				// if tid != "" {
				// 	token.TenantID = tid
				// }
			}
			return nil
		}
		return fmt.Errorf("value must be string: %v", val)
	default:
		token.claims[key] = val
	}
	return nil
}

func (token *NetworkCredentials) GetKind() string {
	return token.kind
}

func (token *NetworkCredentials) Get(key string) (v interface{}) {
	switch key {
	case "credential":
		return token.credential
	}
	return token.claims[key]
}

func (token *NetworkCredentials) GetString(key string) string {
	v := token.Get(key)
	if v == nil {
		return ""
	}
	if v, ok := v.(string); ok {
		return v
	}
	return ""
}

type AuthenticationProvider interface {
	Kind() string
	Do(ctx context.Context, credentials Credentials) (Authentication, error)
}

// type AuthenticationProviderFactory interface {
// 	Create() AuthenticationProvider
// }

// type AuthenticateService interface {
// 	DoAuthenticate(c Credentials) (Principal, error)
// }

type PasswordEncoder interface {
	Hash() string
	Encode(rawPasswd string) (string, error)
	Matches(rawPasswd string, encodedPasswd string) error
}

type PasswordEncoderFactory interface {
	Create() (PasswordEncoder, error)
}
