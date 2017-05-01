package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestKeys(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/keys", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		json.NewEncoder(w).Encode([]string{})
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(server.URL)

	v, err := client.Keys(context.Background())
	if err != nil {
		t.Errorf("query keys: %v", err)
	}
	if !reflect.DeepEqual(v, []string{}) {
		t.Errorf("unexpected response: %v", v)
	}
}

func TestGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/keys/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		json.NewEncoder(w).Encode(struct{}{})
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(server.URL)

	v, err := client.Get(context.Background(), "0")
	if err != nil {
		t.Errorf("query keys: %v", err)
	}
	if !reflect.DeepEqual(v, map[string]interface{}{}) {
		t.Errorf("unexpected response: %v", v)
	}
}

func TestSet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/keys/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("unexpected method: %s", r.Method)
		}

		var v interface{}
		json.NewDecoder(r.Body).Decode(&v)

		if !reflect.DeepEqual(v, map[string]interface{}{}) {
			t.Errorf("unexpected response: %v", v)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(server.URL)

	err := client.Set(context.Background(), "0", bytes.NewReader([]byte("{}")), 0)
	if err != nil {
		t.Errorf("query keys: %v", err)
	}
}

func TestSetTTL(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/keys/", func(w http.ResponseWriter, r *http.Request) {
		ttlStr := r.FormValue("ttl")
		ttl, err := time.ParseDuration(ttlStr)
		if err != nil {
			t.Fatalf("parse duration: %v", err)
		}
		if ttl != 10*time.Second {
			t.Fatalf("unexpected duration: %v", ttl)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(server.URL)

	err := client.Set(context.Background(), "0", nil, 10*time.Second)
	if err != nil {
		t.Errorf("query keys: %v", err)
	}
}

func TestDel(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/keys/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("unexpected method: %s", r.Method)
		}
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(server.URL)

	err := client.Del(context.Background(), "0")
	if err != nil {
		t.Errorf("query keys: %v", err)
	}
}

func TestMalformedURL(t *testing.T) {
	client := NewClient("0://")

	var err error
	_, err = client.Keys(context.Background())
	if err == nil {
		t.Errorf("unexpected pass")
	}

	_, err = client.Get(context.Background(), "0")
	if err == nil {
		t.Errorf("unexpected pass")
	}

	err = client.Set(context.Background(), "0", bytes.NewReader([]byte("{}")), 0)
	if err == nil {
		t.Errorf("unexpected pass")
	}

	err = client.Del(context.Background(), "0")
	if err == nil {
		t.Errorf("unexpected pass")
	}
}

func TestServerDown(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	client := NewClient(server.URL)
	server.Close()

	var err error
	_, err = client.Keys(context.Background())
	if err == nil {
		t.Errorf("unexpected pass")
	}

	_, err = client.Get(context.Background(), "0")
	if err == nil {
		t.Errorf("unexpected pass")
	}

	err = client.Set(context.Background(), "0", bytes.NewReader([]byte("{}")), 0)
	if err == nil {
		t.Errorf("unexpected pass")
	}

	err = client.Del(context.Background(), "0")
	if err == nil {
		t.Errorf("unexpected pass")
	}
}

func TestServerBadResponse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(server.URL)

	var err error
	_, err = client.Keys(context.Background())
	if err == nil {
		t.Errorf("unexpected pass")
	}

	_, err = client.Get(context.Background(), "0")
	if err == nil {
		t.Errorf("unexpected pass")
	}

	err = client.Set(context.Background(), "0", bytes.NewReader([]byte("{}")), 0)
	if err == nil {
		t.Errorf("unexpected pass")
	}

	err = client.Del(context.Background(), "0")
	if err == nil {
		t.Errorf("unexpected pass")
	}
}

func TestServerBadJSONResponse(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{")
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient(server.URL)

	var err error
	_, err = client.Keys(context.Background())
	if err == nil {
		t.Errorf("unexpected pass")
	}

	_, err = client.Get(context.Background(), "0")
	if err == nil {
		t.Errorf("unexpected pass")
	}
}
