package opensearch

import (
	"bytes"
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
    },
	"mappings": {
		"dynamic": false
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

func (c *Client) AddToIndex(ctx context.Context, indexName string, docId string, data []byte) error {
	insertResp, err := c.api.Index(
		ctx,
		opensearchapi.IndexReq{
			Index:      indexName,
			DocumentID: docId,
			Body:       bytes.NewReader(data),
			Params: opensearchapi.IndexParams{
				Refresh: "true",
			},
		},
	)
	if err != nil {
		return err
	}
	fmt.Printf("Created document in %s ID: %s\n", insertResp.Index, insertResp.ID)

	return nil
}
