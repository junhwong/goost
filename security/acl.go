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

var (
	AnyonePermission        = Permission{0, 0} // 任何人
	AnonymousPermission     = Permission{0, 1} // 匿名用户，即未登录
	AuthenticatedPermission = Permission{0, 2} // 认证用户，即已登录
	AdministratorPermission = Permission{0, 3} // 管理员
)
