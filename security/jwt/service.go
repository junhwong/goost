package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/junhwong/goost/security"

	v3 "github.com/gbrlsnchs/jwt/v3"
)

type JwtProvider interface {
	Sign(jti string) (string, error)
	Verify(token string) (jti string, err error)
}

func HS256(secret, issuer, subject string, audience []string, expiration time.Duration) JwtProvider {
	return &JwtProviderImpl{
		secret:     v3.NewHS256([]byte(secret)),
		issuer:     issuer,
		subject:    subject,
		audience:   audience,
		expiration: expiration,
	}
}

type JwtProviderImpl struct {
	secret     *v3.HMACSHA
	issuer     string
	subject    string
	audience   v3.Audience
	expiration time.Duration
}

func (p *JwtProviderImpl) Sign(jti string) (string, error) {
	if jti == "" {
		return "", fmt.Errorf("jti cannot be empty")
	}
	now := time.Now()
	pl := v3.Payload{
		Issuer:         p.issuer,                              // iss(Issuser)：代表这个JWT的签发主体
		Subject:        p.subject,                             // sub(Subject)：代表这个JWT的主体，即它的所有人
		Audience:       p.audience,                            // aud(Audience)：代表这个JWT的接收对象
		ExpirationTime: v3.NumericDate(now.Add(p.expiration)), // exp(Expiration time)：是一个时间戳，代表这个JWT的过期时间
		IssuedAt:       v3.NumericDate(now),                   // iat(Issued at)：是一个时间戳，代表这个JWT的签发时间
		JWTID:          jti,                                   // jti(JWT ID)：是JWT的唯一标识
		// NotBefore:      v3.NumericDate(now.Add(30 * time.Minute)),         // nbf(Not Before)：是一个时间戳，代表这个JWT生效的开始时间，意味着在这个时间之前验证JWT是会失败的
	}
	token, err := v3.Sign(pl, p.secret)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", token), nil
}
func (p *JwtProviderImpl) Verify(token string) (jti string, err error) {
	if token == "" {
		return "", fmt.Errorf("token cannot be empty")
	}
	pl := &v3.Payload{}
	now := time.Now()
	_, err = v3.Verify([]byte(token), p.secret, pl,
		v3.ValidateHeader,
		v3.ValidatePayload(pl,
			v3.IssuedAtValidator(now),
			v3.IssuerValidator(p.issuer),
			v3.ExpirationTimeValidator(now),
			v3.AudienceValidator(p.audience),
		),
	)
	if err != nil {
		ite := &security.InvaildTokenError{Err: err}
		if errors.Is(err, v3.ErrExpValidation) {
			ite.Message = "the token expired"
		}
		return "", ite
	}
	return pl.JWTID, nil
}
