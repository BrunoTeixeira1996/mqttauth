package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

type Hooker struct {
	mqttUsername string
	mqttPassword string
}

type AuthResponse struct {
	Ok bool `json:"ok"`
}

func authHandler(hooker *Hooker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Auth request received: %+v", r)

		if err := r.ParseForm(); err != nil {
			log.Printf("Error parsing form: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		if username == "" || password == "" {
			return
		}

		log.Printf("Got username=%s and password=%s", username, password)

		ok := username == hooker.mqttUsername && password == hooker.mqttPassword
		log.Printf("Auth result: %v", ok)
		if !ok {
			log.Print("Sending status unauthorize")
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

func aclHandler(hooker *Hooker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		topic := r.URL.Query().Get("topic")
		access := r.URL.Query().Get("access") // 1=subscribe, 2=publish

		log.Printf("ACL check: user=%s, topic=%s, access=%s", username, topic, access)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	}
}

func main() {
	mqttUsername := flag.String("mqtt_username", "", "client username")
	mqttPassword := flag.String("mqtt_password", "", "client password")
	listenPort := flag.String("listen", ":9393", "listen port")
	flag.Parse()

	if *mqttUsername == "" || *mqttPassword == "" {
		log.Fatalf("Please provide a username and a password")
	}

	hooker := &Hooker{mqttUsername: *mqttUsername, mqttPassword: *mqttPassword}

	http.HandleFunc("/mqtt/auth", authHandler(hooker))
	http.HandleFunc("/mqtt/acl", aclHandler(hooker))
	log.Printf("Auth server listening on %s", *listenPort)
	log.Fatal(http.ListenAndServe(*listenPort, nil))
}
