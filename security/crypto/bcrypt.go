package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

type BCryptPasswordEncoder struct {
	Cost int
}

//{MD5}BCrypt
//https://www.cnkirito.moe/spring-security-6/
//https://www.cnblogs.com/jpfss/p/11005125.html

func (encoder *BCryptPasswordEncoder) Hash() string {
	return "BCrypt"
}

func (be *BCryptPasswordEncoder) Encode(rawPasswd string) (string, error) {
	data, err := bcrypt.GenerateFromPassword([]byte(rawPasswd), be.Cost)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (be *BCryptPasswordEncoder) Matches(rawPasswd string, encodedPasswd string) error {
	return bcrypt.CompareHashAndPassword([]byte(encodedPasswd), []byte(rawPasswd))
}
