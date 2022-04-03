package security

// 敏感字符串，如：密码
type SensitiveString string

func (ss SensitiveString) String() string {
	if ss == "" {
		return ""
	}
	return "***"
}

func (ss SensitiveString) Raw() string {
	return string(ss)
}

func (ss SensitiveString) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ss.String() + `"`), nil
}

type PasswordEncoder interface {
	Hash() string
	Encode(rawPasswd string) (string, error)
	Matches(rawPasswd string, encodedPasswd string) error
}

type PasswordEncoderFactory interface {
	Create() (PasswordEncoder, error)
}
