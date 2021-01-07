package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

var (
	loginDialogURLBase string = "https://www.facebook.com/v9.0/dialog/oauth"
)

func main() {
	fmt.Println("gogram")

	var redirectURL = os.Getenv("GOGRAM_REDIRECT_URL")
	var clientID = os.Getenv("GOGRAM_CLIENT_ID")
	fmt.Println(redirectURL, clientID)

	var loginDialogURL string = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s", loginDialogURLBase, clientID, redirectURL)
	fmt.Println(loginDialogURL)

	router := mux.NewRouter()

	// callback URL
	router.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		fmt.Fprintf(w, "state: %s\n", strings.Join(queries["state"], ""))
		fmt.Fprintf(w, "code: %s\n", strings.Join(queries["code"], ""))
	})

	// login page
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("login")
		http.Redirect(w, r, loginDialogURL, http.StatusFound)
	})
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(srv.ListenAndServe())
}
