#!/bin/bash
ALB_DNS=$1
echo "Checking health metrics on http://$ALB_DNS/health"
for i in {1..6}; do
  CODE=$(curl -s -o /dev/null -w "%{http_code}" http://$ALB_DNS/health || true)
  if [ "$CODE" -eq 200 ]; then
    echo "App is healthy!"
    exit 0
  fi
  sleep 10
done
echo "Health check failed."
exit 1
