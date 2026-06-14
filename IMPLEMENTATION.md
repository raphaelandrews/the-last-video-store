# THE LAST VIDEO STORE — Retro Video Rental Terminal

> A cyber-secure, RBAC-protected video rental system with a retro 2000s aesthetic.
> Built in Go with a Bubble Tea TUI, BoltDB persistence, and deployable REST API.
> Developed for the **Cybersecurity & Data Structures** university project.

---

## 1. Project Overview

**The Last Video Store** is a full-stack video rental management system styled after the golden age of
Blockbuster (circa 2002). It features a rich terminal user interface powered by
[Charmbracelet Bubble Tea](https://github.com/charmbracelet/bubbletea) and a REST API
backend that runs on [Render](https://render.com).

Users browse a movie catalog, rent tapes, return them, and manage their membership — all
gated by a **6-tier Role-Based Access Control (RBAC)** system enforced through bitmask

The system directly addresses the three core challenges outlined in the project scope:

| Requirement | The Last Video Store Implementation |
|-------------------------|--------------------------|
| **a) Read the file line by line** | BoltDB stores movies persistently; the API serves paginated catalog data; the TUI renders movie cards with autocomplete powered by a custom Trie |
| **b) Allow only authorized users to read** | RBAC bitmask + JWT middleware + brute-force lockout after 5 failed attempts |
| **c) Show the user and the file data** | Profile screen displays tier, active rentals, and full rental history; audit log shows an immutable hash-chained trail of all access events |

### 1.1 System Features at a Glance

**Catalog & Navigation**
- Autocomplete search powered by custom Trie — type 2 letters and see instant DVD/VHS suggestions
- Genre-based filtering via tabs: ALL, ACTION, COMEDY, HORROR, SCIFI, DRAMA
- Format badges on every title: `📀 DVD` / `📼 VHS` — different rental rules per format
- "New Releases" shelf — latest arrivals marked `[NEW]`, restricted by membership plan
- "Staff Picks" curated recommendations by store attendants
- "Last Chance" — titles about to leave the catalog permanently

**Rental & Return**
- Format-aware rental durations — VHS: 3 days, DVD: 5 days (discs are more durable)
- Automatic late fee — $2/day for VHS, $3/day for DVD (higher replacement cost)
- Simultaneous rental limits per membership plan: Bronze=2, Silver=5, Gold=10
- New release waitlist — min-heap ordered by wait time, notified when copies return
- Express return — priority deque processes the most overdue return first
- Full rental history — per-user doubly linked list, navigable in both directions

**Membership Plans**
- **Bronze** (free) — browse catalog, rent up to 2 titles, no new releases
- **Silver** ($9.99/mo) — rent up to 5 titles, access new releases, priority waitlist
- **Gold** ($19.99/mo) — rent up to 10 titles, new releases, reserve upcoming titles
- Each plan gets a color-coded badge + membership card in the Profile page
- Plan upgrades/downgrades processed by store attendants (Employee role)

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
- AES-256-GCM encryption for sensitive data at rest

**Social & Gamification**
- Star ratings — 0 to 5 stars, community average displayed
- Membership plan badge — color-coded card with tier stats
- "Now Showing" — featured title of the week in the header
- Popcorn Points — earned on punctual returns (10 per on-time, -5 per late)

**Administration (Manager / Gerente)**
- Full CRUD on movie catalog — add, edit, remove titles with format (DVD/VHS)
- Full CRUD on users — create, upgrade/downgrade plan, ban
- Audit log viewer — scrollable hash chain with integrity verification
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
- **AES-256-GCM** encryption for sensitive data at rest (audit logs, refresh tokens)

### 2.2 RBAC — 6-tier Hierarchy with Bitmask (Membership Plans + Staff)

Permissions are encoded as bitmask integers checked in O(1) via bitwise `&`.
The system models a real Blockbuster store: customers join a **membership plan** (Bronze / Silver / Gold),
while employees and managers are **store staff** with operational privileges.

```
                    ┌─────────┐
                    │  OWNER  │  0b11111 → All permissions (Dono)
                    └────┬────┘
                    ┌────┴────┐
                    │ MANAGER │  0b01111 → CRUD movies, manage staff (Gerente)
                    └────┬────┘
              ┌─────────┴─────────┐
         ┌────┴────┐        ┌────┴────┐
         │EMPLOYEE │        │  GOLD   │  0b00111 → Rent up to 5, new releases, wishlist
         │(Atend.) │        └────┬────┘
         └────┬────┘        ┌────┴────┐
              │             │ SILVER  │  0b00011 → Rent up to 2, wishlist
              │             └────┬────┘
              │                  │
              └────────┬─────────┘
                  ┌────┴────┐
                  │ BRONZE  │  0b00001 → Browse catalog, wishlist only
                  └─────────┘
```

| Tier | Role (PT-BR) | Bitmask | Max Rentals | New Releases | Wishlist | Audit |
|------|-------------|:-------:|:-----------:|:---:|:---:|:---:|
| Bronze | Cliente Bronze | `0b00001` | 0 | ❌ | ✅ | ❌ |
| Silver | Cliente Prata | `0b00011` | 2 | ❌ | ✅ | ❌ |
| Gold | Cliente Ouro | `0b00111` | 5 | ✅ | ✅ | ❌ |
| Employee | Atendente | `0b00111` | 5 | ✅ | ✅ | ❌ |
| Manager | Gerente | `0b01111` | 10 | ✅ | ✅ | ✅ |
| Owner | Dono | `0b11111` | ∞ | ✅ | ✅ | ✅ |

**Note:** Gold (Cliente Ouro) and Employee (Atendente) share the same bitmask (`0b00111`) for catalog/rental permissions.
The distinction is role-based: Employees are staff who can process returns for **any** customer,
view the full rental ledger, and access the return priority deque. Gold members are customers with
premium rental privileges but no staff capabilities. Enforced via `Role` string field check.

| Permission Bit | Constant | Bronze | Silver | Gold | Employee | Manager | Owner |
|:---:|:---|:---:|:---:|:---:|:---:|:---:|:---:|
| `0b00001` | `PermBrowse` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ |
| `0b00010` | `PermRent` | — | ✅ | ✅ | ✅ | ✅ | ✅ |
| `0b00100` | `PermReserve` | — | — | ✅ | ✅ | ✅ | ✅ |
| `0b01000` | `PermManageUsers` | — | — | — | — | ✅ | ✅ |
| `0b10000` | `PermAdmin` | — | — | — | — | — | ✅ |

When a user attempts an action, the middleware computes:

```go
if session.Permissions & PermRent == 0 {
    return "⛔ ACCESS DENIED — Insufficient clearance"
}
```

### 2.3 Business Rules

**Membership Plans (Planos de Sócio)**

| Plan | Monthly Fee | Max Rentals | New Releases | VHS Late Fee | DVD Late Fee | VHS Duration | DVD Duration |
|------|:----------:|:-----------:|:---:|:---:|:---:|:---:|:---:|
| Bronze | Free | 0 | ❌ | — | — | — | — |
| Silver | $9.99 | 2 | ❌ | $2/day | $3/day | 3 days | 5 days |
| Gold | $19.99 | 5 | ✅ | $2/day | $3/day | 3 days | 5 days |

- Plan upgrades performed by **Attendants (Atendente)** or **Managers (Gerente)** via admin panel
- Plan downgrades also require staff approval
- Membership fees tracked on Profile page (cosmetic — no real payment integration)

**DVD vs VHS Format Rules**

Each movie stores its `Format`: `DVD`, `VHS`, or `Blu-ray`. This affects rental duration,
late fees, and inventory tracking:
- VHS tapes: 3-day rental window (analog tapes wear faster)
- DVDs / Blu-rays: 5-day rental (discs are more durable)
- Late fees: VHS = $2/day, DVD/Blu-ray = $3/day (higher replacement cost)
- Inventory tracked per format (e.g., 3 VHS + 2 DVD copies of The Matrix)
- Format badge on every movie card: `📼 VHS` / `📀 DVD` / `💿 Blu-ray`

**Rental Limits (Limite de Locações Simultâneas)**
- Bronze: 0 active rentals (browse + wishlist only)
- Silver: 2 simultaneous rentals
- Gold / Employee: 5
- Manager: 10
- Owner: unlimited
- Exceeding limit triggers: `"Rental limit reached (X/Y)"` modal

**Late Fees (Multa por Atraso)**
- Auto-calculated on return: `days_overdue × daily_rate` (per format)
- Popcorn Points deducted for late returns (-5 per late)

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
║  ┌─[ALL]──[ACTION]──[COMEDY]──[HORROR]──[SCIFI]──[NEW]──[STAFF PICKS]─┐ ║
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
   │  LOGIN / REGISTER   │
   └─────────┬───────────┘
             ▼
   ┌─────────────────────────────────────────────────────────────┐
   │                       BROWSE CATALOG                         │
   │  (search bar + genre tabs + responsive movie card grid)      │
   └───┬───────────────┬───────────────┬───────────────┬─────────┘
       ▼               ▼               ▼               ▼
┌─────────────┐ ┌─────────────┐ ┌─────────────┐ ┌─────────────┐
│ MOVIE       │ │ MY RENTALS  │ │   PROFILE   │ │ ADMIN PANEL │
│ DETAIL      │ │ (return,    │ │ (stats,     │ │ (manager+)  │
│ (rent,      │ │  history,   │ │  tier badge,│ │             │
│  rate,      │ │  late fees) │ │  popcorn)   │ ├──────┬──────┤
│  waitlist)  │ │             │ │             │ │USERS │MOVIES│
└─────────────┘ └─────────────┘ └─────────────┘ └──┬───┴──┬───┘
                                                   │      │
                                                   └──┬───┘
                                                      ▼
                                               ┌─────────────┐
                                               │ AUDIT LOG   │
                                               │ (hash chain │
                                               │  viewer)    │
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
│   │   └── permissions.go           # Bitmask constants + role hierarchy
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

### Phase 1 — Data Structures (from scratch)

**Goal:** Implement all 8 custom data structures with full test coverage and benchmarks. No third-party collections — everything built with raw Go slices, maps, and nodes.

---

#### Task 1.1: Initialize Go module & project skeleton

- Create project root directory `thelastvideostore/`
- Run `go mod init github.com/thelastvideostore` (or local module path)
- Create directory tree: `internal/ds/{trie,lru,deque,heap,list,bloom,bitmask}`, `internal/crypto`
- Create `cmd/server/` and `cmd/client/` placeholder `main.go` files (just `package main; func main() {}`)
- Verify: `go build ./...` compiles without errors

#### Task 1.2: Implement Bitmask

- Create `internal/ds/bitmask/bitmask.go`
- Define `type Permission uint16`
- Implement functions: `Has(p, flag Permission) bool`, `Set(p, flag Permission) Permission`, `Clear(p, flag Permission) Permission`, `Toggle(p, flag Permission) Permission`
- Define permission constants: `PermBrowse`, `PermRent`, `PermReserve`, `PermManageUsers`, `PermAdmin`
- Define tier constants: `TierBronze = PermBrowse`, `TierSilver = PermBrowse | PermRent`, `TierGold = PermBrowse | PermRent | PermReserve`, `TierEmployee = PermBrowse | PermRent | PermReserve`, `TierManager = PermBrowse | PermRent | PermReserve | PermManageUsers`, `TierOwner = 0b11111`
- Define tier labels map: `map[Permission]string` — "Bronze", "Silver", "Gold", "Employee", "Manager", "Owner"
- Create `internal/ds/bitmask/bitmask_test.go`
- Test table: `struct{ name string; base Permission; flag Permission; wantHas bool; wantSet Permission }`
- Test all combos: Bronze missing Rent, Owner has Admin, set/clear/toggle operations
- Benchmark: `Has`, `Set`, `Clear` — confirm O(1) performance
- Verify: `go test -v -bench=. ./internal/ds/bitmask/`

#### Task 1.3: Implement Doubly Linked List

- Create `internal/ds/list/linkedlist.go`
- Define generic `Node[T any]` struct with `Value T`, `Prev *Node[T]`, `Next *Node[T]`
- Define `List[T any]` struct with `Head *Node[T]`, `Tail *Node[T]`, `Len int`
- Implement: `New[T]()`, `PushBack(v T) *Node[T]`, `PushFront(v T) *Node[T]`, `Remove(node *Node[T]) T`, `PopFront() T`, `PopBack() T`, `Find(pred func(T) bool) *Node[T]`, `Slice() []T`
- Ensure proper nil handling on empty list operations (return zero value + false ok pattern)
- Create `internal/ds/list/linkedlist_test.go`
- Test: push/remove interleaved, pop front/back, find by predicate, iteration via Slice, empty list edge cases
- Benchmark: `PushBack` x 10000, `Find` on 1000-element list
- Verify: `go test -v -bench=. ./internal/ds/list/`

#### Task 1.4: Implement Deque (Ring Buffer)

- Create `internal/ds/deque/deque.go`
- Define `Deque[T any]` struct with ring buffer backing slice: `buf []T`, `head int`, `tail int`, `size int`, `cap int`
- Implement: `New[T](capacity int)`, `PushBack(v T)`, `PushFront(v T)`, `PopFront() T`, `PopBack() T`, `PeekFront() T`, `PeekBack() T`, `Len() int`, `IsEmpty() bool`
- Auto-grow: double capacity when full, handle wrap-around indices correctly
- Return `(T, bool)` for pop/peek on empty deque
- Create `internal/ds/deque/deque_test.go`
- Test: push/pop interleaved mixes, empty deque pops, wrap-around after capacity, grow triggers, peek non-destructive
- Benchmark: push/pop 10000 items from both ends
- Verify: `go test -v -bench=. ./internal/ds/deque/`

#### Task 1.5: Implement Min-Heap

- Create `internal/ds/heap/heap.go`
- Define generic `Heap[T any]` struct with `items []T`, `less func(a, b T) bool` (min-heap by default)
- Implement: `New[T](lessFn)`, `Push(v T)`, `Pop() T`, `Peek() T`, `Len() int`, `IsEmpty() bool`
- Internal: `siftUp(index int)`, `siftDown(index int)`
- Create `internal/ds/heap/heap_test.go`
- Test: sequential push → pop returns sorted order, empty heap pop returns zero value, Peek on empty, priority ordering via custom `less` (e.g., structs ordered by timestamp)
- Benchmark: push 10000 then pop 10000
- Verify: `go test -v -bench=. ./internal/ds/heap/`

#### Task 1.6: Implement Trie (Prefix Tree)

- Create `internal/ds/trie/trie.go`
- Define `TrieNode` struct: `children map[rune]*TrieNode`, `isEnd bool`, `value any` (store movie ID or title)
- Define `Trie` struct: `root *TrieNode`
- Implement: `New()`, `Insert(word string, value any)`, `Search(word string) (any, bool)`, `StartsWith(prefix string) bool`, `Autocomplete(prefix string) []any` (returns all values under prefix subtree via DFS), `Delete(word string) bool`
- Create `internal/ds/trie/trie_test.go`
- Test: insert + search exact match, search missing, startsWith true/false, autocomplete returns all matches for prefix "mat", delete removes exact word but not prefixes, case sensitivity (lowercase only enforced)
- Benchmark: insert 5000 words, autocomplete with 2-char prefix
- Verify: `go test -v -bench=. ./internal/ds/trie/`

#### Task 1.7: Implement LRU Cache

- Create `internal/ds/lru/lru.go`
- Define `entry[K comparable, V any]` struct: `key K`, `value V`
- Define `Cache[K comparable, V any]` struct: `capacity int`, `items map[K]*list.Node[entry[K,V]]`, `order *list.List[entry[K,V]]` (reuses our DoublyLinkedList from Task 1.3)
- Implement: `New[K,V](capacity int)`, `Get(key K) (V, bool)`, `Put(key K, value V)`, `Remove(key K) bool`, `Len() int`, `Contains(key K) bool`
- On Get: move node to front of order (most recently used)
- On Put when full: evict tail of order (least recently used)
- Create `internal/ds/lru/lru_test.go`
- Test: put+get returns value, get missing returns zero, eviction when capacity exceeded (put 4 items capacity 3 → first evicted), update existing key moves to front, remove works, Contains works
- Benchmark: put 10000 items with capacity 1000, get hit rate test
- Verify: `go test -v -bench=. ./internal/ds/lru/`

#### Task 1.8: Implement Bloom Filter

- Create `internal/ds/bloom/bloom.go`
- Define `BloomFilter` struct: `bitset []uint64`, `size uint64` (total bits), `hashCount int`
- Implement: `New(size uint64, hashCount int)`, `Add(data []byte)`, `Contains(data []byte) bool`
- Use double-hashing technique: `h1 = fnv.New64a`, `h2 = fnv.New64` (or murmurhash via hash/fnv + shift)
- Combined hash: `h1 + i*h2` for i from 0 to hashCount-1
- Set/check bits via `bitset[bitIndex/64] & (1 << (bitIndex % 64))`
- Create `internal/ds/bloom/bloom_test.go`
- Test: add string → contains returns true, empty filter contains nothing, multiple adds don't collide, estimate false-positive rate on 1000 items with 10000-bit filter and 3 hash functions
- Benchmark: Add 10000 items, Contains 10000 items
- Verify: `go test -v -bench=. ./internal/ds/bloom/`

#### Task 1.9: Implement Hash Chain (Audit Trail)

- Create `internal/crypto/hashchain.go`
- Define `HashChain` struct: `entries []HashChainEntry`, `lastHash []byte`
- Define `HashChainEntry` struct: `Timestamp int64`, `Action string`, `ActorID string`, `TargetID string`, `Data string`, `Hash []byte`, `PrevHash []byte`
- Implement: `New()`, `Append(action, actorID, targetID, data string) HashChainEntry`, `Verify() bool` (recomputes all hashes from genesis), `GetAll() []HashChainEntry`, `Len() int`
- Hash function: `SHA-256(prevHash || timestamp || action || actorID || targetID || data)`
- Genesis block: first entry has `PrevHash = []byte("GENESIS")`
- Create `internal/crypto/hashchain_test.go`
- Test: genesis entry has correct PrevHash, append links correctly, verify passes on intact chain, verify fails if entry tampered, middle insertion detected
- Benchmark: append 1000 entries + verify
- Verify: `go test -v -bench=. ./internal/crypto/hashchain.go` (or move test to same package)

#### Phase 1 validation:

```bash
go test -v -race -bench=. ./internal/ds/... ./internal/crypto/...
# All tests must pass, all benchmarks must complete, no race conditions
```

---

### Phase 2 — Models & Database Layer

**Goal:** Define all data models as Go structs and implement full CRUD persistence with BoltDB.

---

#### Task 2.1: Install dependencies

- Run `go get github.com/boltdb/bolt` (or `go.etcd.io/bbolt` for maintained fork)
- Run `go get github.com/google/uuid`
- Run `go get golang.org/x/crypto`
- Run `go get github.com/go-chi/chi/v5`
- Run `go get github.com/golang-jwt/jwt/v5`
- Verify: `go mod tidy` succeeds

#### Task 2.2: Create config package

- Create `internal/config/config.go`
- Define `Config` struct: `DBPath string`, `JWTSecret string`, `AESKey string`, `ServerPort string`, `APIBaseURL string`
- Implement `Load() *Config`: reads from env vars with sensible defaults
- Defaults: `DBPath="thelastvideostore.db"`, `ServerPort="8080"`, `APIBaseURL="http://localhost:8080"`
- Add `MustLoad()` variant that panics on missing required vars

#### Task 2.3: Create user model

- Create `internal/models/user.go`
- Define `User` struct: `ID`, `Username`, `PasswordHash`, `Tier` (Permission bitmask), `MaxRentals` int, `RentalCount` int, `Banned` bool, `CreatedAt` int64, `UpdatedAt` int64
- Define JSON tags for API serialization: `json:"id"`, `json:"username"`, `json:"tier"`, `json:"max_rentals"`, `json:"rental_count"`, `json:"banned"` (never expose `PasswordHash`)
- Define `UserResponse` struct (omits password hash for API responses)
- Define helper: `CanRent() bool`, `CanReserve() bool`, `TierName() string`

#### Task 2.4: Create movie model

- Create `internal/models/movie.go`
- Define `Movie` struct: `ID`, `Title`, `Year` int, `Genre` string, `Format` string (VHS, DVD, Blu-ray), `Director`, `Cast` []string, `Synopsis` string, `Rating` float64 (avg), `RatingCount` int, `Available` bool, `CopiesTotal` int, `CopiesAvailable` int, `IsNewRelease` bool, `CoverArt` string (ASCII art placeholder string), `CreatedAt` int64
- Define format constants: `FormatVHS`, `FormatDVD`, `FormatBluRay`
- Define genre constants: `Action`, `Comedy`, `Horror`, `SciFi`, `Drama`, `Thriller`, `Romance`, `Animation`
- Define `MovieResponse` DTO struct for API

#### Task 2.5: Create rental model

- Create `internal/models/rental.go`
- Define `Rental` struct: `ID`, `UserID`, `MovieID`, `MovieFormat` string, `RentedAt` int64, `DueDate` int64, `ReturnedAt` int64 (0 = not returned), `LateFee` float64, `Status` string (active, returned, overdue)
- Define rental status constants: `RentalActive`, `RentalReturned`, `RentalOverdue`
- Define helper: `IsOverdue(now int64) bool`, `CalculateLateFee(now int64) float64` (uses format-specific daily rate: VHS=$2/day, DVD/Blu-ray=$3/day)
- Define helper: `DueDateForFormat(format string, rentedAt int64) int64` (VHS: +3 days, DVD/Blu-ray: +5 days)

#### Task 2.5a: Create wishlist model

- Create `internal/models/wishlist.go`
- Define `WishlistItem` struct: `ID`, `UserID`, `MovieID`, `AddedAt` int64
- Using doubly linked list structure for ordered storage per user

#### Task 2.6: Create audit model

- Create `internal/models/audit.go`
- Define `AuditEntry` struct: `ID`, `Timestamp` int64, `Action` string, `ActorID`, `TargetID`, `Data`, `Hash`, `PrevHash` (mirrors HashChainEntry for DB persistence)
- Define action constants: `ActionLogin`, `ActionLogout`, `ActionRent`, `ActionReturn`, `ActionRegister`, `ActionPromote`, `ActionDemote`, `ActionBan`, `ActionAddMovie`, `ActionEditMovie`, `ActionDeleteMovie`

#### Task 2.7: Create BoltDB store layer

- Create `internal/store/store.go`
- Define `Store` struct wrapping `*bolt.DB`
- Implement `Open(path string) (*Store, error)` — opens BoltDB, creates all buckets: `users`, `movies`, `rentals`, `audit_logs`, `sessions`, `banned`, `movies_by_genre`, `movies_by_title`
- Implement `Close() error`
- Implement helper: `bucketName` constants

#### Task 2.8: Implement user store

- Create `internal/store/users.go`
- Methods on `*Store`:
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

#### Task 2.9: Implement movie store

- Create `internal/store/movies.go`
- Methods on `*Store`:
  - `CreateMovie(movie *models.Movie) error` — stores in `movies` bucket, adds to `movies_by_genre` index, adds to `movies_by_title` index
  - `GetMovieByID(id string) (*models.Movie, error)`
  - `UpdateMovie(movie *models.Movie) error` — updates in all indexes
  - `DeleteMovie(id string) error` — removes from all indexes
  - `ListMovies(genre string, offset, limit int) ([]*models.Movie, int, error)` — paginated; if genre != "" filters by genre index; returns total count
  - `SearchMoviesByPrefix(prefix string, limit int) ([]*models.Movie, error)` — uses `movies_by_title` bucket for prefix scan (BoltDB supports prefix iteration via `Seek`)
  - `GetNewReleases() ([]*models.Movie, error)` — filters by `IsNewRelease = true` (scan all, filter in memory)

#### Task 2.10: Implement rental store

- Create `internal/store/rentals.go`
- Methods on `*Store`:
  - `CreateRental(rental *models.Rental) error` — stores in `rentals` bucket
  - `GetRentalByID(id string) (*models.Rental, error)`
  - `UpdateRental(rental *models.Rental) error`
  - `GetActiveRentalsByUser(userID string) ([]*models.Rental, error)` — iterates all, filters by `UserID` and `Status != returned`
  - `GetRentalHistoryByUser(userID string) ([]*models.Rental, error)` — iterates all, filters by `UserID`, returns all active + returned
  - `GetOverdueRentals() ([]*models.Rental, error)` — iterates all, filters by `DueDate < now && Status == active`
  - `CountActiveRentalsByUser(userID string) (int, error)` — count of non-returned rentals

#### Task 2.11: Implement audit store

- Create `internal/store/audit.go`
- Methods: `AppendAuditEntry(entry *models.AuditEntry) error`, `GetAllAuditEntries() ([]*models.AuditEntry, error)`, `GetAuditEntriesByUser(userID string) ([]*models.AuditEntry, error)`
- Each append: encrypt entry data with AES before storing (call `crypto.AESEncrypt`), store encrypted blob in BoltDB

#### Task 2.12: Implement session store

- Extend `internal/store/store.go` or create `internal/store/sessions.go`
- Methods:
  - `SaveRefreshToken(userID, tokenID string, expiresAt int64) error` — stores in `sessions` bucket
  - `ValidateRefreshToken(userID, tokenID string) (bool, error)` — checks if token exists and not expired
  - `InvalidateRefreshToken(userID, tokenID string) error` — removes from bucket
  - `InvalidateAllUserSessions(userID string) error` — removes all tokens for user
  - `IsTokenRevoked(tokenID string) (bool, error)` — checks `revoked` sub-bucket

#### Task 2.12a: Implement wishlist store

- Create `internal/store/wishlist.go`
- Methods:
  - `AddToWishlist(userID, movieID string) error` — appends to user's wishlist in BoltDB
  - `RemoveFromWishlist(userID, movieID string) error`
  - `GetWishlist(userID string) ([]*models.WishlistItem, error)` — returns ordered items
  - `IsInWishlist(userID, movieID string) (bool, error)`
  - `GetWishlistSize(userID string) (int, error)`
- Wishlist stored as ordered entries per user in `wishlists` bucket

#### Phase 2 validation:

```bash
go test -v ./internal/models/... ./internal/store/...
# Write a store_test.go that opens temp BoltDB, tests CRUD for each entity, closes/deletes temp
```

---

### Phase 3 — Authentication & Security Layer

**Goal:** Implement bcrypt password hashing, JWT token management with refresh rotation, RBAC permission enforcement, brute-force lockout, and AES-256-GCM encryption.

---

#### Task 3.1: Implement password hashing

- Create `internal/auth/password.go`
- Implement `HashPassword(password string) (string, error)` — uses `bcrypt.GenerateFromPassword` with cost 12
- Implement `CheckPassword(hash, password string) bool` — uses `bcrypt.CompareHashAndPassword`
- Create `internal/auth/password_test.go`
- Test: hash produces different string each time, verify same password matches, wrong password doesn't match, empty password rejected

#### Task 3.2: Implement JWT session management

- Create `internal/auth/session.go`
- Define `TokenPair` struct: `AccessToken string`, `RefreshToken string`, `ExpiresAt int64`
- Implement `GenerateTokenPair(userID string, permissions Permission, secret string) (*TokenPair, error)`:
  - Access token: 15-min expiry, claims: `sub=userID`, `perm=permissions`, `iat`, `exp`, `jti` (unique ID)
  - Refresh token: 7-day expiry, claims: `sub=userID`, `jti`, `exp`
  - Sign with HS256
- Implement `ValidateAccessToken(tokenString, secret string) (*Claims, error)` — parses, validates expiry, returns claims
- Implement `ValidateRefreshToken(tokenString, secret string) (*RefreshClaims, error)`
- Create `internal/auth/session_test.go`
- Test: generate valid tokens, validate expired token (manipulate time or create token with -1s expiry), wrong secret fails, malformed token fails

#### Task 3.3: Implement permission enforcement

- Create `internal/auth/permissions.go`
- Re-export bitmask constants from `internal/ds/bitmask` (or import directly — decide which package owns these)
- Define `RequirePermission(userPerms Permission, required Permission) bool` — simple `userPerms & required != 0`
- Define `TierName(perm Permission) string` — returns "Bronze", "Silver", "Gold", "Employee", "Manager", "Owner"
- Define `MaxRentalsForTier(perm Permission) int`:
  - Bronze: 0, Silver: 2, Gold: 5, Employee: 5, Manager: 10, Owner: MaxInt
- Define `CanAccessAdmin(perm Permission) bool` — Manager or Owner

#### Task 3.4: Implement AES-256-GCM encryption

- Create `internal/crypto/aes.go`
- Implement `GenerateAESKey() ([]byte, error)` — 32 random bytes via `crypto/rand`
- Implement `Encrypt(plaintext, key []byte) ([]byte, error)` — AES-256-GCM: generate random nonce, prepend to ciphertext
- Implement `Decrypt(ciphertext, key []byte) ([]byte, error)` — extract nonce, decrypt
- Create `internal/crypto/aes_test.go`
- Test: encrypt then decrypt returns original, different key fails, empty plaintext works, tampered ciphertext detected (GCM auth failure)

#### Task 3.5: Implement brute-force lockout

- Create `internal/auth/lockout.go` (or add to `session.go`)
- Define constants: `MaxAttempts=5`, `LockoutDuration=30*time.Minute`
- Implement `CheckLoginAttempts(store *store.Store, username string) error`:
  - Check if locked → return `ErrAccountLocked` with remaining time
  - Check if attempt count >= MaxAttempts → lock account (save lock expiry), return error
  - If clear → return nil
- Implement `RecordFailedAttempt(store *store.Store, username string) error`
- Implement `RecordSuccessfulLogin(store *store.Store, username string) error` — resets attempts
- Integrate into login handler (Phase 4)

#### Task 3.6: Implement audit log integration

- Create `internal/auth/audit.go` (or add to `internal/crypto/hashchain.go`)
- Connect hash chain to BoltDB store: every state-changing operation appends an entry
- Encrypt audit entries with AES before persisting (call `Encrypt` from Task 3.4)
- Implement `VerifyAuditChain(store *store.Store) (bool, error)` — reads all entries, recomputes hashes, compares

#### Phase 3 validation:

```bash
go test -v ./internal/auth/... ./internal/crypto/...
# JWT round-trip, bcrypt correctness, AES encrypt/decrypt, lockout timer, hash chain integrity
```

---

### Phase 4 — REST API

**Goal:** Build a complete REST API with Chi router, middleware, handlers, DTOs, and server entrypoint.

---

#### Task 4.1: Create DTOs

- Create `api/dto.go`
- Request DTOs:
  - `LoginRequest`: `Username string`, `Password string`
  - `RegisterRequest`: `Username string`, `Password string`
  - `CreateMovieRequest`: `Title string`, `Year int`, `Genre string`, `Director string`, `Cast []string`, `Synopsis string`, `CopiesTotal int`, `IsNewRelease bool`
  - `UpdateMovieRequest`: same fields as Create, all optional via pointers
  - `RentRequest`: `MovieID string`
  - `ReturnRequest`: `RentalID string`
  - `UpdateUserRequest`: `Tier string` (optional), `Banned *bool` (optional)
- Response DTOs (structs with JSON tags):
  - `ErrorResponse`: `Error string`, `Code int`
  - `SuccessResponse`: `Message string`
  - `LoginResponse`: `AccessToken string`, `RefreshToken string`, `User models.UserResponse`
  - `MovieListResponse`: `Movies []models.MovieResponse`, `Total int`, `Page int`, `PageSize int`
- Implement `WriteJSON(w http.ResponseWriter, status int, data interface{})`
- Implement `WriteError(w http.ResponseWriter, status int, message string)`

#### Task 4.2: Implement authentication middleware

- Create `api/middleware.go`
- `AuthMiddleware(secret string, store *store.Store)`:
  - Extract `Authorization: Bearer <token>` header
  - Validate JWT (call `ValidateAccessToken`)
  - Check if token revoked (query `sessions` bucket for revoked JTI)
  - Load user from store by `sub` claim
  - Check user not banned (also check Bloom filter as fast path)
  - Inject `User` + `Permissions` into `context.Context`
  - Return 401 if missing/invalid token || 403 if banned
- `RequirePermission(required Permission)` — middleware factory:
  - Reads permissions from context
  - Calls `bitmask.Has(perms, required)`
  - Returns 403 with `"⛔ ACCESS DENIED — Insufficient clearance"` if check fails
- `RateLimitMiddleware(rate int)` — token bucket:
  - Per-IP counting via in-memory `map[string]*tokenBucket` with mutex
  - 100 req/min default
  - Returns 429 if exceeded
- `CORSMiddleware()` — wraps `chi/cors` with permissive dev defaults
- `LoggingMiddleware()` — logs method, path, status, duration to stdout

#### Task 4.3: Implement auth handlers

- Create `api/auth_handler.go`
- `POST /api/v1/auth/register`:
  - Parse `RegisterRequest`
  - Validate: username 3-20 chars alphanumeric, password 6+ chars
  - Check `UserExists` → 409 if taken
  - Hash password via `auth.HashPassword`
  - Create User with Tier=Bronze (Cliente Bronze), save to store
  - Append to audit hash chain: `ActionRegister`
  - Return 201 with `UserResponse`
- `POST /api/v1/auth/login`:
  - Call `CheckLoginAttempts` → 429 if locked
  - Find user by username → 401 if not found
  - `CheckPassword` → if fail: `RecordFailedAttempt`, return 401
  - `RecordSuccessfulLogin` → resets attempts
  - Check banned (Bloom + DB) → 403
  - Generate `TokenPair` via `auth.GenerateTokenPair`
  - Save refresh token to store (`SaveRefreshToken`)
  - Append to audit: `ActionLogin`
  - Return `LoginResponse` with tokens + user
- `POST /api/v1/auth/refresh`:
  - Accept refresh token from body
  - Validate, check not revoked
  - Invalidate old refresh token (rotation)
  - Generate new token pair
  - Save new refresh token
  - Return new `LoginResponse`
- `POST /api/v1/auth/logout`:
  - Requires JWT auth
  - Invalidate refresh token (or all user sessions)
  - Append to audit: `ActionLogout`
  - Return 200

#### Task 4.4: Implement movie handlers

- Create `api/movie_handler.go`
- `GET /api/v1/movies`:
  - Query params: `genre`, `page` (default 1), `page_size` (default 20)
  - Call `ListMovies(genre, offset, limit)`
  - Return `MovieListResponse`
- `GET /api/v1/movies/search?q=<prefix>`:
  - Requires JWT (any tier)
  - Call `SearchMoviesByPrefix(q, 10)` — uses BoltDB prefix scan on `movies_by_title`
  - Return `[]models.MovieResponse`
- `GET /api/v1/movies/{id}`:
  - Call `GetMovieByID(id)` → 404 if not found
  - Return `MovieResponse`
- `POST /api/v1/movies`:
  - Requires `RequirePermission(PermManageUsers)` (Manager+)
  - Parse `CreateMovieRequest`
  - Validate: title required, year 1900–current, valid genre
  - Create Movie with UUID, `CopiesAvailable = CopiesTotal`, `Available = true`
  - Save to store
  - Append audit: `ActionAddMovie`
  - Return 201 with `MovieResponse`
- `PUT /api/v1/movies/{id}`:
  - Requires Manager+
  - Parse `UpdateMovieRequest`, apply partial updates (only set non-nil fields)
  - Append audit: `ActionEditMovie`
  - Return updated `MovieResponse`
- `DELETE /api/v1/movies/{id}`:
  - Requires Owner only (`PermAdmin`)
  - Delete from store
  - Append audit: `ActionDeleteMovie`
  - Return 200

#### Task 4.5: Implement rental handlers

- Create `api/rental_handler.go`
- `POST /api/v1/rentals/rent`:
  - Requires `RequirePermission(PermRent)` (Bronze cannot rent)
  - Parse `RentRequest`
  - Get user from context
  - Count active rentals → if >= `MaxRentalsForTier(user.Tier)` → 403 "Rental limit reached"
  - Get movie → 404 if not found
  - Check `movie.CopiesAvailable > 0` → 409 "No copies available"
  - Check `movie.IsNewRelease` and user doesn't have `PermReserve` → 403 "Gold plan required for new releases"
  - Create Rental: `DueDate = DueDateForFormat(movie.Format, now)` (VHS: +3d, DVD/Blu-ray: +5d), `Status = active`
  - Decrement `movie.CopiesAvailable`; if 0 → set `Available = false`
  - Increment `user.RentalCount`
  - Update movie + user in store
  - Save rental
  - Append audit: `ActionRent`
  - Return rental with due date
- `POST /api/v1/rentals/return`:
  - Requires `PermRent`
  - Parse `ReturnRequest`
  - Get rental → 404
  - Verify rental belongs to user (or user is Employee+)
  - Set `ReturnedAt = now`
  - If `now > DueDate` → calculate `LateFee = days * 2.00`
  - Increment `movie.CopiesAvailable`; set `Available = true`
  - Decrement `user.RentalCount`
  - Save rental, movie, user
  - Append audit: `ActionReturn`
  - Return rental with late fee if applicable
- `GET /api/v1/rentals/history`:
  - Requires JWT
  - Call `GetRentalHistoryByUser(userID)`
  - Return list of rentals with movie data joined

#### Task 4.6: Implement user handlers (admin)

- Create `api/user_handler.go`
- `GET /api/v1/users`:
  - Requires Manager+
  - Call `ListUsers()`
  - Return list (omit password hashes)
- `POST /api/v1/users`:
  - Requires Manager+
  - Parse `RegisterRequest` + optional tier
  - Same validation as register, but can set initial tier
  - Append audit: action depending on tier
  - Return 201
- `PUT /api/v1/users/{id}`:
  - Requires Manager+
  - Parse `UpdateUserRequest`
  - If tier changed → append audit `ActionPromote`/`ActionDemote`
  - If banned → append audit `ActionBan`, add to Bloom filter
  - Save user
  - Return updated user
- `DELETE /api/v1/users/{id}`:
  - Requires Owner (`PermAdmin`)
  - Delete user from store
  - Append audit
  - Return 200

#### Task 4.7: Implement wishlist handler

- Create `api/wishlist_handler.go`
- `GET /api/v1/wishlist` — requires JWT (Bronze+) — returns user's wishlist with movie details
- `POST /api/v1/wishlist` — requires JWT (Bronze+) — adds movie to wishlist (body: `{movie_id}`)
- `DELETE /api/v1/wishlist/{movieID}` — requires JWT (Bronze+) — removes movie from wishlist
- `GET /api/v1/wishlist/check/{movieID}` — requires JWT — returns `{in_wishlist: true/false}`

#### Task 4.8: Implement audit handler

- Create `api/audit_handler.go`
- `GET /api/v1/audit`:
  - Requires Manager+
  - Query params: `user_id` (optional filter)
  - Call `GetAllAuditEntries` or `GetAuditEntriesByUser`
  - Decrypt entries with AES key
  - Return list

#### Task 4.9: Create router and server entrypoint

- Create `api/router.go`
- Build Chi router:
  - Apply `CORSMiddleware`, `LoggingMiddleware`, `RateLimitMiddleware(100)` globally
  - Group `/api/v1/auth`: register, login (no auth), refresh, logout (JWT)
  - Group `/api/v1/movies`: GET list/search (JWT), POST/PUT/DELETE (JWT + Manager+)
  - Group `/api/v1/rentals`: all require JWT
  - Group `/api/v1/wishlist`: all require JWT (Bronze+)
  - Group `/api/v1/users`: all require JWT + Manager+
  - Group `/api/v1/audit`: JWT + Manager+
  - Health check: `GET /health` returns `{"status":"ok"}`
- Create `cmd/server/main.go`:
  - Load config via `config.Load()`
  - Open BoltDB store via `store.Open(config.DBPath)`
  - Defer store.Close()
  - Create router via `api.NewRouter(store, config)`
  - Start HTTP server on `config.ServerPort`
  - Graceful shutdown on SIGINT/SIGTERM

#### Phase 4 validation:

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

### Phase 5 — TUI Foundation

**Goal:** Set up the Bubble Tea application shell, global state management, theme system, and visual effects.

---

#### Task 5.1: Install TUI dependencies

- Run `go get github.com/charmbracelet/bubbletea`
- Run `go get github.com/charmbracelet/lipgloss`
- Run `go get github.com/charmbracelet/bubbles`
- Verify: `go mod tidy`

#### Task 5.2: Create theme and styles

- Create `tui/styles/theme.go`
- Define color palette constants using `lipgloss.Color`:
  - `Cyan = "#00FFFF"`, `Magenta = "#FF00FF"`, `Yellow = "#FFFF00"`
  - `NeonGreen = "#39FF14"`, `NeonPink = "#FF6EC7"`
  - `Background = "#0A0A2E"` (dark blue), `Surface = "#121240"`, `BorderDim = "#333366"`
  - `Error = "#FF4444"`, `Success = "#44FF44"`, `Warning = "#FFAA00"`
- Define `AppStyle` — full-screen container with background color
- Define `BorderStyle` — lipgloss border with rounded corners, cyan/magenta edge
- Define `TitleStyle` — bold, cyan, large text
- Define `TextStyle`, `DimTextStyle`, `ErrorTextStyle`, `SuccessTextStyle`
- Define tier-specific color map: `TierColors = map[string]lipgloss.Color{...}`

#### Task 5.3: Create visual effects

- Create `tui/styles/effects.go`
- Implement `Scanlines(width, height int) string` — generates alternating lines of semi-transparent `░` pattern over the full terminal dimensions
- Implement `GlitchFrame() string` — returns random character noise `▓▒░█▄▀` (1-3 chars) for temporary glitch overlay; called randomly on page transitions
- Implement `VHSSpinner() []string` — custom spinner frames: `["▌", "▌ ", " ▌", " ▌", "▌ ", "▌", " ▌", " ▌ "]` (tracking artifact)
- Implement `RewindAnimation(tapeName string) string` — returns "◄◄ REWINDING: <tapeName> ... ▌" styled text
- Implement `AccessDeniedOverlay(width, height int) string` — full-screen "⛔ ACCESS DENIED" in red with scanlines

#### Task 5.4: Create global state

- Create `tui/state.go`
- Define `SessionState` struct:
  - `AccessToken string`, `RefreshToken string`
  - `User *models.UserResponse`
  - `Permissions bitmask.Permission`
  - `IsLoggedIn bool`
  - `Cache *lru.Cache[string, interface{}]` — shared cache (capacity 1000)
  - `MovieCache *lru.Cache[string, *models.MovieResponse]` — dedicated movie cache
  - `APIBaseURL string`
- Implement `NewSessionState(apiURL string) *SessionState` — initializes caches
- Implement `Login(tokenPair *auth.TokenPair, user *models.UserResponse)`
- Implement `Logout()` — clears session
- Implement `HasPermission(perm bitmask.Permission) bool`
- Implement `CanAccessAdmin() bool`
- Implement `RefreshAccessToken() error` — calls `/api/v1/auth/refresh` with current refresh token, updates tokens

#### Task 5.5: Create API client

- Create `tui/api_client.go` (or within `state.go`)
- Implement generic `doRequest(method, path string, body interface{}, target interface{}) error`:
  - Uses `net/http` with 10s timeout
  - Sets `Authorization: Bearer <token>` if logged in
  - Sets `Content-Type: application/json`
  - JSON-encodes body if non-nil
  - JSON-decodes response into target
  - On 401 → attempts refresh token → retries once
- Implement typed methods:
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

#### Task 5.6: Create TUI application shell

- Create `tui/app.go`
- Define `Model` struct implementing `bubbletea.Model`:
  - Fields: `currentPage Page`, `session *SessionState`, `width int`, `height int`, `ready bool`, `lastTick time.Time`
- Define `Page` type: `type Page int` with constants:
  - `PageSplash`, `PageLogin`, `PageRegister`, `PageBrowse`, `PageMovieDetail`, `PageMyRentals`, `PageProfile`, `PageAdminUsers`, `PageAdminMovies`, `PageAuditLog`
- Implement `Init() tea.Cmd`:
  - Returns `tea.Batch(tea.EnterAltScreen, tea.ClearScreen, tickCmd())` — ticks for clock/animations
- Implement `Update(msg tea.Msg) (tea.Model, tea.Cmd)`:
  - `tea.WindowSizeMsg` → store width/height, set ready
  - `tea.KeyMsg`:
    - `ctrl+c`, `esc` (on non-modal) → `tea.Quit`
    - Delegate to current page's Update method
  - `tickMsg` → trigger re-render for clock update, return `tickCmd()`
  - Delegate all other messages to current page's Update
- Implement `View() string`:
  - If not ready → "Initializing..."
  - Render: header component + current page View + footer component
- Create `cmd/client/main.go`:
  - Parse CLI flags: `--api-url` (default `http://localhost:8080`), `--debug`
  - Create `SessionState`
  - Create `Model` with `PageSplash`
  - Run `tea.NewProgram(model, tea.WithAltScreen()).Run()`

#### Task 5.5a: Create tick and header

- `tickCmd()` returns `tea.Tick(time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })`
- Header shows: ASCII "THE LAST VIDEO STORE" banner (hardcoded multi-line string with ANSI/lipgloss styling), current time, "NOW SHOWING" if a movie is highlighted

#### Phase 5 validation:

```bash
go build ./cmd/client && ./client
# Splash screen appears
# Clock ticks in header
# Terminal resize recalculates layout
# Ctrl+C quits cleanly
```

---

### Phase 6 — TUI Pages & Components

**Goal:** Build all 10 interactive pages and 9 reusable components.

---

#### Task 6.1: Splash screen

- Create `tui/pages/splash.go`
- 3-second VHS-style intro sequence:
  - Frame 0: "LOADING..." in neon green on black
  - Frame 1: Glitch effect (random characters overlaying)
  - Frame 2: ASCII "THE LAST VIDEO STORE" logo fades in (lipgloss with increasing foreground color brightness)
  - Frame 3: "INSERT MEMBERSHIP CARD █" blinking
- After 3 seconds, auto-transition to `PageLogin`
- Implement as Bubble Tea model: triggers `tea.Tick(3*time.Second, ...)` in `Init()`, transitions on tick message

#### Task 6.2: Login page

- Create `tui/pages/login.go`
- Uses `bubbles/textinput` for username (focused) and password (masked, `EchoMode: textinput.EchoPassword`)
- "TAB" cycles between fields, "ENTER" submits
- Submit: shows spinner via `tea.Batch(spinner.Tick, loginCmd)`
- `loginCmd` calls API client `Login()`:
  - Success → store tokens in session, set user, navigate to `PageBrowse`
  - 401 → show error "Invalid credentials" in red
  - 429 → show "Account locked. Try again in X minutes"
- "ESC" → quit (on splash page) or navigate back
- Link text at bottom: "No card? [R]egister" → pressing R navigates to `PageRegister`

#### Task 6.3: Register page

- Create `tui/pages/register.go`
- Similar to login: username + password + confirm password fields
- Validation: passwords match, username 3+ chars, password 6+ chars
- Submit calls `Register()` API:
  - Success → "Account created! Press any key to login" → navigate to `PageLogin`
  - 409 → "Username already taken"
- "ESC" → back to `PageLogin`

#### Task 6.4: Header component

- Create `tui/components/header.go`
- `HeaderView(width int, session *SessionState) string`:
  - Top border: `══════` full width in cyan
  - Left section: ASCII "THE LAST VIDEO STORE" banner (4 lines)
  - Right section: `🕐 Fri Jun 14 2002  9:48 PM` (dynamic clock)
  - If logged in: `🎫 user@tier BADGE | 🍿 142 pts`
  - Divider: `──── NOW SHOWING: The Matrix ────` (picks random available movie)
  - Bottom border

#### Task 6.5: Footer component

- Create `tui/components/footer.go`
- `FooterView(width int, page Page, session *SessionState) string`:
  - Context-sensitive keybinding bar:
    - Login: `[TAB] switch field  [ENTER] login  [R] register  [ESC] quit`
    - Browse: `[↓↑] navigate  [ENTER] details  [S] search  [R] my rentals  [P] profile  [Q] quit`
    - Admin: adds `[U] users  [M] movies  [A] audit`
    - Movie detail: `[ENTER] rent  [ESC] back  [H] history`
  - Styled with dim text on surface background

#### Task 6.6: Search bar component

- Create `tui/components/searchbar.go`
- Uses `bubbles/textinput` with placeholder "Search movies..."
- On each keypress, debounced (200ms), calls `SearchMovies(prefix)` via API client
- API returns results from BoltDB prefix scan on `movies_by_title`
- Displays dropdown with max 5 suggestions below input
- `ENTER` on suggestion → navigate to `PageMovieDetail` for that movie
- `ESC` clears search and closes dropdown

#### Task 6.7: Movie card component

- Create `tui/components/movie_card.go`
- `MovieCardView(movie *models.MovieResponse, selected bool, width int) string`:
  - Border: magenta if selected, dim if not
  - Content: title (bold, truncated to 20 chars), year `(1999)`, format badge `📼 VHS` / `📀 DVD` / `💿 Blu-ray`, genre badge (colored pill), star rating `★★★★½`, availability `[RENT]` or `[OUT]` or `[NEW]`
  - Fixed-width card: 22x8 characters
  - Selected card gets highlighted border + background

#### Task 6.8: Movie grid component

- Create `tui/components/movie_grid.go`
- `MovieGridView(movies []*models.MovieResponse, selectedIndex int, width, height int) string`:
  - Calculates columns based on terminal width (`width / 22`)
  - Renders grid of `MovieCardView` components
  - Selected index highlighted
  - Handles arrow key navigation: up/down/left/right within grid (wrap or stop at edges)
  - Page up/down for next/prev page of results

#### Task 6.9: Tabs component

- Create `tui/components/tabs.go`
- `TabsView(tabs []string, activeIndex int, width int) string`:
  - Genre tabs: `ALL | ACTION | COMEDY | HORROR | SCIFI | DRAMA | NEW | STAFF PICKS`
  - Active tab: cyan background, bold text
  - Inactive: dim text
  - Styled with lipgloss borders connecting to content below
- `TabWidth` calculates equal widths filling available space

#### Task 6.9a: Wishlist sidebar component

- Create `tui/components/wishlist_sidebar.go`
- `WishlistSidebarView(wishlist []*WishlistItem, width int) string`:
  - Right-side panel showing user's wishlist items
  - Each item: movie title truncated + format badge + availability indicator (`🟢 Available` / `🔴 Rented out`)
  - "Available now!" highlighted entry when a wishlisted title becomes available
  - Quick-rent: press `W` on highlighted movie to rent directly from wishlist
  - Empty state: "Your wishlist is empty. Browse and press [W] to add titles."
  - Styled with dim border, scrollable if many items

#### Task 6.10: Browse page

- Create `tui/pages/browse.go`
- Combines: searchbar + tabs + movie grid + wishlist sidebar
- State: `selectedGenre int`, `movies []MovieResponse`, `selectedMovie int`, `searchMode bool`, `page int`, `showWishlist bool`
- `Init()`: fetches movies for default genre (ALL) via `GetMovies("", 1)`, fetches wishlist
- Wishlist sidebar (right panel): shows user's wishlisted titles, "Available now!" indicator, quick-rent shortcut
- Keybindings:
  - `←→`: switch tabs → refetch movies for genre
  - `↓↑`: navigate grid
  - `ENTER`: navigate to `PageMovieDetail`
  - `/`: focus search bar
  - `W`: add selected movie to wishlist / toggle wishlist sidebar
  - `R`: navigate to `PageMyRentals`
  - `P`: navigate to `PageProfile`
  - If admin: `U` → `PageAdminUsers`, `M` → `PageAdminMovies`, `A` → `PageAuditLog`
- Admin links only visible if `session.HasPermission(PermManageUsers)`

#### Task 6.11: Movie detail page

- Create `tui/pages/movie_detail.go`
- Full-screen movie view:
  - Title (large, bold)
  - `[NEW RELEASE]` or `[AVAILABLE]` or `[RENTED OUT]` badge
  - Year · Genre · Director
  - Star rating: `★★★★½ (4.5/5 from 1,247 ratings)`
  - Synopsis (wrapped text, 3-4 lines)
  - Cast: comma-separated
  - Copies available: `📼 3 of 5 copies available`
- Actions:
  - `ENTER` → rent movie (calls `RentMovie` API):
    - Success → "📼 RENTED! Due: Jun 17 2002" modal, then navigate to browse
    - 403 "Limit reached" → show error modal
    - 403 "Gold plan required" → show promo to upgrade
    - 409 "No copies" → offer "Join waitlist? [Y/N]"
  - `ESC` → back to browse

#### Task 6.12: My rentals page

- Create `tui/pages/my_rentals.go`
- Fetches `GetRentalHistory()` on init
- Lists active rentals at top with:
  - Movie title, rental date, due date, status `🟢 Active` / `🔴 Overdue`, late fee if applicable
- Lists rental history below (returned) with `ReturnedAt` date
- Selected rental can be returned:
  - Press `ENTER` on active rental → confirmation modal "Return The Matrix?"
  - Confirm → calls `ReturnMovie` API
  - Success: "📼 Returned! Late fee: $4.00" or "📼 Returned on time! +10 popcorn points"
  - Movie grid and rental count refresh
- `ESC` → back to browse

#### Task 6.13: Profile page

- Create `tui/pages/profile.go`
- Membership card view:
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
  ╚══════════════════════════╝
  ```
- Membership plan badge in corresponding color (Bronze=#CD7F32, Silver=#C0C0C0, Gold=#FFD700, Employee=magenta, Manager=yellow, Owner=cyan)
- Popcorn points mock calculation: 10 per on-time return, -5 per late
- Stats pulled from rental history, grouped by format (DVD/VHS/Blu-ray)
- `ESC` → back to browse
- `L` → logout → clear session → navigate to `PageLogin`

#### Task 6.14: Modal component

- Create `tui/components/modal.go`
- `ModalView(title, message string, width, height int) string`:
  - Overlay: dimmed background over current page
  - Centered bordered box with title (bold) and message
  - Buttons: `[ENTER] Confirm  [ESC] Cancel`
- `AccessDeniedModal(width, height int) string`:
  - Same layout but red-tinted
  - Title: `⛔ ACCESS DENIED`
  - Message: `Insufficient clearance level`
  - Only shows `[ESC] Dismiss`

#### Task 6.15: Spinner component

- Create `tui/components/spinner.go`
- `VHSSpinnerView() string` — returns current spinner frame
- Used during API calls: login, register, rent, return
- Integrated into pages via `tea.Batch(spinnerTickCmd, apiCallCmd)`

#### Task 6.16: Badge component

- Create `tui/components/badge.go`
- `TierBadgeView(tierName string) string`:
  - Color-coded pill: `[ BRONZE ]` (bronze brown), `[ SILVER ]` (silver gray), `[ GOLD ]` (gold yellow), `[ ATENDENTE ]` (magenta), `[ GERENTE ]` (yellow), `[ DONO ]` (cyan)
  - Styled with lipgloss background + foreground + padding

#### Task 6.17: Admin users page

- Create `tui/pages/admin_users.go`
- Requires `PermManageUsers` — if insufficient, show `AccessDeniedModal`
- Fetches `GetUsers()` from API
- Displays table:
  - Columns: Username | Tier | Rentals | Banned | Actions
  - Each row selectable with `↓↑`
- Actions on selected user:
  - `P` → promote (increment tier, max Owner)
  - `D` → demote (decrement tier, min Bronze)
  - `B` → toggle ban
  - Confirmation modal for each action
- Calls `UpdateUser` API, refreshes list on success

#### Task 6.18: Admin movies page

- Create `tui/pages/admin_movies.go`
- Requires `PermManageUsers`
- Table of all movies: Title | Year | Genre | Copies | Available
- Actions:
  - `A` → add movie form (text inputs for all fields) → calls `CreateMovie`
  - `ENTER` → edit selected movie (populated form) → calls `UpdateMovie`
  - `D` → delete movie (confirmation modal, Owner only) → calls `DeleteMovie`
- Form navigation: TAB between fields, ENTER to submit, ESC to cancel

#### Task 6.19: Audit log page

- Create `tui/pages/audit_log.go`
- Requires `PermManageUsers`
- Fetches `GetAuditEntries()` from API
- Displays scrollable list:
  - Each entry: `[timestamp] ACTION | Actor: user | Target: target | Hash: a1b2c3...`
  - Hash chain verification status at top: `✅ Chain intact (142 entries)` or `❌ Chain broken!`
  - Each entry shows truncated PrevHash → Hash link
- Scroll with `↓↑`, `PgUp/PgDn`, `Home/End`
- `V` → verify chain against API (triggers recomputation)
- `ESC` → back to browse

#### Phase 6 validation:

```bash
go run ./cmd/client
# Full manual walkthrough:
# 1. Splash → Login (as bronze) → Browse (see grid)
# 2. Search "mat" → see Matrix, Matilda, Match Point
# 3. Click Matrix → Movie Detail → Rent → ACCESS DENIED (bronze)
# 4. Logout → Login as gold → Rent Matrix → success
# 5. View My Rentals → see Matrix due date (format-specific)
# 6. Return Matrix → see late fee or on-time confirmation
# 7. Profile → see plan badge (Gold), stats, popcorn points
# 8. Logout → Login as manager (Gerente) → Admin Users → upgrade silver to Gold
# 9. Admin Movies → add a new Blu-ray title
# 10. Audit Log → verify chain intact
# 11. Login as banned → "Account suspended"
# 12. All pages responsive to terminal resize
```

---

### Phase 7 — Seed Data & Integration Testing

**Goal:** Populate the system with realistic data and test all end-to-end flows.

---

#### Task 7.1: Define seed data — movies

- Create `data/movies.json` (or hardcode in `data/seed.go`)
- ~40 real movies from 1980s–2000s:
  - The Matrix (1999), Fight Club (1999), Pulp Fiction (1994), Jurassic Park (1993), The Shawshank Redemption (1994), The Dark Knight (2008), Inception (2010), Forrest Gump (1994), The Godfather (1972), Schindler's List (1993), Goodfellas (1990), The Silence of the Lambs (1991), Se7en (1995),The Usual Suspects (1995), Léon: The Professional (1994), American History X (1998), Saving Private Ryan (1998), The Green Mile (1999), Gladiator (2000), Memento (2000), The Lord of the Rings trilogy (2001-2003), Kill Bill (2003), Eternal Sunshine (2004), The Departed (2006), No Country for Old Men (2007), There Will Be Blood (2007), WALL-E (2008), Inglourious Basterds (2009), District 9 (2009), Blade Runner (1982), Back to the Future (1985), Die Hard (1988), Terminator 2 (1991), Toy Story (1995), The Big Lebowski (1998), The Truman Show (1998), American Beauty (1999), Requiem for a Dream (2000), Spirited Away (2001), City of God (2002)
- Each movie needs: title, year, genre, director, cast (3-5 real names), synopsis (2-3 sentences), rating (3.0-5.0), rating count (100-5000), copies total (2-10), format (mix of DVD/VHS/Blu-ray), is_new_release (first 6 marked true)

#### Task 7.2: Define seed data — users

- Create 7 test users (hardcoded in `data/seed.go`):
  ```
  bronze   / password1  → TierBronze,  MaxRentals=0  (Cliente Bronze — browse + wishlist only)
  silver   / password2  → TierSilver,  MaxRentals=2  (Cliente Prata — rent up to 2, wishlist)
  gold     / password3  → TierGold,    MaxRentals=5  (Cliente Ouro — new releases, waitlist, wishlist)
  employee / password4  → TierEmployee,MaxRentals=5  (Atendente — staff: process any return)
  manager  / password5  → TierManager, MaxRentals=10 (Gerente — manage users, movies, view audit)
  owner    / password6  → TierOwner,   MaxRentals=99 (Dono — all permissions)
  banned   / password7  → TierBronze,  Banned=true   (Blocked account)
  ```
- All passwords hashed with bcrypt
- Add banned user to Bloom filter

#### Task 7.3: Implement seed script

- Create `data/seed.go`
- Package `main` (runnable)
- Functions:
  - `seedMovies(store *store.Store)` — iterates movie list, creates each in BoltDB
  - `seedUsers(store *store.Store)` — iterates user list, hashes passwords, creates in BoltDB
  - `main()`:
    - Load config
    - `os.Remove(config.DBPath)` to start fresh
    - Open store
    - Call seed functions
    - Print summary: "Seeded 40 movies and 7 users."
    - Close store

#### Task 7.4: Integration test flow 1 — Bronze denied

- Login as bronze → Browse → Select movie → Press ENTER to rent
- Expected: Modal "⛔ ACCESS DENIED — Insufficient clearance"
- Verify: No rental created, audit log shows denied attempt

#### Task 7.5: Integration test flow 2 — Silver rental limit

- Login as silver → Rent movie 1 (DVD) → Rent movie 2 (VHS) → Try to rent movie 3
- Expected: Modal "Rental limit reached (2/2)"
- Verify: Only 2 active rentals in DB; format-specific durations applied (3d VHS, 5d DVD)

#### Task 7.6: Integration test flow 3 — Banned user

- Login as banned
- Expected: "Account suspended. Contact store management."
- Verify: Bloom filter check passes (banned flag detected), JWT not issued

#### Task 7.7: Integration test flow 4 — Plan upgrade (Silver → Gold)

- Login as manager (Gerente) → Admin Users → Select silver → Press P to promote
- Expected: Confirmation modal "Upgrade silver to Gold?" → Confirm → User tier updated
- Login as (formerly silver) → Verify can now rent 5 movies, see new releases, join waitlist
- Verify: Audit log entry `ActionPromote` recorded

#### Task 7.8: Integration test flow 5 — Audit chain

- Login as manager → Audit Log → Press V to verify
- Expected: "✅ Chain intact" with entry count
- Tamper test (manual): corrupt one audit entry hash in BoltDB → Verify → "❌ Chain broken at entry #42"

#### Task 7.9: Integration test flow 6 — Employee return with deque

- Login as employee (Atendente) → Return overdue movies for multiple customers
- Expected: Most overdue is processed first (deque pop from back)
- Verify: API returns list sorted by priority; late fees auto-calculated per format

#### Task 7.10: Integration test flow 7 — New release waitlist

- Login as gold → Try to rent sold-out new release
- Expected: "Join waitlist?" modal → Confirm → Added to heap with timestamp
- Verify: Heap peek returns user with oldest timestamp

#### Phase 7 validation:

```bash
go run ./data/seed.go                       # Seeds database
go run ./cmd/server &                        # Start API
go run ./cmd/client                          # Start TUI
# Execute all 7 integration test flows manually
```

---

### Phase 8 — Deployment & Final Polish

**Goal:** Containerize, deploy to Render, cross-compile binaries, write documentation, and polish.

---

#### Task 8.1: Create Makefile

- Create `Makefile`
- Targets:
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

#### Task 8.2: Create Dockerfile

- Create `Dockerfile.server`
- Multi-stage build:
  - Stage 1 (builder): `golang:1.22-alpine`, copy source, `go build -o /server ./cmd/server`
  - Stage 2 (runtime): `alpine:3.19`, copy binary, expose 8080, run
- Ensure static linking: `CGO_ENABLED=0 GOOS=linux`
- Set `JWT_SECRET` and `AES_KEY` as build args with defaults for demo

#### Task 8.3: Create render.yaml

- Create `render.yaml`
- Service definition:
  - Type: `web`
  - Name: `thelastvideostore-api`
  - Runtime: `docker`
  - Dockerfile path: `./Dockerfile.server`
  - Health check path: `/health`
  - Env vars: `JWT_SECRET` (generateValue), `AES_KEY` (generateValue), `DB_PATH=/data/thelastvideostore.db`
  - Disk: mount at `/data` for BoltDB persistence

#### Task 8.4: Cross-compile and test binaries

```bash
make build-all
# Verify Linux binary on local machine:
./bin/thelastvideostore-linux
# Verify Windows binary (if WSL or VM available):
# Copy to Windows, run thelastvideostore.exe
# Or verify with: file bin/thelastvideostore.exe
# Expected: PE32+ executable (console) x86-64
```

#### Task 8.5: Final code cleanup

- Run `make fmt` → all files formatted
- Run `make lint` → zero warnings
- Run `make test` → all tests pass, coverage > 70%
- Remove any debug prints, unused imports, commented-out code
- Ensure no hardcoded secrets (all from env/config)
- Verify sensitive files in .gitignore: `*.db`, `.env`, `bin/`

#### Task 8.6: Create .gitignore

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

#### Task 8.7: Write README.md

- Create `README.md`
- Sections:
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

#### Task 8.8: Demo preparation

- Prepare demo script for presentation:
  1. Start server + seed data
  2. Launch client, show splash screen (VHS-style intro)
  3. Login as bronze → browse catalog → try to rent → ACCESS DENIED
  4. Demo search autocomplete: type "mat" → Trie shows Matrix, Matilda, Match Point
  5. Logout → login as silver → rent 2 movies (1 DVD + 1 VHS, different due dates)
  6. Add a movie to wishlist → show wishlist sidebar
  7. View My Rentals → see format badges, due dates, late fee warnings
  8. Return a movie → see Popcorn Points earned
  9. View Profile → membership card, tier badge, rental stats
  10. Logout → login as manager (Gerente) → upgrade silver to Gold plan → verify new limits
  11. Show Audit Log → verify hash chain integrity
  12. Login as banned → "Account suspended"
  13. Admin movie CRUD: add a new Blu-ray title
  14. Terminal resize demo (responsive movie grid)
  15. Mention cross-platform: show Linux + Windows binaries

#### Task 8.9: Optional polish items

- Add sound effects toggle (beep on rent/return) via `\a` bell character
- Easter egg: Konami code (↑↑↓↓←→←→BA) shows secret "Employee Picks" menu
- ASCII movie posters (hardcoded simple art for top 5 movies)
- On-exit animation: "BE KIND, REWIND" in large ASCII

#### Phase 8 validation:

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
| 1 | **Interface** | Full Bubble Tea TUI with 10 interactive screens, CRT effects, search, grids, modals | Navigate catalog → rent → return → view profile — all in terminal |
| 2 | **Modo de segurança de acesso** | bcrypt + JWT + RBAC bitmask + brute-force lockout + AES-256-GCM + Bloom filter ban list | Show login fail → lockout → ACCESS DENIED modal → manager promotes user → access granted |
| 3 | **Cybersecurity** | 6-layer security: hashing, token auth, bitmask RBAC, encryption at rest, immutable audit via hash chain, input sanitization | Demonstrate hash chain integrity check, AES-encrypted audit entries, Bloom filter banning |
| 4 | **Data Structures** | 8 structures implemented from scratch: Trie, LRU, Deque, MinHeap, DoublyLinkedList, BloomFilter, Bitmask, HashChain | Show `_test.go` passing, explain each structure's role with visual examples |
| 5 | **Read file line by line** | BoltDB stores movies persistently; API reads paginated results; TUI renders each as a card | Browse catalog with search autocomplete (Trie in action) |
| 6 | **Allow only authorized** | JWT middleware + bitmask on every route; 403 returned with structured error | Bronze (Cliente Bronze) tries to rent → "Permissão Negada"; Silver/Gold rents successfully |
| 7 | **Show user and file data** | Profile screen: username, membership plan badge, rental history (linked list), popcorn points | Navigate to Profile, scroll rental history, show plan badge (Bronze/Silver/Gold) |
| 8 | **User registration via file/DB** | BoltDB persistent store + `/auth/register` endpoint | Register new user → login → Bronze plan automatically assigned |
| 9 | **Add/remove access** | Admin user panel: upgrade plan (Silver→Gold), downgrade, ban (add to Bloom filter) | Manager upgrades Bronze → Silver → Gold; new limits immediately active |
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

*Document version 4.1 — The Last Video Store Project Implementation Plan — Last updated: 2026-06-14*
