---
title: Database Schema
tags: [database, postgres, backend]
aliases: [DB Design, Database Design]
---

# Database Schema

## Overview

We use **PostgreSQL** as our primary database.

## Tables

### Users Table

The users table stores account information:

- `id` — UUID primary key
- `email` — unique, not null
- `name` — display name
- `created_at` — timestamp

## Migrations

Database migrations are managed with a simple versioning system.

### Version 1: Initial Schema

The initial migration creates the core tables.

### Version 2: Add Indexes

Added indexes for performance on `email` and `created_at`.

## Related Notes

- [[api-design]] — API that queries this database
- [[infrastructure]] — Where the database runs
