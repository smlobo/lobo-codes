package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestRqliteLogRequestUpdatesWithParameters(t *testing.T) {
	var requests []struct {
		path string
		body []byte
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("read request body: %v", err)
			return
		}
		requests = append(requests, struct {
			path string
			body []byte
		}{r.URL.Path, body})
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/db/query" {
			_, _ = io.WriteString(w, `{"results":[{"rows":[{"id":7}]}]}`)
		} else {
			_, _ = io.WriteString(w, `{"results":[{"rows_affected":1}]}`)
		}
	}))
	defer server.Close()

	previousURL := rqliteURL
	rqliteURL = server.URL
	defer func() { rqliteURL = previousURL }()

	info := &RequestInfo{
		RemoteAddress: "192.0.2.1",
		UserAgent:     `agent' OR 1=1 -- "quoted"`,
	}
	rqliteLogRequest(info, `visitor"logs`, httptest.NewRequest(http.MethodGet, "/", nil))

	if len(requests) != 2 {
		t.Fatalf("got %d requests, want query and update", len(requests))
	}
	for i, want := range [][]interface{}{
		{`SELECT id FROM "visitor""logs" WHERE remote_address=? AND user_agent=?`, info.RemoteAddress, info.UserAgent},
		{`UPDATE "visitor""logs" SET count = count + 1, updated_at = ? WHERE id = ?`, info.UpdatedAt.Format(time.RFC3339Nano), float64(7)},
	} {
		var statements [][]interface{}
		if err := json.Unmarshal(requests[i].body, &statements); err != nil {
			t.Fatalf("decode request %d: %v", i, err)
		}
		if len(statements) != 1 || !reflect.DeepEqual(statements[0], want) {
			t.Errorf("request %d: got %v, want %v", i, statements, want)
		}
	}
	if requests[0].path != "/db/query" || requests[1].path != "/db/execute" {
		t.Errorf("unexpected endpoints: %s, %s", requests[0].path, requests[1].path)
	}
}

func TestRqliteRequestBodyEscapesSQLAsJSON(t *testing.T) {
	query := fmt.Sprintf(`INSERT INTO %s (user_agent) VALUES (?)`, rqliteTable(`visitors`))
	value := "agent \"quoted\"\nline"
	body, err := rqliteRequestBody(query, value)
	if err != nil {
		t.Fatal(err)
	}
	var statements [][]string
	if err := json.Unmarshal(body, &statements); err != nil {
		t.Fatal(err)
	}
	if len(statements) != 1 || !reflect.DeepEqual(statements[0], []string{query, value}) {
		t.Fatalf("got %v", statements)
	}
}
