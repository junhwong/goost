package crypto

import (
	"fmt"
	"strings"

	"github.com/junhwong/goost/security"
	"golang.org/x/crypto/bcrypt"
)

type DelegatingPasswordEncoder struct {
	hash     string
	def      security.PasswordEncoder
	encoders map[string]security.PasswordEncoder
}

func (encoder *DelegatingPasswordEncoder) SetDefault(hash string) error {
	hash = strings.ToUpper(hash)
	if encoder.hash == hash {
		return nil
	}
	e := encoder.encoders[hash]
	if e == nil {
		return fmt.Errorf("hash %q not defined", hash)
	}
	encoder.def = e
	encoder.hash = hash
	return nil
}
func (encoder *DelegatingPasswordEncoder) Hash() string {
	return "delegating"
}
func (encoder *DelegatingPasswordEncoder) Encode(rawPasswd string) (string, error) {
	encodedPasswd, err := encoder.def.Encode(rawPasswd)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("{%s}%s", encoder.def.Hash(), encodedPasswd), nil
}
func (encoder *DelegatingPasswordEncoder) Matches(rawPasswd string, encodedPasswd string) error {
	if !strings.HasPrefix(encodedPasswd, "{") {
		return fmt.Errorf("security/crypto: bad encoded password format")
	}
	i := strings.Index(encodedPasswd, "}")
	if i < 1 {
		return fmt.Errorf("security/crypto: bad encoded password format")
	}
	hash := strings.ToUpper(encodedPasswd[1:i])
	e := encoder.encoders[hash]
	if e == nil {
		return fmt.Errorf("security/crypto: hash %q not defined", hash)
	}

	return e.Matches(rawPasswd, encodedPasswd[i+1:])
}

var (
	ErrMismatchedHashAndPassword = bcrypt.ErrMismatchedHashAndPassword
)

func CreateDelegatingPasswordEncoder(encoders ...security.PasswordEncoder) *DelegatingPasswordEncoder {
	def := &BCryptPasswordEncoder{}
	defHash := strings.ToUpper(def.Hash())
	dpe := &DelegatingPasswordEncoder{
		encoders: map[string]security.PasswordEncoder{defHash: def},
		def:      def,
		hash:     defHash,
	}

	for _, encoder := range encoders {
		if encoder == nil {
			continue
		}
		hash := strings.ToUpper(encoder.Hash())
		if hash == "delegating" {
			panic(fmt.Errorf("hash cannot be %q", encoder.Hash()))
		}
		dpe.encoders[hash] = encoder
		if hash == defHash {
			dpe.def = def
		}
	}

	return dpe
}
