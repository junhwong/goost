package security

import "testing"

func TestAcl(t *testing.T) {
	acl := Acl{}

	acl.Add(Permission{0, 1})
	acl.Add(Permission{0, 1})
	t.Log(acl)
	t.Log(acl.Has(Permission{0, 1}))
	t.Log(acl.Has(Permission{}))
	acl.Remove(Permission{0, 1})
	t.Log(acl.Has(Permission{0, 1}))
	t.Log(acl)
}

func TestA(t *testing.T) {
	acl := Acl{}

	// acl.Add(AnonymousPermission)

	t.Log(acl.Has(AnonymousPermission))

}
