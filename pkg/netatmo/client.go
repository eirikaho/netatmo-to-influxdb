package netatmo

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const baseUrl = "https://api.netatmo.com/"
const reqAuth = baseUrl + "oauth2/token"
const reqStationData = baseUrl + "api/getstationsdata"

type Client struct {
	AccessToken            string `json:"access_token"`
	RefreshToken           string `json:"refresh_token"`
	TokenDurationInSeconds int    `json:"expire_in"`
	clientId               string
	clientSecret           string
	Expires                time.Time
	Error                  string `json:"error"`
}

func NewClient(clientId, clientSecret, refreshToken string) (*Client, error) {
	logrus.Info("Establishing connection to netatmo API ...")
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
	d := url.Values{}
	d.Set("grant_type", "refresh_token")
	d.Set("refresh_token", c.RefreshToken)
	d.Set("client_id", c.clientId)
	d.Set("client_secret", c.clientSecret)
	resp, err := postForm(reqAuth, d)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var cn Client
	err = json.NewDecoder(resp.Body).Decode(&cn)
	if err != nil {
		return err
	}
	if cn.Error != "" {
		return errors.New(cn.Error)
	}

	c.AccessToken = cn.AccessToken
	c.RefreshToken = cn.RefreshToken
	c.Expires = time.Now().Add(time.Second * time.Duration(cn.TokenDurationInSeconds))
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
