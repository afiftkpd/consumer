package es

import (
	"consumer/models"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"testing"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/stretchr/testify/assert"
)

var metaHeaderReValidation = regexp.MustCompile(`^[a-z]{1,}=[a-z0-9\.\-]{1,}(?:,[a-z]{1,}=[a-z0-9\.\-]+)*$`)
var called bool
var defaultURL string

type TransportMock struct {
	Response    *http.Response
	RoundTripFn func(req *http.Request) (*http.Response, error)
}

func (t *TransportMock) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.RoundTripFn(req)
}

func TestClientConfiguration(t *testing.T) {
	defaultURL = "http://localhost:9200"

	// Create Mock for Elasticsearch Server
	mock := TransportMock{
		Response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(strings.NewReader(`{}`)),
			Header:     http.Header{"X-Elastic-Product": []string{"Elasticsearch"}},
		},
	}
	mock.RoundTripFn = func(req *http.Request) (*http.Response, error) { return mock.Response, nil }

	client, err := elasticsearch.NewTypedClient(elasticsearch.Config{
		Transport: &mock,
	})
	if err != nil {
		t.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	t.Run("With empty", func(t *testing.T) {
		c, err := elasticsearch.NewDefaultClient()

		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		u := c.Transport.(*elastictransport.Client).URLs()[0].String()

		if u != defaultURL {
			t.Errorf("Unexpected URL, want=%s, got=%s", defaultURL, u)
		}
	})

	t.Run("With URL from Addresses", func(t *testing.T) {
		c, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{"http://localhost:8080//"}})
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		u := c.Transport.(*elastictransport.Client).URLs()[0].String()

		if u != "http://localhost:8080" {
			t.Errorf("Unexpected URL, want=http://localhost:8080, got=%s", u)
		}
	})

	t.Run("Insert Doc", func(t *testing.T) {
		repo := NewProductRepository(client)

		// data := models.Product
		mock.Response = &http.Response{
			Body: io.NopCloser(strings.NewReader(`{"id": 1, "name": "baju", "price": 1000, "rating": 1, "image_url": "http://google.com"}`)),
		}

		err = repo.Store(context.Background(), models.Product{
			ID:       1,
			Name:     "baju",
			Price:    1000,
			Rating:   1,
			ImageURL: "http://google.com",
		})

		assert.NoError(t, err)
	})

	t.Run("Delete Doc", func(t *testing.T) {
		repo := NewProductRepository(client)

		// data := models.Product
		mock.Response = &http.Response{
			Body: io.NopCloser(strings.NewReader(`{"id": 1, "name": "baju", "price": 1000, "rating": 1, "image_url": "http://google.com"}`)),
		}

		err = repo.Delete(context.Background(), "1")

		assert.NoError(t, err)
	})

	t.Run("Update Doc", func(t *testing.T) {
		repo := NewProductRepository(client)

		// data := models.Product
		mock.Response = &http.Response{
			Body: io.NopCloser(strings.NewReader(`{"id": 1, "name": "baju", "price": 1000, "rating": 1, "image_url": "http://google.com"}`)),
		}

		err = repo.Update(context.Background(), models.Product{
			ID:       1,
			Name:     "baju",
			Price:    1000,
			Rating:   1,
			ImageURL: "http://google.com",
		})

		assert.NoError(t, err)
	})
}
