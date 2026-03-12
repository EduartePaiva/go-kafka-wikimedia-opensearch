package opensearch

import (
	"crypto/tls"
	"net/http"

	osh "github.com/opensearch-project/opensearch-go/v4"
)

type OpenSearchClient struct {
	Client *osh.Client
}

func NewOpenSearchClient() (OpenSearchClient, error) {
	client, err := osh.NewClient(osh.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: []string{"http://localhost:9200"},
	})
	return OpenSearchClient{Client: client}, err
}
