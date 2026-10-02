# Gerenciador de Tarefas — Go + Angular

Sistema web completo com **API REST em Go** e **front-end em Angular**.

```
.
├── backend/                 # API em Go (apenas biblioteca padrão)
│   ├── cmd/server/          # main: configuração, servidor e shutdown gracioso
│   └── internal/
│       ├── task/            # modelo, validação e repositório (memória + arquivo JSON)
│       └── httpapi/         # rotas, handlers, CORS, logs e servidor da SPA
├── frontend/                # Angular 21 (standalone, signals, reactive forms)
│   └── src/app/tasks/       # model, service HTTP e tela de tarefas
├── Dockerfile               # build multi-stage: Angular + Go em uma única imagem
└── Makefile
```

## Pré-requisitos

- Go 1.24+
- Node.js 22.12+ (ou 24+) e npm

## Desenvolvimento

Na primeira vez: `cd frontend && npm install`.

Depois, em dois terminais:

```bash
make backend    # API em http://localhost:8080
make frontend   # App em http://localhost:4200
```

O `ng serve` usa `frontend/proxy.conf.json` para encaminhar `/api` ao backend,
então não é preciso configurar CORS em desenvolvimento.

## Testes

```bash
make test
```

## Produção (Docker)

```bash
docker build -t ccb .
docker run -p 8080:8080 -v ccb-data:/app/data ccb
```

Acesse http://localhost:8080 — o binário Go serve a API e o build do Angular.

## API

| Método | Rota              | Descrição                                               |
|--------|-------------------|---------------------------------------------------------|
| GET    | `/api/health`     | Verificação de saúde                                    |
| GET    | `/api/tasks`      | Lista tarefas (`?status=pending`, `in_progress`, `done`) |
| POST   | `/api/tasks`      | Cria tarefa                                             |
| GET    | `/api/tasks/{id}` | Busca tarefa                                            |
| PUT    | `/api/tasks/{id}` | Atualiza tarefa                                         |
| DELETE | `/api/tasks/{id}` | Remove tarefa                                           |

Corpo de criação/atualização:

```json
{ "title": "Estudar Go", "description": "opcional", "status": "pending" }
```

Erros retornam `{ "error": "mensagem" }` com status 400, 404 ou 422.

## Configuração do backend (variáveis de ambiente)

| Variável      | Padrão                  | Descrição                                    |
|---------------|-------------------------|----------------------------------------------|
| `ADDR`        | `:8080`                 | Endereço de escuta                           |
| `DATA_FILE`   | `data/tasks.json`       | Arquivo de persistência                      |
| `CORS_ORIGIN` | `http://localhost:4200` | Origem liberada no CORS (vazio = desativado) |
| `STATIC_DIR`  | —                       | Pasta do build do Angular a ser servida      |

## Próximos passos sugeridos

- Trocar o repositório em arquivo por um banco (PostgreSQL/SQLite) mantendo a mesma interface do `task.Store`
- Autenticação (JWT) e usuários
- Paginação e busca textual
