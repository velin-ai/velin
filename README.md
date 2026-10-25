# Velin

**Your AI homelab.** Velin is a self-hosted AI chat backend with a web interface: user accounts, API keys, chat history, and an [Ollama](https://ollama.com)-backed model gateway.

English | [Türkçe](README_tr.md)

## Features

- Username/password authentication with short-lived JWT access tokens and rotating refresh tokens
- API key management — keys are shown once and stored as SHA-256 hashes
- Chats and messages stored in PostgreSQL, with per-user ownership checks
- Last 20 messages of a chat sent to the model as context
- Self-hosted model calls through Ollama; no external AI provider required
- Embedded web UI served directly by the Go binary
- English and Turkish UI with an in-app language selector

## Stack

| Layer | Technology |
| --- | --- |
| Backend | Go 1.27, chi router |
| Database | PostgreSQL |
| Auth | JWT (HS256), bcrypt, SHA-256 |
| AI | Ollama |
| Frontend | React 19, TypeScript, Vite |
| Configuration | OS keyring via `go-keyring` |

## Requirements

- Go 1.27 or newer
- Node.js 20.19 or newer and npm
- PostgreSQL 14 or newer
- [Ollama](https://ollama.com) running locally (or reachable on the configured host)
- A working OS keyring: Secret Service on Linux, Keychain on macOS, Credential Manager on Windows

## Setup

### 1. Database

Create an empty PostgreSQL database:

```sh
createdb velin
```

Tables are created automatically on startup.

### 2. Configuration

Velin reads its secrets from the operating system keyring. Store them once with the following commands. Each command writes the value and exits.

```sh
go run ./cmd/velin -set "your-admin-password"
go run ./cmd/velin -db "postgres://user:password@localhost:5432/velin?sslmode=disable"
go run ./cmd/velin -jwt "a-long-random-secret"
```

On headless Linux systems, make sure a Secret Service provider such as `gnome-keyring` is running and unlocked.

### 3. Build the frontend

The Go binary embeds `ui/api/dist` at compile time, so the frontend must be built first:

```sh
cd ui/api
npm install
npm run build
cd ../..
```

If `ui/api/dist` does not exist, `go build` fails with an embed error.

### 4. Run

```sh
go run ./cmd/velin
```

or build a binary:

```sh
go build -o velin ./cmd/velin
./velin
```

Open <http://localhost:8080>.

### Install with `go install`

You can also install the `velin` binary directly from the published module:

```sh
go install github.com/velin-ai/velin/cmd/velin@latest
```

`go install` places the binary in `$(go env GOPATH)/bin` (or `$GOBIN`). Add that directory to `PATH` if it is not there already:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

Then configure Velin once and start it like any other command:

```sh
velin -set "your-admin-password"
velin -db "postgres://user:password@localhost:5432/velin?sslmode=disable"
velin -jwt "a-long-random-secret"
velin
```

The web UI is embedded at compile time. `go install ...@latest` only includes it when `ui/api/dist` is present in the published module version. For a local checkout, build the frontend first and install from the repository instead:

```sh
cd ui/api && npm run build && cd ../..
go install ./cmd/velin
```

### Development mode

Run the API and the Vite dev server side by side:

```sh
# terminal 1
go run ./cmd/velin
```

```sh
# terminal 2
cd ui/api
npm install
npm run dev
```

The frontend runs on <http://localhost:5173> and proxies `/api` requests to `http://localhost:8080`, so UI changes do not require a Go rebuild.

## Ollama

Chat requests are sent to `http://localhost:11434/api/chat`. The endpoint is defined in `internal/api/ai/ollama.go`; adjust it if Ollama runs on another host.

```sh
ollama serve
ollama pull llama3.1
```

The model name is supplied with each chat request, so any model available in Ollama can be used.

## Configuration flags

| Flag | Description | Default |
| --- | --- | --- |
| `-set` | Store or change the admin password in the keyring | — |
| `-db` | Store or change the PostgreSQL connection URL | — |
| `-jwt` | Store or change the JWT signing secret | — |
| `-api-port` | Address the API server listens on | `:8080` |
| `-admin-port` | Reserved for the admin service | `:9090` |

Values are stored in the keyring under the service name `velin`.

## API

All API routes are prefixed with `/api/v1`.

| Method | Path | Authentication | Description |
| --- | --- | --- | --- |
| `POST` | `/auth/register` | Public | Create an account |
| `POST` | `/auth/login` | Public | Sign in and receive access and refresh tokens |
| `POST` | `/auth/refresh` | Public | Rotate a refresh token and issue a new access token |
| `POST` | `/api-key` | JWT | Create an API key |
| `GET` | `/api-keys` | JWT | List your API keys |
| `POST` | `/chats` | JWT | Create a chat |
| `GET` | `/chats` | JWT | List your chats |
| `POST` | `/chats/{chat_id}/messages` | API key | Send a message and receive the model reply |
| `GET` | `/` | Public | Serve the embedded web UI |

JWT-protected routes expect:

```http
Authorization: Bearer <access-token>
```

Chat message routes expect an API key:

```http
Authorization: Bearer vln_live_<key>
```

### Example

```sh
curl -s http://localhost:8080/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo","password":"secret"}'
```

## Project structure

```text
cmd/velin/                 Entrypoint and HTTP router
internal/admin/config/     OS keyring configuration
internal/api/ai/           Ollama client
internal/api/handlers/     Auth, API key and chat handlers
internal/api/middlewares/  JWT and API key authentication
internal/api/models/       Request and response models
internal/store/            PostgreSQL connection and schema
ui/api/                    React + TypeScript frontend
ui/api/src/i18n/           English and Turkish translations
assets/logo.svg            Brand mark
```

## Security notes

- Passwords are hashed with bcrypt (cost 12).
- Access tokens are signed with HS256 and expire after 30 minutes.
- Refresh tokens are valid for 30 days, stored as SHA-256 hashes, and rotated on use.
- API keys are shown only once and stored as SHA-256 hashes.
- Chat and message access is scoped to the authenticated user.
- The frontend keeps tokens in `localStorage`. If you deploy Velin publicly, consider an HTTP-only cookie session layer and additional hardening.

## Internationalization

The UI defaults to English. Turkish is available under **Dashboard → Settings → Language**. The choice is stored in `localStorage` and applied immediately.

Translations live in:

- `ui/api/src/i18n/en.json`
- `ui/api/src/i18n/tr.json`

To add a language, add a JSON file, register it in `ui/api/src/i18n/I18nProvider.tsx` and `ui/api/src/i18n/context.ts`, and add the language label to the existing dictionaries.

## License

[MIT](LICENSE)
