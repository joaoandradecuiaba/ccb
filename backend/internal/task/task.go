// Package task contém o modelo de domínio e o repositório de tarefas.
package task

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// Status possíveis de uma tarefa.
const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusDone       = "done"
)

var (
	ErrNotFound      = errors.New("tarefa não encontrada")
	ErrInvalidTitle  = errors.New("o título é obrigatório (máx. 200 caracteres)")
	ErrInvalidStatus = errors.New("status inválido: use pending, in_progress ou done")
)

// Task representa uma tarefa do sistema.
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Input é o payload aceito na criação e atualização.
type Input struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// Normalize limpa os campos e aplica valores padrão.
func (in *Input) Normalize() {
	in.Title = strings.TrimSpace(in.Title)
	in.Description = strings.TrimSpace(in.Description)
	if in.Status == "" {
		in.Status = StatusPending
	}
}

// Validate verifica se o payload é válido.
func (in Input) Validate() error {
	if in.Title == "" || len([]rune(in.Title)) > 200 {
		return ErrInvalidTitle
	}
	switch in.Status {
	case StatusPending, StatusInProgress, StatusDone:
		return nil
	default:
		return ErrInvalidStatus
	}
}

// Store é um repositório em memória, seguro para uso concorrente,
// com persistência opcional em arquivo JSON.
type Store struct {
	mu     sync.RWMutex
	tasks  map[int64]Task
	nextID int64
	path   string
}

// NewStore cria um repositório. Se path não for vazio, os dados são
// carregados desse arquivo (se existir) e salvos nele a cada alteração.
func NewStore(path string) (*Store, error) {
	s := &Store{tasks: make(map[int64]Task), nextID: 1, path: path}
	if path == "" {
		return s, nil
	}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	var list []Task
	if err := json.Unmarshal(data, &list); err != nil {
		return nil, err
	}
	for _, t := range list {
		s.tasks[t.ID] = t
		if t.ID >= s.nextID {
			s.nextID = t.ID + 1
		}
	}
	return s, nil
}

// List retorna as tarefas ordenadas por ID, filtrando opcionalmente por status.
func (s *Store) List(status string) []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		if status == "" || t.Status == status {
			out = append(out, t)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// Get busca uma tarefa pelo ID.
func (s *Store) Get(id int64) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return t, nil
}

// Create adiciona uma nova tarefa.
func (s *Store) Create(in Input) (Task, error) {
	in.Normalize()
	if err := in.Validate(); err != nil {
		return Task{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	t := Task{
		ID:          s.nextID,
		Title:       in.Title,
		Description: in.Description,
		Status:      in.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.tasks[t.ID] = t
	s.nextID++
	return t, s.persist()
}

// Update substitui os campos editáveis de uma tarefa existente.
func (s *Store) Update(id int64, in Input) (Task, error) {
	in.Normalize()
	if err := in.Validate(); err != nil {
		return Task{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	t.Title = in.Title
	t.Description = in.Description
	t.Status = in.Status
	t.UpdatedAt = time.Now().UTC()
	s.tasks[id] = t
	return t, s.persist()
}

// Delete remove uma tarefa.
func (s *Store) Delete(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return s.persist()
}

// persist grava o estado atual no arquivo (escrita atômica). Deve ser
// chamado com o lock de escrita adquirido.
func (s *Store) persist() error {
	if s.path == "" {
		return nil
	}
	list := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		list = append(list, t)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}
