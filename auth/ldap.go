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
	UserRDN            string
}

func (this LDAPAuth) Init() {

}

func (this LDAPAuth) queryIPAddress(LdapClient *ldap.Conn, username string) ([]string, error) {

	searchRequest := ldap.NewSearchRequest(fmt.Sprintf("%s=%s,%s", this.UserRDN, username, this.DCString),
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=*))",                    // The filter to apply
		[]string{"dn", this.IPAddressAttribute}, // A list attributes to retrieve
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
			return false
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

func (this LDAPAuth) CheckDestinationIP(dstip string, username string, password string) bool {
	if this.IPAddressAttribute == "" {
		return true
	}

	l, err := ldap.Dial("tcp", this.LDAPAddr)

	if err != nil {
		log.Println(err)
		return false
	}

	err = l.Bind(fmt.Sprintf("%s=%s,%s", this.UserRDN, username, this.DCString), password)
	if err != nil {
		log.Println(err)
		return false
	}

	ipList, err := this.queryIPAddress(l, username)
	if err != nil {
		log.Println(err)
		return false
	}
	if len(ipList) == 0 {
		return true // if IPAddressAttribute is not set in LDAP for the user
	}
	return ipAllowed(dstip, ipList)
}

func (this LDAPAuth) Authenticate(username string, password string) bool {
	l, err := ldap.Dial("tcp", this.LDAPAddr)
	if err != nil {
		log.Println(err)
		return false
	}
	defer l.Close()

	err = l.Bind(fmt.Sprintf("%s=%s,%s", this.UserRDN, username, this.DCString), password)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (this LDAPAuth) CheckAdminPanel(Username string, Password string) bool {
	return false
}
