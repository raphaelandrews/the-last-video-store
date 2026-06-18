# The Last Video Store — Features & Workflow

## Overview

A retro VHS-themed video rental store TUI application — built for academic demonstration of **Cybersecurity & Data Structures**. Browse, rent, return, wishlist movies/series/games, order from the snackbar, play games in-store, earn Popcorn Points, purchase premium subscription tiers, and redeem rewards — all through a terminal interface backed by a REST API and BoltDB.

**Two-tier access control**: Premium subscription tiers (Wood→Diamond) govern rental perks and costs. RBAC roles (Bronze→Owner + SnackBar Attendant/Manager + Game Attendant/Manager) govern admin access — both are separate, demonstrating layered security.

---

## Scope Compliance — Requirements Mapping

### Original requirements
> *"Develop a system with an interface and access security. Implement a program using Cybersecurity and Data Structures. Create a system that uses Access Security and a data structure to control access and store data."*

### Security challenges and how they were met

| Challenge | Implementation |
|-----------|---------------|
| **a) Read the file/database line by line** | BoltDB with `First()`/`Next()` cursors — sequential record iteration with middleware-enforced access control |
| **b) Only allow authorized users to read** | JWT + 10-bit RBAC bitmask — every route checks `RequirePermission(flag)`; no token = HTTP 401, no permission = HTTP 403 |
| **c) Display the user and the data** | Header shows `🎫 username | 🏷️ Role | 🍿 pts | 💵 balance`; catalog and audit log display data based on permissions |
| **d) User registration via file/database** | Register screen (`POST /api/v1/auth/register`) + Admin Users panel (`PUT /api/v1/users/{id}`) for promote/demote/ban |
| **e) "Permissão Negada" (Access Denied)** | Dedicated `⛔ ACCESS DENIED` screen when server returns 403 via API middleware — all authorization is server-side, single source of truth |
| **f) Add or remove user access** | Admin Users → keys `P` (promote) / `D` (demote) / `B` (ban) — modify RBAC tier via API |

### Cybersecurity Technologies

| Layer | Technology | Detail |
|--------|-----------|---------|
| Password hashing | bcrypt | Cost 12, per-password salt |
| Authentication | JWT HS256 | Access token (15min) + Refresh token (7 days) with rotation |
| Two-factor auth | TOTP RFC 6238 | HMAC-SHA1, ±30s window, AES-256-GCM encrypted secret |
| Authorization | RBAC Bitmask | 10-bit permission flags, O(1) check |
| Attack protection | Lockout | 5 failed logins = 30min lock; 3 failed TOTP = 10min lock |
| Attack protection | Rate limiting | Per-IP rate limiter on all API routes |
| Integrity | SHA-256 Hash Chain | Immutable audit trail with chain verification |
| Integrity | Bloom Filter | Probabilistic fast banned-user rejection |
| Encryption | AES-256-GCM | TOTP secrets and sensitive data at rest |

### Data Structures (9 custom implementations, zero `container/*`)

| Structure | Application |
|-----------|------------|
| **Bitmask** (10-bit) | RBAC access control — stores permissions in 10 bits across 10 roles |
| **Doubly Linked List** | Rental history — O(1) insert/remove |
| **Deque** (Ring Buffer) | Staff priority return queue |
| **Min-Heap** | Waitlist — ordered by wait time |
| **Trie** (Prefix Tree) | Movie title autocomplete — O(k) |
| **LRU Cache** | Session + movie caching — O(1) |
| **Bloom Filter** | Fast banned-user check — O(k) |
| **Hash Chain** (SHA-256) | Immutable audit trail — O(1) append |
| **Undirected Weighted Graph** | Co-rental recommendations — O(V+E) |

---

## Screens (20 total)

| Screen | Access | Description |
|--------|--------|-------------|
| **Splash** | All | Animated figlet-style banner, `ENTER` to start |
| **Login** | All | Username + password, TOTP challenge if enabled |
| **Register** | All | Create new account (3 fields, validations) |
| **TOTP** | 2FA users | 🔒 6-digit authenticator code entry |
| **Browse** | Authenticated | Tabbed catalog: Movies · Series · Games · SnackBar, pagination, search, staff picks, last chance |
| **Detail** | Authenticated | Synopsis, rating, cast, choose payment method, rent/waitlist |
| **Game Detail** | Authenticated | Game info, platform badge, rent (`R`) or play in-store (`P`), play session timer |
| **Rentals** | Authenticated | Active + history, countdown timer, extend, return with fees |
| **Profile** | Authenticated | Membership card, role badge, subscription tier, money + 🍿 |
| **Tier Shop** | Authenticated | Purchase/upgrade premium subscription (Wood→Diamond) |
| **Wishlist** | Authenticated | Personal waitlist with availability indicators |
| **SnackBar Menu** | Authenticated | 🍿 Concession stand — order drinks, snacks, candy |
| **SnackBar Orders** | Authenticated | View past snackbar orders |
| **SnackBar Manage** | SnackBarManager+ | Restock snackbar items, manage inventory |
| **Rewards** | Authenticated | 🍿 Popcorn Points shop: 26 film-themed collectibles |
| **Inventory** | Authenticated | View owned collectibles |
| **Game Sessions** | GameManager+ | Active in-store play sessions monitor |
| **Access Denied** | Authenticated | ⛔ Full-screen denial with role requirement from server |
| **Admin Movies** | Manager+ | Movie CRUD, staff pick toggle, paginated list |
| **Admin Users** | Supervisor+ | Promote/demote RBAC role, toggle ban |
| **Audit Log** | Supervisor+ | SHA-256 hash chain viewer, chain integrity verification |

---

## Navigation Flow

```
Splash ──ENTER──→ Login ──login──→ Browse ──┬── Movie Detail ──── rent / waitlist
                     ↑                      ├── Game Detail ───── rent / play / end
                     │                      ├── Rentals ────────── return / extend
                 Register                   ├── Profile ────────── logout → Login
                     │                      │   ├── Tier Shop (T)
                     └──────────────────────│   ├── SnackBar (B)
                                            │   ├── Rewards Shop (M)
                                            │   └── Inventory (I)
                                            ├── SnackBar Menu ─── order / orders / manage
                                            │   ├── Orders (O)
                                            │   └── Manage (M, SnackBarManager+)
                                            ├── Wishlist ───────── view / remove
                                            ├── Game Sessions ──── active sessions (GameManager+)
                                            ├── ⛔ Access Denied (server 403)
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
| `Q` | Back to parent screen |

**Q-back hierarchy:**
| From | Returns to |
|------|-----------|
| Detail, Rentals, Wishlist, Game Detail, Game Sessions | Browse |
| Rewards Shop, Tier Shop, Inventory, SnackBar Menu | Profile |
| Admin Movies, Users, Audit Log | Browse |
| Access Denied | Browse |
| Movie Form (`ESC`) | Admin Movies |
| Profile | Browse |
| SnackBar Orders, Manage | SnackBar Menu |

### Browse (Main Catalog)
| Key | Action |
|-----|--------|
| `↑↓` / `J` `K` | Navigate grid |
| `ENTER` / `D` | Open detail (movies/series → Detail, games → Game Detail) |
| `[` / `]` | Media type tabs (Movies / Series / Games / SnackBar) |
| `,` / `.` | Genre subtabs (dynamic per media type) |
| `R` | My Rentals |
| `P` | Profile |
| `V` | View Wishlist |
| `C` | SnackBar |
| `/` | Search (live prefix with Trie backend) |
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
| `ENTER` | Open selected item |
| `ESC` | Cancel search |

### Detail (Movies/Series)
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

### Game Detail
| Key | Action |
|-----|--------|
| `R` | Rent game (takes home, like movie) |
| `P` | Start in-store play session (hourly rate) |
| `E` | End play session |
| `↑↓` / `J` `K` | Navigate related games |
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
| `B` | SnackBar |
| `M` | Rewards Shop |
| `I` | Inventory |
| `Q` | Back to browse |

### SnackBar Menu
| Key | Action |
|-----|--------|
| `↑↓` / `J` `K` | Navigate items |
| `ENTER` | Place order (deducts 💵 from balance) |
| `O` | Order History |
| `M` | Manage (SnackBarManager+) |
| `Q` | Back to Profile |

### SnackBar Manage
| Key | Action |
|-----|--------|
| `↑↓` / `J` `K` | Navigate items |
| `R` | Restock +5 units |
| `Q` | Back to SnackBar Menu |

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
| **💵 Money** | Seed starting balance ($5–$100) | Rentals beyond tier allowance, snackbar orders, game play sessions, buying premium tiers |
| **🍿 Popcorn Points** | Returns (+10 on-time, -5 late), Popcorn Bucket bonus (+5) | Rewards shop, extend rentals (30🍿) |

### Movie-Specific Rental Pricing
| Category | VHS/DVD | Blu-ray |
|----------|---------|---------|
| New releases & 2020+ | $5.99 | $6.99 |
| 2000–2019 | $3.99 | $4.99 |
| Pre-2000 | $2.99 | $3.99 |

---

## RBAC Roles (Admin Access)

10-role 10-bit bitmask — O(1) permission checks:

| Role | Access |
|------|--------|
| **Bronze** | Browse catalog, access snackbar + games |
| **Silver** | + Rent movies/games |
| **Gold** | + New Releases, reservations |
| **Employee** | + Staff tools, browse all sections |
| **Supervisor** | + User management + Audit log |
| **Manager** | + Movie/game/snackbar management |
| **Owner** | Full access |
| **SnackBar Attendant** | Browse + SnackBar access (no movie/game management) |
| **SnackBar Manager** | + SnackBar inventory management |
| **Game Attendant** | Browse + Game access (no movie/snackbar management) |
| **Game Manager** | + Game inventory + play session management |

Permission bits: Browse (1), Rent (2), Reserve (4), ManageUsers (8), Staff (16), Admin (32), SnackBarAccess (64), SnackBarManage (128), GameAccess (256), GameManage (512).

---

## SnackBar Menu

| Category | Items |
|----------|-------|
| **Snacks** | 🍿 Popcorn ($3.99), 🧀 Nachos ($5.99), 🌭 Hot Dog ($4.99), 🍕 Pizza Slice ($4.49), 🥨 Soft Pretzel ($3.99), 🍟 Fries ($3.49), 🍔 Cheeseburger ($6.99) |
| **Candy** | 🍬 Candy Assortment ($2.99), 🍫 Chocolate Bar ($2.49), 🍦 Ice Cream ($3.49) |
| **Drinks** | 🥤 Fountain Soda ($2.99), 💧 Water ($1.49), 🧊 Slushie ($3.99), ☕ Coffee ($2.49), 🥛 Milkshake ($4.99) |

---

## Game Catalog (47 titles)

| Platform | Titles |
|----------|--------|
| **NES** | Super Mario Bros 3, The Legend of Zelda, Metroid, Mega Man 2, Castlevania, Contra, Duck Hunt, Tetris, Punch-Out!! |
| **SNES** | Super Mario World, Zelda: ALTTP, Super Metroid, Chrono Trigger, FF6, Donkey Kong Country, Street Fighter II Turbo, Super Mario Kart, EarthBound |
| **Genesis** | Sonic 2, Streets of Rage 2, Gunstar Heroes, Mortal Kombat II |
| **PS1** | FF VII, Metal Gear Solid, Crash Bandicoot 2, Resident Evil 2, Tony Hawk 2, Castlevania SOTN, Tekken 3 |
| **N64** | Super Mario 64, Zelda OoT, GoldenEye 007, Mario Kart 64, Super Smash Bros, Banjo-Kazooie |
| **PC** | Doom, Half-Life, StarCraft, Age of Empires II, Diablo II, The Sims, RollerCoaster Tycoon 2 |
| **Arcade** | Pac-Man, Space Invaders, Galaga, Street Fighter II, Donkey Kong |

---

## Due Dates (Demo-Friendly Minutes)

| Format | Due In | Late Fee rate |
|--------|--------|--------------|
| VHS | 2 min | ~$0.20/min |
| DVD | 3 min | ~$0.30/min |
| Blu-ray | 4 min | ~$0.40/min |
| Games (cartridge) | 3 min | ~$0.30/min |
| Games (CD) | 4 min | ~$0.40/min |
| **Extend** | +1 min | costs 30🍿 |
| **Play in-store** | Per hour | Deducted from balance on start |

Full lifecycle (rent → overdue → late fee → extend → return) demonstrable in ~5 minutes.

---

## Security Features (Cybersecurity Scope)

| Layer | Feature | Detail |
|-------|---------|--------|
| **Authentication** | bcrypt | Cost factor 12, salt per password |
| **Authentication** | JWT | HS256, 15min access + 7-day refresh rotation |
| **Authentication** | TOTP 2FA | RFC 6238, HMAC-SHA1, 30s window ±1, AES-256-GCM encrypted secrets |
| **Authorization** | RBAC Bitmask | 10-bit permission flags, O(1) check via `p & flag != 0` |
| **Authorization** | Middleware chain | JWT validation → ban check → permission check → handler |
| **Authorization** | Server-side enforcement | API is single source of truth — all denials route through 403 → access denied screen |
| **Attack Protection** | Brute-force lockout | 5 failed logins = 30min lock, 3 failed TOTP = 10min lock |
| **Attack Protection** | Rate limiting | Per-IP rate limiter on all API routes |
| **Integrity** | Hash Chain Audit | SHA-256 Merkle-Damgård chaining, chain verification endpoint |
| **Integrity** | Bloom Filter | O(k) banned-user fast rejection pre-DB lookup |
| **Encryption** | AES-256-GCM | TOTP secrets, audit-sensitive data at rest |

---

## Data Structures (All Custom, No `container/*`)

| Structure | Application | Complexity |
|-----------|------------|:---:|
| **Bitmask** (10-bit) | RBAC permissions | O(1) |
| **Doubly Linked List** | Rental history ordering | O(1) insert/remove |
| **Deque** (Ring Buffer) | Staff return priority queue | O(1) push/pop |
| **Min-Heap** | New release waitlist ordering | O(log n) |
| **Trie** (Prefix Tree) | Movie title autocomplete search | O(k) |
| **LRU Cache** | Session + movie caching | O(1) |
| **Bloom Filter** | Banned user fast check | O(k) |
| **Hash Chain** (SHA-256) | Immutable audit trail | O(1) append |
| **Undirected Weighted Graph** | Co-rental recommendations | O(V+E) |

---

## Popcorn Points Rewards (26 items)

| Item | 🍿 Cost | Effect |
|------|:------:|--------|
| Popcorn Bucket | 50 | +5 bonus points on every future return |
| Blank VHS Tape | 75 | Collectible (inventory) |
| Movie Poster | 100 | Collectible (inventory) |
| Store T-Shirt | 150 | Collectible (inventory) |
| Free Rental Coupon | 200 | +1 free rental (bypasses limit, no fees) |
| Private Screening | 500 | +5 free rentals |
| Tier Upgrade | 1000 | Promote RBAC role one level (up to Gold) |
| Pokemon TCG Booster | 120 | Vintage Jungle expansion — chance of holographic Pikachu |
| Red Pill / Blue Pill Set | 180 | Matrix resin-cast pill keychain pair in velvet pouch |
| Neo's Trench Coat | 800 | Full-length leather-look, Matrix digital rain lining |
| Origami Unicorn | 90 | Hand-folded metallic paper — Blade Runner Gaff tribute |
| Jurassic Park Amber Cane | 350 | Polished resin cane top with faux mosquito inclusion |
| Marlon Brando Cat Plush | 130 | Plush ginger tabby — Godfather opening scene |
| Overlook Carpet Coasters | 60 | Shining hexagonal carpet pattern, set of 4 |
| Big Kahuna Burger Box | 140 | Pulp Fiction burger-shaped tin lunchbox |
| One Ring Replica | 250 | Tungsten LOTR band with elvish inscription |
| Soot Sprite Plushies | 110 | Spirited Away susuwatari, set of 3 |
| Mr. Fusion Prop Replica | 300 | BTTF desktop model — 1.21 gigawatts of style |
| Chestburster Plush | 160 | Surprisingly cute Alien xenomorph hatchling |
| Totem Spinning Top | 200 | Inception brass spinning top in felt-lined case |
| Paper Street Soap Co. Bar | 45 | Fight Club handmade pink soap |
| Indiana Jones Fedora | 280 | Brown felt, adventure-ready |
| Crow Plush (Hitchcock Ed.) | 100 | Surprisingly heavy — The Birds tribute |
| Akira Capsule Patch | 70 | Embroidered pill capsule jacket patch |
| Tarantino Socks | 55 | Feet of various Tarantino characters |
| Mini Monolith | 220 | Solid obsidian-black 2001: A Space Odyssey paperweight |

---

## API Endpoints

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| `POST` | `/api/v1/auth/register` | — | Create account |
| `POST` | `/api/v1/auth/login` | — | Authenticate (TOTP prompt if enabled) |
| `POST` | `/api/v1/auth/login/totp` | Temp | TOTP 2FA verification |
| `POST` | `/api/v1/auth/refresh` | JWT | Rotate refresh token |
| `POST` | `/api/v1/auth/logout` | JWT | Revoke tokens |
| `GET` | `/api/v1/auth/me` | JWT | Get current user state |
| `GET` | `/api/v1/movies` | JWT | List items (paginated, genre + media_type filters) |
| `GET` | `/api/v1/movies/search` | JWT | Trie-based prefix search |
| `GET` | `/api/v1/movies/staff-picks` | JWT | Staff picks |
| `GET` | `/api/v1/movies/last-chance` | JWT | Last copies |
| `GET` | `/api/v1/movies/{id}` | JWT | Detail |
| `POST` | `/api/v1/movies` | Manager+ | Create |
| `PUT` | `/api/v1/movies/{id}` | Manager+ | Update |
| `DELETE` | `/api/v1/movies/{id}` | Manager+ | Delete |
| `POST` | `/api/v1/movies/{id}/staff-pick` | Manager+ | Add staff pick |
| `DELETE` | `/api/v1/movies/{id}/staff-pick` | Manager+ | Remove staff pick |
| `POST` | `/api/v1/rentals/rent` | JWT | Rent (ticket or 💵) |
| `POST` | `/api/v1/rentals/return` | JWT | Return (+🍿, deduct late fees from 💵) |
| `POST` | `/api/v1/rentals/extend` | JWT | Extend due (30🍿) |
| `GET` | `/api/v1/rentals/history` | JWT | Rental history |
| `POST` | `/api/v1/games/play/start` | JWT | Start in-store play session |
| `POST` | `/api/v1/games/play/end` | JWT | End play session |
| `GET` | `/api/v1/games/play/active` | GameManager+ | Active play sessions |
| `GET` | `/api/v1/snackbar` | JWT | SnackBar menu |
| `POST` | `/api/v1/snackbar/order` | JWT | Place order |
| `GET` | `/api/v1/snackbar/orders` | JWT | Order history |
| `POST` | `/api/v1/snackbar/items` | SnackBarManager+ | Add item |
| `PUT` | `/api/v1/snackbar/items/{id}` | SnackBarManager+ | Update item |
| `DELETE` | `/api/v1/snackbar/items/{id}` | SnackBarManager+ | Delete item |
| `POST` | `/api/v1/snackbar/restock` | SnackBarManager+ | Restock |
| `GET` | `/api/v1/tiers` | JWT | List subscription tiers |
| `POST` | `/api/v1/tiers/purchase` | JWT | Buy/upgrade tier |
| `GET`/`POST`/`DELETE` | `/api/v1/wishlist` | JWT | Wishlist CRUD |
| `GET` | `/api/v1/merch` | JWT | Rewards catalog |
| `POST` | `/api/v1/merch/redeem` | JWT | Redeem 🍿 |
| `GET` | `/api/v1/inventory` | JWT | Collectibles |
| `GET` | `/api/v1/users` | Supervisor+ | List users |
| `POST` | `/api/v1/users` | Supervisor+ | Create user |
| `PUT` | `/api/v1/users/{id}` | Supervisor+ | Update role/ban |
| `DELETE` | `/api/v1/users/{id}` | Manager+ | Delete user |
| `POST` | `/api/v1/users/{id}/totp` | Self/Manager+ | TOTP setup |
| `GET` | `/api/v1/audit` | Supervisor+ | Audit entries |
| `GET` | `/api/v1/audit/verify` | Supervisor+ | Verify hash chain |

---

## Seed Data

- **~296 titles** across movies, series, and games — 8 movie genres, 9 series genres, 9 game genres, multiple platforms
- **12 test users** with various subscriptions + RBAC roles (including SnackBar and Game attendants/managers)
- **26 merchandise items** (collectibles, film-themed merch, free rentals, tier upgrade)
- **15 snackbar items** across drinks, snacks, candy
- **47 classic games** across NES, SNES, Genesis, PS1, N64, PC, Arcade
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
