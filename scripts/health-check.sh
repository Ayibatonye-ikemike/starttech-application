#!/bin/bash
ALB_DNS=$1
echo "Checking health metrics on http://$ALB_DNS/health"
for i in {1..30}; do
  CODE=$(curl -s -o /dev/null -w "%{http_code}" http://$ALB_DNS/health || true)
  if [ "$CODE" -eq 200 ]; then
    echo "App is healthy!"
    exit 0
  fi
  echo "Attempt $i: Application is still initializing (HTTP status: $CODE). Retrying in 10s..."
  sleep 10
done
echo "Health check failed after 5 minutes."
exit 1
