package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleProfileBody struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
}

func OAuthIndex(r chi.Router) {

	// r.Handle("/google/callback", http.HandlerFunc(GoogleCallback))
	r.Handle("/google/callback", http.HandlerFunc(GoogleCallBack))
	// r.Handle("/google", http.HandlerFunc(GoogleRequest))
	r.Get("/google", http.HandlerFunc(Google))

}

func Google(w http.ResponseWriter, r *http.Request) {

	godotenv.Load()

	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			"openid",
			"email",
			"profile",
		},
		Endpoint: google.Endpoint,
	}

	url := conf.AuthCodeURL("state")
	http.Redirect(w, r, url, http.StatusFound)
}

func GoogleCallBack(w http.ResponseWriter, r *http.Request) {

	fmt.Println(r.Method)

	fmt.Println("Callback Entered")

	code := r.URL.Query().Get("code")
	conf := &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes: []string{
			"openid",
			"email",
			"profile",
		},
		Endpoint: google.Endpoint,
	}

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		fmt.Println("This is the error")
		w.WriteHeader(500)
		log.Fatalln(err)
	}

	client := conf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?alt=json")
	if err != nil {
		w.WriteHeader(500)
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(500)
		log.Fatalln(err)
		return
	}

	authedPerson := &GoogleProfileBody{}

	if err = json.Unmarshal(bt, &authedPerson); err != nil {
		w.WriteHeader(500)
		log.Fatalln(err)
		return
	}

	if err = ExistsInDatabase(authedPerson.Email); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Failed to authorize"))
		return
	}

	w.Write(bt)
}

//TODO: Check to see if Email exists in database
func ExistsInDatabase(email string) error { return nil }
