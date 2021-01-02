package security

type anonymousAuthentication struct {
}

func (a *anonymousAuthentication) GetCredentials() Credentials { return nil }
func (a *anonymousAuthentication) GetPrincipal() Principal     { return nil }
func (a *anonymousAuthentication) GetDetails() UserDetails     { return nil }
func (a *anonymousAuthentication) IsAuthenticated() bool       { return false }
func (a *anonymousAuthentication) IsGranted(perm Permission) bool {
	return perm[0] == AnonymousPermission[0] && perm[1] == AnonymousPermission[1]
}
