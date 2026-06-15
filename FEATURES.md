# The Last Video Store — Features & Workflow

## Overview

A retro VHS-themed video rental store TUI application. Browse, rent, return, wishlist movies, earn Popcorn Points, purchase premium subscription tiers, and redeem rewards — all through a terminal interface backed by a REST API and BoltDB.

**Two-tier access control**: Premium subscription tiers (Wood→Diamond) govern rental perks and costs. RBAC roles (Employee→Owner) govern admin access. Both are separate — you can be a Wood subscriber with Owner admin rights, or Diamond subscriber with no admin access.

---

## Screens (14 total)

| Screen | Access | Description |
|--------|--------|-------------|
| **Splash** | All | Animated figlet-style banner, `ENTER` to start |
| **Login** | All | Username + password authentication |
| **Register** | All | Create new account (3 fields) |
| **TOTP** | 2FA users | 6-digit authenticator code entry |
| **Browse** | Authenticated | Main catalog grid with pagination, search, and viewing modes |
| **Detail** | Authenticated | Synopsis, rating, cast, rent/waitlist actions, rental cost shown |
| **Rentals** | Authenticated | Active rentals + history, due date countdown, extend, return |
| **Profile** | Authenticated | Membership card, role badge, stats, tier, balance |
| **Tier Shop** | Authenticated | Purchase/upgrade premium subscription tiers |
| **Wishlist** | Authenticated | Personal wishlist viewer with remove |
| **Rewards** | Authenticated | Popcorn Points shop: redeem for free rentals, collectibles |
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
                Register                    │   ├── Tier Shop (T)
                     │                      │   ├── Rewards Shop (M)
                     └──────────────────────│   └── Inventory (I)
                                            ├── Wishlist ─── view / remove
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
| `ENTER` | Rent movie (tier allowance first, then 💵) |
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
| `T` | Tier Shop (subscribe/upgrade) |
| `M` | Rewards Shop |
| `I` | Inventory |
| `Q` | Back to browse |

### Tier Shop
| Key | Action |
|-----|--------|
| `↑↓` / `J` `K` | Navigate tiers |
| `ENTER` | Purchase selected tier |
| `Q` | Back to profile |

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
| Users | `P` | Promote RBAC role |
| Users | `D` | Demote RBAC role |
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

## Premium Subscription Tiers

Users start with **Wood** (free). Purchase higher tiers from Profile → `T`.

| Tier | Price | Free Rentals/mo | Max Concurrent | New Releases | Late Fees |
|------|------:|:--------------:|:------------:|:---:|:---:|
| Wood | Free | 0 | 2 | No | Yes |
| Bronze | $9.99 | 1 | 3 | No | Yes |
| Silver | $19.99 | 3 | 5 | Yes | Yes |
| Gold | $29.99 | 5 | 10 | Yes | **Waived** |
| Diamond | $49.99 | Unlimited | Unlimited | Yes | **Waived** |

- **Free rental allocation** is per billing cycle (renews on purchase/renewal)
- **Paid rentals** cost money from your balance once the free allocation is exhausted
- **Upgrading** costs the full price of the new tier (not the difference)
- Tiers are separate from RBAC admin roles

---

## Money & Popcorn Points (Dual Currency)

| Currency | Earned By | Used For |
|----------|-----------|----------|
| **💵 Money** | Returns (refund of rental cost if on-time), seed starting balance | Renting beyond tier allowance, buying premium tiers |
| **🍿 Popcorn Points** | Returns (+10 on-time, -5 late), Popcorn Bucket bonus (+5) | Rewards shop: free rentals, tier upgrades, collectibles |

### Rental Costs (charged from 💵 balance)
| Format | Cost |
|--------|-----:|
| VHS | $2.99 |
| DVD | $3.99 |
| Blu-ray | $4.99 |

Cost is charged only if the tier's free rental allocation is exhausted. Free rentals from reward coupons also bypass the charge.

---

## RBAC Roles (Admin Access)

| Role | Access |
|------|--------|
| Bronze | Basic member |
| Silver | Basic member |
| Gold | Basic member |
| Employee | Staff tools |
| Supervisor | User management + audit |
| Manager | Movie management + full admin |
| Owner | Unlimited everything |

---

## Rental Rules

| Rule | Detail |
|------|--------|
| VHS due date | 3 days |
| DVD/Blu-ray due date | 5 days |
| VHS late fee | $2/day (deducted from 💵 balance) |
| DVD/Blu-ray late fee | $3/day (deducted from 💵 balance) |
| VHS rewind fee | $1.00 (30% random chance, deducted from 💵) |
| Late fees waived | Gold and Diamond tiers, or free rental coupons |
| Due date display | Countdown: "in N days", "due soon" (≤2 days), "overdue by N days" |
| Extend rental | `E` key, costs 30🍿 for +2 days; overdue rentals become active again |
| Free rentals | Tagged 🎟️ in list, bypass tier limit, waive all rental + late fees |

---

## Popcorn Points Rewards

| Item | 🍿 Cost | Effect |
|------|:------:|--------|
| Popcorn Bucket | 50 | +5 bonus points on every future return |
| Blank VHS Tape | 75 | Collectible (stored in inventory) |
| Movie Poster | 100 | Collectible (stored in inventory) |
| Store T-Shirt | 150 | Collectible (stored in inventory) |
| Free Rental Coupon | 200 | +1 free rental (no fees, bypasses limit) |
| Private Screening | 500 | +5 free rentals |
| Tier Upgrade | 1000 | Promote RBAC role one level (up to Gold) |

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
- **RBAC** — 6-bit permission bitmask per role
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
| `POST` | `/api/v1/rentals/rent` | JWT | Rent (tier allowance or 💵) |
| `POST` | `/api/v1/rentals/return` | JWT | Return (+🍿, ±💵) |
| `POST` | `/api/v1/rentals/extend` | JWT | Extend due date (30🍿, +2d) |
| `GET` | `/api/v1/rentals/history` | JWT | Rental history |
| `GET` | `/api/v1/tiers` | JWT | List subscription tiers |
| `POST` | `/api/v1/tiers/purchase` | JWT | Buy/upgrade tier |
| `GET` | `/api/v1/wishlist` | JWT | View wishlist |
| `POST` | `/api/v1/wishlist` | JWT | Add to wishlist |
| `DELETE` | `/api/v1/wishlist/{movieID}` | JWT | Remove from wishlist |
| `GET` | `/api/v1/merch` | JWT | Rewards catalog |
| `POST` | `/api/v1/merch/redeem` | JWT | Redeem 🍿 for reward |
| `GET` | `/api/v1/inventory` | JWT | Your collectibles |
| `GET` | `/api/v1/users` | Supervisor+ | List users |
| `PUT` | `/api/v1/users/{id}` | Supervisor+ | Update role/ban |
| `GET` | `/api/v1/audit` | Supervisor+ | Audit log |

---

## Seed Data

- **~135 movies** across 8 genres, 3 formats, spanning 1937–2022
- **42 new releases**, **5 staff picks**, varied copy counts (1–5)
- **8 test users** with various subscriptions + roles
- **7 merchandise items** (collectibles, free rentals, tier upgrades)
- **All test users password**: `123`
- **All users start with**: $50–$100 balance + 250🍿 Popcorn Points

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
