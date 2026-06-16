# The Last Video Store — Features & Workflow

## Overview

A retro VHS-themed video rental store TUI application — built for academic demonstration of **Cybersecurity & Data Structures**. Browse, rent, return, wishlist movies, earn Popcorn Points, purchase premium subscription tiers, and redeem rewards — all through a terminal interface backed by a REST API and BoltDB.

**Two-tier access control**: Premium subscription tiers (Wood→Diamond) govern rental perks and costs. RBAC roles (Employee→Owner) govern admin access — both are separate, demonstrating layered security.

---

## Scope Compliance — Requirements Mapping

### Original requirements
> *"Develop a system with an interface and access security. Implement a program using Cybersecurity and Data Structures. Create a system that uses Access Security and a data structure to control access and store data."*

### Security challenges and how they were met

| Challenge | Implementation |
|-----------|---------------|
| **a) Read the file/database line by line** | BoltDB with `First()`/`Next()` cursors — sequential record iteration with middleware-enforced access control |
| **b) Only allow authorized users to read** | JWT + 6-bit RBAC bitmask — every route checks `RequirePermission(flag)`; no token = HTTP 401, no permission = HTTP 403 |
| **c) Display the user and the data** | Header shows `🎫 username | 🏷️ Role | 🍿 pts | 💵 balance`; catalog and audit log display data based on permissions |
| **d) User registration via file/database** | Register screen (`POST /api/v1/auth/register`) + Admin Users panel (`PUT /api/v1/users/{id}`) for promote/demote/ban |
| **e) "Permissão Negada" (Access Denied)** | Dedicated `⛔ ACCESS DENIED` screen when attempting unauthorized routes (e.g., Silver user tries `Ctrl+A` → "Manager role or higher required") |
| **f) Add or remove user access** | Admin Users → keys `P` (promote) / `D` (demote) / `B` (ban) — modify RBAC tier via API |

### Cybersecurity Technologies

| Layer | Technology | Detail |
|--------|-----------|---------|
| Password hashing | bcrypt | Cost 12, per-password salt |
| Authentication | JWT HS256 | Access token (15min) + Refresh token (7 days) with rotation |
| Two-factor auth | TOTP RFC 6238 | HMAC-SHA1, ±30s window, AES-256-GCM encrypted secret |
| Authorization | RBAC Bitmask | 6-bit permission flags, O(1) check |
| Attack protection | Lockout | 5 failed logins = 30min lock; 3 failed TOTP = 10min lock |
| Attack protection | Rate limiting | Per-IP rate limiter on all API routes |
| Integrity | SHA-256 Hash Chain | Immutable audit trail with chain verification |
| Integrity | Bloom Filter | Probabilistic fast banned-user rejection |
| Encryption | AES-256-GCM | TOTP secrets and sensitive data at rest |

### Data Structures (9 custom implementations, zero `container/*`)

| Structure | Application |
|-----------|------------|
| **Bitmask** (6-bit) | RBAC access control — stores permissions in 6 bits |
| **Doubly Linked List** | Rental history — O(1) insert/remove |
| **Deque** (Ring Buffer) | Staff priority return queue |
| **Min-Heap** | Waitlist — ordered by wait time |
| **Trie** (Prefix Tree) | Movie title autocomplete — O(k) |
| **LRU Cache** | Session + movie caching — O(1) |
| **Bloom Filter** | Fast banned-user check — O(k) |
| **Hash Chain** (SHA-256) | Immutable audit trail — O(1) append |
| **Undirected Weighted Graph** | Co-rental recommendations — O(V+E) |

---

## Screens (16 total)

| Screen | Access | Description |
|--------|--------|-------------|
| **Splash** | All | Animated figlet-style banner, `ENTER` to start |
| **Login** | All | Username + password, TOTP challenge if enabled |
| **Register** | All | Create new account (3 fields, validations) |
| **TOTP** | 2FA users | 🔒 6-digit authenticator code entry |
| **Browse** | Authenticated | Catalog grid, pagination, search, staff picks, last chance |
| **Detail** | Authenticated | Synopsis, rating, cast, choose payment method, rent/waitlist |
| **Rentals** | Authenticated | Active + history, countdown timer, extend, return with fees |
| **Profile** | Authenticated | Membership card, role badge, subscription tier, money + 🍿 |
| **Tier Shop** | Authenticated | Purchase/upgrade premium subscription (Wood→Diamond) |
| **Wishlist** | Authenticated | Personal waitlist with availability indicators |
| **Rewards** | Authenticated | 🍿 Popcorn Points shop: redeem for collectibles, free rentals |
| **Inventory** | Authenticated | View owned collectibles |
| **Access Denied** | Authenticated | ⛔ Full-screen denial with role requirement message |
| **Admin Movies** | Manager+ | Movie CRUD, staff pick toggle, paginated list |
| **Admin Users** | Supervisor+ | Promote/demote RBAC role, toggle ban |
| **Audit Log** | Supervisor+ | SHA-256 hash chain viewer, chain integrity verification |

---

## Navigation Flow

```
Splash ──ENTER──→ Login ──login──→ Browse ──┬── Detail ──── rent / waitlist
                     ↑                      ├── Rentals ─── return / extend
                     │                      ├── Profile ─── logout → Login
                Register                    │   ├── Tier Shop (T)
                     │                      │   ├── Rewards Shop (M)
                     └──────────────────────│   └── Inventory (I)
                                            ├── Wishlist ─── view / remove
                                            ├── ⛔ Access Denied (when unauthorized)
                                            ├── Admin Movies (Ctrl+A, Manager+)
                                            ├── Admin Users (Ctrl+U, Supervisor+)
                                            └── Audit Log (Ctrl+G, Supervisor+)
```

---

## Keybindings

### Global
| Key | Action |
|-----|--------|
| `Ctrl+C` / `Ctrl+D` | Quit |
| `Q` | Back to parent screen (see navigation hierarchy below) |

**Q-back hierarchy:**
| From | Returns to |
|------|-----------|
| Detail, Rentals, Wishlist | Browse |
| Rewards Shop, Tier Shop, Inventory | Profile |
| Admin Movies, Users, Audit Log | Browse |
| Access Denied | Browse |
| Movie Form (`ESC`) | Admin Movies |
| Profile | Browse |

### Browse (Main Catalog)
| Key | Action |
|-----|--------|
| `↑↓` / `J` `K` | Navigate grid |
| `ENTER` / `D` | Open movie detail |
| `R` | My Rentals |
| `P` | Profile |
| `V` | View Wishlist |
| `/` | Search movies (live prefix with Trie backend) |
| `[` / `]` | Genre tabs (Action, SciFi, Horror, Comedy, Drama, Thriller, Romance, Animation) |
| `S` | Staff Picks mode |
| `L` | Last Chance mode |
| `A` | All catalog |
| `N` / `B` | Next / Previous page |
| `F5` | Refresh |
| `Ctrl+A` | Admin Movies (Manager+) |
| `Ctrl+U` | Admin Users (Supervisor+) |
| `Ctrl+G` | Audit Log (Supervisor+) |

### Search Mode (`/`)
| Key | Action |
|-----|--------|
| Type | Live prefix search via API |
| `↑↓` / `J` `K` | Navigate results |
| `ENTER` | Open selected movie |
| `ESC` | Cancel search |

### Detail
| Key | Action |
|-----|--------|
| `ENTER` | Open highlighted related movie, or rent current movie |
| `↑↓` / `J` `K` | Navigate franchise/related movies |
| `T` | Use free ticket (🎟️) |
| `M` | Pay with money (💵) |
| `ESC` | Cancel payment choice |
| `W` | Add to waitlist |
| `F5` | Refresh |
| `Q` | Back to Browse |

### Rentals
| Key | Action |
|-----|--------|
| `↑↓` / `J` `K` | Navigate |
| `ENTER` | Return selected rental |
| `E` | Extend due date (+1 min, costs 30🍿) |
| `Q` | Back to browse |

### Profile
| Key | Action |
|-----|--------|
| `L` | Logout |
| `T` | Tier Shop |
| `M` | Rewards Shop |
| `I` | Inventory |
| `Q` | Back to browse |

### Admin Screens
| Screen | Key | Action |
|--------|-----|--------|
| Movies | `A` | Add movie via form |
| Movies | `ENTER` | Edit movie via form |
| Movies | `D` | Delete movie |
| Movies | `S` | Toggle staff pick ★ |
| Movies | `N`/`B` | Page navigation |
| Users | `P` | Promote RBAC role |
| Users | `D` | Demote RBAC role |
| Users | `B` | Toggle ban |
| Audit | `V` | Verify SHA-256 chain integrity |
| Audit | `↑↓` | Navigate entries |

---

## Premium Subscription Tiers

Users start with **Wood** (free). Purchase higher tiers from Profile → `T`.

| Tier | Price | Free Rentals/mo | Max Concurrent | New Releases | Late Fees |
|------|------:|:--------------:|:------------:|:---:|:---:|
| Wood | Free | 0 | 5 | No | Yes |
| Bronze | $9.99 | 1 | 10 | No | Yes |
| Silver | $19.99 | 3 | 15 | Yes | Yes |
| Gold | $29.99 | 5 | 25 | Yes | **Waived** |
| Diamond | $49.99 | 10 | 50 | Yes | **Waived** |

- Tiers are **separate** from RBAC admin roles — a Wood subscriber can have Manager admin access
- Free rental allocation from tier is used first when pressing `T` (ticket)
- Users choose between ticket and money at rental time

---

## Dual Currency System

| Currency | Earned By | Used For |
|----------|-----------|----------|
| **💵 Money** | Seed starting balance ($50–$100) | Renting beyond tier allowance, buying premium tiers |
| **🍿 Popcorn Points** | Returns (+10 on-time, -5 late), Popcorn Bucket bonus (+5) | Rewards shop, extend rentals (30🍿) |

### Movie-Specific Rental Pricing
| Category | VHS/DVD | Blu-ray |
|----------|---------|---------|
| New releases & 2020+ | $5.99 | $6.99 |
| 2000–2019 | $3.99 | $4.99 |
| Pre-2000 | $2.99 | $3.99 |

---

## RBAC Roles (Admin Access)

7-tier 6-bit bitmask — O(1) permission checks:

| Role | Bitmask | Access |
|------|:------:|--------|
| Bronze | `0b000001` | Browse catalog |
| Silver | `0b000011` | Browse + Rent |
| Gold | `0b000111` | + New Releases |
| Employee | `0b010111` | + Staff tools |
| Supervisor | `0b011111` | + User management + Audit |
| Manager | `0b111111` | + Movie management |
| Owner | `0b111111` | Full access |

---

## Due Dates (Demo-Friendly Minutes)

| Format | Due In | Late Fee rate |
|--------|--------|--------------|
| VHS | 2 min | ~$0.20/min |
| DVD | 3 min | ~$0.30/min |
| Blu-ray | 4 min | ~$0.40/min |
| **Extend** | +1 min | costs 30🍿 |

Full lifecycle (rent → overdue → late fee → extend → return) demonstrable in ~5 minutes.

---

## Security Features (Cybersecurity Scope)

| Layer | Feature | Detail |
|-------|---------|--------|
| **Authentication** | bcrypt | Cost factor 12, salt per password |
| **Authentication** | JWT | HS256, 15min access + 7-day refresh rotation |
| **Authentication** | TOTP 2FA | RFC 6238, HMAC-SHA1, 30s window ±1, AES-256-GCM encrypted secrets |
| **Authorization** | RBAC Bitmask | 6-bit permission flags, O(1) check via `p & flag != 0` |
| **Authorization** | Middleware chain | JWT validation → ban check → permission check → handler |
| **Attack Protection** | Brute-force lockout | 5 failed logins = 30min lock, 3 failed TOTP = 10min lock |
| **Attack Protection** | Rate limiting | Per-IP rate limiter on all API routes |
| **Integrity** | Hash Chain Audit | SHA-256 Merkle-Damgård chaining, chain verification endpoint |
| **Integrity** | Bloom Filter | O(k) banned-user fast rejection pre-DB lookup |
| **Encryption** | AES-256-GCM | TOTP secrets, audit-sensitive data at rest |

---

## Data Structures (All Custom, No `container/*`)

| Structure | Application | Complexity |
|-----------|------------|:---:|
| **Bitmask** (6-bit) | RBAC permissions | O(1) |
| **Doubly Linked List** | Rental history ordering | O(1) insert/remove |
| **Deque** (Ring Buffer) | Staff return priority queue | O(1) push/pop |
| **Min-Heap** | New release waitlist ordering | O(log n) |
| **Trie** (Prefix Tree) | Movie title autocomplete search | O(k) |
| **LRU Cache** | Session + movie caching | O(1) |
| **Bloom Filter** | Banned user fast check | O(k) |
| **Hash Chain** (SHA-256) | Immutable audit trail | O(1) append |
| **Undirected Weighted Graph** | Co-rental recommendations | O(V+E) |

---

## Popcorn Points Rewards

| Item | 🍿 Cost | Effect |
|------|:------:|--------|
| Popcorn Bucket | 50 | +5 bonus points on every future return |
| Blank VHS Tape | 75 | Collectible (inventory) |
| Movie Poster | 100 | Collectible (inventory) |
| Store T-Shirt | 150 | Collectible (inventory) |
| Free Rental Coupon | 200 | +1 free rental (bypasses limit, no fees) |
| Private Screening | 500 | +5 free rentals |
| Tier Upgrade | 1000 | Promote RBAC role one level (up to Gold) |

---

## API Endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| `POST` | `/api/v1/auth/register` | — | Create account |
| `POST` | `/api/v1/auth/login` | — | Authenticate (TOTP prompt if enabled) |
| `POST` | `/api/v1/auth/login/totp` | Temp | TOTP 2FA verification |
| `POST` | `/api/v1/auth/refresh` | JWT | Rotate refresh token |
| `GET` | `/api/v1/movies` | JWT | List movies (paginated) |
| `GET` | `/api/v1/movies/search` | JWT | Trie-based prefix search |
| `GET` | `/api/v1/movies/staff-picks` | JWT | Staff picks |
| `GET` | `/api/v1/movies/last-chance` | JWT | Last copies |
| `GET` | `/api/v1/movies/{id}` | JWT | Movie detail |
| `POST` | `/api/v1/movies` | Manager+ | Create movie |
| `PUT` | `/api/v1/movies/{id}` | Manager+ | Update movie |
| `DELETE` | `/api/v1/movies/{id}` | Manager+ | Delete movie |
| `POST` | `/api/v1/movies/{id}/staff-pick` | Manager+ | Toggle staff pick |
| `POST` | `/api/v1/rentals/rent` | JWT | Rent (ticket or 💵, with payment choice) |
| `POST` | `/api/v1/rentals/return` | JWT | Return (+🍿, deduct late fees from 💵) |
| `POST` | `/api/v1/rentals/extend` | JWT | Extend due (30🍿, +1 min) |
| `GET` | `/api/v1/rentals/history` | JWT | Rental history |
| `GET` | `/api/v1/tiers` | JWT | List subscription tiers |
| `POST` | `/api/v1/tiers/purchase` | JWT | Buy/upgrade tier |
| `GET`/`POST`/`DELETE` | `/api/v1/wishlist` | JWT | Wishlist CRUD |
| `GET` | `/api/v1/merch` | JWT | Rewards catalog |
| `POST` | `/api/v1/merch/redeem` | JWT | Redeem 🍿 |
| `GET` | `/api/v1/inventory` | JWT | Collectibles |
| `GET` | `/api/v1/users` | Supervisor+ | List users |
| `PUT` | `/api/v1/users/{id}` | Supervisor+ | Update role/ban |
| `GET` | `/api/v1/audit` | Supervisor+ | Audit entries |
| `GET` | `/api/v1/audit/verify` | Supervisor+ | Verify hash chain integrity |

---

## Recommended Upgrades (Beyond Current Scope)

### To Strengthen Cybersecurity Demonstration
- **🔑 TOTP enrollment flow** — Currently admins set it; add self-service setup/disable for all users with QR code display
- **📊 Security dashboard** — Admin screen showing: login attempts (success/fail), active lockouts, TOTP failure count, rate-limit hits
- **🔐 Session management** — Show active sessions in profile, allow users to revoke other sessions
- **📝 Audit filters** — Filter audit log by action type, date range, or user; export as JSON

### To Strengthen Data Structures Demonstration
- **📊 Co-rental visualization** — ASCII graph showing movie relationships: "Users who rented X also rented Y"
- **⏳ Waitlist demo** — Activate the min-heap waitlist so users queue for unavailable titles and get notified
- **🔄 LRU cache metrics** — Show cache hit/miss ratio on admin dashboard

### Quality of Life
- **📧 Notification system** — Toast-style messages for "Movie X is now available" (from waitlist)
- **🎨 Theme switcher** — Light/dark/VHS green phosphor themes
- **📱 Responsive resize** — Better card grid adaptation on window resize
- **📋 Bulk seed** — Seed multiple rental histories for demo (pre-populated rental data)

---

## Seed Data

- **~147 movies** across 8 genres, 3 formats, spanning 1937–2022 with varied rental prices
- **8 test users** with various subscriptions + RBAC roles
- **7 merchandise items** (collectibles, free rentals, tier upgrade)
- **All passwords**: `123`
- **Starting balances**: $5–$100 + 250🍿

---

## Technology Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.24+ |
| TUI | Bubble Tea + Lipgloss + Bubbles |
| Database | BoltDB (embedded) |
| HTTP Router | Chi v5 |
| Auth | bcrypt + JWT HS256 + TOTP RFC 6238 |
| Encryption | AES-256-GCM |
| Integrity | SHA-256 Hash Chain (Merkle-Damgård) |
| Data Structures | 9 custom implementations (bitmask, list, deque, heap, trie, LRU, bloom, hashchain, graph) |
