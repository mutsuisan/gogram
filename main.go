package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

var (
	loginDialogURLBase    string = "https://www.facebook.com/v9.0/dialog/oauth"
	getAccessTokenURLBase string = "https://graph.facebook.com/v9.0/oauth/access_token"
)

func main() {
	fmt.Println("gogram")

	var redirectURL = os.Getenv("GOGRAM_REDIRECT_URL")
	var clientID = os.Getenv("GOGRAM_CLIENT_ID")
	var clientSecret = os.Getenv("GOGRAM_CLIENT_SECRET")
	fmt.Println(redirectURL, clientID)

	var loginDialogURL string = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s", loginDialogURLBase, clientID, redirectURL)
	fmt.Println(loginDialogURL)

	router := mux.NewRouter()

	// callback URL
	router.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		queries := r.URL.Query()
		fmt.Fprintf(w, "state: %s\n", strings.Join(queries["state"], ""))
		code := strings.Join(queries["code"], "")
		fmt.Fprintf(w, "code: %s\n", code)
		fmt.Fprintln(w, "-----------------------------")

		fmt.Println("fetching access token...")
		var getAccessTokenURL string = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s", getAccessTokenURLBase, clientID, redirectURL, clientSecret, code)
		res, err := http.Get(getAccessTokenURL)
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		bodyStr := string(body)
		fmt.Fprintf(w, bodyStr)
		fmt.Fprintln(w, "\n-----------------------------")
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
