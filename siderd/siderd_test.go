package main

import (
	"bytes"
	"context"
	"github.com/aandryashin/sider/siderd/client"
	"github.com/pborman/uuid"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	key := uuid.New()
	err := cl.Set(context.Background(), key, bytes.NewReader([]byte("{}")), 0)
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	lock.RLock()
	_, ok := storage[key]
	lock.RUnlock()
	if !ok {
		t.Fatalf("key not found")
	}

	lock.Lock()
	defer lock.Unlock()
	delete(storage, key)
}

func TestSetDuplicate(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	key := uuid.New()
	err := cl.Set(context.Background(), key, bytes.NewReader([]byte("{}")), 0)
	if err != nil {
		t.Fatalf("client: %v", err)
	}
	
	err = cl.Set(context.Background(), key, bytes.NewReader([]byte("{}")), 0)
	if err == nil {
		t.Fatalf("possible duplicate key key")
	}
	
	lock.RLock()
	_, ok := storage[key]
	lock.RUnlock()
	if !ok {
		t.Fatalf("key not found")
	}

	lock.Lock()
	defer lock.Unlock()
	delete(storage, key)
}


func TestSetExpiration(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	key := uuid.New()
	err := cl.Set(context.Background(), key, bytes.NewReader([]byte("{}")), 50*time.Millisecond)
	if err != nil {
		t.Fatalf("set key: %v", err)
	}

	<-time.After(100 * time.Millisecond)

	lock.RLock()
	_, ok := storage[key]
	lock.RUnlock()
	if ok {
		t.Fatalf("key is not expired")
	}
}

func TestSetEmptyKey(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	err := cl.Set(context.Background(), "", bytes.NewReader([]byte("{}")), 0)
	if err == nil {
		t.Fatalf("set with empty key")
	}
}

func TestSetNegativeDuration(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	key := uuid.New()
	err := cl.Set(context.Background(), key, bytes.NewReader([]byte("{}")), -50*time.Millisecond)
	if err == nil {
		t.Fatalf("set with negative duration")
	}
}

func TestSetBadJSON(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	key := uuid.New()
	err := cl.Set(context.Background(), key, bytes.NewReader([]byte("{")), 0)
	if err == nil {
		t.Fatalf("set with malformed json")
	}
}

func TestEmptyKeys(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	keys, err := cl.Keys(context.Background())
	if err != nil {
		t.Fatalf("get keys: %v", err)
	}
	if len(keys) != 0 {
		t.Fatalf("non empty list: %v", keys)
	}
}

func TestKeys(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	key := uuid.New()
	err := cl.Set(context.Background(), key, bytes.NewReader([]byte("{}")), 0)
	if err != nil {
		t.Fatalf("set: %v", err)
	}

	keys, err := cl.Keys(context.Background())
	if err != nil {
		t.Fatalf("keys: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("unexpected slice length")
	}
	if keys[0] != key {
		t.Fatalf("key does not exist")
	}
}

func TestEmptyGet(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	_, err := cl.Get(context.Background(), uuid.New())
	if err == nil {
		t.Fatalf("non empty storage")
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	key := uuid.New()
	err := cl.Set(context.Background(), key, bytes.NewReader([]byte("{}")), 0)
	if err != nil {
		t.Fatalf("set: %v", err)
	}

	v, err := cl.Get(context.Background(), key)
	if err != nil {
		t.Fatalf("del: %v", err)
	}
	if !reflect.DeepEqual(v, map[string]interface{}{}) {
		t.Errorf("unexpected response: %v", v)
	}
}

func TestEmptyDel(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	err := cl.Del(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("del unexistent key: %v", err)
	}
}

func TestDel(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	key := uuid.New()
	err := cl.Set(context.Background(), key, bytes.NewReader([]byte("{}")), 0)
	if err != nil {
		t.Fatalf("set: %v", err)
	}

	err = cl.Del(context.Background(), key)
	if err != nil {
		t.Fatalf("del: %v", err)
	}
}

func TestDelDropTimer(t *testing.T) {
	server := httptest.NewServer(handler())
	defer server.Close()

	cl := client.NewClient(server.URL)
	key := uuid.New()
	err := cl.Set(context.Background(), key, bytes.NewReader([]byte("{}")), 120*time.Second)
	if err != nil {
		t.Fatalf("set: %v", err)
	}

	err = cl.Del(context.Background(), key)
	if err != nil {
		t.Fatalf("del: %v", err)
	}
}
