package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/joaoandradecuiaba/ccb/backend/internal/task"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	store, err := task.NewStore("")
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(NewHandler(store, Config{}))
	t.Cleanup(srv.Close)
	return srv
}

func do(t *testing.T, method, url, body string) *http.Response {
	t.Helper()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { resp.Body.Close() })
	return resp
}

func TestTaskCRUD(t *testing.T) {
	srv := newTestServer(t)

	resp := do(t, http.MethodPost, srv.URL+"/api/tasks", `{"title":"  Estudar Go  ","description":"net/http"}`)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create: status %d", resp.StatusCode)
	}
	var created task.Task
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	if created.ID != 1 || created.Title != "Estudar Go" || created.Status != task.StatusPending {
		t.Fatalf("create: tarefa inesperada %+v", created)
	}

	resp = do(t, http.MethodPut, srv.URL+"/api/tasks/1", `{"title":"Estudar Go","status":"done"}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("update: status %d", resp.StatusCode)
	}

	resp = do(t, http.MethodGet, srv.URL+"/api/tasks?status=done", "")
	var list []task.Task
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Status != task.StatusDone {
		t.Fatalf("list: resultado inesperado %+v", list)
	}

	if resp = do(t, http.MethodDelete, srv.URL+"/api/tasks/1", ""); resp.StatusCode != http.StatusNoContent {
		t.Fatalf("delete: status %d", resp.StatusCode)
	}
	if resp = do(t, http.MethodGet, srv.URL+"/api/tasks/1", ""); resp.StatusCode != http.StatusNotFound {
		t.Fatalf("get após delete: status %d", resp.StatusCode)
	}
}

func TestValidation(t *testing.T) {
	srv := newTestServer(t)
	cases := []struct {
		name, method, path, body string
		want                     int
	}{
		{"título vazio", http.MethodPost, "/api/tasks", `{"title":"  "}`, http.StatusUnprocessableEntity},
		{"status inválido", http.MethodPost, "/api/tasks", `{"title":"x","status":"foo"}`, http.StatusUnprocessableEntity},
		{"json inválido", http.MethodPost, "/api/tasks", `{`, http.StatusBadRequest},
		{"campo desconhecido", http.MethodPost, "/api/tasks", `{"title":"x","foo":1}`, http.StatusBadRequest},
		{"id inválido", http.MethodGet, "/api/tasks/abc", "", http.StatusBadRequest},
		{"inexistente", http.MethodPut, "/api/tasks/99", `{"title":"x"}`, http.StatusNotFound},
		{"rota desconhecida", http.MethodGet, "/api/nada", "", http.StatusNotFound},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if resp := do(t, tc.method, srv.URL+tc.path, tc.body); resp.StatusCode != tc.want {
				t.Fatalf("status %d, esperado %d", resp.StatusCode, tc.want)
			}
		})
	}
}
