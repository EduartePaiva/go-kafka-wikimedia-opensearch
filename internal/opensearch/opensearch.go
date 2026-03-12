package opensearch

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"

	osh "github.com/opensearch-project/opensearch-go/v4"
	"github.com/opensearch-project/opensearch-go/v4/opensearchapi"
)

type Client struct {
	api *opensearchapi.Client
}

func NewOpenSearchClient(addresses []string) (*Client, error) {
	client, err := opensearchapi.NewClient(
		opensearchapi.Config{
			Client: osh.Config{
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
				Addresses: addresses,
			},
		},
	)
	return &Client{api: client}, err
}

func (c *Client) CreateIndex(indexName string) error {
	mapping := strings.NewReader(`{
    "settings": {
        "index": {
            "number_of_shards": 1,
            "number_of_replicas": 0
            }
        }
    }`)
	ctx := context.Background()

	res, _ := c.api.Indices.Exists(ctx, opensearchapi.IndicesExistsReq{Indices: []string{indexName}})
	if res.StatusCode == http.StatusOK {
		fmt.Printf("Index \"%s\" already exists.\n", indexName)
		return nil
	}

	createIndexResponse, err := c.api.Indices.Create(
		ctx,
		opensearchapi.IndicesCreateReq{
			Index: indexName,
			Body:  mapping,
		},
	)
	if err != nil {
		return err
	}
	fmt.Printf("Created Index: %s\nShards Acknowledged: %t\n", createIndexResponse.Index, createIndexResponse.ShardsAcknowledged)

	return nil
}
