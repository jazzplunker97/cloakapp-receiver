package db

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"os"
	"sync"
	"time"
)

type APIKey struct {
	Key     string `json:"key"`
	Label   string `json:"label"`
	Created string `json:"created"`
}

type KeyStore struct {
	Keys map[string]APIKey `json:"keys"`
	mu   sync.RWMutex
	path string
}

func NewKeyStore(path string) (*KeyStore, error) {
	ks := &KeyStore{
		Keys: make(map[string]APIKey),
		path: path,
	}
	err := ks.load()
	return ks, err
}

func (ks *KeyStore) load() error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	data, err := os.ReadFile(ks.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return json.Unmarshal(data, &ks.Keys)
}

func (ks *KeyStore) save() error {
	data, err := json.Marshal(ks.Keys)
	if err != nil {
		return err
	}
	return os.WriteFile(ks.path, data, 0644)
}

func (ks *KeyStore) Generate(label string) (string, error) {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	key := hex.EncodeToString(b)
	
	ks.Keys[key] = APIKey{
		Key:     key,
		Label:   label,
		Created: time.Now().Format(time.RFC3339),
	}
	
	return key, ks.save()
}

func (ks *KeyStore) Delete(key string) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	delete(ks.Keys, key)
	return ks.save()
}

func (ks *KeyStore) Validate(key string) bool {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	_, ok := ks.Keys[key]
	return ok
}

func (ks *KeyStore) List() []APIKey {
	ks.mu.RLock()
	defer ks.mu.RUnlock()
	list := make([]APIKey, 0, len(ks.Keys))
	for _, v := range ks.Keys {
		list = append(list, v)
	}
	return list
}
