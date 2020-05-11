package auth

import "github.com/bahusvel/TunnelBeast/boltdb"

type LocalAuth struct {
	Username string
	Password string
}

func (this LocalAuth) Init() {

}

func (this LocalAuth) Authenticate(Username string, Password string) bool {
	if Username == this.Username && Password == this.Password {
		return true
	}

	key := "users/" + Username + "/" + Password
	_, err := boltdb.GetFavourite(key)
	if err == nil {
		return true
	}

	return false
}

func (this LocalAuth) CheckDestinationIP(dstip string, Username string, Password string) bool {
	return true
}

func (this LocalAuth) CheckAdminPanel(Username string, Password string) bool {
	if Username == this.Username && Password == this.Password {
		return true
	}
	return false
}
