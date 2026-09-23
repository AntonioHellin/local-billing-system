# 💼 Local Billing System (`local-billing-system`)

A lightweight, standalone desktop invoicing and billing application built with a high-performance Go backend, embedded React SPA frontend, local SQLite database, and automated PDF export capabilities.

---

## Project Overview

**Local Billing System** provides private, local-first invoicing management for freelancers and small businesses. The application embeds the compiled React single-page frontend directly into the Go executable binary, requiring no external web server or cloud dependencies to run.

### Repository Naming Analysis
- **Recommended Repository Name**: `local-billing-system`
- **Naming Formula**: **Formula A** (`[domain/product]-[core-function]`)
- **Rationale**: Replaces the Spanish identifier `facturacion-local` with a clean kebab-case name specifying domain (`local-billing`) and core function (`system`).

---

## Features

- **Local-First SQLite Storage**: Automatic schema generation and migrations via GORM with pure-Go SQLite driver.
- **Embedded Single-Page Application (SPA)**: Serves React + TailwindCSS interface compiled directly into the binary via `go:embed`.
- **Client & Invoicing API**: RESTful endpoints for creating, querying, and managing clients and invoices.
- **Automated PDF Generation**: High-precision PDF invoice creation with tax calculations, itemized line items, and issuer branding via `gofpdf`.
- **Platform-Independent Exports**: Dynamically discovers safe user home directories across Windows, macOS, and Linux.

---

## Prerequisites

- **Go**: `>= 1.22`
- **Node.js**: `>= 18.0.0`
- **Package Manager**: `npm`

---

## Installation and Build

### 1. Build the Frontend
```bash
cd frontend
npm install
npm run build
cd ..
```

### 2. Compile the Go Application
```bash
go build -o local-billing.exe ./cmd/api
```

### 3. Run the Application
```bash
./local-billing.exe
```
The application will launch and listen at `http://localhost:8080`.

---

## Configuration & Environment Variables

Create a `.env` file based on `.env.example` to customize runtime settings:

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PORT` | Local HTTP server port | `8080` |
| `DB_PATH` | Path to SQLite database file | `facturacion.db` |
| `EXPORT_PATH` | Directory where PDF invoices are saved | `~/BillingExports` |

---

## Defensive Security Architecture

- **Parameterized Queries**: All database lookups and queries utilize GORM parameterized statements to eliminate SQL injection risks.
- **Binary & Database Exclusion**: SQLite database files (`*.db`) and pre-built binaries (`*.exe`) are strictly excluded from git tracking to prevent data leakage and repository bloat.
- **Cross-Platform Path Resolution**: Avoids hardcoded Windows paths (`C:\...`), resolving user directories safely with permission checks.

---

## License

Proprietary. All rights reserved. Not licensed for redistribution, public sublicensing, or resale.
