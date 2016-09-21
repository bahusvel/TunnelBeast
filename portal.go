package main

import (
	"github.com/bahusvel/TunnelBeast/iptables"
	"log"
	"net/http"
	"strings"
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
  padding: 45px;
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
.form button:hover,.form button:active,.form button:focus {
  background: #43A047;
}
.form .message {
  margin: 15px 0 0;
  color: #b3b3b3;
  font-size: 12px;
}
.form .message a {
  color: #4CAF50;
  text-decoration: none;
}
.form .register-form {
  display: none;
}
.container {
  position: relative;
  z-index: 1;
  max-width: 300px;
  margin: 0 auto;
}
.container:before, .container:after {
  content: "";
  display: block;
  clear: both;
}
.container .info {
  margin: 50px auto;
  text-align: center;
}
.container .info h1 {
  margin: 0 0 15px;
  padding: 0;
  font-size: 36px;
  font-weight: 300;
  color: #1a1a1a;
}
.container .info span {
  color: #4d4d4d;
  font-size: 12px;
}
.container .info span a {
  color: #000000;
  text-decoration: none;
}
.container .info span .fa {
  color: #EF3B3A;
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
    </form>
  </div>
</div>
</body>
</html>
`

var connectionTable = map[string]string{}

func AuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Api access", r.RemoteAddr)
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	internalip := r.URL.Query().Get("internalip")
	if username == "" || password == "" || internalip == "" {
		w.Write([]byte("ERROR"))
		return
	}
	if username == "denis" && password == "cp-x2520" {
		clientIP := strings.Split(r.RemoteAddr, ":")[0]
		err := iptables.NewRoute(clientIP, internalip)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		connectionTable[clientIP] = internalip
		w.Write([]byte("OK"))
	} else {
		w.Write([]byte("ERROR"))
	}
}

func SignoutHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Logout request", r.RemoteAddr)
	clientIP := strings.Split(r.RemoteAddr, ":")[0]
	if internalIP, exists := connectionTable[clientIP]; exists {
		delete(connectionTable, clientIP)
		err := iptables.DeleteRoute(clientIP, internalIP)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("OK"))
		return
	}
	w.Write([]byte("Not Logged In"))
}

func PortalEntryHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Portal access", r.RemoteAddr)
	w.Write([]byte(PORTAL_TEMPLATE))
}

func main() {
	err := iptables.Init()
	if err != nil {
		log.Println("Error initializing IPtables", err)
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
