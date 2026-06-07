package test

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/nskforward/gate4/pkg/ssl"
)

func TestHTTPClient(t *testing.T) {
	err := sendHttpRequest(t.Context())
	if err != nil {
		t.Fatal(err)
	}
}

func sendHttpRequest(ctx context.Context) error {
	req, err := http.NewRequest("GET", "https://localhost/api/whoami", nil)
	if err != nil {
		return fmt.Errorf("cannot create http request: %w", err)
	}

	req.WithContext(ctx)
	req.Header.Add("Authorization", "Bearer 45925f18-8934-4d1b-b756-508b4e9dadb2")
	resp, err := httpClient().Do(req)
	if err != nil {
		return fmt.Errorf("cannot execute http request: %w", err)
	}
	defer resp.Body.Close()

	fmt.Println("---- response ----")
	io.Copy(os.Stdout, resp.Body)
	fmt.Println("------------------")

	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("bad status code: %s", resp.Status)
	}

	return nil
}

func httpClient() *http.Client {
	caCert, err := ssl.LoadCertificate("../data/ssl/ca.crt")
	if err != nil {
		panic(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AddCert(caCert)

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs: caCertPool,
			},
		},
	}
}
