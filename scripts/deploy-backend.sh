#!/bin/bash
set -e

ASG_NAME="starttech-backend-asg"
REGION="us-east-1"

echo "Checking for active instance refreshes on $ASG_NAME..."
STATUS=$(aws autoscaling describe-instance-refreshes \
  --auto-scaling-group-name "$ASG_NAME" \
  --region "$REGION" \
  --query "InstanceRefreshes[?Status=='InProgress'].Status" \
  --output text)

if [ "$STATUS" = "InProgress" ]; then
  echo "⚠️ An instance refresh is already running. Waiting for it to complete instead of breaking..."
  while [ "$STATUS" = "InProgress" ]; do
    sleep 15
    STATUS=$(aws autoscaling describe-instance-refreshes \
      --auto-scaling-group-name "$ASG_NAME" \
      --region "$REGION" \
      --query "InstanceRefreshes[?Status=='InProgress'].Status" \
      --output text)
    echo "Still updating servers..."
  done
  echo "✅ Previous deployment finished. Proceeding with your new deployment..."
fi

# FETCH THE INDEPENDENT TEMPLATE ID AND ENFORCE LATEST MAPPING OVERRIDES
TEMPLATE_ID=$(aws autoscaling describe-auto-scaling-groups \
  --auto-scaling-group-names "$ASG_NAME" \
  --region "$REGION" \
  --query "AutoScalingGroups[0].LaunchTemplate.LaunchTemplateId" \
  --output text)

echo "🔄 Overriding ASG definition tracking to map explicit template ID: $TEMPLATE_ID at \$Latest version..."
aws autoscaling update-auto-scaling-group \
  --auto-scaling-group-name "$ASG_NAME" \
  --launch-template "LaunchTemplateId=$TEMPLATE_ID,Version=\$Latest" \
  --region "$REGION"

echo "🚀 Triggering fresh rolling update..."
aws autoscaling start-instance-refresh \
  --auto-scaling-group-name "$ASG_NAME" \
  --preferences '{"MinHealthyPercentage": 50}' \
  --region "$REGION"
