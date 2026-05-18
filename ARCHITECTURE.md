# StartTech System Architecture Documentation

This document provides a comprehensive overview of the technical architecture, network design, and data flows governing the StartTech full-stack application ecosystem.

## 📐 System Topology Overview
The ecosystem uses a highly available, decoupled architecture split into three distinct layers: Content Delivery, Compute Gateway, and Managed Data Tier.

```text
[ Client Browser ] 
       │
       ├──► (HTTPS/443) ──► [ CloudFront CDN ] ──► [ S3 Frontend Bucket ]
       │
       └──► (HTTP/80) ────► [ Application Load Balancer ]
                                   │
                           ┌───────┴───────┐
                           ▼               ▼
                     (Port 8080)     (Port 8080)
                     [ EC2 Node ]    [ EC2 Node ]
                     (ASG Fleet - Public Subnets)
                           │
                           ▼
               [ MongoDB Atlas Cluster ]
```

## 🌐 Component Architectural Breakdown

### 1. Presentation Tier (Frontend Platform)
* **Hosting Base:** AWS S3 configured for static web hosting.
* **Delivery Engine:** AWS CloudFront CDN distribution serving traffic over HTTPS via global edge locations.
* **Deployment Isolation:** Direct public access to S3 is blocked; content is fetched safely using an Origin Access Identity (OAI).

### 2. Application Logic Tier (Compute Gateway)
* **Orchestration:** AWS Auto Scaling Group (ASG) maintaining a highly available fleet of variable EC2 instances (`t3.micro`) spanning multiple Availability Zones.
* **Networking Context:** Instances reside in Public Subnets mapped directly against Internet Gateways to facilitate reliable outward dependency compilation (`docker pull`) and external database transactions.
* **Runtime Sandbox:** Docker containers packaging an independent compilation layer running Golang 1.25.
* **Inbound Access Control:** Servers accept ingress traffic *strictly* on port `8080` forwarded by the application load balancer security group proxy.

### 3. Storage & Data Persistence Tier
* **Database Platform:** MongoDB Atlas Cloud Cluster.
* **Security Mechanics:** Network Access Lists (ACLs) whitelisted globally to accept requests matching validated alphanumeric cryptographic authentication handshake variables.

## 🔄 Core Data & Execution Flows
1. **Frontend Bootstrapping:** User clients query the CloudFront URL edge path to load static index, JavaScript, and style layers from S3.
2. **API Communication:** User interaction triggers REST requests to the Load Balancer DNS over standard HTTP port 80.
3. **Internal Forwarding:** The ALB maps path parameters down to port `8080` across active, healthy target EC2 fleet instances.
4. **Data Operations:** The Go backend application validates authentication headers, establishes runtime connections with MongoDB Atlas, and yields unified JSON payload formats.
