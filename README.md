# Sigma AI Compliance Dashboard

**Submission by George Uche-Umeh - Pastel FDE Assessment**

🔗 **Live Demo:** [pastel-fde-assessment.onrender.com](https://pastel-fde-assessment.onrender.com/)

---

## What This Is

A working web application that demonstrates all three Sigma AI compliance capabilities - built and submitted as part of the Pastel Forward Deployment Engineer assessment.

The app connects directly to Sigma's API and lets you run each capability interactively through a clean browser interface.

---

## What It Demonstrates

### 1 - Transaction Monitoring
Submit a financial transaction and get an instant AI-powered risk verdict.

The dashboard shows:
- **The recommended action** - Approved, Flagged, or Blocked - prominently displayed with a plain-language explanation
- **A risk score** (0–100) with a visual scale so you know what the number means
- **The reason** - which fraud detection rule triggered and why
- **Inline PEP & sanctions screening** on the sender and receiver, automatically included in the transaction check
- A complete breakdown of every data field sent to Sigma

### 2 - PEP & Sanctions Screening
Screen any individual against global watchlists with a single search.

Two separate checks are available - PEP (Politically Exposed Persons) and Sanctions - each returning:
- Match confidence scores, explained in plain English
- Full entity profiles: aliases, countries, political positions held, sanctions list memberships, and source links
- A clear verdict with context on what the match means and what to do next

### 3 - Adverse Media Screening
Search global news sources for negative coverage linked to an individual or entity.

Returns a risk category (High / Medium / Low) with an explanation of what that level implies, the media sources where findings were detected, and important context about how this endpoint works in a real integration (asynchronous, webhook-delivered).

---

## How to Use the Demo

1. Visit [pastel-fde-assessment.onrender.com](https://pastel-fde-assessment.onrender.com/)
2. Select a product from the left sidebar
3. The forms come pre-filled with realistic demo data - click the action button to run the check
4. Results appear inline with full explanations for every field

> The live deployment runs against the **Sigma production API** using the assessment credentials. The sidebar shows the current connection status.

---

## Tech Stack

Built with **Go** on the backend, **HTMX** for real-time form interactions, and **Tailwind CSS** for the UI - no frontend build step required. The application includes a custom Sigma SDK written from scratch.

→ [Technical deep-dive: architecture, SDK design, and engineering decisions](DEEP_DIVE.md)

---

*George Uche-Umeh - Pastel FDE Assessment, March 2026*
