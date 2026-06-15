# 🎬 THE LAST VIDEO STORE

> A cyber-secure, RBAC-protected video rental terminal with retro 2000s aesthetic.
> Built in Go with Bubble Tea TUI, BoltDB persistence, and deployable REST API.

---

## Overview

**The Last Video Store** is a full-stack video rental management system styled after the golden age of
Blockbuster. It features a rich terminal user interface powered by
[Charmbracelet Bubble Tea](https://github.com/charmbracelet/bubbletea) and a REST API
backend that runs on [Render](https://render.com).

Users browse a 40+ movie catalog, rent VHS/DVD/Blu-ray tapes, return them, manage their membership,
and explore co-rental recommendations — all gated by a **7-tier Role-Based Access Control (RBAC)**
system enforced through 6-bit bitmask operations.

## Quick Start

```bash
# Clone and enter
git clone https://github.com/anomalyco/the-last-video-store && cd the-last-video-store

# Seed the database
go run ./data/seed.go

# Start the API server
go run ./cmd/server/

# In another terminal, launch the TUI client
go run ./cmd/client/
# or: go run ./cmd/client/ --api-url http://localhost:8080
```

## Features

- 🎞️ **40+ Movies** — VHS, DVD, Blu-ray with cover art and synopses
- 🔍 **Trie Autocomplete** — Type 2 letters for instant search suggestions
- 📼 **Format-Aware Rentals** — VHS: 3 days, DVD/Blu-ray: 5 days
- 💰 **Late & Rewind Fees** — $2/day VHS, $3/day DVD; $1 VHS rewind fee
- 🏷️ **7 Membership Tiers** — Bronze → Silver → Gold → Employee → Supervisor → Manager → Owner
- 🛡️ **RBAC Bitmask** — O(1) permission checks, Staff bit cleanly separates employees from Gold members
- 🔐 **JWT Auth** — Access (15min) + Refresh (7day) tokens with rotation
- 🔒 **TOTP 2FA** — Optional time-based one-time passwords (Manager+)
- 📋 **Wishlist** — Add titles, get "Available now!" notifications
- ⭐ **Staff Picks & Last Chance** — Curated recommendations + disappearing titles
- 📊 **Graph Recommendations** — "Customers who rented this also rented..."
- 🧾 **Immutable Audit Trail** — SHA-256 hash chain, tamper-detectable
- 🚫 **Brute-force Lockout** — 5 failed attempts = 30-minute account lock
- 🔮 **Bloom Filter** — O(k) banned-user lookup
- 🖥️ **Cross-Platform** — Linux + Windows binaries

## Project Structure

```
thelastvideostore/
├── cmd/
│   ├── server/main.go          # API server entrypoint
│   └── client/main.go          # TUI client entrypoint
├── internal/
│   ├── auth/                   # bcrypt, JWT, RBAC, lockout, TOTP
│   ├── config/                 # Environment configuration
│   ├── crypto/                 # AES-256-GCM, hash chain
│   ├── ds/
│   │   ├── bitmask/            # 6-bit permission bitmask
│   │   ├── list/               # Doubly linked list
│   │   ├── deque/              # Ring buffer deque
│   │   ├── heap/               # Binary min-heap
│   │   ├── trie/               # Prefix tree (autocomplete)
│   │   ├── lru/                # LRU cache
│   │   ├── bloom/              # Bloom filter
│   │   └── graph/              # Undirected weighted graph
│   ├── models/                 # User, Movie, Rental, Wishlist, Audit
│   └── store/                  # BoltDB persistence layer
├── api/                        # Chi REST API (handlers, middleware, DTOs)
├── tui/                        # Bubble Tea TUI (pages, components, styles)
├── data/seed.go                # Database seed script (40+ movies, 8 users)
├── tests/                      # Integration tests
├── Dockerfile.server           # Docker build for Render
├── render.yaml                 # Render Blueprint config
├── Makefile                    # Build/test/run/cross-compile
└── README.md
```

## Test Users

| Username   | Password   | Tier       | Max Rentals |
|-----------|-----------|-----------|:----------:|
| bronze    | password1 | Bronze    | 1 |
| silver    | password2 | Silver    | 2 |
| gold      | password3 | Gold      | 5 |
| employee  | password4 | Employee  | 5 |
| supervisor| password8 | Supervisor| 5 |
| manager   | password5 | Manager   | 10 |
| owner     | password6 | Owner     | ∞ |
| banned    | password7 | Bronze    | (suspended) |

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/v1/auth/register` | — | Register new member (Bronze) |
| `POST` | `/api/v1/auth/login` | — | Login, get JWT (+ TOTP prompt if enabled) |
| `POST` | `/api/v1/auth/login/totp` | Temp | TOTP 2FA step 2 |
| `POST` | `/api/v1/auth/refresh` | JWT | Rotate refresh token |
| `POST` | `/api/v1/auth/logout` | JWT | Logout |
| `GET`  | `/api/v1/movies` | JWT | List movies (paginated, genre filter) |
| `GET`  | `/api/v1/movies/search` | JWT | Prefix search by title |
| `GET`  | `/api/v1/movies/staff-picks` | JWT | Curated picks |
| `GET`  | `/api/v1/movies/last-chance` | JWT | Last copies available |
| `GET`  | `/api/v1/movies/{id}` | JWT | Movie detail |
| `POST` | `/api/v1/movies` | Manager+ | Create movie |
| `PUT`  | `/api/v1/movies/{id}` | Manager+ | Update movie |
| `DELETE`| `/api/v1/movies/{id}` | Manager+ | Delete movie |
| `POST` | `/api/v1/rentals/rent` | JWT | Rent a movie |
| `POST` | `/api/v1/rentals/return` | JWT | Return a movie |
| `GET`  | `/api/v1/rentals/history` | JWT | User's rental history |
| `GET`  | `/api/v1/wishlist` | JWT | View wishlist |
| `POST` | `/api/v1/wishlist` | JWT | Add to wishlist |
| `DELETE`| `/api/v1/wishlist/{movieID}` | JWT | Remove from wishlist |
| `GET`  | `/api/v1/users` | Supervisor+ | List all users |
| `POST` | `/api/v1/users` | Supervisor+ | Create user |
| `PUT`  | `/api/v1/users/{id}` | Supervisor+ | Update user tier/ban |
| `DELETE`| `/api/v1/users/{id}` | Manager+ | Delete user |
| `POST` | `/api/v1/users/{id}/totp` | Self/Manager+ | Enable/disable TOTP 2FA |
| `GET`  | `/api/v1/audit` | Supervisor+ | View audit log |

## Data Structures (All From Scratch)

9 custom Go generics data structures, zero `container/*` usage:

| Structure | Application | Complexity |
|-----------|------------|:---:|
| Bitmask (6-bit) | RBAC permissions | O(1) |
| Doubly Linked List | Rental history, wishlist ordering | O(1) insert/remove |
| Deque (Ring Buffer) | Express return queue | O(1) push/pop |
| Min-Heap | New release waitlist | O(log n) |
| Trie (Prefix Tree) | Movie title autocomplete | O(k) |
| LRU Cache | Recently viewed, session cache | O(1) |
| Bloom Filter | Banned user fast check | O(k) |
| Hash Chain (SHA-256) | Immutable audit trail | O(1) append |
| Undirected Weighted Graph | Co-rental recommendations | O(V+E) |

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.24+ |
| TUI | Bubble Tea + Lipgloss + Bubbles |
| API Router | Chi v5 |
| Database | BoltDB (embedded) |
| Auth | bcrypt + JWT HS256 + TOTP RFC 6238 |
| Encryption | AES-256-GCM |
| Hash Chain | SHA-256 Merkle-style linking |
| Deployment | Docker + Render |

## Makefile Commands

```bash
make build            # Build server + client
make cross-compile    # Cross-compile for Linux + Windows
make test             # Run unit tests
make test-integration # Run integration tests
make seed             # Seed database
make run-server       # Start API server
make run-client       # Start TUI client
make clean            # Clean binaries + DB
make fmt              # Format code
make docker-build     # Build Docker image
make docker-run       # Run Docker container
```

## Deployment

### Render

1. Push to GitHub
2. Create a new Web Service on Render
3. Connect the repository
4. Render auto-detects `render.yaml` or use:
   - **Runtime:** Docker
   - **Dockerfile Path:** `Dockerfile.server`
   - **Port:** 8080
   - **Health Check:** `/health`
5. Required env vars are auto-generated

### Local with Docker

```bash
make docker-build
make docker-run
# Server at http://localhost:8080
```

## Cross-Compilation

```bash
make cross-compile
# Produces:
#   bin/thelastvideostore-linux    (Linux amd64)
#   bin/thelastvideostore.exe      (Windows amd64)
```

Both binaries are statically linked — no runtime dependencies.

## License

MIT — Academic project for Cybersecurity & Data Structures.
