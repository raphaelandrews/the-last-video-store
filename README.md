# 🎬 THE LAST VIDEO STORE

> A cyber-secure, RBAC-protected video rental terminal with retro 2000s aesthetic.
> Built in Go with Bubble Tea TUI, BoltDB persistence, and deployable REST API.

---

## Overview

**The Last Video Store** is a full-stack video rental management system styled after the golden age of
Blockbuster. It features a rich terminal user interface powered by
[Charmbracelet Bubble Tea](https://github.com/charmbracelet/bubbletea) and a REST API
backend that runs on [Render](https://render.com).

Users browse a ~135 movie catalog, rent VHS/DVD/Blu-ray tapes, earn 🍿 Popcorn Points on returns,
redeem rewards (free rentals, collectibles, merch), and purchase **premium subscription tiers**
(Wood → Bronze → Silver → Gold → Diamond) that grant monthly free rental allowances and perks.

Staff access is gated by a **7-role Role-Based Access Control (RBAC)** system enforced through 6-bit bitmask operations, separate from the subscription tier system.

## Quick Start

```bash
# Clone and enter
git clone https://github.com/anomalyco/the-last-video-store && cd the-last-video-store

# Seed the database (~135 movies, 8 users, merchandise catalog)
go run ./data/seed.go

# Start the API server
go run ./cmd/server/

# In another terminal, launch the TUI client
go run ./cmd/client/
# or: go run ./cmd/client/ --api-url http://localhost:8080
```

All 8 test users share the password: `123`

## Features

### Catalog & Browsing
- 🎞️ **~147 Movies** — VHS, DVD, Blu-ray spanning 1937–2022, 8 genres with custom rental prices
- 🔍 **Live Search** — `/` opens search bar, type for instant prefix results
- 📄 **Paginated Browse** — 40 movies per page, `N`/`B` navigation
- 🏷️ **Genre Tabs** — `[`/`]` filter by Action, SciFi, Horror, Comedy, Drama, Thriller, Romance, Animation
- ⭐ **Staff Picks & Last Chance** — `S`/`L` toggle curated + disappearing titles
- 📊 **Movie Recommendations** — Same-genre suggestions + franchise chain (sequels/prequels)

### Rental System
- 📼 **Format-Aware Rentals** — VHS: 2 min, DVD: 3 min, Blu-ray: 4 min (demo-friendly)
- ⏱️ **Due Date Countdown** — "in N min", "due soon", "overdue N min ago"
- 📅 **Extend Rentals** — `E` on rentals page, costs 30🍿 for +1 min
- 💰 **Late & Rewind Fees** — Per-minute late fees; $1 VHS rewind fee (30% random chance)
- 💵 **Paid Rentals** — Movie-specific prices ($2.99–$6.99); free from tier allowance first
- 🎟️ **Free Rentals** — From reward coupons or tier allowance, waive all late fees
- 🎛️ **Payment Choice** — `T` for ticket or `M` for money on ENTER

### Premium Subscription Tiers
| Tier | Price | Free Rentals/mo | Max Concurrent | New Releases | Late Fees |
|------|------:|:--------------:|:------------:|:---:|:---:|
| Wood | Free | 0 | 5 | No | Yes |
| Bronze | $9.99 | 1 | 10 | No | Yes |
| Silver | $19.99 | 3 | 15 | Yes | Yes |
| Gold | $29.99 | 5 | 25 | Yes | **Waived** |
| Diamond | $49.99 | 10 | 50 | Yes | **Waived** |

### 💵 Money & 🍿 Popcorn Points (Dual Currency)
- **Money ($)** — Used to pay for rentals beyond tier allowance and to buy premium tiers
- **Popcorn Points (🍿)** — Loyalty points earned on returns, spent on rewards
- On-time return: +10🍿
- Late return: -5🍿, late fees deducted from money balance
- Popcorn Bucket collectible gives +5🍿 bonus on every return

### Rewards & Collectibles
- 🎁 **Rewards Shop** — `M` from profile, 7 items
- 🎒 **Inventory** — `I` from profile, view owned collectibles
- 🎟️ **Free Rental Coupon** (200🍿) — +1 free rental ticket
- 🍿 **Popcorn Bucket** (50🍿) — +5 bonus points on every future return
- ⬆️ **Tier Upgrade** (1000🍿) — Promote RBAC role one level (up to Gold)
- 🎉 **Private Screening** (500🍿) — +5 free rental tickets

### Admin & Security
- 🏷️ **7 RBAC Roles** — Bronze → Silver → Gold → Employee → Supervisor → Manager → Owner
- 🛡️ **RBAC Bitmask** — O(1) permission checks
- 🔐 **JWT Auth** — Access (15min) + Refresh (7day) tokens with rotation
- 🔒 **TOTP 2FA** — Optional time-based one-time passwords (Manager+)
- 📋 **Wishlist** — Add from detail page, view/remove with `V` key
- 🧾 **Immutable Audit Trail** — SHA-256 hash chain, tamper-detectable
- 🚫 **Brute-force Lockout** — 5 failed attempts = 30min lock
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
│   ├── models/                 # User, Movie, Rental, Wishlist, Audit, Merch, Inventory, Tier
│   └── store/                  # BoltDB persistence layer
├── api/                        # Chi REST API (handlers, middleware, DTOs)
├── tui/                        # Bubble Tea TUI (app, keys, commands, views, pages, components, styles)
├── data/seed.go                # Database seed script (~135 movies, 8 users, merch)
├── tests/                      # Integration tests
├── Dockerfile.server           # Docker build for Render
├── render.yaml                 # Render Blueprint config
├── Makefile                    # Build/test/run/cross-compile
├── FEATURES.md                 # Full feature reference
└── README.md
```

## Test Users

| Username   | Password | Subscription | 🍿 Points | 💵 Balance |
|-----------|---------|-------------|:--------:|:--------:|
| bronze    | 123     | Bronze      | 250      | $50.00   |
| silver    | 123     | Silver      | 250      | $50.00   |
| gold      | 123     | Gold        | 250      | $50.00   |
| employee  | 123     | Gold        | 250      | $50.00   |
| supervisor| 123     | Gold        | 250      | $50.00   |
| manager   | 123     | Diamond     | 250      | $100.00  |
| owner     | 123     | Diamond     | 250      | $100.00  |
| banned    | 123     | Wood        | 250      | $5.00    |

## API Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `POST` | `/api/v1/auth/register` | — | Register new member (Wood subscription) |
| `POST` | `/api/v1/auth/login` | — | Login, get JWT (+ TOTP prompt if enabled) |
| `POST` | `/api/v1/auth/login/totp` | Temp | TOTP 2FA step 2 |
| `POST` | `/api/v1/auth/refresh` | JWT | Rotate refresh token |
| `GET`  | `/api/v1/movies` | JWT | List movies (paginated, 40/page) |
| `GET`  | `/api/v1/movies/search` | JWT | Prefix search by title |
| `GET`  | `/api/v1/movies/staff-picks` | JWT | Curated picks |
| `GET`  | `/api/v1/movies/last-chance` | JWT | Last copies available |
| `GET`  | `/api/v1/movies/{id}` | JWT | Movie detail |
| `POST` | `/api/v1/movies` | Manager+ | Create movie |
| `PUT`  | `/api/v1/movies/{id}` | Manager+ | Update movie |
| `DELETE`| `/api/v1/movies/{id}` | Manager+ | Delete movie |
| `POST` | `/api/v1/rentals/rent` | JWT | Rent a movie (tier allowance or 💵) |
| `POST` | `/api/v1/rentals/return` | JWT | Return a movie (+🍿, ±💵) |
| `POST` | `/api/v1/rentals/extend` | JWT | Extend due date (30🍿, +2 days) |
| `GET`  | `/api/v1/rentals/history` | JWT | Rental history |
| `GET`  | `/api/v1/wishlist` | JWT | View wishlist |
| `POST` | `/api/v1/wishlist` | JWT | Add to wishlist |
| `DELETE`| `/api/v1/wishlist/{movieID}` | JWT | Remove from wishlist |
| `GET`  | `/api/v1/tiers` | JWT | List subscription tiers |
| `POST` | `/api/v1/tiers/purchase` | JWT | Buy/upgrade subscription tier |
| `GET`  | `/api/v1/merch` | JWT | Popcorn Points rewards catalog |
| `POST` | `/api/v1/merch/redeem` | JWT | Redeem 🍿 for a reward |
| `GET`  | `/api/v1/inventory` | JWT | View your collectibles |
| `GET`  | `/api/v1/users` | Supervisor+ | List all users |
| `POST` | `/api/v1/users` | Supervisor+ | Create user |
| `PUT`  | `/api/v1/users/{id}` | Supervisor+ | Update user role/ban |
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
