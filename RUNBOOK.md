# StartTech Operations & Troubleshooting Runbook

This operations document acts as a standardized blueprint for monitoring, maintaining, and recovering the live full-stack production application infrastructure.

## 🚨 Critical Incident: Troubleshooting HTTP 502 Bad Gateway

An HTTP 502 error implies that the Application Load Balancer is healthy, but the backend Go application containers on your EC2 instances are failing to accept connections. Follow this sequential escalation pathway to remediate:

### Phase 1: Local Docker Container Verification
Before debugging cloud network paths, confirm that the application container builds and authenticates cleanly without a localized panic loop. Run this command on your MacBook:
```bash
docker run --name instance-debug-test \
  -p 8080:8080 \
  -e MONGO_URI='mongodb+srv://tonye:Tonye1932014@starttech-cluster.p5oakcb.mongodb.net/?appName=starttech-cluster' \
  -e PORT='8080' \
  -e HOST='0.0.0.0' \
  -e DB_NAME='starttech' \
  -e JWT_SECRET_KEY='StagingStagingSecretKey123' \
  -e USE_REDIS='false' \
  ayibatonye/much-to-do-backend:latest
```
* **Expected Pass Log:** `level=INFO msg="Successfully connected to MongoDB."` followed by `🚀 StartTech Server binding globally`.
* **Failure Remediation:** If the container exits with a SASL authentication failure, the MongoDB cluster password has been modified or expired. Correct the connection string parameters inside your GitHub Actions secrets vault.

### Phase 2: Live AWS Cloud Log Inspection
If the container passes locally, pull the real-time installation and boot logs directly from the active AWS cloud instance using your MacBook terminal:
```bash
# 1. Fetch your active running instance IDs
aws ec2 describe-instances --filters "Name=instance-state-name,Values=running" "Name=instance.group-name,Values=starttech-ec2-sg" --query "Reservations[*].Instances[*].[InstanceId,LaunchTime]" --region us-east-1 --output table

# 3. Extract the unfiltered host console boot buffer logs
aws ec2 get-console-output --instance-id [INSERT_INSTANCE_ID_FROM_ABOVE] --region us-east-1 --output text | tail -n 30
```
* **Check for Core Crashes:** Scan the log output lines. If the final segment terminates at `No such container: app`, the `user_data` script was interrupted before execution. Ensure your infrastructure compute templates are cleanly deployed on AWS with strict failure handling disabled (`set +e`).

---

## ♻️ Manual Emergency Rollback Plan
If an active deployment push triggers a breaking production anomaly and you must revert your EC2 instances back to a known stable application baseline configuration:

1. **Abruptly cancel any stuck rollouts on AWS:**
   ```bash
   aws autoscaling cancel-instance-refresh --auto-scaling-group-name "starttech-backend-asg" --region us-east-1
   ```
2. **Execute your customized localized recovery template script:**
   ```bash
   ./scripts/rollback.sh
   ```
This drops any corrupted active instances and re-provisions clean, operational machines using your verified baseline launch templates automatically.
