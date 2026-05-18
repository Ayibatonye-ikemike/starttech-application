#!/bin/bash
set -e
cd frontend
npm ci
npm run build
aws s3 sync dist/ s3://starttech-frontend-app-unique-bucket-xyz --delete
aws cloudfront create-invalidation --distribution-id $CLOUDFRONT_DIST_ID --paths "/*"
