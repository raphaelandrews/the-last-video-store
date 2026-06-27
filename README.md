# 🎬 The Last Video Store

A retro-styled video rental system with a terminal UI and REST API.

## Features

- **Catalog** — movies, series, games, and snackbar; live search, genre filters, staff picks
- **Rentals** — format-aware durations, late & rewind fees, due-date countdowns, free rental coupons, ticket-or-money payment
- **SnackBar** — order, history, restock
- **Game Arcade** — hourly in-store play sessions
- **Subscriptions** — Wood → Diamond tiers with monthly free-rental allowance
- **Loyalty** — popcorn points on returns, redeemable rewards
- **Admin** — promote/demote, ban/unban, toggle TOTP, tabbed catalog management (movies/series/games), tamper-evident audit log with chain verify
- **Security** — bcrypt passwords, JWT access+refresh tokens, TOTP 2FA, brute-force lockout, RBAC bitmask, hash-chain audit log

## Tech stack

| Layer | Technology |
|---|---|
| Language | Go |
| TUI | Bubble Tea + Lipgloss + Bubbles |
| API | Chi |
| Storage | BoltDB |


## Run locally

```bash
git clone https://github.com/anomalyco/the-last-video-store && cd the-last-video-store
go run ./data/seed.go      # populate DB
go run ./cmd/server/       # API server on :8080
go run ./cmd/client/       # TUI client (in another terminal)
# Default users (all share password "123"): bronze, silver, gold, employee,
# supervisor, manager, owner, bar_attendant, bar_manager, game_attendant,
# game_manager, banned
```


## Deploy to Render

Set env:

| Var | Source |
|---|---|
| `TLVS_JWT_SECRET` | auto-generate |
| `TLVS_AES_KEY` | auto-generate (32 hex chars) |
| `TLVS_SERVER_PORT` | `8080` |
| `TLVS_DB_PATH` | `/app/data/thelastvideostore.db` |
| `TLVS_ENV` | `production` |

