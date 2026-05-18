#!/bin/bash
set -e
aws autoscaling start-instance-refresh \
  --auto-scaling-group-name "starttech-backend-asg" \
  --preferences '{"MinHealthyPercentage": 50}' \
  --region us-east-1
