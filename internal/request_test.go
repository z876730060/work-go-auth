package internal

import (
	"crypto/tls"
	"io"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	req, err := http.NewRequest(http.MethodGet, "https://localhost/debug/pprof", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatal("unexpected status code", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
	if string(body) != "ok" {
		t.Fatal("unexpected body", string(body))
	}
	t.Log("request success")
}
