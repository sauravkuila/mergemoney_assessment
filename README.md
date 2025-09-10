
# Cross-Border Money Transfer System

## Overview

This project implements a robust, secure, and scalable cross-border money transfer system. It enables users to send money internationally with minimal fees and fast settlement times, integrating with third-party payment providers via asynchronous APIs and webhooks.

## Architecture

- **API Layer (`cmd/api/`)**: Exposes RESTful endpoints for user authentication, account linking, money transfer initiation, and webhook handling.
- **Service Layer (`pkg/service/`)**: Contains business logic for user management, transaction processing, and integration with payment providers.
- **DAO Layer (`pkg/dao/`)**: Handles database operations for users, transactions, and account information.
- **DTOs (`pkg/dto/`)**: Defines data transfer objects for API requests/responses and internal service communication.
- **Middleware (`pkg/middleware/`)**: Implements authentication, authorization, and request validation.
- **Logger (`pkg/logger/`)**: Centralized logging for observability and monitoring.
- **Utils (`pkg/utils/`)**: Utility functions for encryption, validation, and other common tasks.
- **Config (`pkg/config/`)**: Configuration management for environment variables and constants.

## Functional Highlights

### 1. User Authentication & Account Linking
- Users log in using their mobile number.
- Third-party API fetches all financial accounts linked to the mobile number.
- Users select the source account for transfers.
- Multi-factor authentication (MFA) is supported.

### 2. Money Transfer Workflow
- Users specify source/destination currency, amount, and recipient details.
- System fetches real-time exchange rates and fees.
- Transfer requests are sent to the payment provider asynchronously.
- Transaction states: Initiated, Pending, In Progress, Completed, Failed.

### 3. Transaction Processing & Asynchronous Status Updates
- Webhook endpoints handle status updates from providers.
- Idempotency and out-of-order update handling ensure consistency.
- Failed transactions trigger reversal/refund logic.
- SMS notifications for key status changes.

### 4. Integration with Third-Party Payment Providers
- Secure API integration for transfer initiation and status updates.
- Supports multiple providers with fallback mechanisms.
- Refunds and reversals are handled gracefully.

## Non-Functional Features

- **Performance**: Optimized for P95 latency < 200ms for transfer requests.
- **Observability**: Centralized logging, monitoring, and alerting.
- **Security**: Encryption for sensitive data, fraud detection, and compliance mechanisms.

## Project Structure

```
mergemoney_assessment/
├── cmd/api/                # API entrypoint
├── pkg/
│   ├── config/             # Configuration management
│   ├── constant/           # Constants and enums
│   ├── dao/                # Data access objects
│   │   └── user/           # User-related DB operations
│   ├── database/           # DB models and connections
│   ├── dto/                # Data transfer objects
│   ├── logger/             # Logging utilities
│   ├── middleware/         # Middleware (auth, validation)
│   ├── server/             # Server/router setup
│   ├── service/            # Business logic
│   │   └── v1/             # Versioned services
│   │       └── login/      # Login service
│   └── utils/              # Utility functions
└── README.md
```

## High-Level Design

- **APIs**: RESTful endpoints for authentication, account linking, transfer initiation, and webhook updates.
- **Database**: Models for users, accounts, transactions, and audit logs.
- **Messaging/Async**: Webhook handlers ensure idempotency and correct transaction state transitions.
- **Security**: Data encryption, MFA, and fraud detection.

## Asynchronous Processing Strategy

- Webhook endpoints are idempotent and handle duplicate/out-of-order updates.
- Transaction state machine ensures consistent state transitions.
- Reconciliation logic detects mismatches between internal and provider records.

## Security & Compliance

- Sensitive data is encrypted at rest and in transit.
- MFA and fraud detection mechanisms are implemented.
- Audit logs for compliance and traceability.

## Bonus Features

- Multiple payment provider support with failover.
- Reconciliation logic for transaction consistency.

## References

- [Sample Architecture Diagram](https://d2908q01vomqb2.cloudfront.net/fc074d501302eb2b93e2554793fcaf50b3bf7291/2023/02/13/adverse_1.png)
- [Sample Data Model](https://media.geeksforgeeks.org/wp-content/uploads/20231215171020/Data-model-design-2.jpg)