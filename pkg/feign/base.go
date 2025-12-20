package feign

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type R[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	Info    map[string]any
}

func do[T any](request *http.Request, key string) (t T, err error) {
	instance, err := fClients.Get(key)
	if err != nil {
		return t, err
	}

	request.URL.Scheme = instance.Scheme
	request.URL.Host = fmt.Sprintf("%s:%d", instance.Ip, instance.Port)
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Accept-Charset", "utf-8")
	response, err := httpClient.Do(request)
	if err != nil {
		return t, err
	}
	defer response.Body.Close()

	all, err := io.ReadAll(response.Body)
	if err != nil {
		return t, err
	}

	var resp R[T]
	err = json.Unmarshal(all, &resp)
	if err != nil {
		return t, err
	}
	t = resp.Data
	return
}
