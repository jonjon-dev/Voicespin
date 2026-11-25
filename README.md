## call-events

A lightweight Go HTTP microservice for tracking call events in memory, accompanied by a Python JSONL summary script, OpenShift/Kubernetes deployment manifests, and a GitLab CI/CD pipeline.

Prerequisites

Go 1.22+

Python 3.10+

Docker / Podman (optional, for container builds)

Local Setup & Run

Initialize and start the Go microservice:

go mod init call-events # if go.mod does not exist yet
go run main.go


The server listens on http://localhost:8080.

Test the endpoints locally:

# Health check
curl -i http://localhost:8080/healthz

# Send a valid event (expects 202 Accepted)
curl -i -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{"session_id": "call-123", "event_type": "call.started"}'

# Retrieve event statistics
curl -s http://localhost:8080/stats


Run automated Go tests (with race detector):

go test -race -v ./...


Run the Python JSONL summary script:

python3 summary.py events.jsonl


AI Usage & Verification Summary

## AI Usage & Verification

* **AI Tool Used:** Gemini was used to generate the initial boilerplate for the Go microservice, Python script, OpenShift/Kubernetes manifests, and GitLab CI/CD pipeline.
* **Manual Verification:** I cross-checked the generated code, CLI flags, YAML fields, and pipeline syntax directly against the official documentation (Go standard library docs, Python `json`/`argparse` docs, OpenShift/Kubernetes specs, and GitLab CI/CD YAML reference) rather than running a live cluster deployment.
* **Fixes & Limitations Found:**
  * **OpenShift `restricted` SCC alignment:** In the initial manifest/Dockerfile, I verified and ensured that the container does not run as `root` (`runAsNonRoot: true`, `allowPrivilegeEscalation: false`, `capabilities.drop: ["ALL"]`, and a non-root `USER` in the `Dockerfile`) so OpenShift's default `restricted` SCC would not block the pod.
  * **GitLab CI/CD variable & rule syntax:** Verified that predefined CI variables (`$CI_COMMIT_SHA`, `$CI_REGISTRY_IMAGE`) and `rules` with `when: manual` matched official GitLab CI syntax so `deploy` only triggers manually on `main`.
  * **Scope limitation:** The Go service stores counters in memory (`map` + `sync.RWMutex`), so data resets if a pod restarts—suitable for this task, though production would need an external database or Redis.
