---
title: API Design
tags: [api, design, backend]
aliases: [API Design Document, REST API Spec]
---

# API Design

## Overview

This document describes the REST API design for our service.

## Authentication

All endpoints require a Bearer token:

```bash
curl -H "Authorization: Bearer <token>" https://api.example.com/v1/
```

## Endpoints

### GET /users

Returns a list of users:

```go
func ListUsers(w http.ResponseWriter, r *http.Request) {
    users, err := db.Query("SELECT * FROM users")
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(users)
}
```

### POST /users

Creates a new user. Configuration example:

```yaml
api:
  host: "0.0.0.0"
  port: 8080
  timeout: 30s
  database:
    driver: postgres
    url: ${DATABASE_URL}
```

## Error Handling

Errors are returned in a standard format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid email format",
    "details": [
      {"field": "email", "issue": "not a valid email address"}
    ]
  }
}
```

## See Also

- [[database]] - Database schema
- [[infrastructure]] - Deployment setup
