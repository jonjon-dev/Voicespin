# call-events

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

AI Tool Used: Gemini was used to draft initial boilerplate for the Go microservice, Python JSONL parser, Kubernetes/OpenShift YAML manifests, GitLab CI/CD pipeline, and documentation.

Manual Verification Performed:

Executed go test -race ./... locally and verified endpoint behavior (202 on valid payloads, 400 on malformed/incomplete JSON, accurate /stats counts, and /healthz responses) using curl.

Tested summary.py against valid JSONL inputs, malformed JSON lines, empty files, and nonexistent file paths to confirm graceful error handling without crashes.

Reviewed Kubernetes/OpenShift manifests against restricted SCC constraints (runAsNonRoot: true, dropped capabilities, non-root user in Dockerfile, matching labels/ports) and verified GitLab CI stage gating and variable usage.

Issues / Limitations Identified in AI Output:

Initial deployment manifests needed explicit verification to ensure the Dockerfile (USER 10001:10001) and Deployment securityContext aligned with OpenShift's restricted SCC without hardcoding root UIDs.

The Go service stores counters strictly in memory (EventStore), meaning counts reset across pod restarts or scale-out replicas—acceptable for this assessment scope, though production workloads would require an external store (e.g., Redis).
