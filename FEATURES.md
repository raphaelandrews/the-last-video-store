# The Last Video Store — Features & Workflow

## Overview

A retro VHS-themed video rental store TUI application. Browse, rent, return, wishlist movies, earn Popcorn Points, and redeem rewards — all through a terminal interface backed by a REST API and BoltDB. Supports 7 membership tiers, TOTP 2FA, admin management, staff picks, and an immutable audit trail.

---

## Screens (13 total)

| Screen | Access | Description |
|--------|--------|-------------|
| **Splash** | All | Animated figlet-style banner, `ENTER` to start |
| **Login** | All | Username + password authentication |
| **Register** | All | Create new account (3 fields) |
| **TOTP** | 2FA users | 6-digit authenticator code entry |
| **Browse** | Authenticated | Main catalog grid with pagination, search, and viewing modes |
| **Detail** | Authenticated | Synopsis, rating, cast, rent/waitlist actions |
| **Rentals** | Authenticated | Active rentals + history, due date countdown, extend, return |
| **Profile** | Authenticated | Membership card, tier badge, stats, rewards shop access |
| **Wishlist** | Authenticated | Personal wishlist viewer with remove |
| **Rewards** | Authenticated | Popcorn Points shop: redeem for free rentals, tier upgrades, collectibles |
| **Inventory** | Authenticated | View owned collectibles (popcorn bucket, poster, etc.) |
| **Admin Movies** | Manager+ | Movie CRUD, staff pick toggle |
| **Admin Users** | Supervisor+ | Promote/demote/ban users |
| **Audit Log** | Supervisor+ | Immutable hash chain log viewer |

---

## Navigation Flow

```
Splash ──ENTER──→ Login ──login──→ Browse ──┬── Detail ──── rent / add to waitlist
                     ↑                      ├── Rentals ─── return / extend
                     │                      ├── Profile ─── logout → Login
                Register                    │   ├── Rewards Shop (M)
                     │                      │   └── Inventory (I)
                     └──────────────────────├── Wishlist ─── view / remove
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
| `Q` | Back to Browse (from sub-screens) |

### Browse (Main Catalog)
| Key | Action |
|-----|--------|
| `↑↓` / `J` `K` | Navigate grid |
| `ENTER` / `D` | Open movie detail |
| `R` | My Rentals |
| `P` | Profile |
| `V` | View Wishlist |
| `/` | Search movies (live prefix search) |
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
| `ENTER` | Rent movie (if available) |
| `W` | Add to waitlist |
| `F5` | Refresh availability |
| `Q` | Back to browse |

### Rentals
| Key | Action |
|-----|--------|
| `↑↓` / `J` `K` | Navigate |
| `ENTER` | Return selected rental |
| `E` | Extend due date (+2 days, costs 30🍿) |
| `Q` | Back to browse |

### Profile
| Key | Action |
|-----|--------|
| `L` | Logout |
| `M` | Rewards Shop |
| `I` | Inventory |
| `Q` | Back to browse |

### Rewards Shop
| Key | Action |
|-----|--------|
| `↑↓` / `J` `K` | Navigate items |
| `ENTER` | Redeem selected reward |
| `Q` | Back to profile |

### Wishlist
| Key | Action |
|-----|--------|
| `↑↓` / `J` `K` | Navigate items |
| `ENTER` | Item info |
| `D` / `DELETE` | Remove item |
| `Q` | Back to browse |

### Admin Screens
| Screen | Key | Action |
|--------|-----|--------|
| Movies | `A` | Add movie |
| Movies | `ENTER` | Edit movie |
| Movies | `D` | Delete movie |
| Movies | `S` | Toggle staff pick |
| Users | `P` | Promote tier |
| Users | `D` | Demote tier |
| Users | `B` | Toggle ban |
| Audit | `V` | Verify chain integrity |

---

## Browse Modes

| Mode | Key | Source |
|------|-----|--------|
| **All** | `A` | Full catalog, paginated at 40 per page |
| **Staff Picks** | `S` | Manager-curated recommendations |
| **Last Chance** | `L` | Movies with 1 copy remaining |

---

## Membership Tiers

| Tier | Max Rentals | New Releases | Admin Access |
|------|:----------:|:---:|:---:|
| Bronze | 1 | No | — |
| Silver | 2 | No | — |
| Gold | 5 | Yes | — |
| Employee | 5 | Yes | Staff |
| Supervisor | 5 | Yes | Users + Audit |
| Manager | 10 | Yes | Movies + All |
| Owner | ∞ | Yes | Full access |

---

## Rental Rules

| Rule | Detail |
|------|--------|
| VHS due date | 3 days |
| DVD/Blu-ray due date | 5 days |
| VHS late fee | $2/day |
| DVD/Blu-ray late fee | $3/day |
| VHS rewind fee | $1.00 (30% random chance) |
| Due date display | Countdown: "in N days", "due soon" (≤2 days), "overdue by N days" |
| Extend rental | `E` key, costs 30🍿 for +2 days; overdue rentals become active again |
| Free rentals | Bypass tier limit, waive all late fees, tagged with 🎟️ in list |

---

## Popcorn Points System

| Action | Points |
|--------|:------:|
| On-time return (no fees) | +10 |
| Late return | -5 |
| Popcorn Bucket bonus (per return) | +5 |
| Private Screening bonus | +5 free rental tokens |

### Rewards Catalog

| Item | 🍿 Cost | Effect |
|------|:------:|--------|
| Popcorn Bucket | 50 | +5 bonus points on every future return |
| Blank VHS Tape | 75 | Collectible (stored in inventory) |
| Movie Poster | 100 | Collectible (stored in inventory) |
| Store T-Shirt | 150 | Collectible (stored in inventory) |
| Free Rental Coupon | 200 | +1 free rental (no late fees, bypasses limit) |
| Private Screening | 500 | +5 free rentals |
| Tier Upgrade | 1000 | Permanent tier promotion (up to Gold), increases max rentals |

---

## Wishlist Workflow

1. Browse → `V` → view your wishlist
2. Detail → `W` → add current movie
3. Wishlist shows availability (✓ available / ✗ out)
4. Wishlist → `D` → remove item via API

---

## Authentication & Security

- **bcrypt** password hashing (cost 12)
- **JWT** access tokens + refresh token rotation
- **TOTP 2FA** — HMAC-SHA1, 6-digit codes, AES-256-GCM encrypted secrets
- **RBAC** — 6-bit permission bitmask per tier
- **Brute-force lockout** — 5 failed login attempts = 30min lock
- **TOTP lockout** — 3 failed codes = 10min lock
- **Audit Trail** — SHA-256 hash chain, immutable, append-only

---

## Data Structures (Custom Implementations)

| Structure | Purpose |
|-----------|---------|
| Bitmask | RBAC permissions |
| Doubly Linked List | Rental history ordering |
| Deque (Ring Buffer) | Staff return priority queue |
| Min-Heap | Waitlist ordering |
| Trie | Movie title prefix search |
| LRU Cache | Session + movie caching |
| Bloom Filter | Fast banned-user lookup |
| Hash Chain | Immutable audit trail |
| Graph | Co-rental recommendations |

---

## API Endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| `POST` | `/api/v1/auth/login` | — | Authenticate (step 1) |
| `POST` | `/api/v1/auth/login/totp` | Temp | TOTP verification (step 2) |
| `POST` | `/api/v1/auth/register` | — | Create account |
| `POST` | `/api/v1/auth/refresh` | JWT | Refresh JWT |
| `GET` | `/api/v1/movies` | JWT | List movies (paginated) |
| `GET` | `/api/v1/movies/search` | JWT | Prefix search |
| `GET` | `/api/v1/movies/staff-picks` | JWT | Staff picks |
| `GET` | `/api/v1/movies/last-chance` | JWT | Last chance |
| `GET` | `/api/v1/movies/{id}` | JWT | Single movie |
| `POST` | `/api/v1/rentals/rent` | JWT | Rent a movie |
| `POST` | `/api/v1/rentals/return` | JWT | Return a rental |
| `POST` | `/api/v1/rentals/extend` | JWT | Extend due date (30🍿, +2 days) |
| `GET` | `/api/v1/rentals/history` | JWT | Rental history |
| `GET` | `/api/v1/wishlist` | JWT | View wishlist |
| `POST` | `/api/v1/wishlist` | JWT | Add to wishlist |
| `DELETE` | `/api/v1/wishlist/{movieID}` | JWT | Remove from wishlist |
| `GET` | `/api/v1/merch` | JWT | Rewards catalog |
| `POST` | `/api/v1/merch/redeem` | JWT | Redeem popcorn points |
| `GET` | `/api/v1/inventory` | JWT | Your collectibles |
| `GET` | `/api/v1/users` | Supervisor+ | List users |
| `PUT` | `/api/v1/users/{id}` | Supervisor+ | Update user |
| `GET` | `/api/v1/audit` | Supervisor+ | Audit log |

---

## Seed Data

- **~135 movies** across 8 genres, 3 formats, spanning 1937–2022
- **42 new releases**, **5 staff picks**, varied copy counts (1–5)
- **8 test users** (bronze, silver, gold, employee, supervisor, manager, owner, banned)
- **7 merchandise items** (collectibles, free rentals, tier upgrade)
- **All test users password**: `123`
- **All users start with**: 250🍿 Popcorn Points

---

## Technology Stack

| Layer | Technology |
|-------|-----------|
| Language | Go |
| TUI Framework | Bubble Tea + Lipgloss + Bubbles |
| Database | BoltDB (embedded key-value) |
| HTTP Router | Chi |
| Auth | bcrypt, JWT (HS256), TOTP (RFC 6238) |
| Crypto | AES-256-GCM, SHA-256 |
| Transport | REST JSON API |
