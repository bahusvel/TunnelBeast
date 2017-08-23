package auth

type AuthProvider interface {
	Init()
	Authenticate(Username string, Password string) bool
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
	return false
}
