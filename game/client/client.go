package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/joohnes/wp-sea/game/app"
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
	urlOppDesc   = "/api/game/desc"
	urlRefresh   = "/api/game/refresh"
	urlList      = "/api/list"
	urlStats     = "/api/stats"
	tokenHeader  = "X-Auth-Token"
	errAuthToken = "no auth token"
)

type Stats struct {
	Games  int    `json:"games"`
	Nick   string `json:"nick"`
	Points int    `json:"points"`
	Rank   int    `json:"rank"`
	Wins   int    `json:"wins"`
}

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
		"coords":      coords,
		"desc":        desc,
		"nick":        nick,
		"target_nick": targetOpponent,
		"wpbot":       wpbot,
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
	if resp.StatusCode != 200 {
		return errors.New("could not init game")
	}

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
	if resp.StatusCode != 200 {
		return []string{}, errors.New("could not retrieve board")
	}
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
	if resp.StatusCode != 200 {
		return nil, errors.New("could not retrieve status")
	}
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
	if resp.StatusCode != 200 {
		return "", errors.New("could not shoot")
	}
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
	if resp.StatusCode != 200 {
		return errors.New("could not resign")
	}
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

	if resp.StatusCode != 200 {
		return "", "", errors.New("could not retrieve opp description")
	}
	var body map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return "", "", err
	}
	opp_nick := fmt.Sprintf("%v", body["opponent"])
	opp_desc := fmt.Sprintf("%v", body["opp_desc"])
	return opp_nick, opp_desc, nil
}

func (c *Client) Refresh() error {
	urlPath, err := url.JoinPath(c.serverAddr, urlRefresh)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodGet, urlPath, http.NoBody)
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
	if resp.StatusCode != 200 {
		return errors.New("could not refresh")
	}
	return nil
}

func (c *Client) PlayerList() ([]map[string]string, error) {
	urlPath, err := url.JoinPath(c.serverAddr, urlList)
	if err != nil {
		return []map[string]string{}, err
	}

	req, err := http.NewRequest(http.MethodGet, urlPath, http.NoBody)
	if err != nil {
		return []map[string]string{}, err
	}

	req.Header = http.Header{
		tokenHeader: []string{c.token},
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return []map[string]string{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return []map[string]string{}, errors.New("could not get player list")
	}
	var body []map[string]string
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return []map[string]string{}, err
	}
	return body, nil
}

func (c *Client) Stats() (map[string][]int, error) {
	urlPath, err := url.JoinPath(c.serverAddr, urlStats)
	if err != nil {
		return map[string][]int{}, err
	}

	req, err := http.NewRequest(http.MethodGet, urlPath, http.NoBody)
	if err != nil {
		return map[string][]int{}, err
	}

	req.Header = http.Header{
		tokenHeader: []string{c.token},
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return map[string][]int{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return map[string][]int{}, errors.New("player not found")
	}
	var body map[string][]Stats
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return map[string][]int{}, err
	}
	stats := make(map[string][]int)
	for _, x := range body["stats"] {
		stats[x.Nick] = []int{x.Games, x.Points, x.Rank, x.Wins}
	}

	return stats, nil
}

func (c *Client) StatsPlayer(nick string) ([]int, error) {
	urlPath, err := url.JoinPath(c.serverAddr, urlStats)
	if err != nil {
		return []int{}, err
	}
	if nick == "" {
		fmt.Println("Please enter a nick!")
		return []int{}, err
	}
	urlPath, err = url.JoinPath(urlPath, nick)
	if err != nil {
		return []int{}, err
	}
	req, err := http.NewRequest(http.MethodGet, urlPath, http.NoBody)
	if err != nil {
		return []int{}, err
	}

	req.Header = http.Header{
		tokenHeader: []string{c.token},
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return []int{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return []int{}, errors.New("player not found")
	}
	var body map[string]Stats
	err = json.NewDecoder(resp.Body).Decode(&body)
	if err != nil {
		return []int{}, err
	}
	stats := []int{
		body["stats"].Games,
		body["stats"].Points,
		body["stats"].Rank,
		body["stats"].Wins,
	}

	return stats, nil
}
