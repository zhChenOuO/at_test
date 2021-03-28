package igmail

import (
	"amazing_talker/internal/errors"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Config ...
type Config struct {
	Path     string `mapstructure:"path"`
	FileName string `mapstructure:"file_name"`
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"

	if dir := os.Getenv("PROJ_DIR"); dir != "" {
		tokFile = dir + "/" + tokFile
	}

	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok, err := getTokenFromWeb(config)
		if err != nil {
			return nil, err
		}
		if err := saveToken(tokFile, tok); err != nil {
			return nil, err
		}
	}
	return config.Client(context.Background(), tok), nil
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, errors.NewWithMessagef(errors.ErrInternalError, "Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		return nil, errors.NewWithMessagef(errors.ErrInternalError, "Unable to retrieve token from web: %v", err)
	}
	return tok, nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return errors.NewWithMessagef(errors.ErrInternalError, "Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
	return nil
}

// NewGmailService ...
func NewGmailService(cfg *Config) (*gmail.Service, error) {
	filename := "credentials.json"
	if cfg.FileName != "" {
		filename = cfg.FileName
	}
	if dir := os.Getenv("PROJ_DIR"); dir != "" {
		filename = dir + "/" + cfg.Path + "/" + filename
	}
	log.Debug().Msg(filename)
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.NewWithMessagef(errors.ErrInternalError, "Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		return nil, errors.NewWithMessagef(errors.ErrInternalError, "Unable to parse client secret file to config: %v", err)
	}
	client, err := getClient(config)
	if err != nil {
		return nil, err
	}

	srv, err := gmail.New(client)
	if err != nil {
		return nil, errors.NewWithMessagef(errors.ErrInternalError, "Unable to retrieve Gmail client: %v", err)
	}

	return srv, nil
}
