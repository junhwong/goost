package authentication

import (
	"fmt"
	"strings"

	"github.com/junhwong/goost/security"
	"github.com/junhwong/goost/web"
)

func SetClaimsWithContext(ctx web.Context, credentials security.NetworkCredentials) error {
	if err := credentials.SetClaim("host", ctx.Request.Host); err != nil {
		return err
	}
	if err := credentials.SetClaim("http_method", ctx.Request.Method); err != nil {
		return err
	}
	if err := credentials.SetClaim("remote_address", ctx.ClientIP()); err != nil {
		return err
	}
	if err := credentials.SetClaim("referer", ctx.GetHeader("Referer")); err != nil {
		return err
	}
	return nil
}

//git clone github_jun:junhwong/duzee-go.git
//https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#Authentication_schemes
//https://docs.aws.amazon.com/AmazonS3/latest/API/sigv4-auth-using-authorization-header.html
func GetTokenFromHeader(ctx web.Context) (security.Credentials, error) {
	s := ctx.GetHeader("Authorization")
	if s == "" {
		return nil, nil //fmt.Errorf("authentication is required")
	}
	arr := strings.Split(s, " ")
	if len(arr) != 2 {
		return nil, fmt.Errorf("Bad Authentication format: %s", s)
	}
	credentials := security.NewNetworkCredentials(arr[0]) // 一般是 Bearer, https://cloud.tencent.com/developer/article/1586937
	if err := credentials.SetClaim("credential", arr[1]); err != nil {
		return nil, err
	}

	if err := SetClaimsWithContext(ctx, *credentials); err != nil {
		return nil, err
	}

	return credentials, nil

}
