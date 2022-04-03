package security

// 授权器
type Authorizer interface {
	// 载入主体权限,如果没有载入过
	LoadPermissionsIfAbsent(Principal) error
	// 判断主体是否有该权限
	IsPermitted(Principal, Permission) bool
}
