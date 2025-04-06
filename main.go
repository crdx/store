package store

import (
	"context"
	"fmt"
	"strings"

	"github.com/carlmjohnson/requests"
)

type Store struct {
	baseUrl  string
	apiToken string
}

// New instantiates and returns a [*Store] using baseUrl as the base URL and apiToken as the API
// token.
func New(baseUrl string, apiToken string) *Store {
	// Ensure the baseUrl ends with a trailing slash so that relative paths are handled correctly
	// by the Path method of carlmjohnson/requests.
	baseUrl = strings.TrimRight(baseUrl, "/") + "/"

	return &Store{
		baseUrl:  baseUrl,
		apiToken: apiToken,
	}
}

// Set sets the value of a key and returns a human-readable string stating which action the server
// carried out.
func (self Store) Set(key, value string) (string, error) {
	var res baseResponse

	req := setRequest{
		Value: value,
	}

	err := self.httpClient().
		Path(key).
		BodyJSON(&req).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}

	if !res.Success {
		return "", fmt.Errorf("set: %s", res.Message)
	}

	return res.Message, nil
}

// Append appends a line to the current value of a key and returns a human-readable string stating which
// action the server carried out.
func (self Store) Append(key, value string) (string, error) {
	currentValue, err := self.Get(key)
	if err != nil {
		return "", err
	}

	var newValue string

	if currentValue == "" {
		newValue = value
	} else {
		lines := strings.Split(currentValue, "\n")
		lines = append(lines, value)
		newValue = strings.Join(lines, "\n")
	}

	return self.Set(key, newValue)
}

// GetOrDefault returns the value of a key. If the value is empty then defaultValue is returned.
func (self Store) GetOrDefault(key, defaultValue string) (string, error) {
	value, err := self.Get(key)
	if err != nil {
		return "", err
	}

	if value == "" {
		return defaultValue, nil
	}

	return value, nil
}

// Get returns the value of a key.
func (self Store) Get(key string) (string, error) {
	var res getResponse

	err := self.httpClient().
		Path(key).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}

	if !res.Success {
		return "", fmt.Errorf("get: %s", res.Message)
	}

	return strings.TrimSuffix(res.Value, "\n"), nil
}

// Delete deletes a key and returns a human-readable string stating which action the server carried
// out.
func (self Store) Delete(key string) (string, error) {
	var res baseResponse

	err := self.httpClient().
		Delete().
		Path(key).
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		return "", err
	}

	if !res.Success {
		return "", fmt.Errorf("delete: %s", res.Message)
	}

	return res.Message, nil
}

// List returns a slice containing each of the keys currently stored.
func (self Store) List() ([]string, error) {
	var res listResponse

	err := self.httpClient().
		ToJSON(&res).
		Fetch(context.Background())
	if err != nil {
		return nil, err
	}

	if !res.Success {
		return nil, fmt.Errorf("list: %s", res.Message)
	}

	var items []string
	for _, item := range res.Items {
		items = append(items, item.K)
	}

	return items, nil
}

// —————————————————————————————————————————————————————————————————————————————————————————————————

func (self Store) httpClient() *requests.Builder {
	return requests.URL(self.baseUrl).Bearer(self.apiToken)
}
