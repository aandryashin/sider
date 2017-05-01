package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Client struct {
	Endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{endpoint}
}

func (c *Client) Keys(ctx context.Context) ([]string, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/keys", c.Endpoint), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
	}
	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("keys: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("keys: %v", err)
		}
		return nil, fmt.Errorf("keys: %s", string(msg))
	}
	var keys []string
	err = json.NewDecoder(resp.Body).Decode(&keys)
	if err != nil {
		return nil, fmt.Errorf("decode response: %v", err)
	}
	return keys, nil
}

func (c *Client) Get(ctx context.Context, key string) (interface{}, error) {
	r, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/keys/%s", c.Endpoint, key), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %v", err)
	}
	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("get: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("keys: %v", err)
		}
		return nil, fmt.Errorf("get: %s", string(msg))
	}
	var value interface{}
	err = json.NewDecoder(resp.Body).Decode(&value)
	if err != nil {
		return nil, fmt.Errorf("decode response: %v", err)
	}
	return value, nil
}

func (c *Client) Set(ctx context.Context, key string, body io.Reader, ttl time.Duration) error {
	var u string
	if ttl != 0 {
		u = fmt.Sprintf("%s/keys/%s?ttl=%v", c.Endpoint, key, ttl)
	} else {
		u = fmt.Sprintf("%s/keys/%s", c.Endpoint, key)
	}
	r, err := http.NewRequest(http.MethodPost, u, body)
	if err != nil {
		return fmt.Errorf("new request: %v", err)
	}
	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("set: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("set: %v", err)
		}
		return fmt.Errorf("set: %s", string(msg))
	}
	return nil
}

func (c *Client) Del(ctx context.Context, key string) error {
	r, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/keys/%s", c.Endpoint, key), nil)
	if err != nil {
		return fmt.Errorf("new request: %v", err)
	}
	resp, err := http.DefaultClient.Do(r.WithContext(ctx))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("del: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		msg, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("del: %v", err)
		}
		return fmt.Errorf("del: %s", string(msg))
	}
	return nil
}
