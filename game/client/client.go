package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	urlInitGame  = "/api/game"
	urlGetBoard  = "/api/game/board"
	urlGetStatus = "/api/game"
	tokenHeader  = "x-auth-token"
	errAuthToken = "no auth token"
)

type Board struct {
	Board []string
}

type client struct {
	client     *http.Client
	serverAddr string
	token      string
}

func New(addr string, t time.Duration) *client {
	return &client{
		client: &http.Client{
			Timeout: t,
		},
		serverAddr: addr,
	}
}

func (c *client) PrintToken() {
	fmt.Println("Token: ", c.token)
}

func (c *client) InitGame(coords []string, desc, nick, target_opponent string, wpbot bool) error {
	params := map[string]any{
		"coords":          coords,
		"desc":            desc,
		"nick":            nick,
		"target_opponent": target_opponent,
		"wpbot":           wpbot,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	url, err := url.JoinPath(c.serverAddr, urlInitGame)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.token = string(resp.Header.Get(tokenHeader))
	if c.token != string(resp.Header.Get(tokenHeader)) {
		return errors.New(errAuthToken)
	}
	return nil
}

func (c *client) Board() ([]string, error) {
	url, err := url.JoinPath(c.serverAddr, urlGetBoard)
	if err != nil {
		return []string{}, err
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []string{}, err
	}
	req.Header = http.Header{
		tokenHeader: []string{c.token},
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return []string{}, err
	}
	defer resp.Body.Close()

	body := Board{}
	json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return []string{}, err
	}

	fmt.Print(body)
	return []string{}, err
}

// func (c *client) Status() (*StatusResponse, error) {
// 	url, err := url.JoinPath(c.serverAddr, urlGetStatus)
// 	if err != nil {
// 		return err
// 	}
// 	req, err := http.NewRequest(http.MethodGet, url, nil)
// 	if err != nil {
// 		return err
// 	}
// 	req.Header = http.Header{
// 		tokenHeader: []string{c.token},
// 	}
// 	resp, err := c.client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Print(string(body))
// 	return nil
// }
