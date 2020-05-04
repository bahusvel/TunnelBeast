package auth

import "github.com/bahusvel/TunnelBeast/boltdb"

type AuthProvider interface {
	Init()
	CheckDestinationIP(dstip string, Username string, Password string) bool
	Authenticate(Username string, Password string) bool
	CheckAdminPanel(Username string, Password string) bool
}

type TestAuth struct {
	Username string
	Password string
}

func (this TestAuth) Init() {

}

func (this TestAuth) Authenticate(Username string, Password string) bool {
	if Username == this.Username && Password == this.Password {
		return true
	}

	key := "users/" + Username + "/" + Password
	_, err := boltdb.GetRecord(key)
	if err == nil {
		return true
	}

	return false
}

func (this TestAuth) CheckDestinationIP(dstip string, Username string, Password string) bool {
	return true
}

func (this TestAuth) CheckAdminPanel(Username string, Password string) bool {
	if Username == this.Username && Password == this.Password {
		return true
	}
	return false
}
