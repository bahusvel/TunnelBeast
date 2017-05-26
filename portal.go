package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bahusvel/TunnelBeast/auth"
	"github.com/bahusvel/TunnelBeast/config"
	"github.com/bahusvel/TunnelBeast/iptables"
)

const PORTAL_TEMPLATE = `
<html>
<head>
<style>
@import url(https://fonts.googleapis.com/css?family=Roboto:300);

.login-page {
  width: 360px;
  padding: 8% 0 0;
  margin: auto;
}
form {
  position: relative;
  z-index: 1;
  background: #FFFFFF;
  max-width: 360px;
  margin: 0 auto 100px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 0 20px 0 rgba(0, 0, 0, 0.2), 0 5px 5px 0 rgba(0, 0, 0, 0.24);
}
input {
  font-family: "Roboto", sans-serif;
  outline: 0;
  background: #f2f2f2;
  width: 100%;
  border: 0;
  margin: 0 0 15px;
  padding: 15px;
  box-sizing: border-box;
  font-size: 14px;
}
button {
  font-family: "Roboto", sans-serif;
  text-transform: uppercase;
  outline: 0;
  background: #4CAF50;
  width: 100%;
  border: 0;
  padding: 15px;
  color: #FFFFFF;
  font-size: 14px;
  -webkit-transition: all 0.3 ease;
  transition: all 0.3 ease;
  cursor: pointer;
}
body {
  background: #76b852; /* fallback for old browsers */
  background: -webkit-linear-gradient(right, #76b852, #8DC26F);
  background: -moz-linear-gradient(right, #76b852, #8DC26F);
  background: -o-linear-gradient(right, #76b852, #8DC26F);
  background: linear-gradient(to left, #76b852, #8DC26F);
  font-family: "Roboto", sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
</head>
<body>
<div class="login-page" >
    <form class="login-form" action="/auth" method="get">
      <input type="text" name="username" placeholder="username"/>
      <input type="password" name="password" placeholder="password"/>
	  <input type="text" name="internalip" placeholder="internalip"/>
      <button type="submit">login</button>
	  <p>Powered by <a href="https://github.com/bahusvel/TunnelBeast">TunnelBeast</a></p>
    </form>
  </div>
</div>
</body>
</html>
`

const LOGOUT_TEMPLATE = `
<html>
<head>
<style>
@import url(https://fonts.googleapis.com/css?family=Roboto:300);

.login-page {
  width: 360px;
  padding: 8% 0 0;
  margin: auto;
}
form {
  position: relative;
  z-index: 1;
  background: #FFFFFF;
  max-width: 360px;
  margin: 0 auto 100px;
  padding: 20px;
  text-align: center;
  box-shadow: 0 0 20px 0 rgba(0, 0, 0, 0.2), 0 5px 5px 0 rgba(0, 0, 0, 0.24);
}
button {
  font-family: "Roboto", sans-serif;
  text-transform: uppercase;
  outline: 0;
  background: #4CAF50;
  width: 100%;
  border: 0;
  padding: 15px;
  color: #FFFFFF;
  font-size: 14px;
  -webkit-transition: all 0.3 ease;
  transition: all 0.3 ease;
  cursor: pointer;
}
body {
  background: #76b852; /* fallback for old browsers */
  background: -webkit-linear-gradient(right, #76b852, #8DC26F);
  background: -moz-linear-gradient(right, #76b852, #8DC26F);
  background: -o-linear-gradient(right, #76b852, #8DC26F);
  background: linear-gradient(to left, #76b852, #8DC26F);
  font-family: "Roboto", sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
</head>
<body>
<div class="login-page" >
    <form class="login-form" action="/logout" method="get">
	  <p>You are already logged in</p>
      <button type="submit">logout</button>
	  <p>Powered by <a href="https://github.com/bahusvel/TunnelBeast">TunnelBeast</a></p>
    </form>
  </div>
</div>
</body>
</html>
`

var connectionTable = map[string]string{}
var authProvider auth.AuthProvider

func AuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Api access", r.RemoteAddr)
	w.Header().Set("Cache-Control", "no-cache")
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	internalip := r.URL.Query().Get("internalip")
	if username == "" || password == "" || internalip == "" {
		w.Write([]byte("ERROR"))
		return
	}
	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	if _, exists := connectionTable[clientIP]; exists {
		w.Write([]byte("LOGGEDIN"))
		return
	}
	if authProvider.Authenticate(username, password, internalip) {
		err := iptables.NewRoute(clientIP, internalip)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		connectionTable[clientIP] = internalip
		http.Redirect(w, r, "/", 302)
	} else {
		w.Write([]byte("ERROR"))
	}
}

func SignoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Logout request", r.RemoteAddr)
	w.Header().Set("Cache-Control", "no-cache")
	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	if internalIP, exists := connectionTable[clientIP]; exists {
		delete(connectionTable, clientIP)
		err := iptables.DeleteRoute(clientIP, internalIP)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		http.Redirect(w, r, "/", 302)
		return
	}
	w.Write([]byte("Not Logged In"))
}

func PortalEntryHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Portal access", r.RemoteAddr)
	w.Header().Set("Cache-Control", "no-cache")
	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	if _, exists := connectionTable[clientIP]; exists {
		w.Write([]byte(LOGOUT_TEMPLATE))
	} else {
		w.Write([]byte(PORTAL_TEMPLATE))
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Usage: TunnelBeast /path/to/config.yml")
	}
	conf := config.Configuration{}
	config.LoadConfig(os.Args[1], &conf)
	log.Printf("%+v\n", conf)

	authProvider = conf.AuthProvider
	authProvider.Init()

	err := iptables.Init(conf.ListenDev)
	if err != nil {
		log.Println("Error initializing iptables", err)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", PortalEntryHandler)
	mux.HandleFunc("/auth", AuthenticationHandler)
	mux.HandleFunc("/logout", SignoutHandler)

	port80 := &http.Server{Addr: ":80", Handler: mux}
	port666 := &http.Server{Addr: ":666", Handler: mux}
	go func() {
		errInternal := port80.ListenAndServe()
		if errInternal != nil {
			log.Println(errInternal)
		}
	}()
	err = port666.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
