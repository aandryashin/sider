package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type node struct {
	done chan struct{}
	data interface{}
}

var (
	storage map[string]*node = make(map[string]*node)
	lock    sync.RWMutex
)

func set(w http.ResponseWriter, r *http.Request, key string, ttl time.Duration) {
	var data interface{}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Parse request: %v", err), http.StatusBadRequest)
		return
	}
	
	lock.RLock()
	_, ok := storage[key]
	lock.RUnlock()
	
	if ok {
		http.Error(w, "Key already exists", http.StatusConflict)
		return
	}

	lock.Lock()
	defer lock.Unlock()
	var done chan struct{}
	if ttl != 0 {
		log.Printf("Key: [%s] expires in: [%v].\n", key, ttl)
		done = make(chan struct{})
		go func(key string, done chan struct{}) {
			select {
			case <-time.After(ttl):
				lock.Lock()
				defer lock.Unlock()
				delete(storage, key)
				log.Printf("Expired: [%s].\n", key)
			case <-done:
				log.Printf("Drop timer: [%s].\n", key)
			}
		}(key, done)
	}

	storage[key] = &node{done: done, data: data}
	log.Printf("Set: [%s].\n", key)
}

func get(w http.ResponseWriter, r *http.Request, key string, ttl time.Duration) {
	lock.RLock()
	defer lock.RUnlock()
	v, ok := storage[key]
	if !ok {
		http.Error(w, fmt.Sprintf("Key [%s] not found.", key), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(v.data)
	log.Printf("Get: [%s].\n", key)
}

func del(w http.ResponseWriter, r *http.Request, key string, ttl time.Duration) {
	lock.Lock()
	defer lock.Unlock()
	v, ok := storage[key]
	if !ok {
		return
	}
	if v.done != nil {
		close(v.done)
	}
	delete(storage, key)
	log.Printf("Del: [%s].\n", key)
}

func list(w http.ResponseWriter, r *http.Request) {
	var list []string = []string{}
	lock.RLock()
	defer lock.RUnlock()
	for k, _ := range storage {
		list = append(list, k)
		select {
		case <-r.Context().Done():
			return
		default:
		}
	}
	json.NewEncoder(w).Encode(list)
}
