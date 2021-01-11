package main

import (
	"encoding/json"
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
	accountPageURLBase    string = "https://graph.facebook.com/v9.0/me/accounts"
	businessPageURLBase   string = "https://graph.facebook.com/v9.0/me"
)

// GraphAuth stores instagram graph API credentials
type GraphAuth struct {
	RedirectURL          string
	FBAccessToken        string
	FBPageName           string
	FBPageID             string
	InstagramAccessToken string
	ClientID             string
	ClientSecret         string
}

// AccessTokenResponse response of second token call
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
}

// AccountResponse response of account call
type AccountResponse struct {
	Data []InstagramAccountAuth `json:"data"`
}

// InstagramAccountAuth response of fb account page
type InstagramAccountAuth struct {
	AccessToken string `json:"access_token"`
	Category    string `json:"category"`
	Name        string `json:"name"`
	ID          string `json:"id"`
}

// ExchangeTokenResponse response of second token call
type ExchangeTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
}

func main() {
	fmt.Println("gogram")
	var auth GraphAuth

	// initialize auth
	auth.RedirectURL = os.Getenv("GOGRAM_REDIRECT_URL")
	auth.ClientID = os.Getenv("GOGRAM_CLIENT_ID")
	auth.ClientSecret = os.Getenv("GOGRAM_CLIENT_SECRET")

	var loginDialogURL string = fmt.Sprintf("%s?client_id=%s&redirect_uri=%s", loginDialogURLBase, auth.ClientID, auth.RedirectURL)
	fmt.Println(loginDialogURL)

	router := mux.NewRouter()

	// callback URL
	router.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		// code を取得
		code := strings.Join(q["code"], "")

		var getAccessTokenURL string = fmt.Sprintf(
			"%s?client_id=%s&redirect_uri=%s&client_secret=%s&code=%s",
			getAccessTokenURLBase, auth.ClientID, auth.RedirectURL, auth.ClientSecret, code)

		res, err := http.Get(getAccessTokenURL)
		if err != nil {
			panic(err)
		}

		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		// unmarshal access_token
		var accessTokenResponse AccessTokenResponse
		err = json.Unmarshal(b, &accessTokenResponse)

		// save access_token
		auth.FBAccessToken = accessTokenResponse.AccessToken
		bodyStr := string(b)
		fmt.Fprintf(w, bodyStr)
	})

	// user
	router.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		if auth.FBAccessToken == "" {
			fmt.Fprintf(w, "please call /login first")
			return
		}

		var accountPageURL string = fmt.Sprintf("%s?access_token=%s", accountPageURLBase, auth.FBAccessToken)
		res, err := http.Get(accountPageURL)
		if err != nil {
			panic(err)
		}
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		var accountResponse AccountResponse
		err = json.Unmarshal(b, &accountResponse)

		auth.InstagramAccessToken = accountResponse.Data[0].AccessToken
		auth.FBPageName = accountResponse.Data[0].Name
		auth.FBPageID = accountResponse.Data[0].ID
		bodyStr := string(b)
		fmt.Fprintf(w, bodyStr)
	})

	// businness
	router.HandleFunc("/business", func(w http.ResponseWriter, r *http.Request) {
		if auth.InstagramAccessToken == "" {
			fmt.Fprintf(w, "please call /user first")
			return
		}

		var businessPageURL string = fmt.Sprintf("%s/?fields=instagram_business_account&access_token=%s", businessPageURLBase, auth.InstagramAccessToken)
		res, err := http.Get(businessPageURL)
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		bodyStr := string(body)
		fmt.Fprintf(w, bodyStr)
	})

	// login page
	router.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, loginDialogURL, http.StatusFound)
	})

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
	}

	log.Fatal(srv.ListenAndServe())
}
