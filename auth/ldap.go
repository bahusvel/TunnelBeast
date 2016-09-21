package auth

import (
	"log"

	"gopkg.in/ldap.v2"
)

type LDAPAuth struct {
	LDAPAddr string
	DCString string
}

func (this LDAPAuth) Init() {

}

func (this LDAPAuth) Authenticate(Username string, Password string, InternalIP string) bool {
	l, err := ldap.Dial("tcp", this.LDAPAddr)
	if err != nil {
		log.Println(err)
		return false
	}
	defer l.Close()

	err = l.Bind("cn="+Username+","+this.DCString, Password)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
