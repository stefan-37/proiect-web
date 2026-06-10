> This README was written by [Claude](https://claude.ai).

# Gym Management API

A RESTful backend for managing a gym: users, trainers, admins, subscription plans, and classes. Written in Go with the Gin web framework, GORM over PostgreSQL, JWT-based auth, and shipped behind an Nginx reverse proxy via Docker Compose.

## Stack

- **Language:** Go 1.26
- **Web framework:** [Gin](https://github.com/gin-gonic/gin)
- **ORM:** [GORM](https://gorm.io/) with the PostgreSQL driver
- **Database:** PostgreSQL 17 (alpine)
- **Auth:** JWT (`github.com/golang-jwt/jwt/v5`) stored in an HTTP cookie named `key`
- **Email:** SMTP via `net/smtp` (Gmail) — sends a notification email on user signup and the new password on password reset (requires `EMAIL` / `EMAIL_PASSWORD` env vars; no-ops if unset)
- **Reverse proxy:** Nginx (with rate-limited `/login` endpoints)
- **Containerization:** Docker + Docker Compose, multi-stage build on `distroless` base

## Project Layout

```
backend/
├── main.go                # Entry point — auto-migrates models, seeds plans, starts router
├── config/                # JWT secret / app config
├── db/                    # GORM singleton + connection pool setup
├── models/                # User, Admin, Trainer, Subscription, UserSubscription, Class, Person
├── repository/            # Data-access layer (one file per model)
├── service/               # Business logic (auth, user, admin, trainer, subscription, class, mail)
├── handler/               # Gin HTTP handlers, thin wrappers over services
├── middleware/            # JWT auth middleware (role-aware)
├── router/                # Route definitions
├── seed/                  # plans.json + seeder for default subscription tiers
├── Docker-Compose.yaml    # backend + postgres + nginx
├── dockerfile             # Multi-stage build
└── nginx.conf             # Reverse proxy + rate limiting on /login
```

## Domain Model

- **User** — gym member; can subscribe to a plan. Signup requires `name`, `email`, `password`, and `phone` (all non-empty).
- **Admin** — manages trainers and subscription plans.
- **Trainer** — created under an admin; can create classes.
- **Subscription** — a plan tier (`Basic`, `Premium`) with a price, owned by an admin.
- **UserSubscription** — links a user to a subscription with `StartedAt` / `ExpiresAt`.
- **Class** — scheduled session created by a trainer. Tracks `Capacity` and a live `Users` count; bookings are rejected once `Users` reaches `Capacity`.
- **BookingSituation** — links a user to a class they booked. Created via `/user/book`, which checks the class's remaining capacity and increments the class's `Users` count atomically (the class row is locked for the duration so concurrent bookings can't overbook). Listed via `/user/reservations` and cancelled via `/user/reservation`, which deletes the booking and decrements the class's `Users` count in the same locked transaction so the freed slot becomes bookable again.
- **Person** — a gym-attendance record; one row per visit, holding the member's id (`PersonID`), a `CheckedIn` timestamp, and a nullable `CheckedOut` timestamp. A row with `CheckedOut = NULL` means the member is currently in the gym.

```mermaid
classDiagram
    direction LR

    class Admin {
        int id
        string email
        string password
    }
    class User {
        int id
        string name
        string email
        string password
        string phone
    }
    class Subscription {
        int id
        string type
        float price
        int adminID
    }
    class UserSubscription {
        int id
        int UserID
        int SubscriptionID
        time start
        time end
        int adminID
    }
    class Trainer {
        int id
        string name
        string email
        string password
        string description
        int adminID
    }
    class Classes {
        int id
        date date
        time time
        string name
        string description
        int trainerID
        int capacity
        int adminID
    }
    class BookingSituation {
        int id
        int UserID
        int ClassID
        string status
        int adminID
    }
    class Person {
        int id
        int PersonID
        time checkedIn
        time checkedOut
    }

    User "1" --> "*" Person
    User "0..1" --> "*" UserSubscription
    Subscription "1" --> "*" UserSubscription
    Trainer "1" --> "*" Classes
    User "1" --> "*" BookingSituation
    Classes "1" --> "*" BookingSituation
```

Two background goroutines run in [main.go](backend/main.go): one calls `service.CheckSubscriptions` every 24 hours to keep subscription state up to date, and one re-seeds admins and trainers from JSON files every 5 minutes.

## API

All authenticated routes require a JWT cookie (`key`) with a `role` claim matching the route group.

### Public

| Method | Path             | Description                            |
| ------ | ---------------- | -------------------------------------- |
| GET    | `/subscriptions` | List available subscription plans      |
| GET    | `/classes`       | List scheduled classes                 |
| GET    | `/personsCount`  | Count members currently checked in     |
| POST   | `/forgot-password` | Reset a user's password by email — generates a new random password, emails it, and stores its hash (no auth required) |

### `/user`

| Method | Path                  | Auth      | Description                       |
| ------ | --------------------- | --------- | --------------------------------- |
| POST   | `/user/signup`        | public    | Register a new user (sends a notification email) |
| POST   | `/user/login`         | public    | Login, sets JWT cookie            |
| POST   | `/user/logout`        | user      | Logout, clears the JWT cookie     |
| GET    | `/user/get`           | user      | Get current user                  |
| PUT    | `/user/update`        | user      | Update current user               |
| POST   | `/user/reset-password`| user      | Reset current user's password — generates a new random password, emails it, and stores its hash |
| DELETE | `/user/delete`        | user      | Delete current user               |
| POST   | `/user/subscribe`     | user      | Subscribe to a plan               |
| GET    | `/user/subscription`  | user      | List the user's subscriptions     |
| POST   | `/user/checkin`       | user      | Check in to the gym               |
| POST   | `/user/checkout`      | user      | Check out of the gym              |
| POST   | `/user/book`          | user      | Book a class (capacity-enforced)  |
| GET    | `/user/reservations`  | user      | List the user's class reservations |
| DELETE | `/user/reservation`   | user      | Cancel a class reservation (frees the slot) |
| GET    | `/user/payments`      | user      | Payment history: returns `{ payments, subscriptions, classes, reservations }` — the user's purchases and reservations plus the full plan and class lists so the client can resolve `subscription_id`/`class_id` to a type/price/name |

### `/admin`

| Method | Path             | Auth  | Description           |
| ------ | ---------------- | ----- | --------------------- |
| POST   | `/admin/login`   | public| Login                 |
| GET    | `/admin/get`     | admin | Get current admin     |
| PUT    | `/admin/update`  | admin | Update current admin  |
| DELETE | `/admin/delete`  | admin | Delete current admin  |
| POST   | `/admin/class`   | admin | Create a class (no trainer; `TrainerID` 0) |

Admin accounts are seeded from `seed/admins.json` — there is no public signup endpoint.

### `/trainer`

| Method | Path                | Auth    | Description                  |
| ------ | ------------------- | ------- | ---------------------------- |
| POST   | `/trainer/login`    | public  | Login                        |
| GET    | `/trainer/get`      | trainer | Get current trainer          |
| PUT    | `/trainer/update`   | trainer | Update current trainer       |
| DELETE | `/trainer/delete`   | trainer | Delete current trainer       |
| POST   | `/trainer/class`    | trainer | Create a class               |

Trainer accounts are seeded from `seed/trainers.json` — there is no public signup endpoint.

## Running

### With Docker Compose (recommended)

From [backend/](backend/):

```bash
docker compose -f Docker-Compose.yaml up --build
```

This brings up:
- `postgres` on `:5432`
- `backend` (Go API) on internal `:8080`
- `nginx` on `:80`, proxying to the backend; all routes are burst-limited (burst 30 nodelay), with a tighter burst (10 nodelay) on `*/login` paths

The API is reachable at `http://localhost/`.

### Locally (without Docker)

You'll need a PostgreSQL instance reachable at the DSN hard-coded in [backend/db/db.go](backend/db/db.go):

```
postgresql://postgres:password@postgres:5432/mydatabase
```

Then:

```bash
cd backend
go mod download
go run .
```

The server listens on `:8080`.

## Seeding

On startup, [seed/plans.json](backend/seed/plans.json) is loaded into the `subscriptions` table — two default tiers (Basic / Premium).

A second background goroutine re-seeds [seed/admins.json](backend/seed/admins.json) and [seed/trainers.json](backend/seed/trainers.json) every 5 minutes. This is how admin and trainer accounts are provisioned — both files are bind-mounted into the container via Docker Compose. New entries are inserted idempotently (existing rows by email are never overwritten).

## Performance

A self-contained smoke / load test lives at [test_api.js](test_api.js) (Node 18+, no deps). It runs the full CRUD lifecycle for `/user`, `/admin`, `/trainer` plus auth-negative, cross-role, subscription, and class checks, and reports per-endpoint latency stats and end-to-end throughput.

```bash
node test_api.js [iterations] [concurrency]   # default 20 1
```

Reference numbers from `node test_api.js 200 50` against the full Docker Compose stack on this branch:

| Metric | Value |
|---|---|
| Hot-loop throughput | ~2900 req/s |
| Hot requests | 1600 |
| `user_get` p50 / p99 | ~15 ms / ~50 ms |
| `user_update` p50 / p99 | ~20 ms / ~80 ms |
| `user_login` p50 / p99 | ~180 ms / ~400 ms |
| `user_signup` p50 / p99 | ~210 ms / ~520 ms |
| Auth / cross-role / subscription / class checks | all 100% pass |

Things to know when reading those numbers:

- **Login and signup are bcrypt-bound.** [`bcrypt.DefaultCost`](https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-constants) = 10, deliberately ~60–100 ms per call. That's protection against offline cracking, not a perf bug. See [repository/user.go](backend/repository/user.go), [repository/admin.go](backend/repository/admin.go), [repository/trainer.go](backend/repository/trainer.go).
- **At concurrency > cores, bcrypt queues** — mean signup/login latency scales roughly with `(concurrency / CPU_cores) × single_call_cost`.
- **DB pool** is 20 max-open / 10 idle in [db/db.go](backend/db/db.go). Concurrency above ~20 will queue on non-bcrypt paths.
- **Nginx burst-limits all routes** (burst 30 nodelay) and applies a tighter burst (10 nodelay) on `*/login` paths — see [nginx.conf](backend/nginx.conf). Both zones currently run at a very high rate ceiling, so the burst caps are what enforce the limit in practice.

