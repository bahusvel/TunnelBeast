package auth

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"testing"
)

func TestIPFromLDAP(t *testing.T) {
	this := LDAPAuth{LDAPAddr: "192.168.1.90:389", DCString: "dc=unitecloud,dc=net", IPAddressAttribute: "postalAddress", UserObjectClass: "simpleSecurityObject"}

	l, err := ldap.Dial("tcp", this.LDAPAddr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()
	ips, err := this.queryIPAddress(l, "guest")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ips)
}
