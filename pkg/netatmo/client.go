package netatmo

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const baseUrl = "https://api.netatmo.com/"
const reqAuth = baseUrl + "oauth2/token"
const reqStationData = baseUrl + "api/getstationsdata"

type Client struct {
	clientId     string
	clientSecret string
	Expires      time.Time
	Error        string `json:"error"`
	AuthResponse
}

type AuthResponse struct {
	AccessToken            string `json:"access_token"`
	RefreshToken           string `json:"refresh_token"`
	TokenDurationInSeconds int    `json:"expire_in"`
	Error                  string `json:"error"`
}

func NewClient(clientId, clientSecret string) (*Client, error) {
	logrus.Info("Establishing connection to netatmo API ...")
	refreshToken, err := ReadRefreshToken()
	if err != nil {
		return nil, err
	}
	d := url.Values{}
	d.Set("grant_type", "refresh_token")
	d.Set("client_id", clientId)
	d.Set("client_secret", clientSecret)
	d.Set("refresh_token", refreshToken)
	resp, err := postForm(reqAuth, d)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var c Client
	err = json.NewDecoder(resp.Body).Decode(&c)
	if err != nil {
		return nil, err
	}
	if c.Error != "" {
		return nil, errors.New(c.Error)
	}

	c.Expires = time.Now().Add(time.Second * time.Duration(c.TokenDurationInSeconds))
	c.clientId = clientId
	c.clientSecret = clientSecret
	logrus.Info("Connected!")

	err = WriteRefreshToken(c.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (c *Client) EnsureValidToken() error {
	if time.Now().After(c.Expires) {
		err := c.RefreshTokens()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) RefreshTokens() error {
	logrus.Info("Refreshing access token ...")
	refreshToken, err := ReadRefreshToken()
	if err != nil {
		return err
	}
	logrus.Info("Read refresh token from file: " + refreshToken)
	d := url.Values{}
	d.Set("grant_type", "refresh_token")
	d.Set("refresh_token", refreshToken)
	d.Set("client_id", c.clientId)
	d.Set("client_secret", c.clientSecret)
	resp, err := postForm(reqAuth, d)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var ar AuthResponse
	err = json.NewDecoder(resp.Body).Decode(&ar)
	if err != nil {
		return err
	}
	if ar.Error != "" {
		return errors.New(ar.Error)
	}

	c.AccessToken = ar.AccessToken
	err = WriteRefreshToken(ar.RefreshToken)
	if err != nil {
		return err
	}
	c.Expires = time.Now().Add(time.Second * time.Duration(ar.TokenDurationInSeconds))
	logrus.Info("Refresh OK")
	return nil
}

func (c *Client) GetStationData() (*DeviceListResponseBody, error) {
	logrus.Info("Getting station data ...")
	d := url.Values{}
	err := c.EnsureValidToken()
	if err != nil {
		return nil, err
	}
	d.Set("access_token", c.AccessToken)
	resp, err := postForm(reqStationData, d)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var dlr DeviceListResponse
	err = json.NewDecoder(resp.Body).Decode(&dlr)
	if err != nil {
		return nil, err
	}

	return &dlr.Body, nil
}

func postForm(url string, params url.Values) (*http.Response, error) {
	c := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")

	return c.Do(req)
}

func ReadRefreshToken() (string, error) {
	logrus.Info("Reading refresh token from file")
	bytes, err := os.ReadFile("/var/lib/netatmo/refresh_token.txt")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytes)), nil
}

func WriteRefreshToken(token string) error {
	if _, err := os.Stat("/var/lib/netatmo/refresh_token.txt"); os.IsNotExist(err) {
		err := os.MkdirAll("/var/lib/netatmo", 0600)
		if err != nil {
			return err
		}
	}
	logrus.Info("Writing refresh token to file")
	return os.WriteFile("/var/lib/netatmo/refresh_token.txt", []byte(token), 0644)
}
