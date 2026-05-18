# StartTech Full-Stack Application Delivery (CI/CD)

This repository contains the full-stack application codebase alongside the advanced dual-track continuous deployment pipelines built to deliver secure, zero-downtime rolling environments.

## ⚙️ Repository Layout Checklist
* **`/Client`**: React.js / Vite frontend user interface environment.
* **`/Server/MuchToDo`**: Highly optimized Golang REST API server tier.
* **`/.github/workflows`**: Native GitHub Actions pipeline tracks governing automated testing, image publishing, rolling deployments, and health verifications.
* **`/scripts`**: Automation scripts handling deployments, custom loop health checks, and emergency rollbacks.

## 🚢 Dual-Track CI/CD Pipeline Mechanics

### 1. Frontend Delivery Flow (`frontend-ci-cd.yml`)
* **Trigger:** Triggered automatically on pushes to the `main` branch modifying files inside the `Client/**` path.
* **Actions:** Runs validation audits, builds static compilation files via Node.js, syncs compilation assets directly to the target S3 bucket, and forces a global CloudFront cache invalidation.

### 2. Backend Delivery Flow (`backend-ci-cd.yml`)
* **Trigger:** Triggered automatically on pushes to the `main` branch modifying files inside the `Server/**` path.
* **Actions:** Runs the Go testing suite, compiles a lightweight multi-stage Docker image layer using Go `1.25`, and publishes it to Docker Hub.
* **Deployment Execution:** Triggers an AWS Auto Scaling **Instance Refresh** rolling update to progressively swap cluster nodes, waits out an intentional initialization cooldown window, and triggers an advanced custom health check script.

## 🛠️ Operations & Operations Runbook

### Troubleshooting a Stuck 502 Bad Gateway
If the pipeline health loop fails on an incoming deployment track, execute these sequential checks directly from your local terminal terminal:

1. **Verify Local Docker Container Integrity:**
   ```bash
   docker run --name local-test -p 8080:8080 -e MONGO_URI="[Your-URI]" -e PORT="8080" -e HOST="0.0.0.0" ayibatonye/much-to-do-backend:latest
   ```
2. **Check AWS Instance Build Progress:**
   ```bash
   aws autoscaling describe-instance-refreshes --auto-scaling-group-name "starttech-backend-asg" --region us-east-1
   ```

### Executing a Manual Rollback
If a production deployment goes sideways and you need to safely restore your cluster environment to a known baseline profile instantly:
```bash
./scripts/rollback.sh
```
