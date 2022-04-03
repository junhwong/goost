package security

import "encoding/json"

type ClaimsSetter interface {
	Set(string, interface{}) ClaimsSetter
}
type ClaimsGetter interface {
	Get(string) interface{}
	GetAll() map[string]interface{}
}

type BasicClaims struct {
	Claims map[string]interface{}
}

func (c *BasicClaims) Set(k string, v interface{}) ClaimsSetter {
	if c.Claims == nil {
		c.Claims = map[string]interface{}{}
	}
	if v == nil {
		delete(c.Claims, k)
	} else {
		c.Claims[k] = v
	}
	return c
}

func (c *BasicClaims) Get(k string) interface{} {
	if c.Claims == nil {
		return nil
	}
	return c.Claims[k]
}

func (c *BasicClaims) ToJSON(m ...func(v interface{}) ([]byte, error)) ([]byte, error) {
	var marshal func(v interface{}) ([]byte, error)
	for _, it := range m {
		if it != nil {
			marshal = it
		}
	}
	if marshal == nil {
		marshal = json.Marshal
	}
	return marshal(c.Claims)
}

func (c *BasicClaims) GetAll() map[string]interface{} {
	return c.Claims
}
