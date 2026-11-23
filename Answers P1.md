1. The problem is that, the Pod is asking to run as UID 0 which is root but the provider won't allow it,
Delete the runAsUser: 0 from the deployment YAML so the Pod wont run as root.

2. The ports don't match, change the target port and service port to match.

4. use df -H followed by various df/du/find checks (אני לא זוכר בדיוק איזה פקודות, בפועל הייתי נעזר ב AI) to try to find the culprit. if the space was taken up by random files that are misplaced just use rm ,commonly, logs are the problem, in that case use tranucate, never rm, that way the actual logging process closes and frees the space.

5. I don't know what SIP is, I used AI for this question:
- **NAT or firewall blocking the audio ports:** SIP sets up the call fine, but the actual audio (RTP) uses separate UDP ports that might be blocked by a firewall or sent to the wrong local IP behind NAT. To check, I’d capture packets with Wireshark or `tcpdump` to see if RTP traffic is flowing both ways and check the firewall or NAT logs.
- **Asymmetric routing or missing media relay:** If the endpoints are on separate networks or behind double NAT, the audio stream can get stuck without an SBC or media proxy bridging them. I’d inspect the IPs and ports advertised in the SDP (`INVITE` and `200 OK`) and trace the route to see if that media path is actually reachable.
- **Codec mismatch or local hardware issue:** One side might fail to decode the audio due to an unsupported codec, or it could just be a muted mic or faulty headset on one end. I’d check the negotiated codecs in the call logs and run an echo test or swap headsets to isolate the faulty.
