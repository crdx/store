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

func New(baseUrl string, apiToken string) *Store {
	// Ensure the baseUrl ends with a trailing slash so that relative paths are handled correctly
	// by the Path method of carlmjohnson/requests.
	baseUrl = strings.TrimRight(baseUrl, "/") + "/"

	return &Store{
		baseUrl:  baseUrl,
		apiToken: apiToken,
	}
}

func (self Store) httpClient() *requests.Builder {
	return requests.URL(self.baseUrl).Bearer(self.apiToken)
}

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
