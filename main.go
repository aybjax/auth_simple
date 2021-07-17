package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))
var users = map[string]string{
	"aybjax": "aybjax",
	"admin": "password",
}

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session.id")

	if (session.Values["authenticated"] != nil) &&
		(session.Values["authenticated"]) != false {
			w.Write([]byte(time.Now().String()))
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request)  {
	session, _ := store.Get(r, "session.id")
	log.Println(session)
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Please pass the data as URL form encoded",
					http.StatusBadRequest)

		return
	}

	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")

	if originalPassword, ok := users[username]; ok {
		if password == originalPassword {
			session.Values["authenticated"] = true;
			session.Save(r, w)
		} else {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)

			return
		}
	} else {
		http.Error(w, "User not found", http.StatusNotFound)

		return
	}
}

func Logouthandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session.id")

	session.Values["authenticated"] = false
	session.Save(r, w)

	w.Write([]byte(""))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/healthcheck", HealthcheckHandler)
	r.HandleFunc("/logout", Logouthandler)

	http.Handle("/", r)
	srv := &http.Server {
		Handler: r,
		Addr: "127.0.0.1:8080",
	}

	log.Fatal(srv.ListenAndServe())
}

func init() {
	os.Setenv("SESSION_SECRET", "MY_SESSION_SECRET")
}