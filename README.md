**Threat Event Processing System (TEPS)**

This project contains two Go services:

**1. event-file-reader**

**2. ThreatEventProcessingService**

These services work together to generate, retrieve, store, cache, and manage cybersecurity event data using PostgreSQL and Redis.

**Project Structure**

```
go.work
event-file-reader/           # Generates random event data and serves it via an API
    ├── data/events.json     # Auto-generated file with random events
    └── main.go              # Starts the file reader service

ThreatEventProcessingService/
    ├── teps/main.go     # Main entry point for TEPS
    ├── internal/
    │   ├── async/           # Scheduler implementations
    │   ├── handlers/        # HTTP handlers
    │   ├── model/           # Data models
    │   ├── repository/      # DB & cache access
    │   ├── schedulers/      # Scheduler setup
    │   └── service/         # Business logic
    ├── .env                 # Configuration file
    └── go.mod

```

**Service 1: event-file-reader**

Description:
Generates 200 random events and saves them as data/events.json on startup.

Exposes a GET API on port 9090 at:

http://localhost:9090/events

Purpose: This API is consumed by TEPS to fetch events daily and persist them in the database.

Run:
```
cd event-file-reader
go run main.go
```

**Service 2: ThreatEventProcessingService**

Description:
Exposes CRUD APIs for Event data.

Uses PostgreSQL for persistence.

Uses Redis for caching GET calls.

Runs two background schedulers:

1. Fetch Scheduler: Calls the reader service and inserts events into DB daily (according to fetch api cron defined in .env)

2. Cleanup Scheduler: Deletes events older than 24 hours from DB & cache, and uploads them to S3 (at the end of the day)

Run:
```
cd ThreatEventProcessingService
go run teps/main.go
```

**APIs (TEPS)**

Base URL:
http://localhost:8080

Endpoints:

**POST** _/events_ — Create event

Sample Request Body:

```
{
    "id":1,
    "source": "Firewall-1",
    "threat_type": "SQL Injection",
    "detected_at": "2025-05-09T14:00:00Z",
    "processed_at": "2025-05-09T14:05:00Z",
    "details": "Suspicious payload detected in login form"
  }
```

**GET** /events/:id — Get event (uses Redis cache)

**PUT** _/events/_ — Update event

Sample Request Body
```
{
    "id":1,
    "source": "Firewall-1-Edited",
    "threat_type": "SQL Injection",
    "detected_at": "2025-05-09T14:00:00Z",
    "processed_at": "2025-05-09T14:05:00Z",
    "details": "Suspicious payload detected in login form"
  }
```

**DELETE** _/events/:id_ — Delete event

**Environment Configuration**
Both services load configuration from a .env file.

Example .env for ThreatEventProcessingService:

**Setup Instructions**

Prerequisites:
Go 1.21+

PostgreSQL (running locally)

Redis (running locally)

**Steps:**

1. Clone the repo and navigate to root

2. Create the following Postgres table:

**sql**

```
CREATE TABLE events (
    id SERIAL PRIMARY KEY,
    source TEXT,
    threat_type TEXT,
    detected_at TIMESTAMP,
    processed_at TIMESTAMP,
    details TEXT
);
```

3. Start both services:

# Terminal 1
```
cd event-file-reader
go run main.go
```

# Terminal 2
```
cd ThreatEventProcessingService
go run teps/main.go
```

go.work
This project uses a go.work file to link both modules, allowing cross-package imports and shared model references.

Make sure you're at the root level when running commands so the go.work file can resolve dependencies properly.

**Generated File**
**event-file-reader/data/events.json**: This file is overwritten on every startup with new random events.

**Note:**__ Upload S3 part is commented for now, bucket names and key can be modifed in .env file
