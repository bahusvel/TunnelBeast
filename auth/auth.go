package auth

type AuthProvider interface {
	Init()
	CheckDestinationIP(dstip string, Username string, Password string) bool
	Authenticate(Username string, Password string) bool
}

type TestAuth struct {
	Users map[string]string
}

func (this TestAuth) Init() {

}

func (this TestAuth) Authenticate(Username string, Password string) bool {
	if password, exists := this.Users[Username]; exists && password == Password {
		return true
	}
	return false
}

func (this TestAuth) CheckDestinationIP(dstip string, Username string, Password string) bool {
	return true
}
