package main

import (
	"github.com/bahusvel/TunnelBeast/iptables"
	"log"
	"net/http"
	"strings"
)

const PORTAL_TEMPLATE = `
<div class="login-page" >
    <form class="login-form" action="/auth" method="get">
      <input type="text" name="username" placeholder="username"/>
      <input type="password" name="password" placeholder="password"/>
	  <input type="text" name="internalip" placeholder="internalip"/>
      <button type="submit">login</button>
    </form>
  </div>
</div>
`

func AuthenticationHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	internalip := r.URL.Query().Get("internalip")
	if username == "" || password == "" || internalip == "" {
		w.Write([]byte("ERROR"))
		return
	}
	if username == "denis" && password == "cp-x2520" {
		clientIPPort := r.RemoteAddr
		clientIP := strings.Split(clientIPPort, ":")[0]
		log.Println(clientIP)
		err := iptables.NewRoute(clientIP, internalip)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("OK"))
	} else {
		w.Write([]byte("ERROR"))
	}
}

func PortalEntryHandler(w http.ResponseWriter, r *http.Request) {
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

	server := &http.Server{Addr: "0.0.0.0:8080", Handler: mux}

	err = server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}
