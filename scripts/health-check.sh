#!/bin/bash
set -e

ALB_DNS=$1
if [ -z "$ALB_DNS" ]; then
    echo "❌ Error: Missing Application Load Balancer DNS URL."
    exit 1
fi

echo "🔍 Starting deployment verification against: http://$ALB_DNS/health"
MAX_ATTEMPTS=30
DELAY_SECONDS=15

for ((i=1; i<=MAX_ATTEMPTS; i++)); do
    # Fetch HTTP status code and response body
    RESPONSE=$(curl -s -w "\n%{http_code}" "http://$ALB_DNS/health" || true)
    HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
    BODY=$(echo "$RESPONSE" | sed '$d')

    if [ "$HTTP_STATUS" -eq 200 ]; then
        echo "🟢 SUCCESS: Application health check passed with status 200!"
        echo "📥 Response Payload: $BODY"
        exit 0
    elif [ "$HTTP_STATUS" -eq 502 ]; then
        echo "⚠️ Attempt $i/$MAX_ATTEMPTS: Received HTTP 502 (Bad Gateway). ALB cannot reach Go app container yet."
    else
        echo "⚠️ Attempt $i/$MAX_ATTEMPTS: Server returned HTTP status $HTTP_STATUS."
    fi

    sleep $DELAY_SECONDS
done

echo "🚨 CRITICAL: Deployment health verification timed out after $((MAX_ATTEMPTS * DELAY_SECONDS / 60)) minutes."
exit 1
