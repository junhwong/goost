package security

import (
	"fmt"
)

// Permission 表示特定域对象的权限。
type Permission [2]uint

func (perm Permission) String() string {
	return fmt.Sprintf("%d:%d", perm[0], perm[1])
}

// Acl 表示域对象的访问控制列表（access control list）。
type Acl map[uint]uint

func (acl Acl) Add(perm Permission) error {
	if acl == nil {
		return nil
	}
	index, mask := perm[0], perm[1]
	acl[index] |= mask
	return nil
}

func (acl Acl) Remove(perm Permission) error {
	if acl == nil {
		return nil
	}
	index, mask := perm[0], perm[1]
	acl[index] = acl[index] & ^mask
	if acl[index] == 0 {
		delete(acl, index)
	}
	return nil
}

// Has 判断 acl 中是否包含给定 perm。
func (acl Acl) Has(perm Permission) bool {
	if acl == nil {
		return false
	}
	index, mask := perm[0], perm[1]
	return (acl[index] & mask) == mask
}

//https://docs.spring.io/spring-security/site/docs/5.4.1/api/
//https://docs.microsoft.com/zh-cn/dotnet/api/system.security.principal.iprincipal?view=netcore-3.1
//https://docs.microsoft.com/zh-cn/dotnet/api/system.security?view=net-5.0
//https://segmentfault.com/a/1190000021552296?utm_source=sf-related
//https://www.jianshu.com/p/a773806f9831
//https://connect2id.com/products/nimbus-oauth-openid-connect-sdk
//https://www.keycloak.org/
//https://docs.microsoft.com/zh-cn/dotnet/api/system.security.claims.claimtypes?view=netcore-3.1
//https://developer.github.com/apps/building-oauth-apps/authorizing-oauth-apps/
//https://docs.microsoft.com/zh-cn/dotnet/api/system.net.authenticationmanager?view=netcore-3.1
type Attributes interface {
}

var (
	AnyonePermission        = Permission{0, 0} // 任何人
	AnonymousPermission     = Permission{0, 1} // 匿名用户，即未登录
	AuthenticatedPermission = Permission{0, 2} // 认证用户，即已登录
	AdministratorPermission = Permission{0, 3} // 管理员
)
