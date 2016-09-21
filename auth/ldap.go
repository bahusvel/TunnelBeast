package auth

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"log"
)

type LDAPAuth struct {
	LDAPAddr           string
	DCString           string
	IPAddressAttribute string
	UserObjectClass    string
}

func (this LDAPAuth) Init() {

}

func (this LDAPAuth) queryIPAddress(LdapClient *ldap.Conn, Username string) ([]string, error) {

	searchRequest := ldap.NewSearchRequest(fmt.Sprintf("cn=%s,%s", Username, this.DCString),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass="+this.UserObjectClass+"))",   // The filter to apply
		[]string{"dn", "cn", this.IPAddressAttribute}, // A list attributes to retrieve
		nil,
	)

	sr, err := LdapClient.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}
	if len(sr.Entries) != 1 {
		return []string{}, fmt.Errorf("No entries, or too many entries")
	}
	return sr.Entries[0].GetAttributeValues(this.IPAddressAttribute), nil
}

func ipAllowed(ip string, whitelist []string) bool {
	for _, testIp := range whitelist {
		if ip == testIp || testIp == "*" {
			return true
		}
	}
	return false
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
	ipList, err := this.queryIPAddress(l, Username)
	if err != nil {
		log.Println(err)
		return false
	}
	return ipAllowed(InternalIP, ipList)
}
