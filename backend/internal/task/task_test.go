package task

import (
	"path/filepath"
	"testing"
)

func TestStorePersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")

	s, err := NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.Create(Input{Title: "a"}); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Create(Input{Title: "b", Status: StatusInProgress}); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(1); err != nil {
		t.Fatal(err)
	}

	reloaded, err := NewStore(path)
	if err != nil {
		t.Fatal(err)
	}
	list := reloaded.List("")
	if len(list) != 1 || list[0].Title != "b" {
		t.Fatalf("dados recarregados inesperados: %+v", list)
	}
	created, err := reloaded.Create(Input{Title: "c"})
	if err != nil {
		t.Fatal(err)
	}
	if created.ID != 3 {
		t.Fatalf("próximo ID = %d, esperado 3", created.ID)
	}
}
