package utils

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	GET_METHOD  = "GET"
	POST_METHOD = "POST"
)

type LoraHTTPClient struct {
	loraserver      string
	loraserver_user string
	loraserver_pwd  string
	token           string
	expired         int64
}

func NewLoraHTTPClient(url, user, password string) *LoraHTTPClient {
	return &LoraHTTPClient{
		loraserver:      url,
		loraserver_user: user,
		loraserver_pwd:  password,
	}
}

func (c *LoraHTTPClient) Get(api_path string, params, header map[string]string, body map[string]interface{}) ([]byte, error) {
	return c.send(GET_METHOD, api_path, params, header, body)
}

func (c *LoraHTTPClient) Post(api_path string, params, header map[string]string, body map[string]interface{}) ([]byte, error) {
	return c.send(POST_METHOD, api_path, params, header, body)
}

func (c *LoraHTTPClient) send(method, api_path string, params, header map[string]string, body map[string]interface{}) ([]byte, error) {
	uri := fmt.Sprintf("%s%s", c.loraserver, api_path)
	// add jwt token
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}
	if header == nil {
		header = make(map[string]string)
	}
	header["Grpc-Metadata-Authorization"] = fmt.Sprintf("Bearer %s", token)
	header["Content-Type"] = "application/json; charset=utf-8"
	return sendhttpRequest(method, uri, params, header, body)
}

func sendhttpRequest(method, uri string, params, header map[string]string, body map[string]interface{}) ([]byte, error) {
	var (
		req       *http.Request
		resp      *http.Response
		client    http.Client
		send_data string
		err       error
	)

	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		uri = fmt.Sprintf("%s?%s", uri, q.Encode())
	}

	if len(body) > 0 {
		send_body, json_err := json.Marshal(body)
		if json_err != nil {
			return nil, json_err
		}
		send_data = string(send_body)
	}

	client.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	req, err = http.NewRequest(method, uri, strings.NewReader(send_data))
	if err != nil {
		return nil, err
	}
	defer func() {
		if req.Body != nil {
			req.Body.Close()
		}
	}()

	for k, v := range header {
		req.Header.Add(k, v)
	}
	// log.Debug("http client header:", req.Header)

	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("error http code:%d ", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

func (c *LoraHTTPClient) getToken() (string, error) {
	now := time.Now().Add(time.Minute * -5).Unix()
	if now > c.expired {
		uri := fmt.Sprintf("%s/api/internal/login", c.loraserver)
		header := map[string]string{
			"Content-Type": "application/json; charset=utf-8",
		}
		resp_data, err := sendhttpRequest(POST_METHOD, uri, nil, header, map[string]interface{}{
			"username": c.loraserver_user,
			"password": c.loraserver_pwd,
		})
		if err != nil {
			return "", err
		}
		token := struct {
			JWT string `json:"jwt"`
		}{}

		if err := json.Unmarshal(resp_data, &token); err != nil {
			return "", err
		}
		c.token = token.JWT
		jwts := strings.Split(token.JWT, ".")
		if len(jwts) == 3 {
			raw_claims, err := base64.RawURLEncoding.DecodeString(jwts[1])
			if err == nil {
				claims := make(map[string]interface{})
				if err := json.Unmarshal(raw_claims, &claims); err == nil {
					if exp, ok := claims["exp"].(float64); ok {
						c.expired = int64(exp)
					}
				}

			}
		}
	}
	return c.token, nil
}
