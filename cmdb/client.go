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

	"github.com/Sirupsen/logrus"
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
	url := c.cmdbEndpoint
	url.Path = "/v2/items/system/" + systemCode
	req := http.Request{
		Method: "GET",
		URL:    url,
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

func (c *Client) UpdateRelationship(relationshipType string, contact Contact, system System) error {
	url := c.cmdbEndpoint
	url.Path = "/v2/relationships/system/" + system.SystemCode + "/" + relationshipType
	resp, err := c.makeRequest("GET", url, nil)
	logrus.Infof("Get relationship [%s]: %d", url, resp.StatusCode)
	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Error(err)
			return err
		}

		var relList []Relationship
		err = json.Unmarshal(body, &relList)
		if err != nil {
			logrus.Error(err)
			return err
		}

		for _, r := range relList {
			url := c.cmdbEndpoint
			url.Path = "/v2/relationships/" + r.SubjectType + "/" + r.SubjectID + "/" + r.RelationshipType + "/" + r.ObjectType + "/" + r.ObjectID
			resp, err := c.makeRequest("DELETE", url, nil)
			logrus.Infof("Delete relationship [%s]: %d", url, resp.StatusCode)
			if err != nil {
				logrus.Error(err)
				return err
			}
		}
	}
	if err != nil {
		logrus.Error(err)
		return err
	}

	for _, e := range contact.Entries {
		logrus.Infof("contact: %v", e)
		r := Relationship{
			SubjectType:      "system",
			SubjectID:        system.SystemCode,
			RelationshipType: relationshipType,
			ObjectType:       "contact",
			ObjectID:         e.DataItemID,
		}

		url := c.cmdbEndpoint
		url.Path = "/v2/relationships/" + r.SubjectType + "/" + r.SubjectID + "/" + r.RelationshipType + "/" + r.ObjectType + "/" + r.ObjectID
		resp, err := c.makeRequest("PUT", url, nil)
		logrus.Infof("Create relationship [%s]: %d", url, resp.StatusCode)
		if err != nil {
			logrus.Errorf("Error adding relationship: %s", err)
			return err
		}
	}

	return err
}

func (c *Client) makeRequest(method string, url *url.URL, body []byte) (*http.Response, error) {
	req := http.Request{
		Method: method,
		URL:    url,
		Header: http.Header{
			"X-API-KEY": {c.apiKey},
		},
	}

	if len(body) > 0 {
		req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		req.ContentLength = int64(len(body))
		req.Header.Add("Content-Type", "application/json")
	}
	return c.httpClient.Do(&req)
}

func (c *Client) UpdateSystemAttributes(system SystemAttributes) error {
	if system.SystemCode == "" {
		return errors.New("No system code set")
	}

	body, err := json.Marshal(system)
	if err != nil {
		return err
	}

	url := c.cmdbEndpoint
	url.Path = "/v2/items/system/" + system.SystemCode

	resp, err := c.makeRequest("PUT", url, body)
	if resp.StatusCode != http.StatusOK {
		logrus.Error(resp.StatusCode, err)
		return err
	}
	return err
}

func (c *Client) UpdateSystem(system System) error {

	err := c.UpdateRelationship("primaryContact", system.PrimaryContact, system)
	if err != nil {
		logrus.Errorf("Error updating primaryContact, %s", err)
		return err
	}
	err = c.UpdateRelationship("secondaryContact", system.SecondaryContact, system)
	if err != nil {
		logrus.Errorf("Error updating secondaryContact, %s", err)
		return err
	}
	err = c.UpdateRelationship("programme", system.Programme, system)
	if err != nil {
		logrus.Errorf("Error updating programme, %s", err)
		return err
	}
	err = c.UpdateRelationship("productOwner", system.ProductOwner, system)
	if err != nil {
		logrus.Errorf("Error updating productOwner, %s", err)
		return err
	}
	err = c.UpdateRelationship("technicalLead", system.TechnicalLead, system)
	if err != nil {
		logrus.Errorf("Error updating technicalLead, %s", err)
		return err
	}

	err = c.UpdateSystemAttributes(system.ToSystemAttributes())
	if err != nil {
		logrus.Error(err)
		return err
	}
	return err

}
