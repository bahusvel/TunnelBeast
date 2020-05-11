package auth

import (
	"fmt"
	"log"
	"testing"

	"gopkg.in/ldap.v2"
)

func TestIPFromLDAP(t *testing.T) {
	this := LDAPAuth{LDAPAddr: "ldap.unitecloud.net:389", DCString: "dc=unitecloud,dc=net", IPAddressAttribute: "ipHostNumber"}

	l, err := ldap.Dial("tcp", this.LDAPAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	err = l.Bind("cn="+"tunnel2"+","+this.DCString, "cp-x2520")
	if err != nil {
		log.Println(err)
		return
	}

	ips, err := this.queryIPAddress(l, "tunnel2")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ips)
}
