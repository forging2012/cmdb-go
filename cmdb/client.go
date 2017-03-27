package cmdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	httpClient   *http.Client
	cmdbEndpoint *url.URL
	apiKey       string
}

func NewClient(cmdbEndpoint string, apiKey string) (Client, error) {

	tr := &http.Transport{
		MaxIdleConnsPerHost: 32,
		Dial: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
	}
	c := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	reqUrl, err := url.Parse(cmdbEndpoint)
	if err != nil {
		return Client{}, err
	}

	return Client{
		httpClient:   c,
		cmdbEndpoint: reqUrl,
		apiKey:       apiKey,
	}, nil
}

func (c *Client) GetSystem(systemCode string) (System, error) {
	c.cmdbEndpoint.Path = "/v2/items/system/" + systemCode
	req := http.Request{
		Method: "GET",
		URL:    c.cmdbEndpoint,
		Header: http.Header{
			"X-API-KEY": {c.apiKey},
		},
	}
	resp, err := c.httpClient.Do(&req)
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	if err != nil {
		return System{}, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return System{}, err
	}

	var system System
	err = json.Unmarshal(body, &system)
	if err != nil {
		return System{}, err
	}
	return system, err
}

func (c *Client) UpdateSystemAttributes(system SystemAttributes) error {
	if system.SystemCode == "" {
		return errors.New("No system code set")
	}

	body, err := json.Marshal(system)
	if err != nil {
		return err
	}

	c.cmdbEndpoint.Path = "/v2/items/system/" + system.SystemCode
	req := http.Request{
		Method: "PUT",
		URL:    c.cmdbEndpoint,
		Header: http.Header{
			"X-API-KEY":    {c.apiKey},
			"Content-Type": {"application/json"},
		},
		Body:          ioutil.NopCloser(bytes.NewBuffer(body)),
		ContentLength: int64(len(body)),
	}
	resp, err := c.httpClient.Do(&req)
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	return err
}

func (c *Client) UpdateSystem(system System) error {

	return c.UpdateSystemAttributes(system.ToSystemAttributes())

}
