# Gym Management API

A RESTful backend for managing a gym: users, trainers, admins, subscription plans, and classes. Written in Go with the Gin web framework, GORM over PostgreSQL, JWT-based auth, and shipped behind an Nginx reverse proxy via Docker Compose.

## Stack

- **Language:** Go 1.26
- **Web framework:** [Gin](https://github.com/gin-gonic/gin)
- **ORM:** [GORM](https://gorm.io/) with the PostgreSQL driver
- **Database:** PostgreSQL 17 (alpine)
- **Auth:** JWT (`github.com/golang-jwt/jwt/v5`) stored in an HTTP cookie named `key`
- **Reverse proxy:** Nginx (with rate-limited `/login` endpoints)
- **Containerization:** Docker + Docker Compose, multi-stage build on `distroless` base

## Project Layout

```
backend/
├── main.go                # Entry point — auto-migrates models, seeds plans, starts router
├── config/                # JWT secret / app config
├── db/                    # GORM singleton + connection pool setup
├── models/                # User, Admin, Trainer, Subscription, UserSubscription, Class
├── repository/            # Data-access layer (one file per model)
├── service/               # Business logic (auth, user, admin, trainer, subscription, class)
├── handler/               # Gin HTTP handlers, thin wrappers over services
├── middleware/            # JWT auth middleware (role-aware)
├── router/                # Route definitions
├── seed/                  # plans.json + seeder for default subscription tiers
├── Docker-Compose.yaml    # backend + postgres + nginx
├── dockerfile             # Multi-stage build
└── nginx.conf             # Reverse proxy + rate limiting on /login
```

## Domain Model

- **User** — gym member; can subscribe to a plan.
- **Admin** — manages trainers and subscription plans.
- **Trainer** — created under an admin; can create classes.
- **Subscription** — a plan tier (`Basic`, `Premium`) with a price, owned by an admin.
- **UserSubscription** — links a user to a subscription with `StartedAt` / `ExpiresAt`.
- **Class** — scheduled session created by a trainer, with capacity tracking.

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

    User "0..1" --> "*" UserSubscription
    Subscription "1" --> "*" UserSubscription
    Trainer "1" --> "*" Classes
    User "1" --> "*" BookingSituation
    Classes "1" --> "*" BookingSituation
```

A background goroutine in [main.go](backend/main.go) runs `service.CheckSubscriptions` every 24 hours to keep subscription state up to date.

## API

All authenticated routes require a JWT cookie (`key`) with a `role` claim matching the route group.

### Public

| Method | Path             | Description                            |
| ------ | ---------------- | -------------------------------------- |
| GET    | `/subscriptions` | List available subscription plans      |
| GET    | `/classes`       | List scheduled classes                 |

### `/user`

| Method | Path                  | Auth      | Description                       |
| ------ | --------------------- | --------- | --------------------------------- |
| POST   | `/user/signup`        | public    | Register a new user               |
| POST   | `/user/login`         | public    | Login, sets JWT cookie            |
| GET    | `/user/get`           | user      | Get current user                  |
| PUT    | `/user/update`        | user      | Update current user               |
| DELETE | `/user/delete`        | user      | Delete current user               |
| POST   | `/user/subscribe`     | user      | Subscribe to a plan               |
| GET    | `/user/subscription`  | user      | List the user's subscriptions     |

### `/admin`

| Method | Path             | Auth  | Description           |
| ------ | ---------------- | ----- | --------------------- |
| POST   | `/admin/signup`  | public| Register an admin     |
| POST   | `/admin/login`   | public| Login                 |
| GET    | `/admin/get`     | admin | Get current admin     |
| PUT    | `/admin/update`  | admin | Update current admin  |
| DELETE | `/admin/delete`  | admin | Delete current admin  |

### `/trainer`

| Method | Path                | Auth    | Description                  |
| ------ | ------------------- | ------- | ---------------------------- |
| POST   | `/trainer/signup`   | public  | Register a trainer           |
| POST   | `/trainer/login`    | public  | Login                        |
| GET    | `/trainer/get`      | trainer | Get current trainer          |
| PUT    | `/trainer/update`   | trainer | Update current trainer       |
| DELETE | `/trainer/delete`   | trainer | Delete current trainer       |
| POST   | `/trainer/class`    | trainer | Create a class               |

## Running

### With Docker Compose (recommended)

From [backend/](backend/):

```bash
docker compose -f Docker-Compose.yaml up --build
```

This brings up:
- `postgres` on `:5432`
- `backend` (Go API) on internal `:8080`
- `nginx` on `:80`, proxying to the backend and rate-limiting `/login` endpoints (1000 r/s, burst 10)

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
- **Nginx rate-limits `/login`** at 10 r/m (burst 10) per IP — see [nginx.conf](backend/nginx.conf). Hammering a single IP past the burst returns 503s.

