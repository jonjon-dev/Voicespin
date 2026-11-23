# Verification & Troubleshooting

## How to Verify the Deployment

1. **Check rollout & pod status:**
   ```bash
   oc rollout status deployment/call-events
   oc get pods -l app=call-events -o wide
   ```
   Confirm 2 pods are `Running` and `READY` (`1/1`).

2. **Check Service, Endpoints, and Route:**
   ```bash
   oc get svc call-events
   oc get endpoints call-events
   oc get route call-events
   ```
   Verify `endpoints` lists both pod IPs on port `8080`.

3. **Test health, event ingestion, and stats over HTTPS:**
   ```bash
   HOST=$(oc get route call-events -o jsonpath='{.spec.host}')
   curl -i -X GET "https://${HOST}/healthz"
   curl -i -X POST "https://${HOST}/events" \
     -H "Content-Type: application/json" \
     -d '{"session_id":"call-123","event_type":"call.started"}'
   curl -s "https://${HOST}/stats"
   ```

## Troubleshooting If Deployment Fails

* **`Pending` pods:** Check node capacity or resource quotas with `oc describe pod <pod-name>`.
* **`ImagePullBackOff` / `ErrImagePull`:** Verify the image name/tag, pull secrets, and registry network access.
* **`CrashLoopBackOff` / failing probes:** Inspect `oc logs <pod-name> --previous` and `oc describe pod <pod-name>`; confirm the app listens on `0.0.0.0:8080` and `/healthz` returns `200`.
* **SCC / permission errors:** Ensure non-root settings, `allowPrivilegeEscalation: false`, and dropped capabilities are set.
* **No traffic via Route:** Check that the `Service` selector matches pod labels, `targetPort` is `8080`, and `oc get endpoints call-events` lists pod IPs.