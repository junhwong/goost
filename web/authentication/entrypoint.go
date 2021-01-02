package authentication

import (
	"strings"

	"github.com/junhwong/goost/data"
	"github.com/junhwong/goost/errors"
	"github.com/junhwong/goost/security"
	"github.com/junhwong/goost/security/jwt"
	"github.com/junhwong/goost/web"
	"golang.org/x/crypto/bcrypt"
)

type TokenProvider interface {
	GetIdentity(credentials security.Credentials) (string, error)
}

//JwtTokenEnhancer implements TokenEnhancer enhance
//https://blog.csdn.net/zhenghongcs/article/details/107241168?utm_medium=distribute.pc_relevant_bbs_down.none-task-blog-baidujs-2.nonecase&depth_1-utm_source=distribute.pc_relevant_bbs_down.none-task-blog-baidujs-2.nonecase
//https://www.cnblogs.com/harrychinese/p/SpringBoot_security_basics.html

// see also: [WWW-Authenticate](https://tools.ietf.org/html/rfc6750#section-3)
type BearerTokenAuthenticationEntryPoint struct {
	Provider           TokenProvider
	UserDetailsService security.UserDetailsService
}

func (p *BearerTokenAuthenticationEntryPoint) Commence(ctx web.Context) (security.Authentication, error) {
	credentials, err := GetTokenFromHeader(ctx)
	if credentials == nil || err != nil {
		return nil, err
	}
	id, err := p.Provider.GetIdentity(credentials)
	if err != nil {
		var ite *security.InvaildTokenError
		if e, ok := err.(*security.InvaildTokenError); ok {
			ite = e
		}
		if ite == nil {
			ite = &security.InvaildTokenError{Err: err}
		}
		ite.TokenType = credentials.GetKind()
		return nil, ite
	}

	details, err := p.UserDetailsService.LoadUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &AuthenticationImpl{
		credentials: credentials,
		details:     details,
	}, nil
}

type JwtAuthenticationProvider struct {
	Service jwt.JwtProvider
}

func (p *JwtAuthenticationProvider) GetIdentity(credentials security.Credentials) (string, error) {
	tokenString := credentials.GetString("credential")
	return p.Service.Verify(tokenString)
}

//login_url_
//authentication_filter_base
//ClaimAccessor

type UserPasswordToken struct {
	ClientID     string `json:"client_id" form:"client_id"`         // 客户端标识
	ClientSecret string `json:"client_secret" form:"client_secret"` // 客户端密钥
	User         string `json:"username" form:"username"`           // 登录名。不能以 `ROLE_`,`SCOPE_`开头，不能包含 `:`,`,`字符。
	Password     string `json:"password" form:"password"`           // 密码
	Scope        string `json:"scope" form:"scope"`                 // 密码
}

type LoginUrlAuthenticationEntryPoint struct {
	Path               string
	UserDetailsService security.UserDetailsService
	PasswordEncoder    security.PasswordEncoder
}
type UserPasswordForm struct {
	ClientID     string `json:"client_id" form:"client_id"` // 客户端标识
	ClientSecret string `json:"client_secret" form:"client_secret"`
	Username     string `json:"username" form:"username"` // 登录名。(kind<uint8|string>:user)
	Password     string `json:"password" form:"password"` // 密码
}

func (p *LoginUrlAuthenticationEntryPoint) Commence(ctx web.Context) (security.Authentication, error) {
	if !strings.EqualFold(ctx.Request.URL.Path, p.Path) {
		return nil, nil
	}
	if !strings.EqualFold(ctx.Request.Method, "POST") {
		return nil, nil
	}
	// TODO  防止恶意用户名、密码爆猜
	form := UserPasswordForm{}
	if err := ctx.Bind(&form); err != nil {
		return nil, err // 400
	}

	if form.Username == "" {
		return nil, &errors.ArgumentError{Field: "username", Code: "missing_field"}
	}

	details, err := p.UserDetailsService.LoadUser(ctx, form.ClientID, form.Username)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return nil, errors.NewInvalidArgumentError("username,password", "", err)
		}
		return nil, err
	}

	err = p.PasswordEncoder.Matches(form.Password, details.GetPassword())
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, errors.NewInvalidArgumentError("username,password", "", err)
		}
		return nil, err
	}

	return &AuthenticationImpl{
		credentials: nil,
		details:     details,
	}, nil
}

type AuthenticationImpl struct {
	credentials security.Credentials
	details     security.UserDetails
}

func (a *AuthenticationImpl) GetCredentials() security.Credentials {
	return a.credentials
}
func (a *AuthenticationImpl) GetDetails() security.UserDetails {
	return a.details
}
func (a *AuthenticationImpl) GetPrincipal() security.Principal {
	return a.details
}
func (a *AuthenticationImpl) IsAuthenticated() bool {
	// TODO: 过期时间
	return true
}
func (a *AuthenticationImpl) IsGranted(perm security.Permission) bool {
	// TODO:
	return false
}
