---
title: Infrastructure
tags: [infrastructure, devops, kubernetes]
aliases: [Infra, Deployment Setup]
---

# Infrastructure

## Overview

Our infrastructure runs on **Kubernetes** with Helm charts.

## Components

- **API Servers** — 3 replicas, auto-scaling
- **Database** — PostgreSQL with streaming replication
- **Cache** — Redis cluster for session storage
- **CDN** — CloudFront for static assets

## Monitoring

We use Grafana dashboards for observability.

> This is a blockquote with some deployment notes.

## Related

- [[api-design]] — API served by this infrastructure
- [[database]] — Database running on this infrastructure
