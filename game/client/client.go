package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"game/game/app"
	"net/http"
	"net/url"
	"time"
)

const (
	urlInitGame  = "/api/game"
	urlGetBoard  = "/api/game/board"
	urlGetStatus = "/api/game"
	urlFire      = "/api/game/fire"
	urlResign    = "/api/game/abandon"
	urlOppDesc   = "api/game/desc"
	tokenHeader  = "X-Auth-Token"
	errAuthToken = "no auth token"
)

type Board struct {
	Board []string
}

type Client struct {
	client     *http.Client
	serverAddr string
	token      string
}

func New(addr string, t time.Duration) *Client {
	return &Client{
		client: &http.Client{
			Timeout: t,
		},
		serverAddr: addr,
	}
}

func (c *Client) PrintToken() {
	fmt.Println("Token: ", c.token)
}

func (c *Client) InitGame(coords []string, desc, nick, targetOpponent string, wpbot bool) error {
	fmt.Println("Connecting to server...")
	params := map[string]any{
		"coords":          coords,
		"desc":            desc,
		"nick":            nick,
		"target_opponent": targetOpponent,
		"wpbot":           wpbot,
	}

	data, err := json.Marshal(params)
	if err != nil {
		return err
	}
	urlPath, err := url.JoinPath(c.serverAddr, urlInitGame)
	if err != nil {
		return err
	}

	resp, err := http.Post(urlPath, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c.token = resp.Header.Get(tokenHeader)
	if c.token != resp.Header.Get(tokenHeader) {
		return errors.New(errAuthToken)
	}
	return nil
}

func (c *Client) Board() ([]string, error) {
	urlPath, err := url.JoinPath(c.serverAddr, urlGetBoard)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, urlPath, nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		tokenHeader: []string{c.token},
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body := Board{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return body.Board, err
}

func (c *Client) Status() (*app.StatusResponse, error) {
	urlPath, err := url.JoinPath(c.serverAddr, urlGetStatus)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, urlPath, http.NoBody)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		tokenHeader: []string{c.token},
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body := app.StatusResponse{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	return &body, nil
}

func (c *Client) Shoot(coord string) (string, error) {
	urlPath, err := url.JoinPath(c.serverAddr, urlFire)
	if err != nil {
		return "", err
	}
	params := map[string]string{
		"coord": coord,
	}
	data, err := json.Marshal(params)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(http.MethodPost, urlPath, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	req.Header = http.Header{
		tokenHeader: []string{c.token},
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var body map[string]string

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", err
	}

	return body["result"], nil
}

func (c *Client) Resign() error {
	urlPath, err := url.JoinPath(c.serverAddr, urlResign)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, urlPath, http.NoBody)
	if err != nil {
		return err
	}

	req.Header = http.Header{
		tokenHeader: []string{c.token},
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Println("You have successfully resigned the game!")
	return nil
}

func (c *Client) GetOppDesc() (string, string, error) {
	urlPath, err := url.JoinPath(c.serverAddr, urlOppDesc)
	if err != nil {
		return "", "", err
	}

	req, err := http.NewRequest(http.MethodGet, urlPath, http.NoBody)
	if err != nil {
		return "", "", err
	}

	req.Header = http.Header{
		tokenHeader: []string{c.token},
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	var body map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", "", err
	}
	opp_nick := fmt.Sprintf("%v", body["opponent"])
	opp_desc := fmt.Sprintf("%v", body["opp_desc"])
	return opp_nick, opp_desc, nil
}
