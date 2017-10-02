package auth

import (
	"fmt"
	"log"
    "net"
    "gopkg.in/ldap.v2"
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
	ipAddr := net.ParseIP(ip)
    for _, testIp := range whitelist {
		_, ipv4Net, err := net.ParseCIDR(testIp)
        if err != nil {
            log.Println(err)
        }
        if ipv4Net.Contains(ipAddr) {
			return true
		}
	}
	return false
}

func (this LDAPAuth) CheckSourceIP(srcip string) bool {
	return true
}

func (this LDAPAuth) CheckDestinationIP(dstip string, Username string) bool {
	if this.IPAddressAttribute == "" {
		return true
	}
    l, err := ldap.Dial("tcp", this.LDAPAddr)
	ipList, err := this.queryIPAddress(l, Username)
	if err != nil {
		log.Println(err)
		return false
	}
	return ipAllowed(dstip, ipList)
	return true
}

func (this LDAPAuth) Authenticate(Username string, Password string) bool {
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
