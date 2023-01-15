package porkbun

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, pattern, filename string) *Client {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	mux.HandleFunc(pattern, func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(rw, "invalid method: "+req.Method, http.StatusBadRequest)
			return
		}

		all, _ := io.ReadAll(req.Body)
		if !strings.HasPrefix(string(all), `{"apikey":"key","secretapikey":"secret"`) {
			http.Error(rw, `{"status": "ERROR","message": "invalid auth"}`, http.StatusOK)
			return
		}

		data, err := os.ReadFile(fmt.Sprintf("./fixtures/%s.json", filename))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		_, _ = rw.Write(data)
	})

	client := New("secret", "key")
	client.BaseURL, _ = url.Parse(server.URL)

	return client
}

func TestClient_Ping(t *testing.T) {
	client := setup(t, "/ping", "ping")

	ping, err := client.Ping(context.Background())
	require.NoError(t, err)

	assert.Equal(t, "2a02:842b:5da:c101:4b81:e1b5:83f7:3e7c", ping)
}

func TestClient_Ping_error(t *testing.T) {
	client := setup(t, "/ping", "error")

	_, err := client.Ping(context.Background())
	require.Error(t, err)
}

func TestClient_CreateRecord(t *testing.T) {
	client := setup(t, "/dns/create/example.com", "create")

	record := Record{
		Type:    "TXT",
		Content: "foobar",
		TTL:     DefaultTTL,
		Prio:    "1",
	}

	id, err := client.CreateRecord(context.Background(), "example.com", record)
	require.NoError(t, err)

	assert.Equal(t, 106926659, id)
}

func TestClient_CreateRecord_error(t *testing.T) {
	client := setup(t, "/dns/create/example.com", "error")

	record := Record{
		Type:    "TXT",
		Content: "foobar",
		TTL:     DefaultTTL,
		Prio:    "1",
	}

	_, err := client.CreateRecord(context.Background(), "example.com", record)
	require.Error(t, err)
}

func TestClient_EditRecord(t *testing.T) {
	client := setup(t, "/dns/edit/example.com/666", "edit")

	record := Record{
		Type:    "TXT",
		Content: "foobar",
		TTL:     DefaultTTL,
		Prio:    "1",
	}

	err := client.EditRecord(context.Background(), "example.com", 666, record)
	require.NoError(t, err)
}

func TestClient_EditRecord_error(t *testing.T) {
	client := setup(t, "/dns/edit/example.com/666", "error")

	record := Record{
		Type:    "TXT",
		Content: "foobar",
		TTL:     DefaultTTL,
		Prio:    "1",
	}

	err := client.EditRecord(context.Background(), "example.com", 666, record)
	require.Error(t, err)
}

func TestClient_DeleteRecord(t *testing.T) {
	client := setup(t, "/dns/delete/example.com/666", "edit")

	err := client.DeleteRecord(context.Background(), "example.com", 666)
	require.NoError(t, err)
}

func TestClient_DeleteRecord_error(t *testing.T) {
	client := setup(t, "/dns/delete/example.com/666", "error")

	err := client.DeleteRecord(context.Background(), "example.com", 666)
	require.Error(t, err)
}

func TestClient_RetrieveRecords(t *testing.T) {
	client := setup(t, "/dns/retrieve/example.com", "retrieve")

	records, err := client.RetrieveRecords(context.Background(), "example.com")
	require.NoError(t, err)

	expected := []Record{
		{
			ID:      "106926652",
			Name:    "borseth.ink",
			Type:    "A",
			Content: "1.1.1.1",
			TTL:     "300",
			Prio:    "0",
			Notes:   "",
		},
		{
			ID:      "106926659",
			Name:    "www.borseth.ink",
			Type:    "A",
			Content: "1.1.1.1",
			TTL:     "300",
			Prio:    "0",
			Notes:   "",
		},
	}

	assert.Equal(t, expected, records)
}

func TestClient_RetrieveRecords_error(t *testing.T) {
	client := setup(t, "/dns/retrieve/example.com", "error")

	_, err := client.RetrieveRecords(context.Background(), "example.com")
	require.Error(t, err)
}
