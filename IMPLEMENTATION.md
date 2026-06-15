# THE LAST VIDEO STORE — Retro Video Rental Terminal

> A cyber-secure, RBAC-protected video rental system with a retro 2000s aesthetic.
> Built in Go with a Bubble Tea TUI, BoltDB persistence, and deployable REST API.
> Developed for the **Cybersecurity & Data Structures** university project.

---

## 1. Project Overview

**The Last Video Store** is a full-stack video rental management system styled after the golden age of
Blockbuster. It features a rich terminal user interface powered by
[Charmbracelet Bubble Tea](https://github.com/charmbracelet/bubbletea) and a REST API
backend that runs on [Render](https://render.com).

Users browse a movie catalog, rent tapes, return them, and manage their membership — all
gated by a **7-tier Role-Based Access Control (RBAC)** system enforced through bitmask

The system directly addresses the three core challenges outlined in the project scope:

| Requirement | The Last Video Store Implementation |
|-------------------------|--------------------------|
| **a) Read the file line by line** | BoltDB stores movies persistently; the API serves paginated catalog data; the TUI renders movie cards with autocomplete powered by a custom Trie |
| **b) Allow only authorized users to read** | RBAC bitmask + JWT middleware + brute-force lockout after 5 failed attempts |
| **c) Show the user and the file data** | Profile screen displays tier, active rentals, and full rental history; audit log shows an immutable hash-chained trail of all access events |

### 1.1 System Features at a Glance

**Catalog & Navigation**
- Autocomplete search powered by custom Trie — type 2 letters and see instant DVD/VHS suggestions
- Genre-based filtering via tabs: ALL, ACTION, COMEDY, HORROR, SCIFI, DRAMA, NEW, STAFF PICKS, LAST CHANCE
- Format badges on every title: `📀 DVD` / `📼 VHS` / `💿 Blu-ray` — different rental rules per format
- "New Releases" shelf — latest arrivals marked `[NEW]`, restricted by membership plan (Gold+)
- "Staff Picks" — Manager-curated recommendations, toggled from Admin Movies panel
- "Last Chance" — titles with only 1 copy remaining, about to leave the catalog permanently
- Co-rental recommendations on Movie Detail: "Customers who rented this also rented..." (powered by custom Graph DS)

**Rental & Return**
- Format-aware rental durations — VHS: 3 days, DVD/Blu-ray: 5 days
- Format-specific late fees — $2/day for VHS, $3/day for DVD/Blu-ray
- Rewind fee — $1.00 if VHS tape returned unrewound (30% random chance, set at rental time)
- Simultaneous rental limits per membership plan: Bronze=1, Silver=2, Gold/Employee/Supervisor=5, Manager=10, Owner=∞
- New release waitlist — min-heap ordered by wait time, notified when copies return
- Express return — priority deque processes the most overdue return first
- Full rental history — per-user doubly linked list, navigable in both directions

**Membership Plans**
- **Bronze** (free) — browse catalog, rent up to 1 title, wishlist, no new releases
- **Silver** ($9.99/mo) — rent up to 2 titles, wishlist, no new releases
- **Gold** ($19.99/mo) — rent up to 5 titles, new releases, waitlist, wishlist
- Each plan gets a color-coded badge + membership card in the Profile page
- Plan upgrades/downgrades processed by Supervisors or Managers via admin panel

**Wishlist**
- Every member (Bronze+) can add/remove titles to a personal wishlist
- Wishlist stored server-side in BoltDB, displayed on the Browse page sidebar
- "Available now!" notification when a wishlisted title has copies back in stock
- Quick-rent from wishlist — select and rent in one keystroke

**Security**
- Brute-force lockout — 5 failed attempts = 30-minute account lock
- JWT sessions with rotating refresh tokens (15-min access / 7-day refresh)
- Immutable audit trail — hash chain with SHA-256 linking, tamper-detectable
- Bloom filter for banned members — O(k) lookup before any operation
- bcrypt password hashing (cost 12) — passwords never stored in plaintext
- AES-256-GCM encryption for sensitive data at rest (audit logs, TOTP secrets)
- Optional TOTP 2FA — Manager+ can enable per-account, HMAC-SHA1 time-based codes

**Social & Gamification**
- Star ratings — 0 to 5 stars, community average displayed
- Membership plan badge — color-coded card with tier stats (7 tiers)
- "Now Showing" — featured title of the week in the header
- Popcorn Points — earned on punctual returns (10 per on-time, -5 per late)

**Administration (Supervisor / Gerente)**
- Full CRUD on movie catalog — add, edit, remove titles with format (DVD/VHS/Blu-ray) — Manager+ only
- Full CRUD on users — create, upgrade/downgrade plan, ban — Supervisor+
- Audit log viewer — scrollable hash chain with integrity verification — Supervisor+
- Staff return processing — Employee+ can process returns for any customer via priority deque
- Live dashboard — copies rented, overdue count, daily revenue per format

---

## 2. Security & Access Control

### 2.1 Authentication Flow

```
┌──────────┐     POST /auth/login      ┌──────────────┐
│  Client  │ ─────────────────────────► │   API Server │
│  (TUI)   │ ◄─────── JWT ──────────── │   (chi)      │
│          │     (access + refresh)     │              │
│          │                           │  bcrypt      │
│          │  Every request includes    │  compare     │
│          │  Authorization: Bearer     │  against     │
│          │  <token>                   │  BoltDB      │
└──────────┘                           └──────────────┘
```

- **bcrypt** (cost factor 12) for password hashing — never store plaintext
- **JWT** (HS256) access tokens with 15-minute expiry
- **Refresh tokens** (opaque, stored in BoltDB) with 7-day expiry, rotated on each use
- **Brute-force protection**: 5 failed login attempts → 30-minute lockout (per IP + per user)
- **TOTP 2FA** (optional): Time-based one-time password via HMAC-SHA1 — Manager+ can enable per-account; adds second authentication step on login
- **AES-256-GCM** encryption for sensitive data at rest (audit logs, refresh tokens, TOTP secrets)

### 2.2 RBAC — 7-tier Hierarchy with Bitmask (Membership Plans + Staff)

Permissions are encoded as 6-bit integers checked in O(1) via bitwise `&`.
The system models a real Blockbuster store: customers join a **membership plan** (Bronze / Silver / Gold),
while store staff (Employee / Supervisor / Manager / Owner) have operational privileges.
Gold and Employee no longer share the same bitmask — `PermStaff` cleanly distinguishes staff from premium customers.

```
                    ┌─────────┐
                    │  OWNER  │  0b111111 → Supreme authority (Dono)
                    └────┬────┘
                    ┌────┴────┐
                    │ MANAGER │  0b111111 → Full CRUD + staff + audit (Gerente)
                    └────┬────┘
                    ┌────┴────┐
                    │SUPERVIS.│  0b011111 → Manage users + staff + audit
                    └────┬────┘
                    ┌────┴────┐
                    │EMPLOYEE │  0b010111 → Staff: process any return
                    │(Atend.) │
                    └────┬────┘
               ┌─────────┴──────────────┐
          ┌────┴────┐              ┌────┴────┐
          │  GOLD   │ 0b000111     │ SILVER  │ 0b000011
          │Cl. Ouro │              │Cl. Prata│
          └─────────┘              └────┬────┘
                                        │
                                   ┌────┴────┐
                                   │ BRONZE  │ 0b000001
                                   └─────────┘
```

| Tier | Role (PT-BR) | Bitmask | Max Rentals | New Releases | Wishlist | Audit | Staff |
|------|-------------|:-------:|:-----------:|:---:|:---:|:---:|:---:|
| Bronze | Cliente Bronze | `0b000001` | 1 | ❌ | ✅ | ❌ | — |
| Silver | Cliente Prata | `0b000011` | 2 | ❌ | ✅ | ❌ | — |
| Gold | Cliente Ouro | `0b000111` | 5 | ✅ | ✅ | ❌ | — |
| Employee | Atendente | `0b010111` | 5 | ✅ | ✅ | ❌ | ✅ |
| Supervisor | Supervisor | `0b011111` | 5 | ✅ | ✅ | ✅ | ✅ |
| Manager | Gerente | `0b111111` | 10 | ✅ | ✅ | ✅ | ✅ |
| Owner | Dono | `0b111111` | ∞ | ✅ | ✅ | ✅ | ✅ |

**Key distinction:** The 6th bit `PermStaff` (0b010000) cleanly separates employees from Gold members.
No string-based role checks needed — all permission logic is pure bitmask arithmetic.
Supervisor is a new intermediate role with `PermManageUsers` + `PermStaff` but without `PermAdmin` (cannot CRUD movies).
Manager and Owner share the same bitmask (`0b111111`); the distinction is symbolic — Owner is the store proprietor.

| Permission Bit | Constant | Bronze | Silver | Gold | Employee | Supervisor | Manager | Owner |
|:---:|:---|:---:|:---:|:---:|:---:|:---:|:---:|:---:|
| `0b000001` | `PermBrowse` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `0b000010` | `PermRent` | — | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `0b000100` | `PermReserve` | — | — | ✅ | ✅ | ✅ | ✅ | ✅ |
| `0b001000` | `PermManageUsers` | — | — | — | — | ✅ | ✅ | ✅ |
| `0b010000` | `PermStaff` | — | — | — | ✅ | ✅ | ✅ | ✅ |
| `0b100000` | `PermAdmin` | — | — | — | — | — | ✅ | ✅ |

When a user attempts an action, the middleware computes:

```go
if session.Permissions & PermRent == 0 {
    return "⛔ ACCESS DENIED — Insufficient clearance"
}
```

### 2.3 Business Rules

**Membership Plans (Planos de Sócio)**

| Plan | Monthly Fee | Max Rentals | New Releases | VHS Late Fee | DVD/BD Late Fee | VHS Duration | DVD/BD Duration | Rewind Fee |
|------|:----------:|:-----------:|:---:|:---:|:---:|:---:|:---:|:---:|
| Bronze | Free | 1 | ❌ | $2/day | $3/day | 3 days | 5 days | $1.00 |
| Silver | $9.99 | 2 | ❌ | $2/day | $3/day | 3 days | 5 days | $1.00 |
| Gold | $19.99 | 5 | ✅ | $2/day | $3/day | 3 days | 5 days | $1.00 |

- Plan upgrades performed by **Supervisors** (Supervisor) or **Managers (Gerente)** via admin panel
- Plan downgrades also require staff approval
- Membership fees tracked on Profile page (cosmetic — no real payment integration)

**DVD vs VHS Format Rules**

Each movie stores its `Format`: `DVD`, `VHS`, or `Blu-ray`. This affects rental duration,
late fees, rewind fees, and inventory tracking:
- VHS tapes: 3-day rental window (analog tapes wear faster)
- DVDs / Blu-rays: 5-day rental (discs are more durable)
- Late fees: VHS = $2/day, DVD/Blu-ray = $3/day (higher replacement cost)
- **Rewind fee:** Returning a VHS tape unrewound costs $1.00 (30% random chance on return, set at rental time)
- Inventory tracked per format (e.g., 3 VHS + 2 DVD copies of The Matrix)
- Format badge on every movie card: `📼 VHS` / `📀 DVD` / `💿 Blu-ray`

**Rental Limits (Limite de Locações Simultâneas)**
- Bronze: 1 simultaneous rental
- Silver: 2 simultaneous rentals
- Gold / Employee / Supervisor: 5
- Manager: 10
- Owner: unlimited
- Exceeding limit triggers: `"Rental limit reached (X/Y)"` modal

**Late Fees (Multa por Atraso)**
- Auto-calculated on return: `days_overdue × daily_rate` (per format: VHS=$2/day, DVD/Blu-ray=$3/day)
- Popcorn Points deducted for late returns (-5 per late)
- **Rewind fee:** $1.00 added if VHS tape was returned unrewound (`NeedsRewind` flag on rental)

**Waitlist (Fila de Espera para Lançamentos)**
- Gold members can join waitlist when all copies of a new release are rented
- Min-heap ordered by earliest join timestamp — oldest waiter gets first notification
- Per-title, per-format queue (separate waitlists for DVD vs VHS of same movie)

### 2.4 Additional Security Layers

- **Bloom Filter** for banned-user membership check — O(k) lookup with <0.1% false-positive rate
- **Hash Chain** (blockchain-lite) for immutable audit trail — each log entry contains `SHA-256(prev_hash + entry_data)`
- **Session management** — active sessions tracked in BoltDB; tokens invalidated on logout
- **Rate limiting** — 100 requests/min per IP via token bucket in middleware
- **Input sanitization** — all user input validated and sanitized before storage

---

## 3. Data Structures

Every data structure is **implemented from scratch in pure Go** (no `container/*` usage).
Each includes table-driven unit tests and benchmarks.

| Structure | File | Application | Complexity |
|-----------|------|-------------|:---:|
| **Trie (Prefix Tree)** | `internal/ds/trie/trie.go` | Movie title autocomplete in search bar + wishlist title lookup — user types "mat" → sees "Matrix", "Matilda", "Match Point" | Insert: O(k), Search: O(k) |
| **LRU Cache** | `internal/ds/lru/lru.go` | Recently viewed movies + active session cache + wishlist status cache — avoids repeated BoltDB reads | Get/Put: O(1) |
| **Deque** | `internal/ds/deque/deque.go` | Return processing queue — FIFO for order of arrival, pop-back to prioritize most overdue | Push/Pop either end: O(1) |
| **Min-Heap** | `internal/ds/heap/heap.go` | New-release waitlist per movie+format — customer with longest wait time rises to top | Insert: O(log n), Extract: O(log n) |
| **Doubly Linked List** | `internal/ds/list/linkedlist.go` | Rental history per user + wishlist items per user — navigable in both directions | Insert/Remove: O(1) at known node |
| **Bloom Filter** | `internal/ds/bloom/bloom.go` | Banned-user fast check before any DB query | Add/Contains: O(k) |
| **Bitmask** | `internal/ds/bitmask/bitmask.go` | RBAC permission flags — single integer comparison | Check: O(1) |
| **Hash Chain** | `internal/crypto/hashchain.go` | Immutable audit log — each entry links to previous via SHA-256 | Append: O(1), Verify: O(n) |
| **Graph (Undirected)** | `internal/ds/graph/graph.go` | Co-rental recommendation engine — edge weight = # of users who rented both movies; BFS for similarity search | AddEdge: O(1), GetRecommendations: O(V+E) |

---

## 4. System Architecture

```
┌───────────────────────────────────────────────────────────────────┐
│                        USER'S TERMINAL                           │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    TUI CLIENT (Bubble Tea)                   │ │
│  │  ┌─────────┐ ┌──────────┐ ┌──────────┐ ┌────────────────┐  │ │
│  │  │ Splash  │ │  Login   │ │  Browse  │ │  Admin Panel   │  │ │
│  │  │ Screen  │→│ Register │→│ Catalog  │→│  Users/Movies  │  │ │
│  │  └─────────┘ └──────────┘ └──────────┘ └────────────────┘  │ │
│  │       │              │            │              │           │ │
│  │       └──────────────┴────────────┴──────────────┘           │ │
│  │                          │  HTTPS + JWT                      │ │
│  └──────────────────────────┼───────────────────────────────────┘ │
└─────────────────────────────┼─────────────────────────────────────┘
                              │
                              ▼
┌───────────────────────────────────────────────────────────────────┐
│                   RENDER CLOUD (Docker)                           │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    API SERVER (Go + Chi)                     │ │
│  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌───────────────┐ │ │
│  │  │   Auth   │ │  Movies  │ │ Rentals  │ │  Audit Log    │ │ │
│  │  │ Handler  │ │ Handler  │ │ Handler  │ │  Handler      │ │ │
│  │  └────┬─────┘ └────┬─────┘ └────┬─────┘ └───────┬───────┘ │ │
│  │       │             │            │               │          │ │
│  │       └─────────────┴────────────┴───────────────┘          │ │
│  │                          │                                   │ │
│  │                    ┌─────┴─────┐                             │ │
│  │                    │  BoltDB   │                             │ │
│  │                    │ (persist) │                             │ │
│  │                    └───────────┘                             │ │
│  └─────────────────────────────────────────────────────────────┘ │
└───────────────────────────────────────────────────────────────────┘
```

### Cross-Compilation Targets

```bash
GOOS=linux   GOARCH=amd64 go build -o thelastvideostore-linux   ./cmd/client
GOOS=windows GOARCH=amd64 go build -o thelastvideostore.exe      ./cmd/client
```

Both binaries ship as single statically-linked executables — no runtime dependencies.

### 4.3 TUI Interface Mockup — Browse Catalog (Main Screen)

```
╔══════════════════════════════════════════════════════════════════════════╗
║ ██▄ █ █ █ █▄▀ █▄▄ █ █ ▄▀▀ ▀█▀ ██▀   █ ▄▀▄ █▄▀ ▄▀▀ ▀█▀ ██▀ ██▀     ║
║ █▄█ ▀▄▀ █ █▀▄ █▄█ ▀▄▀ ▀▄▄  █  █▄▄   █ ▀▄▀ █▀▄ ▀▄▄  █  █▄▄ █▄▄     ║
╠══════════════════════════════════════════════════════════════════════════╣
║  ★ NOW SHOWING: The Matrix (1999)  │  🕐 Fri Jun 14 2002  9:48 PM      ║
╠══════════════════════════════════════════════════════════════════════════╣
║                                                                          ║
║  ┌──────────────────────────────┐  ┌──────────────────────────────────┐ ║
║  │ 🔍 Search: mat_             │  │ 📼 Your Rentals: 2 / 5           │ ║
║  │ ─────────────────────────── │  │ • The Matrix — due Jun 16        │ ║
║  │ ▸ The Matrix (1999) ★★★★½  │  │ • Pulp Fiction — due Jun 15      │ ║
║  │   Matilda (1996) ★★★★☆      │  │                                  │ ║
║  │   Match Point (2005) ★★★☆☆  │  │  🍿 Popcorn Points: 142         │ ║
║  └──────────────────────────────┘  └──────────────────────────────────┘ ║
║                                                                          ║
║  ┌─[ALL]──[ACTION]──[COMEDY]──[HORROR]──[SCIFI]──[NEW]──[STAFF ★]──[⏳ LAST]─┐ ║
║  │                                                                      │ ║
║  │  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐               │ ║
║  │  │ ░░░░░░░░ │ │ ░░░░░░░░ │ │ ░░░░░░░░ │ │ ░░░░░░░░ │               │ ║
║  │  │ MATRIX   │ │ PULP     │ │ FIGHT    │ │ JURASSIC │               │ ║
║  │  │ ★★★★½   │ │ FICTION  │ │ CLUB     │ │ PARK     │               │ ║
║  │  │ [RENT]   │ │ ★★★★★   │ │ ★★★★☆   │ │ ★★★½☆   │               │ ║
║  │  │          │ │ [OUT]    │ │ [RENT]   │ │ [RENT]   │               │ ║
║  │  └──────────┘ └──────────┘ └──────────┘ └──────────┘               │ ║
║  │                                                                      │ ║
║  └──────────────────────────────────────────────────────────────────────┘ ║
╠══════════════════════════════════════════════════════════════════════════╣
║ │[TAB] sections [↓↑] navigate [ENTER] select [R] rent [Q] quit        │ ║
╚══════════════════════════════════════════════════════════════════════════╝
```

**Visual effects implemented in the TUI:**
- CRT scanlines overlay with 2% opacity on the background
- Random glitch effect on page transitions (1 frame of scrambled characters)
- "Rewinding..." animation when returning a tape (ASCII spinner with blinking text)
- Neon color palette: cyan `#00FFFF`, magenta `#FF00FF`, yellow `#FFFF00` on deep navy background `#0A0A2E`
- Movie cards with lipgloss gradient borders
- Context-sensitive footer showing available keybindings per screen

### 4.4 Navigation Flow

```
        ┌───────────┐
        │  SPLASH   │  3s VHS-style animated intro
        └─────┬─────┘
              ▼
   ┌─────────────────────┐
   │  LOGIN / REGISTER   │  (+ TOTP prompt if 2FA enabled)
   └─────────┬───────────┘
             ▼
   ┌───────────────────────────────────────────────────────────────────┐
   │                          BROWSE CATALOG                            │
   │  (search bar + genre tabs + responsive movie card grid +           │
   │   wishlist sidebar + staff picks / last chance tabs)               │
   └────┬──────────────┬──────────────┬──────────────┬──────────────────┘
        ▼              ▼              ▼              ▼
 ┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────────┐
 │ MOVIE       │ │ MY RENTALS  │ │   PROFILE   │ │   ADMIN PANEL   │
 │ DETAIL      │ │ (return,    │ │ (stats,     │ │                 │
 │ (rent,      │ │  history,   │ │  tier badge,│ ├────────┬────────┤
 │  waitlist,  │ │  late fees, │ │  popcorn,   │ │ USERS  │ MOVIES │
 │  recommend) │ │  rewind fee)│ │  TOTP)      │ │(suprv+)|(mgr+)  │
 └─────────────┘ └─────────────┘ └─────────────┘ └───┬────┴────┬───┘
                                                     │         │
                                                     └────┬────┘
                                                          ▼
                                                   ┌─────────────┐
                                                   │ AUDIT LOG   │
                                                   │ (hash chain │
                                                   │  viewer,    │
                                                   │  suprv+)    │
                                                   └─────────────┘
```

---

## 5. Project Directory Structure

```
thelastvideostore/
│
├── cmd/
│   ├── server/
│   │   └── main.go                  # API server entrypoint (Render deploy target)
│   └── client/
│       └── main.go                  # TUI client entrypoint (local execution)
│
├── internal/
│   ├── auth/
│   │   ├── password.go              # bcrypt hash + verify
│   │   ├── session.go               # JWT create, validate, refresh
│   │   ├── permissions.go           # Bitmask constants + role hierarchy
│   │   ├── totp.go                  # TOTP 2FA implementation (HMAC-SHA1, time-based)
│   │   └── lockout.go               # Brute-force lockout tracking
│   │
│   ├── crypto/
│   │   ├── aes.go                   # AES-256-GCM encrypt/decrypt
│   │   └── hashchain.go             # Blockchain-like immutable audit trail
│   │
│   ├── ds/
│   │   ├── trie/
│   │   │   ├── trie.go              # Prefix tree implementation
│   │   │   └── trie_test.go
│   │   ├── lru/
│   │   │   ├── lru.go               # LRU cache (hashmap + doubly linked list)
│   │   │   └── lru_test.go
│   │   ├── deque/
│   │   │   ├── deque.go             # Double-ended queue (ring buffer)
│   │   │   └── deque_test.go
│   │   ├── heap/
│   │   │   ├── heap.go              # Binary min-heap
│   │   │   └── heap_test.go
│   │   ├── list/
│   │   │   ├── linkedlist.go        # Doubly linked list
│   │   │   └── linkedlist_test.go
│   │   ├── bloom/
│   │   │   ├── bloom.go             # Bloom filter (multi-hash)
│   │   │   └── bloom_test.go
│   │   └── bitmask/
│   │       ├── bitmask.go           # Generic bitmask operations
│   │       └── bitmask_test.go
│   │   └── graph/
│   │       ├── graph.go             # Undirected weighted graph (co-rental recommendations)
│   │       └── graph_test.go
│   │
│   ├── models/
│   │   ├── user.go                  # User struct, roles, permissions
│   │   ├── movie.go                 # Movie struct, genres, format (DVD/VHS/Blu-ray)
│   │   ├── rental.go                # Rental struct, due dates, late fees
│   │   ├── wishlist.go              # Wishlist item struct per user
│   │   └── audit.go                 # Audit log entry struct
│   │
│   ├── store/
│   │   ├── store.go                 # BoltDB initialization + bucket setup
│   │   ├── users.go                 # User CRUD operations
│   │   ├── movies.go                # Movie CRUD + search + filter
│   │   ├── rentals.go               # Rental CRUD + history queries
│   │   ├── wishlist.go              # Wishlist add/remove/fetch per user
│   │   ├── sessions.go              # Refresh token storage + revocation
│   │   └── audit.go                 # Audit log append + verify
│   │
│   └── config/
│       └── config.go                # Environment variables, constants, defaults
│
├── api/
│   ├── router.go                    # Chi router setup + route registration
│   ├── middleware.go                # JWT auth, rate-limit, logging, RBAC
│   ├── auth_handler.go              # POST /auth/login, /auth/register, /auth/refresh, /auth/logout
│   ├── movie_handler.go             # GET/POST/PUT/DELETE /movies, GET /movies/search
│   ├── rental_handler.go            # POST /rentals/rent, POST /rentals/return, GET /rentals/history
│   ├── wishlist_handler.go          # GET/POST/DELETE /wishlist
│   ├── user_handler.go              # GET/POST/PUT/DELETE /users (admin only)
│   ├── audit_handler.go             # GET /audit (manager+ only)
│   └── dto.go                       # Request/Response data transfer objects
│
├── tui/
│   ├── app.go                       # Bubbletea Model — top-level application
│   ├── state.go                     # Global session state (token, user, cache)
│   ├── api_client.go                # HTTP client with JWT auto-refresh
│   │
│   ├── styles/
│   │   ├── theme.go                 # Lipgloss color palette, borders, typography
│   │   └── effects.go               # CRT scanlines, glitch frames, rewind animation
│   │
│   ├── pages/
│   │   ├── splash.go                # Animated VHS-style intro (3 seconds)
│   │   ├── login.go                 # Username + password form
│   │   ├── register.go              # Registration form with tier selection
│   │   ├── browse.go                # Main catalog: grid + search + genre tabs
│   │   ├── movie_detail.go          # Synopsis, rating, rent/reserve buttons
│   │   ├── my_rentals.go            # Active rentals, due dates, return flow
│   │   ├── profile.go               # User stats, tier badge, rental history
│   │   ├── admin_users.go           # User management: promote, demote, ban
│   │   ├── admin_movies.go          # Movie management: add, edit, remove
│   │   └── audit_log.go             # Immutable audit chain viewer
│   │
│   └── components/
│       ├── header.go                # ASCII banner "THE LAST VIDEO STORE" + clock + now showing
│       ├── footer.go                # Keybinding hints + status bar
│       ├── searchbar.go             # Input field wired to Trie autocomplete
│       ├── movie_card.go            # Single movie card with poster, stars, format badge
│       ├── movie_grid.go            # Responsive grid layout of movie cards
│       ├── wishlist_sidebar.go      # Personal wishlist panel (add/remove/notify)
│       ├── tabs.go                  # Genre/category tab bar
│       ├── modal.go                 # Confirmation dialogs + "ACCESS DENIED"
│       ├── spinner.go               # Custom VHS-tracking spinner
│       └── badge.go                 # Membership plan badge (Bronze/Silver/Gold, colored)
│
├── data/
│   └── seed.go                      # Seed catalog: ~40 real movies (1980s–2000s)
│
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile.server                # Backend container for Render deployment
├── render.yaml                      # Render Blueprint deploy configuration
└── README.md
```

---

## 6. Implementation Phases

---

### Phase 1 — Data Structures (from scratch)  [x]

**Goal:** Implement all 9 custom data structures with full test coverage and benchmarks. No third-party collections — everything built with raw Go slices, maps, and nodes.

---

#### Task 1.1: Initialize Go module & project skeleton  [x]

- [x] Create project root directory `thelastvideostore/`
- [x] Run `go mod init github.com/thelastvideostore` (or local module path)
- [x] Create directory tree: `internal/ds/{trie,lru,deque,heap,list,bloom,bitmask,graph}`, `internal/crypto`
- [x] Create `cmd/server/` and `cmd/client/` placeholder `main.go` files (just `package main; func main() {}`)
- [x] Verify: `go build ./...` compiles without errors

#### Task 1.2: Implement Bitmask (6-bit)  [x]

- [x] Create `internal/ds/bitmask/bitmask.go`
- [x] Define `type Permission uint16`
- [x] Implement functions: `Has(p, flag Permission) bool`, `Set(p, flag Permission) Permission`, `Clear(p, flag Permission) Permission`, `Toggle(p, flag Permission) Permission`
- [x] Define 6 permission constants: `PermBrowse = 0b000001`, `PermRent = 0b000010`, `PermReserve = 0b000100`, `PermManageUsers = 0b001000`, `PermStaff = 0b010000`, `PermAdmin = 0b100000`
- [x] Define tier constants: `TierBronze = PermBrowse`, `TierSilver = PermBrowse | PermRent`, `TierGold = PermBrowse | PermRent | PermReserve`, `TierEmployee = PermBrowse | PermRent | PermReserve | PermStaff`, `TierSupervisor = PermBrowse | PermRent | PermReserve | PermManageUsers | PermStaff`, `TierManager = PermBrowse | PermRent | PermReserve | PermManageUsers | PermStaff | PermAdmin`, `TierOwner = TierManager` (same bits, symbolic difference)
- [x] Define tier labels map: `map[Permission]string` — "Bronze", "Silver", "Gold", "Employee", "Supervisor", "Manager", "Owner"
- [x] Create `internal/ds/bitmask/bitmask_test.go`
- [x] Test table: `struct{ name string; base Permission; flag Permission; wantHas bool; wantSet Permission }`
- [x] Test all combos: Bronze missing Rent, Owner has Admin, set/clear/toggle operations, Gold != Employee (Staff bit distinguishes them)
- [x] Benchmark: `Has`, `Set`, `Clear` — confirm O(1) performance
- [x] Verify: `go test -v -bench=. ./internal/ds/bitmask/`

#### Task 1.3: Implement Doubly Linked List  [x]

- [x] Create `internal/ds/list/linkedlist.go`
- [x] Define generic `Node[T any]` struct with `Value T`, `Prev *Node[T]`, `Next *Node[T]`
- [x] Define `List[T any]` struct with `Head *Node[T]`, `Tail *Node[T]`, `Len int`
- [x] Implement: `New[T]()`, `PushBack(v T) *Node[T]`, `PushFront(v T) *Node[T]`, `Remove(node *Node[T]) T`, `PopFront() T`, `PopBack() T`, `Find(pred func(T) bool) *Node[T]`, `Slice() []T`
- [x] Ensure proper nil handling on empty list operations (return zero value + false ok pattern)
- [x] Create `internal/ds/list/linkedlist_test.go`
- [x] Test: push/remove interleaved, pop front/back, find by predicate, iteration via Slice, empty list edge cases
- [x] Benchmark: `PushBack` x 10000, `Find` on 1000-element list
- [x] Verify: `go test -v -bench=. ./internal/ds/list/`

#### Task 1.4: Implement Deque (Ring Buffer)  [x]

- [x] Create `internal/ds/deque/deque.go`
- [x] Define `Deque[T any]` struct with ring buffer backing slice: `buf []T`, `head int`, `tail int`, `size int`, `cap int`
- [x] Implement: `New[T](capacity int)`, `PushBack(v T)`, `PushFront(v T)`, `PopFront() T`, `PopBack() T`, `PeekFront() T`, `PeekBack() T`, `Len() int`, `IsEmpty() bool`
- [x] Auto-grow: double capacity when full, handle wrap-around indices correctly
- [x] Return `(T, bool)` for pop/peek on empty deque
- [x] Create `internal/ds/deque/deque_test.go`
- [x] Test: push/pop interleaved mixes, empty deque pops, wrap-around after capacity, grow triggers, peek non-destructive
- [x] Benchmark: push/pop 10000 items from both ends
- [x] Verify: `go test -v -bench=. ./internal/ds/deque/`

#### Task 1.5: Implement Min-Heap  [x]

- [x] Create `internal/ds/heap/heap.go`
- [x] Define generic `Heap[T any]` struct with `items []T`, `less func(a, b T) bool` (min-heap by default)
- [x] Implement: `New[T](lessFn)`, `Push(v T)`, `Pop() T`, `Peek() T`, `Len() int`, `IsEmpty() bool`
- [x] Internal: `siftUp(index int)`, `siftDown(index int)`
- [x] Create `internal/ds/heap/heap_test.go`
- [x] Test: sequential push → pop returns sorted order, empty heap pop returns zero value, Peek on empty, priority ordering via custom `less` (e.g., structs ordered by timestamp)
- [x] Benchmark: push 10000 then pop 10000
- [x] Verify: `go test -v -bench=. ./internal/ds/heap/`

#### Task 1.6: Implement Trie (Prefix Tree)  [x]

- [x] Create `internal/ds/trie/trie.go`
- [x] Define `TrieNode` struct: `children map[rune]*TrieNode`, `isEnd bool`, `value any` (store movie ID or title)
- [x] Define `Trie` struct: `root *TrieNode`
- [x] Implement: `New()`, `Insert(word string, value any)`, `Search(word string) (any, bool)`, `StartsWith(prefix string) bool`, `Autocomplete(prefix string) []any` (returns all values under prefix subtree via DFS), `Delete(word string) bool`
- [x] Create `internal/ds/trie/trie_test.go`
- [x] Test: insert + search exact match, search missing, startsWith true/false, autocomplete returns all matches for prefix "mat", delete removes exact word but not prefixes, case sensitivity (lowercase only enforced)
- [x] Benchmark: insert 5000 words, autocomplete with 2-char prefix
- [x] Verify: `go test -v -bench=. ./internal/ds/trie/`

#### Task 1.7: Implement LRU Cache  [x]

- [x] Create `internal/ds/lru/lru.go`
- [x] Define `entry[K comparable, V any]` struct: `key K`, `value V`
- [x] Define `Cache[K comparable, V any]` struct: `capacity int`, `items map[K]*list.Node[entry[K,V]]`, `order *list.List[entry[K,V]]` (reuses our DoublyLinkedList from Task 1.3)
- [x] Implement: `New[K,V](capacity int)`, `Get(key K) (V, bool)`, `Put(key K, value V)`, `Remove(key K) bool`, `Len() int`, `Contains(key K) bool`
- [x] On Get: move node to front of order (most recently used)
- [x] On Put when full: evict tail of order (least recently used)
- [x] Create `internal/ds/lru/lru_test.go`
- [x] Test: put+get returns value, get missing returns zero, eviction when capacity exceeded (put 4 items capacity 3 → first evicted), update existing key moves to front, remove works, Contains works
- [x] Benchmark: put 10000 items with capacity 1000, get hit rate test
- [x] Verify: `go test -v -bench=. ./internal/ds/lru/`

#### Task 1.8: Implement Bloom Filter  [x]

- [x] Create `internal/ds/bloom/bloom.go`
- [x] Define `BloomFilter` struct: `bitset []uint64`, `size uint64` (total bits), `hashCount int`
- [x] Implement: `New(size uint64, hashCount int)`, `Add(data []byte)`, `Contains(data []byte) bool`
- [x] Use double-hashing technique: `h1 = fnv.New64a`, `h2 = fnv.New64` (or murmurhash via hash/fnv + shift)
- [x] Combined hash: `h1 + i*h2` for i from 0 to hashCount-1
- [x] Set/check bits via `bitset[bitIndex/64] & (1 << (bitIndex % 64))`
- [x] Create `internal/ds/bloom/bloom_test.go`
- [x] Test: add string → contains returns true, empty filter contains nothing, multiple adds don't collide, estimate false-positive rate on 1000 items with 10000-bit filter and 3 hash functions
- [x] Benchmark: Add 10000 items, Contains 10000 items
- [x] Verify: `go test -v -bench=. ./internal/ds/bloom/`

#### Task 1.9: Implement Hash Chain (Audit Trail)  [x]

- [x] Create `internal/crypto/hashchain.go`
- [x] Define `HashChain` struct: `entries []HashChainEntry`, `lastHash []byte`
- [x] Define `HashChainEntry` struct: `Timestamp int64`, `Action string`, `ActorID string`, `TargetID string`, `Data string`, `Hash []byte`, `PrevHash []byte`
- [x] Implement: `New()`, `Append(action, actorID, targetID, data string) HashChainEntry`, `Verify() bool` (recomputes all hashes from genesis), `GetAll() []HashChainEntry`, `Len() int`
- [x] Hash function: `SHA-256(prevHash || timestamp || action || actorID || targetID || data)`
- [x] Genesis block: first entry has `PrevHash = []byte("GENESIS")`
- [x] Create `internal/crypto/hashchain_test.go`
- [x] Test: genesis entry has correct PrevHash, append links correctly, verify passes on intact chain, verify fails if entry tampered, middle insertion detected
- [x] Benchmark: append 1000 entries + verify
- [x] Verify: `go test -v -bench=. ./internal/crypto/hashchain.go` (or move test to same package)

#### Task 1.10: Implement Undirected Weighted Graph (Co-rental Recommendations)  [x]

- [x] Create `internal/ds/graph/graph.go`
- [x] Define `Graph` struct: `vertices map[string]*Vertex`
- [x] Define `Vertex` struct: `ID string`, `Edges map[string]int` (neighbor ID → co-rental weight), `Data interface{}`
- [x] Implement: `New()`, `AddVertex(id string)`, `AddEdge(v1, v2 string)`, `IncrementEdge(v1, v2 string)` (increments edge weight by 1), `GetNeighbors(id string) map[string]int`, `BFS(start string) []string` (visit order), `GetRecommendations(id string, k int) []string` (top-k neighbors by edge weight, excluding self and already-connected zero-weight), `HasVertex(id string) bool`, `VertexCount() int`, `EdgeCount() int`
- [x] Application: when building co-rental graph, `IncrementEdge(movieA, movieB)` for every rental pair, building up co-rental counts
- [x] Create `internal/ds/graph/graph_test.go`
- [x] Test: add vertices and edges, increment weights, get neighbors sorted by weight, BFS visits correct order, GetRecommendations returns top-k by weight, empty graph operations return gracefully
- [x] Benchmark: build graph with 100 vertices and 500 edges, GetRecommendations on dense vertex
- [x] Verify: `go test -v -bench=. ./internal/ds/graph/`

#### Phase 1 validation:  [x]

```bash
go test -v -race -bench=. ./internal/ds/... ./internal/crypto/...
# All tests must pass, all benchmarks must complete, no race conditions
```

---

### Phase 2 — Models & Database Layer  [x]

**Goal:** Define all data models as Go structs and implement full CRUD persistence with BoltDB.

---

#### Task 2.1: Install dependencies  [x]

- [x] Run `go get github.com/boltdb/bolt` (or `go.etcd.io/bbolt` for maintained fork)
- [x] Run `go get github.com/google/uuid`
- [x] Run `go get golang.org/x/crypto`
- [x] Run `go get github.com/go-chi/chi/v5`
- [x] Run `go get github.com/golang-jwt/jwt/v5`
- [x] Verify: `go mod tidy` succeeds

#### Task 2.2: Create config package  [x]

- [x] Create `internal/config/config.go`
- [x] Define `Config` struct: `DBPath string`, `JWTSecret string`, `AESKey string`, `ServerPort string`, `APIBaseURL string`
- [x] Implement `Load() *Config`: reads from env vars with sensible defaults
- [x] Defaults: `DBPath="thelastvideostore.db"`, `ServerPort="8080"`, `APIBaseURL="http://localhost:8080"`
- [x] Add `MustLoad()` variant that panics on missing required vars

#### Task 2.3: Create user model  [x]

- [x] Create `internal/models/user.go`
- [x] Define `User` struct: `ID`, `Username`, `PasswordHash`, `Tier` (Permission bitmask), `MaxRentals` int, `RentalCount` int, `Banned` bool, `TOTPEnabled` bool, `TOTPSecret` string (AES encrypted at rest), `CreatedAt` int64, `UpdatedAt` int64
- [x] Define JSON tags for API serialization: `json:"id"`, `json:"username"`, `json:"tier"`, `json:"max_rentals"`, `json:"rental_count"`, `json:"banned"`, `json:"totp_enabled"` (never expose `PasswordHash` or `TOTPSecret`)
- [x] Define `UserResponse` struct (omits password hash and TOTP secret for API responses)
- [x] Define helper: `CanRent() bool`, `CanReserve() bool`, `TierName() string`, `HasStaffAccess() bool` (checks `PermStaff` bit)

#### Task 2.4: Create movie model  [x]

- [x] Create `internal/models/movie.go`
- [x] Define `Movie` struct: `ID`, `Title`, `Year` int, `Genre` string, `Format` string (VHS, DVD, Blu-ray), `Director`, `Cast` []string, `Synopsis` string, `Rating` float64 (avg), `RatingCount` int, `Available` bool, `CopiesTotal` int, `CopiesAvailable` int, `IsNewRelease` bool, `CoverArt` string (ASCII art placeholder string), `CreatedAt` int64
- [x] Define format constants: `FormatVHS`, `FormatDVD`, `FormatBluRay`
- [x] Define genre constants: `Action`, `Comedy`, `Horror`, `SciFi`, `Drama`, `Thriller`, `Romance`, `Animation`
- [x] Define `MovieResponse` DTO struct for API

#### Task 2.5: Create rental model  [x]

- [x] Create `internal/models/rental.go`
- [x] Define `Rental` struct: `ID`, `UserID`, `MovieID`, `MovieFormat` string, `RentedAt` int64, `DueDate` int64, `ReturnedAt` int64 (0 = not returned), `LateFee` float64, `RewindFee` float64, `NeedsRewind` bool (30% chance for VHS, set at rental time), `Status` string (active, returned, overdue)
- [x] Define rental status constants: `RentalActive`, `RentalReturned`, `RentalOverdue`
- [x] Define helper: `IsOverdue(now int64) bool`, `CalculateLateFee(now int64) float64` (uses format-specific daily rate: VHS=$2/day, DVD/Blu-ray=$3/day)
- [x] Define helper: `DueDateForFormat(format string, rentedAt int64) int64` (VHS: +3 days, DVD/Blu-ray: +5 days)
- [x] Define helper: `CalculateRewindFee() float64` — returns $1.00 if `NeedsRewind && MovieFormat == FormatVHS`, else $0.00
- [x] Define helper: `TotalFee() float64` — `LateFee + RewindFee`

#### Task 2.5a: Create wishlist model  [x]

- [x] Create `internal/models/wishlist.go`
- [x] Define `WishlistItem` struct: `ID`, `UserID`, `MovieID`, `AddedAt` int64
- [x] Using doubly linked list structure for ordered storage per user

#### Task 2.6: Create audit model  [x]

- [x] Create `internal/models/audit.go`
- [x] Define `AuditEntry` struct: `ID`, `Timestamp` int64, `Action` string, `ActorID`, `TargetID`, `Data`, `Hash`, `PrevHash` (mirrors HashChainEntry for DB persistence)
- [x] Define action constants: `ActionLogin`, `ActionLogout`, `ActionRent`, `ActionReturn`, `ActionRegister`, `ActionPromote`, `ActionDemote`, `ActionBan`, `ActionAddMovie`, `ActionEditMovie`, `ActionDeleteMovie`, `ActionTOTPEnabled`, `ActionTOTPDisabled`, `ActionAddToWishlist`, `ActionRemoveFromWishlist`, `ActionAddStaffPick`, `ActionRemoveStaffPick`

#### Task 2.7: Create BoltDB store layer  [x]

- [x] Create `internal/store/store.go`
- [x] Define `Store` struct wrapping `*bolt.DB`
- [x] Implement `Open(path string) (*Store, error)` — opens BoltDB, creates all buckets: `users`, `movies`, `rentals`, `audit_logs`, `sessions`, `banned`, `movies_by_genre`, `movies_by_title`, `staff_picks`, `wishlists`, `totp_secrets`
- [x] Implement `Close() error`
- [x] Implement helper: `bucketName` constants

#### Task 2.8: Implement user store  [x]

- [x] Create `internal/store/users.go`
- [x] Methods on `*Store`:
  - `CreateUser(user *models.User) error` — serializes to JSON, stores in `users` bucket by ID, also indexes by username in same bucket via `username:<name>` key
  - `GetUserByID(id string) (*models.User, error)`
  - `GetUserByUsername(username string) (*models.User, error)` — looks up `username:<name>` key to get ID, then fetches user
  - `UpdateUser(user *models.User) error` — overwrites entire user JSON
  - `DeleteUser(id string) error` — removes user and username index
  - `ListUsers() ([]*models.User, error)` — iterates all keys, skips username index keys
  - `UserExists(username string) bool`
  - `IncrementFailedAttempts(username string) (int, error)` — stores attempts in `sessions` bucket
  - `ResetFailedAttempts(username string) error`
  - `LockUserUntil(username string, until int64) error` — stores lock expiry in `sessions` bucket
  - `IsUserLocked(username string) (bool, error)` — checks lock expiry

#### Task 2.9: Implement movie store  [x]

- [x] Create `internal/store/movies.go`
- [x] Methods on `*Store`:
  - `CreateMovie(movie *models.Movie) error` — stores in `movies` bucket, adds to `movies_by_genre` index, adds to `movies_by_title` index
  - `GetMovieByID(id string) (*models.Movie, error)`
  - `UpdateMovie(movie *models.Movie) error` — updates in all indexes
  - `DeleteMovie(id string) error` — removes from all indexes
  - `ListMovies(genre string, offset, limit int) ([]*models.Movie, int, error)` — paginated; if genre != "" filters by genre index; returns total count
  - `SearchMoviesByPrefix(prefix string, limit int) ([]*models.Movie, error)` — uses `movies_by_title` bucket for prefix scan (BoltDB supports prefix iteration via `Seek`)
  - `GetNewReleases() ([]*models.Movie, error)` — filters by `IsNewRelease = true` (scan all, filter in memory)

#### Task 2.10: Implement rental store  [x]

- [x] Create `internal/store/rentals.go`
- [x] Methods on `*Store`:
  - `CreateRental(rental *models.Rental) error` — stores in `rentals` bucket
  - `GetRentalByID(id string) (*models.Rental, error)`
  - `UpdateRental(rental *models.Rental) error`
  - `GetActiveRentalsByUser(userID string) ([]*models.Rental, error)` — iterates all, filters by `UserID` and `Status != returned`
  - `GetRentalHistoryByUser(userID string) ([]*models.Rental, error)` — iterates all, filters by `UserID`, returns all active + returned
  - `GetOverdueRentals() ([]*models.Rental, error)` — iterates all, filters by `DueDate < now && Status == active`
  - `CountActiveRentalsByUser(userID string) (int, error)` — count of non-returned rentals

#### Task 2.11: Implement audit store  [x]

- [x] Create `internal/store/audit.go`
- [x] Methods: `AppendAuditEntry(entry *models.AuditEntry) error`, `GetAllAuditEntries() ([]*models.AuditEntry, error)`, `GetAuditEntriesByUser(userID string) ([]*models.AuditEntry, error)`
- [x] Each append: encrypt entry data with AES before storing (call `crypto.AESEncrypt`), store encrypted blob in BoltDB

#### Task 2.12: Implement session store  [x]

- [x] Extend `internal/store/store.go` or create `internal/store/sessions.go`
- [x] Methods:
  - `SaveRefreshToken(userID, tokenID string, expiresAt int64) error` — stores in `sessions` bucket
  - `ValidateRefreshToken(userID, tokenID string) (bool, error)` — checks if token exists and not expired
  - `InvalidateRefreshToken(userID, tokenID string) error` — removes from bucket
  - `InvalidateAllUserSessions(userID string) error` — removes all tokens for user
  - `IsTokenRevoked(tokenID string) (bool, error)` — checks `revoked` sub-bucket

#### Task 2.12a: Implement wishlist store  [x]

- [x] Create `internal/store/wishlist.go`
- [x] Methods:
  - `AddToWishlist(userID, movieID string) error` — appends to user's wishlist in BoltDB
  - `RemoveFromWishlist(userID, movieID string) error`
  - `GetWishlist(userID string) ([]*models.WishlistItem, error)` — returns ordered items
  - `IsInWishlist(userID, movieID string) (bool, error)`
  - `GetWishlistSize(userID string) (int, error)`
- [x] Wishlist stored as ordered entries per user in `wishlists` bucket

#### Task 2.12b: Implement staff picks & last chance store  [x]

- [x] Extend `internal/store/movies.go` or create `internal/store/staffpicks.go`
- [x] Methods:
  - `AddStaffPick(movieID string) error` — adds movie ID to `staff_picks` bucket
  - `RemoveStaffPick(movieID string) error`
  - `GetStaffPicks() ([]*models.Movie, error)` — resolves IDs to full movie objects
  - `IsStaffPick(movieID string) bool`
  - `GetLastChanceMovies() ([]*models.Movie, error)` — queries movies where `CopiesAvailable == 1 && !IsNewRelease`
- [x] Staff Picks bucket keyed by movie ID, value = timestamp of when it was picked

#### Task 2.12c: Implement TOTP store operations  [x]

- [x] Extend `internal/store/store.go` or create `internal/store/totp.go`
- [x] Methods:
  - `SaveTOTPSecret(userID string, encryptedSecret []byte) error` — stores in `totp_secrets` bucket, AES encrypted
  - `GetTOTPSecret(userID string) ([]byte, error)` — retrieves encrypted secret
  - `DeleteTOTPSecret(userID string) error`
  - `IncrementTOTPFailures(userID string) (int, error)` — TOTP-specific lockout after 3 failures
  - `ResetTOTPFailures(userID string) error`
  - `LockTOTPUserUntil(userID string, until int64) error` — 10-minute TOTP lockout

#### Phase 2 validation:  [x]

```bash
go test -v ./internal/models/... ./internal/store/...
# Write a store_test.go that opens temp BoltDB, tests CRUD for each entity, closes/deletes temp
```

---

### Phase 3 — Authentication & Security Layer  [x]

**Goal:** Implement bcrypt password hashing, JWT token management with refresh rotation, RBAC permission enforcement, brute-force lockout, and AES-256-GCM encryption.

---

#### Task 3.1: Implement password hashing  [x]

- [x] Create `internal/auth/password.go`
- [x] Implement `HashPassword(password string) (string, error)` — uses `bcrypt.GenerateFromPassword` with cost 12
- [x] Implement `CheckPassword(hash, password string) bool` — uses `bcrypt.CompareHashAndPassword`
- [x] Create `internal/auth/password_test.go`
- [x] Test: hash produces different string each time, verify same password matches, wrong password doesn't match, empty password rejected

#### Task 3.2: Implement JWT session management  [x]

- [x] Create `internal/auth/session.go`
- [x] Define `TokenPair` struct: `AccessToken string`, `RefreshToken string`, `ExpiresAt int64`
- [x] Implement `GenerateTokenPair(userID string, permissions Permission, secret string) (*TokenPair, error)`:
  - Access token: 15-min expiry, claims: `sub=userID`, `perm=permissions`, `iat`, `exp`, `jti` (unique ID)
  - Refresh token: 7-day expiry, claims: `sub=userID`, `jti`, `exp`
  - Sign with HS256
- [x] Implement `ValidateAccessToken(tokenString, secret string) (*Claims, error)` — parses, validates expiry, returns claims
- [x] Implement `ValidateRefreshToken(tokenString, secret string) (*RefreshClaims, error)`
- [x] Create `internal/auth/session_test.go`
- [x] Test: generate valid tokens, validate expired token (manipulate time or create token with -1s expiry), wrong secret fails, malformed token fails

#### Task 3.3: Implement permission enforcement  [x]

- [x] Create `internal/auth/permissions.go`
- [x] Re-export bitmask constants from `internal/ds/bitmask` (or import directly — decide which package owns these)
- [x] Define `RequirePermission(userPerms Permission, required Permission) bool` — simple `userPerms & required != 0`
- [x] Define `TierName(perm Permission) string` — returns "Bronze", "Silver", "Gold", "Employee", "Supervisor", "Manager", "Owner"
- [x] Define `MaxRentalsForTier(perm Permission) int`:
  - Bronze: 1, Silver: 2, Gold: 5, Employee: 5, Supervisor: 5, Manager: 10, Owner: MaxInt
- [x] Define `CanAccessAdmin(perm Permission) bool` — Manager or Owner
- [x] Define `IsStaff(perm Permission) bool` — checks `PermStaff` bit (Employee, Supervisor, Manager, Owner)
- [x] Define `CanManageUsers(perm Permission) bool` — checks `PermManageUsers` bit (Supervisor, Manager, Owner)

#### Task 3.4: Implement AES-256-GCM encryption  [x]

- [x] Create `internal/crypto/aes.go`
- [x] Implement `GenerateAESKey() ([]byte, error)` — 32 random bytes via `crypto/rand`
- [x] Implement `Encrypt(plaintext, key []byte) ([]byte, error)` — AES-256-GCM: generate random nonce, prepend to ciphertext
- [x] Implement `Decrypt(ciphertext, key []byte) ([]byte, error)` — extract nonce, decrypt
- [x] Create `internal/crypto/aes_test.go`
- [x] Test: encrypt then decrypt returns original, different key fails, empty plaintext works, tampered ciphertext detected (GCM auth failure)

#### Task 3.5: Implement brute-force lockout  [x]

- [x] Create `internal/auth/lockout.go` (or add to `session.go`)
- [x] Define constants: `MaxAttempts=5`, `LockoutDuration=30*time.Minute`
- [x] Implement `CheckLoginAttempts(store *store.Store, username string) error`:
  - Check if locked → return `ErrAccountLocked` with remaining time
  - Check if attempt count >= MaxAttempts → lock account (save lock expiry), return error
  - If clear → return nil
- [x] Implement `RecordFailedAttempt(store *store.Store, username string) error`
- [x] Implement `RecordSuccessfulLogin(store *store.Store, username string) error` — resets attempts
- [x] Integrate into login handler (Phase 4)

#### Task 3.6: Implement audit log integration  [x]

- [x] Create `internal/auth/audit.go` (or add to `internal/crypto/hashchain.go`)
- [x] Connect hash chain to BoltDB store: every state-changing operation appends an entry
- [x] Encrypt audit entries with AES before persisting (call `Encrypt` from Task 3.4)
- [x] Implement `VerifyAuditChain(store *store.Store) (bool, error)` — reads all entries, recomputes hashes, compares

#### Task 3.7: Implement TOTP 2FA (optional, Manager+ feature)  [x]

- [x] Create `internal/auth/totp.go`
- [x] Implement using only Go stdlib (`crypto/hmac`, `crypto/sha1`, `crypto/rand`, `encoding/base32`, `time`)
- [x] `GenerateTOTPSecret() (string, error)` — generates 20 random bytes, returns base32-encoded string (e.g., `"JBSWY3DPEHPK3PXP"`)
- [x] `GenerateTOTPCode(secret string, t time.Time) (string, error)` — HMAC-SHA1 of (counter = unix/30), returns 6-digit code per RFC 6238
- [x] `ValidateTOTPCode(secret string, code string) bool` — checks current code ± 1 interval (30s skew tolerance)
- [x] `GenerateTOTPURL(issuer, accountName, secret string) string` — returns `otpauth://totp/...` URL for QR generation
- [x] Create `internal/auth/totp_test.go`
- [x] Test: generate secret, generate code, validate same code, reject expired code, reject wrong code, consistent output for same time step
- [x] Integration: on login, if user has `TOTPEnabled`, after password validation prompt for TOTP code before issuing tokens
- [x] Profile page: Manager+ can enable/disable TOTP, view setup key, verify setup with one test code

#### Phase 3 validation:  [x]

```bash
go test -v ./internal/auth/... ./internal/crypto/...
# JWT round-trip, bcrypt correctness, AES encrypt/decrypt, lockout timer, hash chain integrity
```

---

### Phase 4 — REST API  [x]

**Goal:** Build a complete REST API with Chi router, middleware, handlers, DTOs, and server entrypoint.

---

#### Task 4.1: Create DTOs  [x]

- [x] Create `api/dto.go`
- [x] Request DTOs:
  - `LoginRequest`: `Username string`, `Password string`
  - `RegisterRequest`: `Username string`, `Password string`
  - `CreateMovieRequest`: `Title string`, `Year int`, `Genre string`, `Director string`, `Cast []string`, `Synopsis string`, `CopiesTotal int`, `IsNewRelease bool`
  - `UpdateMovieRequest`: same fields as Create, all optional via pointers
  - `RentRequest`: `MovieID string`
  - `ReturnRequest`: `RentalID string`
  - `UpdateUserRequest`: `Tier string` (optional), `Banned *bool` (optional)
- [x] Response DTOs (structs with JSON tags):
  - `ErrorResponse`: `Error string`, `Code int`
  - `SuccessResponse`: `Message string`
  - `LoginResponse`: `AccessToken string`, `RefreshToken string`, `User models.UserResponse`
  - `MovieListResponse`: `Movies []models.MovieResponse`, `Total int`, `Page int`, `PageSize int`
- [x] Implement `WriteJSON(w http.ResponseWriter, status int, data interface{})`
- [x] Implement `WriteError(w http.ResponseWriter, status int, message string)`

#### Task 4.2: Implement authentication middleware  [x]

- [x] Create `api/middleware.go`
- [x] `AuthMiddleware(secret string, store *store.Store)`:
  - Extract `Authorization: Bearer <token>` header
  - Validate JWT (call `ValidateAccessToken`)
  - Check if token revoked (query `sessions` bucket for revoked JTI)
  - Load user from store by `sub` claim
  - Check user not banned (also check Bloom filter as fast path)
  - Inject `User` + `Permissions` into `context.Context`
  - Return 401 if missing/invalid token || 403 if banned
- [x] `RequirePermission(required Permission)` — middleware factory:
  - Reads permissions from context
  - Calls `bitmask.Has(perms, required)`
  - Returns 403 with `"⛔ ACCESS DENIED — Insufficient clearance"` if check fails
- [x] `RequireStaff()` — middleware requiring `PermStaff`:
  - Used on return-processing endpoints where staff can return any customer's rentals
- [x] `TOTPMiddleware()` — if user has `TOTPEnabled`, validates TOTP header `X-TOTP-Code` on login step 2
- [x] `RateLimitMiddleware(rate int)` — token bucket:
  - Per-IP counting via in-memory `map[string]*tokenBucket` with mutex
  - 100 req/min default
  - Returns 429 if exceeded
- [x] `CORSMiddleware()` — wraps `chi/cors` with permissive dev defaults
- [x] `LoggingMiddleware()` — logs method, path, status, duration to stdout

#### Task 4.3: Implement auth handlers  [x]

- [x] Create `api/auth_handler.go`
- [x] `POST /api/v1/auth/register`:
  - Parse `RegisterRequest`
  - Validate: username 3-20 chars alphanumeric, password 6+ chars
  - Check `UserExists` → 409 if taken
  - Hash password via `auth.HashPassword`
  - Create User with Tier=Bronze (Cliente Bronze), save to store
  - Append to audit hash chain: `ActionRegister`
  - Return 201 with `UserResponse`
- [x] `POST /api/v1/auth/login`:
  - Call `CheckLoginAttempts` → 429 if locked
  - Find user by username → 401 if not found
  - `CheckPassword` → if fail: `RecordFailedAttempt`, return 401
  - If user has `TOTPEnabled == true`:
    - Require `X-TOTP-Code` header
    - Decrypt TOTP secret from store
    - Validate code via `auth.ValidateTOTPCode` → if fail: return 401 `"Invalid TOTP code"`
  - `RecordSuccessfulLogin` → resets attempts
  - Check banned (Bloom + DB) → 403
  - Generate `TokenPair` via `auth.GenerateTokenPair`
  - Save refresh token to store (`SaveRefreshToken`)
  - Append to audit: `ActionLogin`
  - Return `LoginResponse` with tokens + user (include `totp_required: true/false` during initial password auth so the TUI knows to prompt)
- [x] `POST /api/v1/auth/login/totp`:
  - Temporary session token from step 1 (5-min expiry, no full access)
  - Accept `{code: "123456"}`
  - Validate TOTP → if valid, issue real token pair
  - If invalid → increment TOTP failure counter (lock after 3 TOTP failures)
- [x] `POST /api/v1/auth/refresh`:
  - Accept refresh token from body
  - Validate, check not revoked
  - Invalidate old refresh token (rotation)
  - Generate new token pair
  - Save new refresh token
  - Return new `LoginResponse`
- [x] `POST /api/v1/auth/logout`:
  - Requires JWT auth
  - Invalidate refresh token (or all user sessions)
  - Append to audit: `ActionLogout`
  - Return 200

#### Task 4.4: Implement movie handlers  [x]

- [x] Create `api/movie_handler.go`
- [x] `GET /api/v1/movies`:
  - Query params: `genre`, `page` (default 1), `page_size` (default 20)
  - Call `ListMovies(genre, offset, limit)`
  - Return `MovieListResponse`
- [x] `GET /api/v1/movies/search?q=<prefix>`:
  - Requires JWT (any tier)
  - Call `SearchMoviesByPrefix(q, 10)` — uses BoltDB prefix scan on `movies_by_title`
  - Return `[]models.MovieResponse`
- [x] `GET /api/v1/movies/staff-picks`:
  - Requires JWT
  - Returns movies curated by Manager/Owner via dedicated BoltDB bucket `staff_picks`
  - Store picks as movie IDs; resolve to full movie objects on read
- [x] `GET /api/v1/movies/last-chance`:
  - Requires JWT
  - Returns movies where `CopiesAvailable == 1 && !IsNewRelease` — titles about to leave catalog
- [x] `GET /api/v1/movies/{id}`:
  - Call `GetMovieByID(id)` → 404 if not found
  - Return `MovieResponse`
- [x] `POST /api/v1/movies`:
  - Requires `RequirePermission(PermAdmin)` (Manager+)
  - Parse `CreateMovieRequest`
  - Validate: title required, year 1900–current, valid genre
  - Create Movie with UUID, `CopiesAvailable = CopiesTotal`, `Available = true`
  - Save to store
  - Append audit: `ActionAddMovie`
  - Return 201 with `MovieResponse`
- [x] `PUT /api/v1/movies/{id}`:
  - Requires `PermAdmin` (Manager+)
  - Parse `UpdateMovieRequest`, apply partial updates (only set non-nil fields)
  - Append audit: `ActionEditMovie`
  - Return updated `MovieResponse`
- [x] `DELETE /api/v1/movies/{id}`:
  - Requires `PermAdmin` (Manager+)
  - Delete from store
  - Append audit: `ActionDeleteMovie`
  - Return 200
- [x] `POST /api/v1/movies/{id}/staff-pick`:
  - Requires `PermAdmin` (Manager+)
  - Adds movie ID to `staff_picks` bucket
  - Return 200 with `{staff_pick: true}`
- [x] `DELETE /api/v1/movies/{id}/staff-pick`:
  - Requires `PermAdmin` (Manager+)
  - Removes movie ID from `staff_picks` bucket

#### Task 4.5: Implement rental handlers  [x]

- [x] Create `api/rental_handler.go`
- [x] `POST /api/v1/rentals/rent`:
  - Requires `RequirePermission(PermRent)`
  - Parse `RentRequest`
  - Get user from context
  - Count active rentals → if >= `MaxRentalsForTier(user.Tier)` → 403 "Rental limit reached"
  - Get movie → 404 if not found
  - Check `movie.CopiesAvailable > 0` → 409 "No copies available"
  - Check `movie.IsNewRelease` and user doesn't have `PermReserve` → 403 "Gold plan required for new releases"
  - Create Rental: `DueDate = DueDateForFormat(movie.Format, now)` (VHS: +3d, DVD/Blu-ray: +5d), `Status = active`
  - If `movie.Format == FormatVHS`: randomly set `NeedsRewind = true` (30% chance) to simulate unrewound tape
  - Decrement `movie.CopiesAvailable`; if 0 → set `Available = false`
  - Increment `user.RentalCount`
  - Update movie + user in store
  - Save rental
  - Append audit: `ActionRent`
  - Return rental with due date (+ rewind flag if VHS)
- [x] `POST /api/v1/rentals/return`:
  - Requires `PermRent`
  - Parse `ReturnRequest`
  - Get rental → 404
  - Verify rental belongs to user (or user has `PermStaff` — Employee, Supervisor, Manager, Owner)
  - Set `ReturnedAt = now`
  - If `now > DueDate` → calculate `LateFee = daysLate × dailyRate` where `dailyRate = $2.00` for VHS, `$3.00` for DVD/Blu-ray
  - If `NeedsRewind` → set `RewindFee = 1.00` (VHS rewind fee)
  - Calculate `TotalFee = LateFee + RewindFee`; if >0, display breakdown to user
  - Increment `movie.CopiesAvailable`; set `Available = true`
  - Decrement `user.RentalCount`
  - Save rental, movie, user
  - Append audit: `ActionReturn`
  - Return rental with fee breakdown (late fee + rewind fee)
- [x] `GET /api/v1/rentals/history`:
  - Requires JWT
  - Call `GetRentalHistoryByUser(userID)`
  - Return list of rentals with movie data joined

#### Task 4.6: Implement user handlers (admin)  [x]

- [x] Create `api/user_handler.go`
- [x] `GET /api/v1/users`:
  - Requires Supervisor+ (`PermManageUsers`)
  - Call `ListUsers()`
  - Return list (omit password hashes)
- [x] `POST /api/v1/users`:
  - Requires Supervisor+ (`PermManageUsers`)
  - Parse `RegisterRequest` + optional tier
  - Same validation as register, but can set initial tier
  - Append audit: action depending on tier
  - Return 201
- [x] `PUT /api/v1/users/{id}`:
  - Requires Supervisor+ (`PermManageUsers`)
  - Parse `UpdateUserRequest`
  - If tier changed → append audit `ActionPromote`/`ActionDemote`
  - If banned → append audit `ActionBan`, add to Bloom filter
  - Save user
  - Return updated user
- [x] `DELETE /api/v1/users/{id}`:
  - Requires Manager+ (`PermAdmin`)
  - Delete user from store
  - Append audit
  - Return 200
- [x] `POST /api/v1/users/{id}/totp`:
  - Requires the user themselves OR Manager+
  - `{enabled: true}` → generate TOTP secret, store AES-encrypted, return secret & otpauth URL for QR setup
  - `{enabled: false}` → clear TOTP secret, disable
  - Append audit: `ActionTOTPEnabled` / `ActionTOTPDisabled`

#### Task 4.7: Implement wishlist handler  [x]

- [x] Create `api/wishlist_handler.go`
- [x] `GET /api/v1/wishlist` — requires JWT (Bronze+) — returns user's wishlist with movie details
- [x] `POST /api/v1/wishlist` — requires JWT (Bronze+) — adds movie to wishlist (body: `{movie_id}`)
- [x] `DELETE /api/v1/wishlist/{movieID}` — requires JWT (Bronze+) — removes movie from wishlist
- [x] `GET /api/v1/wishlist/check/{movieID}` — requires JWT — returns `{in_wishlist: true/false}`

#### Task 4.8: Implement audit handler  [x]

- [x] Create `api/audit_handler.go`
- [x] `GET /api/v1/audit`:
  - Requires Manager+
  - Query params: `user_id` (optional filter)
  - Call `GetAllAuditEntries` or `GetAuditEntriesByUser`
  - Decrypt entries with AES key
  - Return list

#### Task 4.9: Create router and server entrypoint  [x]

- [x] Create `api/router.go`
- [x] Build Chi router:
  - Apply `CORSMiddleware`, `LoggingMiddleware`, `RateLimitMiddleware(100)` globally
  - Group `/api/v1/auth`: register (no auth), login (no auth), login/totp (temporary session), refresh (JWT), logout (JWT)
  - Group `/api/v1/movies`: GET list/search/staff-picks/last-chance (JWT), POST/PUT/DELETE/staff-pick (JWT + PermAdmin)
  - Group `/api/v1/rentals`: all require JWT
  - Group `/api/v1/wishlist`: all require JWT (Bronze+)
  - Group `/api/v1/users`: GET/POST/PUT (JWT + Supervisor+), DELETE (JWT + PermAdmin), TOTP endpoints (JWT, self or Manager+)
  - Group `/api/v1/audit`: JWT + Supervisor+
  - Group `/api/v1/recommendations/{movieID}`: JWT — returns co-rental recommendations from Graph DS
  - Health check: `GET /health` returns `{"status":"ok"}`
- [x] Create `cmd/server/main.go`:
  - Load config via `config.Load()`
  - Open BoltDB store via `store.Open(config.DBPath)`
  - Defer store.Close()
  - Create router via `api.NewRouter(store, config)`
  - Start HTTP server on `config.ServerPort`
  - Graceful shutdown on SIGINT/SIGTERM

#### Phase 4 validation:  [x]

```bash
# Start server
go run ./cmd/server &
# Test endpoints with curl/httpie
curl -X POST http://localhost:8080/api/v1/auth/register -d '{"username":"test","password":"test123"}'
curl -X POST http://localhost:8080/api/v1/auth/login -d '{"username":"test","password":"test123"}'
# Export token, test protected routes
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/movies
curl -H "Authorization: Bearer $TOKEN" http://localhost:8080/api/v1/rentals/rent -d '{"movie_id":"..."}'
# Test 403 on guest rental attempt
# Test 401 on missing token
# Test rate limiting (100 rapid requests)
```

---

### Phase 5 — TUI Foundation  [x]

**Goal:** Set up the Bubble Tea application shell, global state management, theme system, and visual effects.

---

#### Task 5.1: Install TUI dependencies  [x]

- [x] Run `go get github.com/charmbracelet/bubbletea`
- [x] Run `go get github.com/charmbracelet/lipgloss`
- [x] Run `go get github.com/charmbracelet/bubbles`
- [x] Verify: `go mod tidy`

#### Task 5.2: Create theme and styles  [x]

- [x] Create `tui/styles/theme.go`
- [x] Define color palette constants using `lipgloss.Color`:
  - `Cyan = "#00FFFF"`, `Magenta = "#FF00FF"`, `Yellow = "#FFFF00"`
  - `NeonGreen = "#39FF14"`, `NeonPink = "#FF6EC7"`
  - `Background = "#0A0A2E"` (dark blue), `Surface = "#121240"`, `BorderDim = "#333366"`
  - `Error = "#FF4444"`, `Success = "#44FF44"`, `Warning = "#FFAA00"`
- [x] Define `AppStyle` — full-screen container with background color
- [x] Define `BorderStyle` — lipgloss border with rounded corners, cyan/magenta edge
- [x] Define `TitleStyle` — bold, cyan, large text
- [x] Define `TextStyle`, `DimTextStyle`, `ErrorTextStyle`, `SuccessTextStyle`
- [x] Define tier-specific color map: `TierColors = map[string]lipgloss.Color{...}`

#### Task 5.3: Create visual effects  [x]

- [x] Create `tui/styles/effects.go`
- [x] Implement `Scanlines(width, height int) string` — generates alternating lines of semi-transparent `░` pattern over the full terminal dimensions
- [x] Implement `GlitchFrame() string` — returns random character noise `▓▒░█▄▀` (1-3 chars) for temporary glitch overlay; called randomly on page transitions
- [x] Implement `VHSSpinner() []string` — custom spinner frames: `["▌", "▌ ", " ▌", " ▌", "▌ ", "▌", " ▌", " ▌ "]` (tracking artifact)
- [x] Implement `RewindAnimation(tapeName string) string` — returns "◄◄ REWINDING: <tapeName> ... ▌" styled text
- [x] Implement `AccessDeniedOverlay(width, height int) string` — full-screen "⛔ ACCESS DENIED" in red with scanlines

#### Task 5.4: Create global state  [x]

- [x] Create `tui/state.go`
- [x] Define `SessionState` struct:
  - `AccessToken string`, `RefreshToken string`
  - `User *models.UserResponse`
  - `Permissions bitmask.Permission`
  - `IsLoggedIn bool`
  - `Cache *lru.Cache[string, interface{}]` — shared cache (capacity 1000)
  - `MovieCache *lru.Cache[string, *models.MovieResponse]` — dedicated movie cache
  - `APIBaseURL string`
- [x] Implement `NewSessionState(apiURL string) *SessionState` — initializes caches
- [x] Implement `Login(tokenPair *auth.TokenPair, user *models.UserResponse)`
- [x] Implement `Logout()` — clears session
- [x] Implement `HasPermission(perm bitmask.Permission) bool`
- [x] Implement `CanAccessAdmin() bool`
- [x] Implement `RefreshAccessToken() error` — calls `/api/v1/auth/refresh` with current refresh token, updates tokens

#### Task 5.5: Create API client  [x]

- [x] Create `tui/api_client.go` (or within `state.go`)
- [x] Implement generic `doRequest(method, path string, body interface{}, target interface{}) error`:
  - Uses `net/http` with 10s timeout
  - Sets `Authorization: Bearer <token>` if logged in
  - Sets `Content-Type: application/json`
  - JSON-encodes body if non-nil
  - JSON-decodes response into target
  - On 401 → attempts refresh token → retries once
- [x] Implement typed methods:
  - `Login(username, password string) (*LoginResponse, error)`
  - `Register(username, password string) (*UserResponse, error)`
  - `SearchMovies(query string) ([]*MovieResponse, error)`
  - `GetMovies(genre string, page int) (*MovieListResponse, error)`
  - `GetMovie(id string) (*MovieResponse, error)`
  - `RentMovie(movieID string) (*RentalResponse, error)`
  - `ReturnMovie(rentalID string) (*RentalResponse, error)`
  - `GetRentalHistory() ([]*RentalResponse, error)`
  - `GetUsers() ([]*UserResponse, error)`
  - `UpdateUser(id string, req *UpdateUserRequest) (*UserResponse, error)`
  - `GetAuditEntries() ([]*AuditEntry, error)`

#### Task 5.6: Create TUI application shell  [x]

- [x] Create `tui/app.go`
- [x] Define `Model` struct implementing `bubbletea.Model`:
  - Fields: `currentPage Page`, `session *SessionState`, `width int`, `height int`, `ready bool`, `lastTick time.Time`
- [x] Define `Page` type: `type Page int` with constants:
  - `PageSplash`, `PageLogin`, `PageRegister`, `PageBrowse`, `PageMovieDetail`, `PageMyRentals`, `PageProfile`, `PageAdminUsers`, `PageAdminMovies`, `PageAuditLog`
- [x] Implement `Init() tea.Cmd`:
  - Returns `tea.Batch(tea.EnterAltScreen, tea.ClearScreen, tickCmd())` — ticks for clock/animations
- [x] Implement `Update(msg tea.Msg) (tea.Model, tea.Cmd)`:
  - `tea.WindowSizeMsg` → store width/height, set ready
  - `tea.KeyMsg`:
    - `ctrl+c`, `esc` (on non-modal) → `tea.Quit`
    - Delegate to current page's Update method
  - `tickMsg` → trigger re-render for clock update, return `tickCmd()`
  - Delegate all other messages to current page's Update
- [x] Implement `View() string`:
  - If not ready → "Initializing..."
  - Render: header component + current page View + footer component
- [x] Create `cmd/client/main.go`:
  - Parse CLI flags: `--api-url` (default `http://localhost:8080`), `--debug`
  - Create `SessionState`
  - Create `Model` with `PageSplash`
  - Run `tea.NewProgram(model, tea.WithAltScreen()).Run()`

#### Task 5.5a: Create tick and header  [x]

- [x] `tickCmd()` returns `tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })`
- [x] Header shows: ASCII "THE LAST VIDEO STORE" banner (hardcoded multi-line string with ANSI/lipgloss styling), current time, "NOW SHOWING" if a movie is highlighted

#### Phase 5 validation:  [x]

```bash
go build ./cmd/client && ./client
# Splash screen appears
# Clock ticks in header
# Terminal resize recalculates layout
# Ctrl+C quits cleanly
```

---

### Phase 6 — TUI Pages & Components  [x]

**Goal:** Build all 10 interactive pages and 9 reusable components.

---

#### Task 6.1: Splash screen  [x]

- [x] Create `tui/pages/splash.go`
- [x] 3-second VHS-style intro sequence:
  - Frame 0: "LOADING..." in neon green on black
  - Frame 1: Glitch effect (random characters overlaying)
  - Frame 2: ASCII "THE LAST VIDEO STORE" logo fades in (lipgloss with increasing foreground color brightness)
  - Frame 3: "INSERT MEMBERSHIP CARD █" blinking
- [x] After 3 seconds, auto-transition to `PageLogin`
- [x] Implement as Bubble Tea model: triggers `tea.Tick(3*time.Second, ...)` in `Init()`, transitions on tick message

#### Task 6.2: Login page  [x]

- [x] Create `tui/pages/login.go`
- [x] Uses `bubbles/textinput` for username (focused) and password (masked, `EchoMode: textinput.EchoPassword`)
- [x] "TAB" cycles between fields, "ENTER" submits
- [x] Submit: shows spinner via `tea.Batch(spinner.Tick, loginCmd)`
- [x] `loginCmd` calls API client `Login()`:
  - Success → store tokens in session, set user, navigate to `PageBrowse`
  - 401 → show error "Invalid credentials" in red
  - 429 → show "Account locked. Try again in X minutes"
- [x] "ESC" → quit (on splash page) or navigate back
- [x] Link text at bottom: "No card? [R]egister" → pressing R navigates to `PageRegister`

#### Task 6.3: Register page  [x]

- [x] Create `tui/pages/register.go`
- [x] Similar to login: username + password + confirm password fields
- [x] Validation: passwords match, username 3+ chars, password 6+ chars
- [x] Submit calls `Register()` API:
  - Success → "Account created! Press any key to login" → navigate to `PageLogin`
  - 409 → "Username already taken"
- [x] "ESC" → back to `PageLogin`

#### Task 6.4: Header component  [x]

- [x] Create `tui/components/header.go`
- [x] `HeaderView(width int, session *SessionState) string`:
  - Top border: `══════` full width in cyan
  - Left section: ASCII "THE LAST VIDEO STORE" banner (4 lines)
  - Right section: `🕐 Fri Jun 14 2002  9:48 PM` (dynamic clock)
  - If logged in: `🎫 user@tier BADGE | 🍿 142 pts`
  - Divider: `──── NOW SHOWING: The Matrix ────` (picks random available movie)
  - Bottom border

#### Task 6.5: Footer component  [x]

- [x] Create `tui/components/footer.go`
- [x] `FooterView(width int, page Page, session *SessionState) string`:
  - Context-sensitive keybinding bar:
    - Login: `[TAB] switch field  [ENTER] login  [R] register  [ESC] quit` (+ `[T] TOTP` if prompted)
    - Browse: `[↓↑] navigate  [ENTER] details  [S] search  [R] my rentals  [W] wishlist  [P] profile  [Q] quit`
    - Admin (Supervisor+): adds `[U] users  [M] movies  [A] audit`
    - Admin (Manager+): adds `[S] staff picks` toggle on selected movie
    - Movie detail: `[ENTER] rent  [ESC] back  [W] add to wishlist` (shows recommendations panel if available)
    - Profile: `[T] toggle TOTP  [L] logout  [ESC] back`
  - Styled with dim text on surface background

#### Task 6.6: Search bar component  [x]

- [x] Create `tui/components/searchbar.go`
- [x] Uses `bubbles/textinput` with placeholder "Search movies..."
- [x] On each keypress, debounced (200ms), calls `SearchMovies(prefix)` via API client
- [x] API returns results from BoltDB prefix scan on `movies_by_title`
- [x] Displays dropdown with max 5 suggestions below input
- [x] `ENTER` on suggestion → navigate to `PageMovieDetail` for that movie
- [x] `ESC` clears search and closes dropdown

#### Task 6.7: Movie card component  [x]

- [x] Create `tui/components/movie_card.go`
- [x] `MovieCardView(movie *models.MovieResponse, selected bool, width int) string`:
  - Border: magenta if selected, dim if not
  - Content: title (bold, truncated to 20 chars), year `(1999)`, format badge `📼 VHS` / `📀 DVD` / `💿 Blu-ray`, genre badge (colored pill), star rating `★★★★½`, availability `[RENT]` or `[OUT]` or `[NEW]`
  - Fixed-width card: 22x8 characters
  - Selected card gets highlighted border + background

#### Task 6.8: Movie grid component  [x]

- [x] Create `tui/components/movie_grid.go`
- [x] `MovieGridView(movies []*models.MovieResponse, selectedIndex int, width, height int) string`:
  - Calculates columns based on terminal width (`width / 22`)
  - Renders grid of `MovieCardView` components
  - Selected index highlighted
  - Handles arrow key navigation: up/down/left/right within grid (wrap or stop at edges)
  - Page up/down for next/prev page of results

#### Task 6.9: Tabs component  [x]

- [x] Create `tui/components/tabs.go`
- [x] `TabsView(tabs []string, activeIndex int, width int) string`:
  - Genre tabs: `ALL | ACTION | COMEDY | HORROR | SCIFI | DRAMA | NEW | STAFF PICKS | LAST CHANCE`
  - Active tab: cyan background, bold text
  - Inactive: dim text
  - Styled with lipgloss borders connecting to content below
- [x] `TabWidth` calculates equal widths filling available space

#### Task 6.9a: Wishlist sidebar component  [x]

- [x] Create `tui/components/wishlist_sidebar.go`
- [x] `WishlistSidebarView(wishlist []*WishlistItem, width int) string`:
  - Right-side panel showing user's wishlist items
  - Each item: movie title truncated + format badge + availability indicator (`🟢 Available` / `🔴 Rented out`)
  - "Available now!" highlighted entry when a wishlisted title becomes available
  - Quick-rent: press `W` on highlighted movie to rent directly from wishlist
  - Empty state: "Your wishlist is empty. Browse and press [W] to add titles."
  - Styled with dim border, scrollable if many items

#### Task 6.10: Browse page  [x]

- [x] Create `tui/pages/browse.go`
- [x] Combines: searchbar + tabs + movie grid + wishlist sidebar
- [x] State: `selectedGenre int`, `movies []MovieResponse`, `selectedMovie int`, `searchMode bool`, `page int`, `showWishlist bool`
- [x] `Init()`: fetches movies for default genre (ALL) via `GetMovies("", 1)`, fetches wishlist
- [x] Wishlist sidebar (right panel): shows user's wishlisted titles, "Available now!" indicator, quick-rent shortcut
- [x] Keybindings:
  - `←→`: switch tabs → refetch movies for genre
  - `↓↑`: navigate grid
  - `ENTER`: navigate to `PageMovieDetail`
  - `/`: focus search bar
  - `W`: add selected movie to wishlist / toggle wishlist sidebar
  - `R`: navigate to `PageMyRentals`
  - `P`: navigate to `PageProfile`
  - Supervisor+: `U` → `PageAdminUsers`, `A` → `PageAuditLog`
  - Manager+: `U` → `PageAdminUsers`, `M` → `PageAdminMovies`, `A` → `PageAuditLog`
- [x] Admin links only visible if `session.HasPermission(PermManageUsers)` or `session.HasPermission(PermAdmin)`

#### Task 6.11: Movie detail page  [x]

- [x] Create `tui/pages/movie_detail.go`
- [x] Full-screen movie view:
  - Title (large, bold)
  - `[NEW RELEASE]` or `[AVAILABLE]` or `[RENTED OUT]` badge
  - Year · Genre · Director · Format badge
  - Star rating: `★★★★½ (4.5/5 from 1,247 ratings)`
  - Synopsis (wrapped text, 3-4 lines)
  - Cast: comma-separated
  - Copies available: `📼 3 of 5 copies available`
  - **Co-rental recommendations** (bottom panel): "Customers who rented this also:" — 3-5 recommended titles from Graph DS, ordered by co-rental weight
- [x] Actions:
  - `ENTER` → rent movie (calls `RentMovie` API):
    - Success → "📼 RENTED! Due: Jun 17 2002" modal, then navigate to browse
    - 403 "Limit reached" → show error modal
    - 403 "Gold plan required" → show promo to upgrade
    - 409 "No copies" → offer "Join waitlist? [Y/N]"
  - `W` → add/remove from wishlist
  - `ESC` → back to browse

#### Task 6.12: My rentals page  [x]

- [x] Create `tui/pages/my_rentals.go`
- [x] Fetches `GetRentalHistory()` on init
- [x] Lists active rentals at top with:
  - Movie title, rental date, due date, format badge, status `🟢 Active` / `🔴 Overdue`, `🔄 VHS` rewind indicator if applicable
- [x] Lists rental history below (returned) with `ReturnedAt` date, late fee, rewind fee (if any)
- [x] Selected rental can be returned:
  - Press `ENTER` on active rental → confirmation modal "Return The Matrix?"
  - Confirm → calls `ReturnMovie` API
  - Success: "📼 Returned! Late fee: $4.00 (VHS: 2 days × $2/day) + Rewind fee: $1.00" or "📼 Returned on time! +10 popcorn points"
  - Movie grid and rental count refresh
- [x] `ESC` → back to browse

#### Task 6.13: Profile page  [x]

- [x] Create `tui/pages/profile.go`
- [x] Membership card view:
  ```
  ╔══════════════════════════╗
  ║  THE LAST VIDEO STORE MEMBERSHIP    ║
  ║                          ║
  ║  Username: gold_member   ║
  ║  Plan:   ★★ GOLD ★★     ║
  ║  Member since: 2002-06-01 ║
  ║                          ║
  ║  📀 DVD Rentals: 2/5     ║
  ║  📼 VHS Rentals: 1/5     ║
  ║  🍿 Popcorn Points: 142  ║
  ║  ⏱ Total movies: 27     ║
  ║  $ Late fees paid: $6   ║
  ║  $ Rewind fees: $2      ║
  ║  🔒 2FA: Enabled         ║
  ╚══════════════════════════╝
  ```
- [x] Membership plan badge in corresponding color (Bronze=#CD7F32, Silver=#C0C0C0, Gold=#FFD700, Employee=magenta, Supervisor=orange, Manager=yellow, Owner=cyan)
- [x] Popcorn points mock calculation: 10 per on-time return, -5 per late
- [x] Stats pulled from rental history, grouped by format (DVD/VHS/Blu-ray)
- [x] Rewind fees tracked separately from late fees
- [x] `L` → logout → clear session → navigate to `PageLogin`
- [x] `T` → toggle TOTP setup (Manager+ only) — shows secret key, otpauth URL, verification prompt
- [x] `ESC` → back to browse

#### Task 6.14: Modal component  [x]

- [x] Create `tui/components/modal.go`
- [x] `ModalView(title, message string, width, height int) string`:
  - Overlay: dimmed background over current page
  - Centered bordered box with title (bold) and message
  - Buttons: `[ENTER] Confirm  [ESC] Cancel`
- [x] `AccessDeniedModal(width, height int) string`:
  - Same layout but red-tinted
  - Title: `⛔ ACCESS DENIED`
  - Message: `Insufficient clearance level`
  - Only shows `[ESC] Dismiss`

#### Task 6.15: Spinner component  [x]

- [x] Create `tui/components/spinner.go`
- [x] `VHSSpinnerView() string` — returns current spinner frame
- [x] Used during API calls: login, register, rent, return
- [x] Integrated into pages via `tea.Batch(spinnerTickCmd, apiCallCmd)`

#### Task 6.16: Badge component  [x]

- [x] Create `tui/components/badge.go`
- [x] `TierBadgeView(tierName string) string`:
- [x] Color-coded pill: `[ BRONZE ]` (bronze brown), `[ SILVER ]` (silver gray), `[ GOLD ]` (gold yellow), `[ ATENDENTE ]` (magenta), `[ SUPERVISOR ]` (orange), `[ GERENTE ]` (yellow), `[ DONO ]` (cyan)
- [x] Styled with lipgloss background + foreground + padding

#### Task 6.17: Admin users page  [x]

- [x] Create `tui/pages/admin_users.go`
- [x] Requires `PermManageUsers` (Supervisor+) — if insufficient, show `AccessDeniedModal`
- [x] Fetches `GetUsers()` from API
- [x] Displays table:
  - Columns: Username | Tier | Rentals | Banned | TOTP | Actions
  - Each row selectable with `↓↑`
- [x] Actions on selected user:
  - `P` → promote (increment tier, max Owner)
  - `D` → demote (decrement tier, min Bronze)
  - `B` → toggle ban
  - Confirmation modal for each action
- [x] Calls `UpdateUser` API, refreshes list on success
- [x] TOTP status column: `🔒` if enabled, `—` if disabled

#### Task 6.18: Admin movies page  [x]

- [x] Create `tui/pages/admin_movies.go`
- [x] Requires `PermAdmin` (Manager+) — if insufficient (Supervisor), show `AccessDeniedModal`
- [x] Table of all movies: Title | Year | Genre | Copies | Available | Staff Pick
- [x] Actions:
  - `A` → add movie form (text inputs for all fields) → calls `CreateMovie`
  - `ENTER` → edit selected movie (populated form) → calls `UpdateMovie`
  - `D` → delete movie (confirmation modal) → calls `DeleteMovie`
  - `S` → toggle Staff Pick (Manager+)
- [x] Form navigation: TAB between fields, ENTER to submit, ESC to cancel

#### Task 6.19: Audit log page  [x]

- [x] Create `tui/pages/audit_log.go`
- [x] Requires `PermManageUsers`
- [x] Fetches `GetAuditEntries()` from API
- [x] Displays scrollable list:
  - Each entry: `[timestamp] ACTION | Actor: user | Target: target | Hash: a1b2c3...`
  - Hash chain verification status at top: `✅ Chain intact (142 entries)` or `❌ Chain broken!`
  - Each entry shows truncated PrevHash → Hash link
- [x] Scroll with `↓↑`, `PgUp/PgDn`, `Home/End`
- [x] `V` → verify chain against API (triggers recomputation)
- [x] `ESC` → back to browse

#### Phase 6 validation:  [x]  [x]

```bash
go run ./cmd/client
# Full manual walkthrough:
# 1. Splash → Login (as bronze) → Browse (see grid)
# 2. Search "mat" → see Matrix, Matilda, Match Point
# 3. Click Matrix → Movie Detail → Rent → success (Bronze can rent 1 title)
# 4. Try to rent a New Release as Bronze → "Gold plan required for new releases"
# 5. Logout → Login as gold → Rent Matrix → success
# 6. Add movie to wishlist → sidebar shows "Available now!" when back in stock
# 7. View My Rentals → see Matrix due date (format-specific)
# 8. Return Matrix → see late fee or on-time confirmation (+ popcorn or -5 late)
# 9. Profile → see plan badge (Gold), stats, popcorn points, TOTP setup option
# 10. Logout → Login as supervisor → Admin Users → upgrade silver to Gold
# 11. Logout → Login as manager (Gerente) → Admin Movies → add a new Blu-ray title
# 12. Manager: set Staff Picks, browse → see on "Staff Picks" tab
# 13. Audit Log → verify chain intact
# 14. Login as banned → "Account suspended"
# 15. Enable TOTP on manager account → login with 2FA → success
# 16. All pages responsive to terminal resize
```

---

### Phase 7 — Seed Data & Integration Testing  [x]

**Goal:** Populate the system with realistic data and test all end-to-end flows.

---

#### Task 7.1: Define seed data — movies  [x]

- [x] Create `data/movies.json` (or hardcode in `data/seed.go`)
- [x] ~40 real movies from 1980s–2000s:
  - The Matrix (1999), Fight Club (1999), Pulp Fiction (1994), Jurassic Park (1993), The Shawshank Redemption (1994), The Dark Knight (2008), Inception (2010), Forrest Gump (1994), The Godfather (1972), Schindler's List (1993), Goodfellas (1990), The Silence of the Lambs (1991), Se7en (1995),The Usual Suspects (1995), Léon: The Professional (1994), American History X (1998), Saving Private Ryan (1998), The Green Mile (1999), Gladiator (2000), Memento (2000), The Lord of the Rings trilogy (2001-2003), Kill Bill (2003), Eternal Sunshine (2004), The Departed (2006), No Country for Old Men (2007), There Will Be Blood (2007), WALL-E (2008), Inglourious Basterds (2009), District 9 (2009), Blade Runner (1982), Back to the Future (1985), Die Hard (1988), Terminator 2 (1991), Toy Story (1995), The Big Lebowski (1998), The Truman Show (1998), American Beauty (1999), Requiem for a Dream (2000), Spirited Away (2001), City of God (2002)
- [x] Each movie needs: title, year, genre, director, cast (3-5 real names), synopsis (2-3 sentences), rating (3.0-5.0), rating count (100-5000), copies total (2-10), format (mix of DVD/VHS/Blu-ray), is_new_release (first 6 marked true)

#### Task 7.2: Define seed data — users  [x]

- [x] Create 8 test users (hardcoded in `data/seed.go`):
  ```
  bronze    / password1  → TierBronze,     MaxRentals=1  (Cliente Bronze — browse + 1 rental + wishlist)
  silver    / password2  → TierSilver,     MaxRentals=2  (Cliente Prata — rent up to 2, wishlist)
  gold      / password3  → TierGold,       MaxRentals=5  (Cliente Ouro — new releases, waitlist, wishlist)
  employee  / password4  → TierEmployee,   MaxRentals=5  (Atendente — staff: process any return)
  supervisor/ password8  → TierSupervisor, MaxRentals=5  (Supervisor — manage users, view audit)
  manager   / password5  → TierManager,    MaxRentals=10 (Gerente — CRUD movies + manage users + view audit)
  owner     / password6  → TierOwner,      MaxRentals=99 (Dono — all permissions)
  banned    / password7  → TierBronze,     Banned=true   (Blocked account)
  ```
- [x] All passwords hashed with bcrypt
- [x] Add banned user to Bloom filter

#### Task 7.3: Implement seed script  [x]

- [x] Create `data/seed.go`
- [x] Package `main` (runnable)
- [x] Functions:
  - `seedMovies(store *store.Store)` — iterates movie list, creates each in BoltDB
  - `seedUsers(store *store.Store)` — iterates user list, hashes passwords, creates in BoltDB
  - `main()`:
    - Load config
    - `os.Remove(config.DBPath)` to start fresh
    - Open store
    - Call seed functions
    - Print summary: "Seeded 40 movies and 8 users."
    - Close store

#### Task 7.4: Integration test flow 1 — Bronze rental limit  [x]

- [x] Login as bronze → Browse → Rent a movie (DVD) → Success (1 rental allowed)
- [x] Try to rent a second movie → Expected: Modal "Rental limit reached (1/1)"
- [x] Try to rent a New Release → Expected: "Gold plan required for new releases"
- [x] Verify: Only 1 active rental in DB

#### Task 7.5: Integration test flow 2 — Silver rental limit  [x]

- [x] Login as silver → Rent movie 1 (DVD) → Rent movie 2 (VHS) → Try to rent movie 3
- [x] Expected: Modal "Rental limit reached (2/2)"
- [x] Verify: Only 2 active rentals in DB; format-specific durations applied (3d VHS, 5d DVD)

#### Task 7.6: Integration test flow 3 — Banned user  [x]

- [x] Login as banned
- [x] Expected: "Account suspended. Contact store management."
- [x] Verify: Bloom filter check passes (banned flag detected), JWT not issued

#### Task 7.7: Integration test flow 4 — Plan upgrade by Supervisor (Silver → Gold)  [x]

- [x] Login as supervisor → Admin Users → Select silver → Press P to promote
- [x] Expected: Confirmation modal "Upgrade silver to Gold?" → Confirm → User tier updated
- [x] Try to access Admin Movies as supervisor → ACCESS DENIED (no PermAdmin)
- [x] Login as (formerly silver) → Verify can now rent 5 movies, see new releases, join waitlist
- [x] Verify: Audit log entry `ActionPromote` recorded

#### Task 7.8: Integration test flow 5 — Audit chain  [x]

- [x] Login as supervisor → Audit Log → Press V to verify
- [x] Expected: "✅ Chain intact" with entry count
- [x] Tamper test (manual): corrupt one audit entry hash in BoltDB → Verify → "❌ Chain broken at entry #42"

#### Task 7.9: Integration test flow 6 — Employee return with deque & rewind fee  [x]

- [x] Login as employee (Atendente) → Return overdue movies for multiple customers
- [x] Expected: Most overdue is processed first (deque pop from back)
- [x] One VHS rental has `NeedsRewind=true` → on return, $1.00 rewind fee added
- [x] Verify: API returns list sorted by priority; late fees auto-calculated per format ($2/day VHS, $3/day DVD); rewind fee shown separately

#### Task 7.10: Integration test flow 7 — New release waitlist  [x]

- [x] Login as gold → Try to rent sold-out new release
- [x] Expected: "Join waitlist?" modal → Confirm → Added to heap with timestamp
- [x] Verify: Heap peek returns user with oldest timestamp

#### Task 7.11: Integration test flow 8 — Co-rental recommendations (Graph)  [x]

- [x] Login as gold → Rent The Matrix + The Matrix Reloaded (both action/sci-fi)
- [x] Navigate to The Matrix detail page → See "Customers who rented this also:" section
- [x] Verify: Graph edge between The Matrix and The Matrix Reloaded incremented; recommendations shown ordered by co-rental weight

#### Task 7.12: Integration test flow 9 — TOTP 2FA setup & login  [x]

- [x] Login as manager → Profile → Enable TOTP
- [x] Expected: ASCII art display of TOTP secret + otpauth URL (for QR scanning)
- [x] Logout → Login as manager → After password, prompted for TOTP code
- [x] Enter correct code → Access granted
- Enter wrong code → "Invalid TOTP code" → After 3 failures → "TOTP locked for 10 minutes"

#### Phase 7 validation:  [x]  [x]

```bash
go run ./data/seed.go                       # Seeds database
go run ./cmd/server &                        # Start API
go run ./cmd/client                          # Start TUI
# Execute all 9 integration test flows manually
```

---

### Phase 8 — Deployment & Final Polish  [x]

**Goal:** Containerize, deploy to Render, cross-compile binaries, write documentation, and polish.

---

#### Task 8.1: Create Makefile  [x]

- [x] Create `Makefile`
- [x] Targets:
  - `build-server`: cross-compile server binary for linux/amd64
  - `build-client-linux`: cross-compile client binary for linux/amd64
  - `build-client-windows`: cross-compile client binary for windows/amd64
  - `build-all`: all three builds
  - `test`: run all tests with race detector and coverage
  - `test-cover`: generate coverage HTML report
  - `lint`: run golangci-lint
  - `seed`: run seed script
  - `run-server`: run server locally
  - `run-client`: run client locally
  - `clean`: remove binaries and test DB
  - `tidy`: go mod tidy
  - `fmt`: go fmt ./...

#### Task 8.2: Create Dockerfile  [x]

- [x] Create `Dockerfile.server`
- [x] Multi-stage build:
  - Stage 1 (builder): `golang:1.22-alpine`, copy source, `go build -o /server ./cmd/server`
  - Stage 2 (runtime): `alpine:3.19`, copy binary, expose 8080, run
- [x] Ensure static linking: `CGO_ENABLED=0 GOOS=linux`
- [x] Set `JWT_SECRET` and `AES_KEY` as build args with defaults for demo

#### Task 8.3: Create render.yaml  [x]

- [x] Create `render.yaml`
- [x] Service definition:
  - Type: `web`
  - Name: `thelastvideostore-api`
  - Runtime: `docker`
  - Dockerfile path: `./Dockerfile.server`
  - Health check path: `/health`
  - Env vars: `JWT_SECRET` (generateValue), `AES_KEY` (generateValue), `DB_PATH=/data/thelastvideostore.db`
  - Disk: mount at `/data` for BoltDB persistence

#### Task 8.4: Cross-compile and test binaries  [x]

```bash
make build-all
# Verify Linux binary on local machine:
./bin/thelastvideostore-linux
# Verify Windows binary (if WSL or VM available):
# Copy to Windows, run thelastvideostore.exe
# Or verify with: file bin/thelastvideostore.exe
# Expected: PE32+ executable (console) x86-64
```

#### Task 8.5: Final code cleanup  [x]

- [x] Run `make fmt` → all files formatted
- [x] Run `make lint` → zero warnings
- [x] Run `make test` → all tests pass, coverage > 70%
- [x] Remove any debug prints, unused imports, commented-out code
- [x] Ensure no hardcoded secrets (all from env/config)
- [x] Verify sensitive files in .gitignore: `*.db`, `.env`, `bin/`

#### Task 8.6: Create .gitignore  [x]

```
# Binaries
bin/
*.exe

# Database
*.db

# Environment
.env
.env.local

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Coverage
coverage.out
coverage.html

# Dependencies
vendor/
```

#### Task 8.7: Write README.md  [x]

- [x] Create `README.md`
- [x] Sections:
  - Title + ASCII banner
  - Project description (2 paragraphs)
  - Architecture diagram (ASCII from Section 4)
  - Features list
  - Tech stack badges
  - Quick start (clone, seed, run-server, run-client)
  - Test users table
  - Screenshots (placeholders or generated via terminal capture tool)
  - API documentation (endpoints table)
  - Cross-compilation instructions
  - Deployment to Render instructions
  - Project structure tree
  - Academic context (links back to IMPLEMENTATION.md)

#### Task 8.8: Demo preparation  [x]

- [x] Prepare demo script for presentation:
  1. Start server + seed data
  2. Launch client, show splash screen (VHS-style intro)
  3. Login as bronze → browse catalog → rent 1 movie (Bronze now rents!) → try second → "Rental limit reached (1/1)"
  4. Demo search autocomplete: type "mat" → Trie shows Matrix, Matilda, Match Point
  5. Logout → login as silver → rent 2 movies (1 DVD + 1 VHS, different due dates)
  6. Add a movie to wishlist → show wishlist sidebar with "Available now!" notification
  7. View My Rentals → see format badges, due dates, late fee warnings
  8. Return a VHS movie → "REWIND FEE: $1.00" (30% chance), Popcorn Points breakdown
  9. View Profile → membership card, Gold tier badge, rental stats, TOTP setup option
  10. Logout → login as supervisor → Admin Users → upgrade silver to Gold plan
  11. Try Supervisor to access Admin Movies → ACCESS DENIED (no PermAdmin — Manager only)
  12. Logout → login as manager → manage catalog, set Staff Picks, browse Staff Picks tab
  13. Show Audit Log → verify hash chain integrity
  14. Movie detail → show co-rental recommendations (Graph data structure in action)
  15. Enable TOTP on manager account → logout → login with 2FA → success
  16. Login as banned → "Account suspended"
  17. Terminal resize demo (responsive movie grid + wishlist sidebar)
  18. Mention cross-platform: show Linux + Windows binaries

#### Task 8.9: Optional polish items  [x]

- [x] Add sound effects toggle (beep on rent/return) via `\a` bell character
- [x] Easter egg: Konami code (↑↑↓↓←→←→BA) shows secret "Employee Picks" menu
- [x] ASCII movie posters (hardcoded simple art for top 5 movies)
- [x] On-exit animation: "BE KIND, REWIND" in large ASCII
- [x] TOTP QR code rendered as ASCII QR in terminal (via qrcode-terminal-go or custom block characters)
- Co-rental graph visualization: render small ASCII graph on Movie Detail showing connected titles

#### Phase 8 validation:  [x]  [x]

```bash
make build-all
make test
make lint
# Docker build
docker build -f Dockerfile.server -t thelastvideostore-server .
docker run -p 8080:8080 -e JWT_SECRET=test -e AES_KEY=0123456789abcdef0123456789abcdef thelastvideostore-server
# curl http://localhost:8080/health → {"status":"ok"}
# Deploy to Render via render.yaml or gh render
```

---

## 7. Requirements Traceability Matrix

| # | Requirement | The Last Video Store Implementation | Evidence in Presentation |
|---|------------------------|--------------------------|--------------------------|
| 1 | **Interface** | Full Bubble Tea TUI with 10+ interactive screens, CRT effects, search, grids, modals, wishlist sidebar | Navigate catalog → rent → return → view profile — all in terminal |
| 2 | **Modo de segurança de acesso** | bcrypt + JWT + 7-tier RBAC bitmask + brute-force lockout + AES-256-GCM + Bloom filter ban list + optional TOTP 2FA | Show login fail → lockout → ACCESS DENIED modal → Supervisor promotes user → access granted → TOTP 2FA demo |
| 3 | **Cybersecurity** | 7-layer security: hashing, token auth, bitmask RBAC, encryption at rest, immutable audit via hash chain, input sanitization, TOTP 2FA | Demonstrate hash chain integrity check, AES-encrypted audit entries, Bloom filter banning, TOTP setup + verification |
| 4 | **Data Structures** | 9 structures implemented from scratch: Trie, LRU, Deque, MinHeap, DoublyLinkedList, BloomFilter, Bitmask, HashChain, Graph | Show `_test.go` passing, explain each structure's role (Graph for co-rental recommendations) |
| 5 | **Read file line by line** | BoltDB stores movies persistently; API reads paginated results; TUI renders each as a card | Browse catalog with search autocomplete (Trie in action) |
| 6 | **Allow only authorized** | JWT middleware + 6-bit bitmask on every route; PermStaff distinguishes Employee from Gold cleanly | Bronze rents 1 movie → reaches limit; New Release requires Gold; Supervisor manages users but NOT movies |
| 7 | **Show user and file data** | Profile screen: username, 7-tier plan badge, rental history (linked list), popcorn points, rewind fees, TOTP status | Navigate to Profile, scroll rental history, show plan badge (Bronze/Silver/Gold/Employee/Supervisor/Manager/Owner) |
| 8 | **User registration via file/DB** | BoltDB persistent store + `/auth/register` endpoint | Register new user → login → Bronze plan automatically assigned (MaxRentals=1) |
| 9 | **Add/remove access** | Admin user panel: upgrade/downgrade plan (Supervisor+), ban (add to Bloom filter), toggle TOTP | Supervisor upgrades Bronze → Silver → Gold; Manager CRUD movies + set Staff Picks |
| 10 | **Cross-platform** | Go cross-compilation: Linux + Windows binaries | Show both binaries; run on Linux, optionally demo on Windows |

---

## 8. Technology Stack Summary

| Layer | Technology | Justification |
|-------|-----------|---------------|
| Language | Go 1.22+ | Compiled, cross-platform, excellent cryptography stdlib, fast |
| TUI | Bubble Tea + Lipgloss + Bubbles | Mature Elm-architecture TUI, extensive styling, terminal-native |
| HTTP Router | Chi v5 | Lightweight, idiomatic, middleware-friendly |
| Database | BoltDB | Embedded, zero-config, ACID, pure Go — no external DB to install |
| Auth | golang-jwt v5 + bcrypt | Industry standard, well-audited libraries |
| Crypto | golang.org/x/crypto + crypto/aes + crypto/sha256 | Go stdlib — no third-party crypto dependencies |
| Deploy | Render + Docker | Free tier, Go-native buildpack, managed HTTPS |
| Linting | golangci-lint | Comprehensive Go linter aggregator |
| Version Control | Git + GitHub | Standard |

*Document version 5.0 — The Last Video Store Project Implementation Plan — Last updated: 2026-06-14*
*
*Major v5 changes: 7-tier RBAC (6-bit bitmask + Supervisor + PermStaff), Bronze=1 rental, Graph DS, TOTP 2FA, Rewind Fee, Staff Picks/Last Chance endpoints*
