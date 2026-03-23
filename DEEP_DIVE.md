# Technical Deep Dive - Sigma AI Compliance Dashboard

← [Back to README](README.md)

This document covers the architecture, engineering decisions, and implementation details behind the Sigma AI Compliance Dashboard. Intended for engineers reviewing the codebase.

**Pastel FDE Assessment - George Uche-Umeh**

A full-stack compliance demonstration platform that integrates three of Pastel's Sigma AI product capabilities via a clean, browser-based interface. Built in Go with a focus on clean architecture, correct API usage, and clear data presentation for both technical and non-technical audiences.

> The live deployment is available at: **https://pastel-fde-assessment.onrender.com/**


---

## Table of Contents

- [Technical Deep Dive - Sigma AI Compliance Dashboard](#technical-deep-dive---sigma-ai-compliance-dashboard)
  - [Table of Contents](#table-of-contents)
  - [What This Application Does](#what-this-application-does)
  - [Technology Stack \& Why Go](#technology-stack--why-go)
    - [Language: Go](#language-go)
    - [Frontend: HTMX + Tailwind CSS](#frontend-htmx--tailwind-css)
    - [Dependency Injection: `samber/do`](#dependency-injection-samberdo)
    - [Logging: `zerolog`](#logging-zerolog)
  - [Repository Structure](#repository-structure)
  - [Architecture Overview](#architecture-overview)
    - [Request Flow (Transaction Example)](#request-flow-transaction-example)
  - [The Sigma SDK (`pkg/sigma`)](#the-sigma-sdk-pkgsigma)
    - [Why a Custom SDK?](#why-a-custom-sdk)
    - [`SigmaClient` Interface (`interface.go`)](#sigmaclient-interface-interfacego)
    - [Live Client (`client.go`)](#live-client-clientgo)
    - [Typed Request Models (`transaction.go`, `aml.go`, `adverse_media.go`)](#typed-request-models-transactiongo-amlgo-adverse_mediago)
    - [Mock Client (`mock.go`)](#mock-client-mockgo)
  - [Configuration \& Environment Variables](#configuration--environment-variables)
  - [Running Locally](#running-locally)
  - [Running with Docker](#running-with-docker)
    - [Makefile](#makefile)
  - [Feature Walkthrough](#feature-walkthrough)
    - [Task 1: Transaction Monitoring](#task-1-transaction-monitoring)
    - [Task 2: PEP \& Sanctions Screening](#task-2-pep--sanctions-screening)
    - [Task 3: Adverse Media Screening](#task-3-adverse-media-screening)
  - [Design Decisions \& Going Above and Beyond](#design-decisions--going-above-and-beyond)
    - [1. A Complete SDK, Not Inline HTTP Calls](#1-a-complete-sdk-not-inline-http-calls)
    - [2. Interface-Driven Design for Testability](#2-interface-driven-design-for-testability)
    - [3. Optional Fields as Pointer Types](#3-optional-fields-as-pointer-types)
    - [4. Information Architecture on Result Pages](#4-information-architecture-on-result-pages)
    - [5. Full Responsiveness](#5-full-responsiveness)
    - [6. Dual Base URL Routing](#6-dual-base-url-routing)
    - [7. Graceful Shutdown](#7-graceful-shutdown)
  - [Mock Mode vs Live API](#mock-mode-vs-live-api)

---

## What This Application Does

This platform allows a user to interact with three Sigma AI compliance APIs through a modern web dashboard, see exactly what data is sent to Sigma, receive and understand the full Sigma response, and get plain-English explanations of what every field and every result means.

The three capabilities demonstrated are:

| Capability                    | What It Does                                                                                                                                                                                                                                                                |
| :---------------------------- | :-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Transaction Monitoring**    | Submits a financial transaction to Sigma's AI engine for real-time fraud risk analysis. Returns a recommended action (approve / flag / block), a numeric risk score, the triggering rule, and inline PEP/sanctions checks on the transaction parties.                       |
| **PEP & Sanctions Screening** | Screens an individual's name against global Politically Exposed Persons (PEP) lists and international sanctions databases (OFAC, UN, EU, UK, and others). Returns match confidence scores and full entity profiles for any hits found.                                      |
| **Adverse Media Screening**   | Searches global news and media archives for negative coverage tied to an individual or entity covering fraud allegations, criminal investigations, money laundering, and financial misconduct. Operates asynchronously, with results delivered via webhook in production. |

---

## Technology Stack & Why Go

### Language: Go

Go was selected as the primary language because it is my most fluent server-side language. For a technical assessment with a tight deadline, working in the language I know best let me focus on what matters most: getting the API integration right, building a clean architecture, and presenting data clearly, rather than fighting the language.

Go also offers real practical advantages for this kind of integration work:

- **Compiled and fast;** startup time is near-instant; the binary runs without a runtime environment
- **Excellent standard library;** `net/http`, `html/template`, and `encoding/json` cover everything needed without external frameworks
- **Strong typing;** the type system enforces correct API payloads at compile time; any field mismatch is caught before the code runs
- **Small, readable code;** request/response models are self-documenting structs with JSON struct tags

### Frontend: HTMX + Tailwind CSS

Rather than building a separate Single-Page Application and a REST API, the frontend uses **HTMX** to progressively enhance the Go-rendered HTML. When a user submits a form, HTMX fires an HTTP POST to the Go server and swaps only the result section of the page with the server's response - no page reload, no JavaScript framework, no build step.

**Why this matters for an FDE context:** HTMX lets a backend-focused engineer (like a typical FDE) build rich, interactive UIs without JavaScript expertise. The result is a simpler, more maintainable codebase with zero frontend build tooling.

**Tailwind CSS** (via CDN) handles all styling. The design uses a dark-themed, responsive layout with a fixed sidebar for navigation, product-specific colour accents (blue for transactions, amber for PEP/sanctions, rose for adverse media), and a mobile-first approach throughout.

### Dependency Injection: `samber/do`

Rather than passing dependencies through function arguments across every layer, the application uses a small DI container (`samber/do/v2`). This lets the application decide at startup whether to wire up the live Sigma client or the mock client - based on a single environment variable - without touching a single line of business logic.

### Logging: `zerolog`

Structured JSON logging via `zerolog`, readable in development mode with coloured console output.

---

## Repository Structure

```
pastel-fde-assessment/
│
├── main.go                    # Entry point - wires everything together
│
├── internal/
│   ├── config/
│   │   └── config.go          # Reads environment variables into a Config struct
│   └── di/
│       └── container.go       # Dependency injection wiring (client, services)
│
├── pkg/
│   └── sigma/                 # Custom Sigma SDK (see section below)
│       ├── interface.go       # SigmaClient interface - the contract
│       ├── client.go          # Live HTTP client implementation
│       ├── mock.go            # Mock client for UI testing
│       ├── transaction.go     # Request/response models for Transaction Monitoring
│       ├── aml.go             # Request/response models for PEP & Sanctions
│       ├── adverse_media.go   # Request/response models for Adverse Media
│       └── errors.go          # Typed API errors with sentinel error values
│
├── services/
│   ├── transaction_service.go # Maps form input → SDK request → result struct
│   ├── screening_service.go   # PEP/sanctions service layer
│   └── adverse_media_service.go
│
├── handlers/
│   ├── home_handler.go        # Serves the dashboard home page
│   ├── transaction_handler.go # Parses form data + calls service + renders result
│   ├── screening_handler.go   # PEP and sanctions handlers
│   └── adverse_media_handler.go
│
├── templates/                 # Go HTML templates (one file per view)
│   ├── layout.html            # Global chrome: sidebar, mobile header, CSS
│   ├── home.html              # Dashboard landing page
│   ├── transaction.html       # Transaction form
│   ├── transaction_result.html
│   ├── screening.html         # PEP & sanctions form
│   ├── screening_result.html
│   ├── adverse_media.html     # Adverse media form
│   └── adverse_media_result.html
│
├── Dockerfile                 # Multi-stage Docker build
├── docker-compose.yml         # One-command local container setup
├── Makefile                   # Common development commands
└── assessment.md              # Original assessment brief
```

---

## Architecture Overview

The application follows a clean, layered architecture. Each layer has one responsibility and talks only to the layer directly below it.

```
Browser (HTMX forms)
        │  HTTP POST
        ▼
  Handler Layer        (handlers/)
  Parse → Validate → Call service → Render template
        │
        ▼
  Service Layer        (services/)
  Map form input → SDK request → Call client → Return result
        │
        ▼
  SDK / Client Layer   (pkg/sigma/)
  Build HTTP request → Call Sigma API → Decode response → Typed error
        │
        ▼
  Sigma API            (external)
  sigmaprod.sabipay.com  /  sigmaaml.sabipay.com
```

### Request Flow (Transaction Example)

1. User fills the transaction form and clicks **Submit Transaction for AI Analysis**
2. HTMX fires `POST /transaction/submit` with the form body
3. `TransactionHandler.Submit` parses every form field (strings, floats, booleans, timestamps)
4. A `TransactionInput` struct is populated and passed to `TransactionService.SubmitTransaction`
5. The service maps the input into Sigma's `SubmitTransactionRequest` structure - setting optional pointer fields to `nil` when empty, preventing spurious data from being sent
6. The SDK client serialises the struct to JSON, attaches `apiKey` and `apiSecret` headers, and POSTs to `api/v1/transaction-monitoring/instant`
7. The response is decoded into a typed `TransactionResponse`, returned up the chain
8. The handler renders `transaction_result.html` with the result, which HTMX swaps into the page

---

## The Sigma SDK (`pkg/sigma`)

One of the deliberate architectural decisions was to write a **purpose-built SDK** for the Sigma API rather than making raw `http.Get` calls inside handlers. This is the `pkg/sigma` package.

### Why a Custom SDK?

Raw API calls scattered through handlers are fragile, hard to test, and hard to maintain. A SDK encapsulates:

- Authentication (header injection on every request)
- URL construction (two base URLs: one for transaction monitoring, one for AML)
- Request serialisation and response deserialisation  
- Error handling with typed sentinel errors
- A clean, testable interface

### `SigmaClient` Interface (`interface.go`)

```go
type SigmaClient interface {
    SubmitTransaction(ctx context.Context, req *SubmitTransactionRequest) (*TransactionResponse, error)
    CheckPEP(ctx context.Context, req *ScreeningRequest) (*ScreeningResponse, error)
    CheckSanction(ctx context.Context, req *ScreeningRequest) (*ScreeningResponse, error)
    CheckAdverseMedia(ctx context.Context, req *AdverseMediaRequest) (*AdverseMediaResponse, error)
}
```

All services depend on this interface, not on the concrete `*Client`. This means the live client and the mock client are interchangeable - the business logic is completely unaware of which one it's talking to.

### Live Client (`client.go`)

The live client makes real HTTP requests to the Sigma API. Key design decisions:

- **Two base URLs** - Transaction Monitoring uses `sigmaprod.sabipay.com`; PEP, Sanctions, and Adverse Media use `sigmaaml.sabipay.com`. The client has separate `doSigmaRequest` and `doAMLRequest` methods that route to the correct base.
- **Custom header auth** - Sigma uses `apiKey` and `apiSecret` as request headers (not `Authorization: Bearer`). The client sets these on every outgoing request.
- **10-second timeout** - prevents the server from hanging indefinitely on a slow Sigma response.
- **Structured error handling** - non-2xx responses are decoded into `SigmaAPIError`, which wraps typed sentinel errors (`ErrUnauthorized`, `ErrRateLimited`, etc.) so callers can use `errors.Is` to handle specific failure modes.
- **Body draining** - after reading the response, the remaining body bytes are always discarded and the body closed, preventing connection leaks.

### Typed Request Models (`transaction.go`, `aml.go`, `adverse_media.go`)

Every field the Sigma API accepts is represented as a named Go type. Optional fields use pointer types (`*string`, `*bool`, `*float64`) so that `omitempty` in the JSON struct tag causes them to be omitted entirely when not provided - rather than sending zero-value noise like `"email": ""` or `"balanceBefore": 0`.

Enum values (transaction type, channel, account type, severity level, action type) are defined as named string constants:

```go
const (
    TxTypeDebit  TransactionType = "debit"
    TxTypeCredit TransactionType = "credit"
)
```

This makes incorrect values impossible to express at the call site.

### Mock Client (`mock.go`)

The mock client implements `SigmaClient` with randomised but realistic responses - realistic enough to exercise every part of the result templates. It generates:

- Random risk scores, rule names, reason codes and messages
- Random PEP/sanctions match results with realistic entity profiles (names, aliases, positions, political party affiliations, birth dates, associated countries)
- Random adverse media findings with categorised risk levels and real source publication names

This was built specifically to allow full UI testing without burning through live API quota or requiring a live internet connection during development.

---

## Configuration & Environment Variables

All configuration is read from environment variables at startup. Every variable has a sensible default.

| Variable             | Default                          | Description                                               |
| :------------------- | :------------------------------- | :-------------------------------------------------------- |
| `SIGMA_API_KEY`      | *(assessment key)*               | Sigma API authentication key                              |
| `SIGMA_API_SECRET`   | *(assessment secret)*            | Sigma API authentication secret                           |
| `SIGMA_BASE_URL`     | `https://sigmaprod.sabipay.com/` | Base URL for Transaction Monitoring                       |
| `SIGMA_AML_BASE_URL` | `https://sigmaaml.sabipay.com/`  | Base URL for PEP, Sanctions, Adverse Media                |
| `PORT`               | `80`                             | HTTP port the server listens on                           |
| `USE_MOCK`           | `true`                           | When `true`, uses the mock client instead of the live API |

---

## Running Locally

**Prerequisites:** Go 1.21 or later.

```bash
# Clone the repository
git clone https://github.com/UcGeorge/pastel-fde-assessment.git
cd pastel-fde-assessment

# Run with the mock client (no API credentials needed)
USE_MOCK=true go run main.go

# Run against the live Sigma API
USE_MOCK=false \
  SIGMA_API_KEY=your-key \
  SIGMA_API_SECRET=your-secret \
  go run main.go
```

The server starts on port 80 by default. Open `http://localhost` in your browser.

To use a different port:

```bash
PORT=8080 USE_MOCK=true go run main.go
# Then open http://localhost:8080
```

---

## Running with Docker

A multi-stage Dockerfile produces a minimal image. The app binary is compiled in the builder stage and copied to a small `alpine` base image.

```bash
# Build and start with docker-compose
docker-compose up --build

# Or build and run manually
docker build -t sigma-dashboard .
docker run -p 8080:80 \
  -e USE_MOCK=false \
  -e SIGMA_API_KEY=your-key \
  -e SIGMA_API_SECRET=your-secret \
  sigma-dashboard
```

The `docker-compose.yml` reads credentials from `.env.docker` for convenience.

### Makefile

```bash
make build    # Compile the binary
make run      # Run with mock mode
make docker   # Build and run the Docker container
```

---

## Feature Walkthrough

### Task 1: Transaction Monitoring

**Form page:** `/transaction`

The transaction form is the most comprehensive in the application. It mirrors the full `SubmitTransactionRequest` model from the Sigma SDK and is organized into logical sections:

**Transaction Details** *(required)* - Reference ID, amount, currency, date/time, channel (card payment, bank transfer, ATM, POS, etc.), type (debit or credit), completion status, and whether funds left the platform.

**Anonymized User Context** *(required)* - A hashed/tokenized user identifier, ban status, KYC verification status. Critically, **Sigma does not require or want real PII** - it works with anonymized signals, which is why this section uses a `uniqueId` (your internal hashed identifier) rather than a name or account number.

**Optional Metadata** - Account flags (dormant, internal, staff, cheque), transaction context (sender/receiver account numbers, balance before, narration, session token), user profile (account type, age, city, country), device fingerprint (device ID, OS, manufacturer), GPS coordinates, third-party counterparty data, declared account limits, inline screening names, and beneficiary information.

**How results are displayed:**

The result page shows a prominent action banner - green for *approved*, amber for *flagged*, red for *rejected* - with a plain-English explanation of what each outcome means and what action to take. Below the banner, a risk score scale (0–33 low, 34–66 medium, 67–100 high) puts the numeric score in context. The full submitted transaction data and the complete Sigma response (transaction ID, rule result, triggering rule name and ID, reason code and message, inline screening results) are displayed in organized field cards.

---

### Task 2: PEP & Sanctions Screening

**Form page:** `/screening`

The screening page screens a named individual against two distinct databases - PEP and Sanctions - using separate API calls. Two buttons trigger each check independently, sharing the same form inputs.

Before reaching the form, the page presents information cards explaining:

- **What is a PEP?** - An individual in a prominent public function (heads of state, senior politicians, central bank governors). PEPs carry elevated financial crime risk due to their position and potential exposure to bribery or corruption.
- **What are Sanctions?** - Lists maintained by OFAC, UN, EU, UK HMT, and others, identifying parties prohibited from financial transactions. Transacting with a sanctioned entity is a serious regulatory violation.

The **Match Threshold** field (0.0–1.0) controls the minimum confidence required for a match to be returned. A lower threshold returns more potential matches (higher recall, lower precision); a higher threshold returns fewer but more certain matches.

**How results are displayed:**

An amber banner appears when matches are found, explicitly stating what a match means ("this does not automatically mean the subject is the same person - review the confidence score and profile data"). Each matched entity is shown in a card with:

- **Search Score** - how closely the name string matched
- **Confidence Score** - holistic match probability combining name, location, date of birth, and other corroborating signals
- Entity details: aliases, associated countries, known addresses, political positions (with dates held), political party affiliations, education, sanctions list memberships (in red badges), risk tags, and clickable source dataset links

---

### Task 3: Adverse Media Screening

**Form page:** `/adverse-media`

The adverse media page searches global news and media sources for negative coverage linked to a named subject.

An important notice is displayed prominently before the form: **this endpoint operates asynchronously**. After submitting a request, Sigma queues a deep media search. In a production integration, results are pushed to a configured webhook URL once processing completes - which may take several minutes. The demo displays the initial submission acknowledgement and any data returned in the immediate response.

**How results are displayed:**

A status banner reflects whether findings were returned and the overall risk category (High / Medium / Low), each with a plain-English explanation of what that category implies and what action is appropriate. The full Sigma response fields (request ID, processing status, query, business profile, timestamps, sources) are displayed with descriptions. A persistent notice explains the async delivery model and quotes the current job status.

---

## Design Decisions & Going Above and Beyond

### 1. A Complete SDK, Not Inline HTTP Calls

The `pkg/sigma` package is a self-contained SDK with typed models, enums, a clean interface, structured error handling, and full documentation comments. This goes meaningfully beyond making `http.Post` calls from handler code.

### 2. Interface-Driven Design for Testability

The `SigmaClient` interface means the application can be tested or demoed without touching the live API. The mock client is sophisticated enough to exercise every part of the UI, including multi-entity PEP results with positions, aliases, and sanctions.

### 3. Optional Fields as Pointer Types

The Sigma transaction payload has many optional fields. Rather than sending empty strings or zero values that the API would have to ignore, the service layer maps each optional form input through helper functions (`ptrStr`, `ptrFloat`, `ptrBoolForm`, `ptrTime`) that return `nil` when the value is empty - so the JSON encoder's `omitempty` tag omits them entirely from the payload. This is the correct way to interact with an optional-field API.

### 4. Information Architecture on Result Pages

Every result page explains not just what the data is, but what it means and what to do with it. Risk scores have a visual scale. Confidence scores have inline explanations. Sanctions matches are highlighted in red. The action recommendation comes with a plain-English interpretation.

### 5. Full Responsiveness

The layout adapts from a 5-column desktop grid down to a single-column mobile stack throughout every form and result page. The sidebar collapses to an off-canvas drawer on small screens, with a hamburger menu in the mobile header.

### 6. Dual Base URL Routing

Sigma exposes compliance capabilities through two separate base URLs. Transaction monitoring uses `sigmaprod.sabipay.com`; AML capabilities (PEP, sanctions, adverse media) use `sigmaaml.sabipay.com`. The SDK client routes each method call to the correct base transparently, so the service layer never has to know about this distinction.

### 7. Graceful Shutdown

The server catches `SIGINT` and `SIGTERM` signals and triggers a 10-second graceful shutdown using Go's built-in `http.Server.Shutdown`, allowing in-flight requests to complete before the process exits.

---

## Mock Mode vs Live API

The application ships in mock mode by default (`USE_MOCK=true`). This behaviour is intentional:

| Mode             | Behaviour                                                                                                                                                          |
| :--------------- | :----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `USE_MOCK=true`  | No Sigma API calls are made. All responses are generated by the `MockClient` with realistic randomised data. Useful for demoing the UI without network dependency. |
| `USE_MOCK=false` | All API calls are made to the live Sigma endpoints using the configured credentials. Requires `SIGMA_API_KEY` and `SIGMA_API_SECRET` to be set.                    |

The mock/live switch happens entirely inside the DI container (`internal/di/container.go`). No handler or service code branches on this setting - they all receive a `SigmaClient` and never know which concrete implementation they're using.

---

*Built by George Uche-Umeh - Pastel FDE Assessment, March 2026*
