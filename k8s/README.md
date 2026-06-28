# Kubernetes manifests for the gym API

Translation of `backend/Docker-Compose.yaml` to Kubernetes. Designed for a local
single-node cluster (**kind** or **minikube**) on Windows.

## What maps to what

| Compose | Kubernetes (file) |
|---|---|
| `backend` service | Deployment + Service `backend` (`03-backend.yaml`) |
| `postgres` service + `postgres_data` volume | StatefulSet + PVC + Service `postgres` (`01-postgres.yaml`) |
| `nginx` service | Deployment + Service + nginx.conf ConfigMap (`04-nginx.yaml`) |
| `environment:` secrets | Secret `app-secrets` (`00-secret.yaml`) |
| `seed/*.json` bind-mounts | ConfigMap `seed-files` (`02-seed-configmap.yaml`) |
| `healthcheck` / `depends_on` | readiness/liveness probes + initContainer |
| `ports: "80:80"` | NodePort 30080 on the nginx Service |

The Service names `postgres` and `backend` are load-bearing: the Go DSN and
`nginx.conf` address those hostnames, so **no app code changes are needed.**

## Build the backend image and load it into the cluster

Kubernetes pulls images by name — it can't `build:` like Compose. Build locally,
then load the image into the node (no registry needed).

```powershell
cd backend
docker build -t gym-backend:local .

# kind:
kind load docker-image gym-backend:local

# OR minikube:
minikube image load gym-backend:local
```

## Apply everything

```powershell
kubectl apply -f k8s/
kubectl get pods -w        # wait until all are Running / Ready
```

## Reach it

```powershell
# kind: forward the nginx Service to localhost
kubectl port-forward service/nginx 8080:80
#   -> http://localhost:8080

# minikube: get the node URL for the NodePort
minikube service nginx --url
```

## Notes / gotchas

- **Single replica** for the backend on purpose: the seeder and 24h
  subscription-checker goroutines run in-process. Scaling up duplicates them.
  Move them to a Kubernetes `CronJob` before raising `replicas`.
- **Health checks are TCP-only** (no Go changes). All three pods have liveness +
  readiness probes on their port; a failing **liveness** probe restarts the pod,
  a failing **readiness** probe just pulls it out of its Service. The backend also
  has a `startupProbe` so a slow first boot (AutoMigrate + seeding) isn't killed.
  A TCP probe confirms the port is open, not that the DB is reachable — if you
  later want DB-aware readiness, add a `GET /readyz` handler in Go that pings
  Postgres and switch that probe to `httpGet`.
- **Frontend static files** are not on disk in this branch, so `location /` 404s.
  API routes still work. Restore `frontend/FrontEnd/` and bake it into a custom
  nginx image (or another ConfigMap) to serve the UI.
- **Secrets**: `00-secret.yaml` ships dev values (`key: "test"`). Replace before
  any non-local use; don't commit real credentials.
