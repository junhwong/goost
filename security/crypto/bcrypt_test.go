package crypto

import "testing"

func TestBCryptPasswordEncoder(t *testing.T) {
	encoder := &BCryptPasswordEncoder{}
	s, err := encoder.Encode("123456dfgdgfdgdgdftgertexvx#@%#%@#%@!@@#%*(&")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s, len(s))
	err = encoder.Matches("123456", s)
	if err != nil {
		t.Fatal(err)
	}
	err = encoder.Matches("1234567", s)
	if err != nil {
		t.Fatal(err)
	}

}

func TestCreateDelegatingPasswordEncoder(t *testing.T) {
	encoder := CreateDelegatingPasswordEncoder()
	s, err := encoder.Encode("123456")
	t.Log(s, err)
	err = encoder.Matches("123456", s)
	t.Log(err)
}
