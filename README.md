# Sigma Compliance Dashboard

A production-grade web application demonstrating [Pastel's Sigma AI](https://sigma-docs.pastel.africa) compliance capabilities: **Transaction Monitoring**, **PEP & Sanctions Screening**, and **Adverse Media Screening**.

Built with **Go**, **HTMX**, and **Tailwind CSS**. Uses **dependency injection** via [`samber/do`](https://github.com/samber/do) for clean, testable architecture.

## Quick Start

### Prerequisites
- Go 1.21+ installed

### Run Locally
```bash
git clone https://github.com/UcGeorge/pastel-fde-assessment.git
cd pastel-fde-assessment
go run .
```

Open [http://localhost:8080](http://localhost:8080) in your browser.

### Environment Variables (Optional)

| Variable | Default | Description |
|---|---|---|
| `SIGMA_API_KEY` | *(provided)* | Sigma API key |
| `SIGMA_API_SECRET` | *(provided)* | Sigma API secret |
| `SIGMA_BASE_URL` | `https://sigmaprod.sabipay.com/` | Transaction monitoring API |
| `SIGMA_AML_BASE_URL` | `https://sigmaaml.sabipay.com/` | AML/Screening API |
| `PORT` | `8080` | HTTP server port |

## Features

### 1. Transaction Monitoring
Submit a transaction to Sigma's AI engine for real-time fraud analysis. The app displays:
- Transaction details submitted (reference, amount, currency, sender, receiver, channel, type)
- Full Sigma response with risk score and recommended action (approve/block)
- Risk severity badge color-coded (green/amber/red)

### 2. PEP & Sanctions Screening
Screen individuals against global Politically Exposed Persons lists and international sanctions databases (OFAC, UN, EU). The app displays:
- Individual details submitted (name, threshold, country)
- Match count with confidence scores
- Entity details: aliases, positions, sanctions lists, data sources

### 3. Adverse Media Screening
Search news and media sources for negative coverage. The app displays:
- Query parameters submitted
- Request status and tracking ID
- Full API response with webhook-based processing status

## Architecture

```
├── main.go              # Entry point: DI setup, routing, graceful shutdown
├── internal/
│   ├── config/          # Environment-based configuration
│   └── di/              # samber/do DI container wiring
├── pkg/sigma/           # Sigma API SDK (client, models, errors)
├── services/            # Business logic layer
├── handlers/            # HTTP handlers (pages + HTMX endpoints)
└── templates/           # Go html/template with Tailwind + HTMX
```

**Key Design Decisions:**
- **Dependency Injection**: `samber/do/v2` for clean service wiring and testability
- **SDK Pattern**: Reusable `pkg/sigma` client with dual base URL support (transaction monitoring vs AML)
- **HTMX**: Dynamic form submissions without JavaScript frameworks; server-rendered HTML partials
- **Error Handling**: Custom `SigmaAPIError` type with sentinel errors for `errors.Is()` support
- **Graceful Shutdown**: Signal handling for clean server termination

## Notes

- The API returns `200 OK` with body `{"message":"Access Denied"}` for some endpoints — this is handled gracefully in the UI
- The Sigma documentation has a typo: `sacntion` instead of `sanction` in the sanctions endpoint path — the SDK uses the documented path as-is
- Adverse media screening uses a webhook-based pattern; the initial response shows a `pending` status with a request ID

## Author
George Uche-Umeh — [engineering@pastel.africa](mailto:engineering@pastel.africa)
