package prosper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
)

type ProsperClient struct {
	ClientId     string
	ClientSecret string
	Filter       string
	BaseUrl      string
}

type ProsperToken struct {
	Token     string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
}

func NewProsperClient() (ProsperClient, error) {
	p := ProsperClient{}
	p.ClientId = viper.GetString("client_id")
	p.ClientSecret = viper.GetString("client_secret")
	p.Filter = viper.GetString("filter")
	p.BaseUrl = "https://api.prosper.com/v1"
	return p, nil
}

func (p *ProsperClient) GetToken() ProsperToken {
	payload := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s", p.ClientId, p.ClientSecret)
	url := fmt.Sprintf("%s%s", p.BaseUrl, "/security/oauth/token")
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(payload))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var token ProsperToken
	if err := json.Unmarshal(responseData, &token); err != nil {
		log.Fatal(err)
	}
	return token
}

func (p *ProsperClient) GetListings(filter string, token string) string {
	if filter == "" {
		filter = "sort_by=effective_yield&amount_remaining_max=1000"
	}

	url := fmt.Sprintf("%s%s?%s", p.BaseUrl, "/search/listings", filter)
	req, err := http.NewRequest("GET", url, bytes.NewBufferString(filter))

	// Setting Headers for authentication
	authHeader := fmt.Sprintf("bearer %s", token)
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return string(responseData)
}
