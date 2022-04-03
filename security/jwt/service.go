package jwt

import (
	"fmt"
	"time"

	"github.com/junhwong/goost/idworker"
	"github.com/junhwong/goost/security"
	"github.com/spf13/cast"

	v3 "github.com/gbrlsnchs/jwt/v3"
)

type Provider interface {
	NewClaims(timeout time.Duration) ClaimsBuilder
	Sign(interface{}) (string, error)
	Verify(string, time.Duration) (ClaimsGetter, error)
}

type providerImpl struct {
	alg v3.Algorithm
	// issuer     string
	// subject    string
	// audience   v3.Audience
	// expiration time.Duration
}

func (p *providerImpl) Sign(payload interface{}) (string, error) {
	if payload == nil {
		return "", fmt.Errorf("payload cannot be empty")
	}

	token, err := v3.Sign(payload, p.alg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", token), nil
}

// func (p *providerImpl) Verify(token string, target interface{}) (err error) {
// 	if token == "" {
// 		return fmt.Errorf("token cannot be empty")
// 	}
// 	// pl := &v3.Payload{}
// 	// now := time.Now()
// 	_, err = v3.Verify([]byte(token), p.alg, target,
// 		v3.ValidateHeader,
// 		// v3.ValidatePayload(pl,
// 		// 	v3.IssuedAtValidator(now),
// 		// 	v3.IssuerValidator(p.issuer),
// 		// 	v3.ExpirationTimeValidator(now),
// 		// 	v3.AudienceValidator(p.audience),
// 		// ),
// 	)
// 	return
// 	// if err != nil {
// 	// 	ite := &security.InvaildTokenError{Err: err}
// 	// 	if errors.Is(err, v3.ErrExpValidation) {
// 	// 		ite.Message = "the token expired"
// 	// 	}
// 	// 	return "", ite
// 	// }
// 	// return pl.JWTID, nil
// }

func (p *providerImpl) Verify(token string, timeout time.Duration) (ClaimsGetter, error) {
	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}
	pl := p.newClaims(timeout)
	pl.Set("jti", nil)
	// now := time.Now()
	_, err := v3.Verify([]byte(token), p.alg, &pl.Claims,
		v3.ValidateHeader,
		// v3.VerifyOption
		// v3.ValidatePayload(pl,
		// v3.IssuedAtValidator(now),
		// 	v3.IssuerValidator(p.issuer),
		// 	v3.ExpirationTimeValidator(now),
		// 	v3.AudienceValidator(p.audience),
		// ),
	)
	if err != nil {
		return nil, &security.InvaildTokenError{Err: err}
	}
	err = expValidator()(pl)
	if err != nil {
		return nil, &security.InvaildTokenError{Err: err}
	}

	return pl, nil
	// if err != nil {
	// 	ite := &security.InvaildTokenError{Err: err}
	// 	if errors.Is(err, v3.ErrExpValidation) {
	// 		ite.Message = "the token expired"
	// 	}
	// 	return "", ite
	// }
	// return pl.JWTID, nil
}

type Validator func(pl ClaimsGetter) error

// ExpireAtValidator validates the "iat" claim.
func expValidator() Validator {
	return func(pl ClaimsGetter) error {
		if NumericDate(time.Now()).After(pl.ExpireAt()) {
			return v3.ErrExpValidation
		}
		return nil
	}
}

type ProviderOption interface {
	applyProviderOption(*providerImpl) error
}

func NewProvider(options ...ProviderOption) func() (Provider, error) {
	return func() (Provider, error) {
		impl := &providerImpl{}
		var err error
		for _, opt := range options {
			if opt != nil {
				err = opt.applyProviderOption(impl)
				if err != nil {
					break
				}
			}
		}
		return impl, err
	}
}

type providerOptionFunc func(*providerImpl) error

func (f providerOptionFunc) applyProviderOption(opt *providerImpl) error {
	return f(opt)
}
func WithHS256(secret string) providerOptionFunc {
	return func(pi *providerImpl) error {
		pi.alg = v3.NewHS256([]byte(secret))
		return nil
	}
}

type (
	Time     = v3.Time
	Audience = v3.Audience
)

var (
	NumericDate = v3.NumericDate
)

type ClaimsBuilder interface {
	security.ClaimsSetter
	SetIssuedAt(time.Time) ClaimsBuilder
	SetAudience(...string) ClaimsBuilder
	SetSubject(string) ClaimsBuilder
	SetIssuer(string) ClaimsBuilder
	Sign() (string, error)
}
type ClaimsGetter interface {
	security.ClaimsGetter
	IssuedAt() time.Time
	ExpireAt() time.Time
	Audience() []string
	Subject() string
	Issuer() string
	ID() string
}

type claimsImpl struct {
	security.BasicClaims
	timeout  time.Duration
	provider Provider
}

func (p *providerImpl) newClaims(timeout time.Duration) *claimsImpl {
	c := &claimsImpl{
		provider: p,
		timeout:  timeout,
	}
	c.Set("jti", fmt.Sprint(idworker.NextId()))
	return c
}
func (p *providerImpl) NewClaims(timeout time.Duration) ClaimsBuilder {
	c := &claimsImpl{
		provider: p,
		timeout:  timeout,
	}
	c.Set("jti", fmt.Sprint(idworker.NextId()))
	return c
}

func (c *claimsImpl) SetIssuedAt(v time.Time) ClaimsBuilder {
	t := NumericDate(v)
	c.Set("iat", t.Unix())
	c.Set("exp", t.Add(c.timeout).Unix())
	return c
}
func (c *claimsImpl) SetAudience(v ...string) ClaimsBuilder {
	c.Set("aud", Audience(v))
	return c
}
func (c *claimsImpl) SetSubject(v string) ClaimsBuilder {
	c.Set("sub", v)
	return c
}
func (c *claimsImpl) SetIssuer(v string) ClaimsBuilder {
	c.Set("iss", v)
	return c
}
func (c *claimsImpl) Sign() (string, error) {
	//c.ToJSON()
	return c.provider.Sign(c.GetAll())
}

func (c *claimsImpl) IssuedAt() time.Time {
	v := cast.ToInt64(c.Get("iat"))
	return NumericDate(time.Unix(v, 0)).Time
}
func (c *claimsImpl) ExpireAt() time.Time {
	v := cast.ToInt64(c.Get("exp"))
	if v > 0 {
		return NumericDate(time.Unix(v, 0)).Time
	}
	v = cast.ToInt64(c.Get("iat"))
	return NumericDate(time.Unix(v+int64(c.timeout.Seconds()), 0)).Time
}
func (c *claimsImpl) Audience() []string {
	v := cast.ToStringSlice(c.Get("aud"))
	// v := cast.ToString(c.Get("aud"))
	return v
}
func (c *claimsImpl) Subject() string {
	v := cast.ToString(c.Get("sub"))
	return v
}
func (c *claimsImpl) Issuer() string {
	v := cast.ToString(c.Get("iss"))
	return v
}
func (c *claimsImpl) ID() string {
	v := cast.ToString(c.Get("jti"))
	return v
}
