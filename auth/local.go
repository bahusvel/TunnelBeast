package auth

import (
	"log"

	"github.com/bahusvel/TunnelBeast/boltdb"
)

type LocalAuth struct {
	Username string
	Password string
}

func (this LocalAuth) Init() {
	//add admin user to DB
	err := boltdb.AddUser(this.Username, this.Password)
	if err == boltdb.ErrExists {
		return
	}
	if err != nil {
		log.Fatalln(err)
	}
}

func (this LocalAuth) Authenticate(Username string, Password string) bool {
	if Username == this.Username && Password == this.Password {
		return true
	}

	return boltdb.Authenticate(Username, Password)
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
