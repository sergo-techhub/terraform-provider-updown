# Terraform Provider for updown.io

[![CI](https://github.com/sergo-techhub/terraform-provider-updown/actions/workflows/ci.yml/badge.svg)](https://github.com/sergo-techhub/terraform-provider-updown/actions/workflows/ci.yml)
[![Release](https://github.com/sergo-techhub/terraform-provider-updown/actions/workflows/release.yml/badge.svg)](https://github.com/sergo-techhub/terraform-provider-updown/actions/workflows/release.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/sergo-techhub/terraform-provider-updown.svg)](https://pkg.go.dev/github.com/sergo-techhub/terraform-provider-updown)

A Terraform provider for [updown.io](https://updown.io) - a simple and affordable website monitoring service.

## Fork Notice

This project was forked from [mvisonneau/terraform-provider-updown](https://github.com/mvisonneau/terraform-provider-updown) by [SERGO GmbH](https://github.com/sergo-techhub).

We are actively modernizing and updating this provider to support the current [updown.io API](https://updown.io/api) implementation, including:

- Support for all check types (`http`, `https`, `icmp`, `tcp`, `tcps`)
- HTTP verb configuration (`GET`, `HEAD`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`)
- HTTP body for POST/PUT/PATCH requests
- Modern Go version (1.24+)
- Latest Terraform Plugin SDK (v2.38.1)

## Contributing

We'd love to see you contribute! Whether it's:

- Reporting bugs or suggesting features via [Issues](https://github.com/sergo-techhub/terraform-provider-updown/issues)
- Submitting [Pull Requests](https://github.com/sergo-techhub/terraform-provider-updown/pulls) with improvements or fixes
- Improving documentation

All contributions are welcome!

## Resources

| Type | Name | Description |
|------|------|-------------|
| **data** | `updown_nodes` | Returns the list of monitoring nodes IPv4 and IPv6 addresses |
| **resource** | `updown_check` | Creates and manages a check |
| **resource** | `updown_recipient` | Creates and manages a recipient |
| **resource** | `updown_status_page` | Creates and manages a status page |
| **resource** | `updown_webhook` | Creates a webhook _(DEPRECATED - use recipients instead)_ |

## Installation

### From Source (Development)

```bash
git clone https://github.com/sergo-techhub/terraform-provider-updown.git
cd terraform-provider-updown
go build -o terraform-provider-updown
```

Then configure Terraform to use the local provider:

```hcl
# ~/.terraformrc
provider_installation {
  dev_overrides {
    "sergo-techhub/updown" = "/path/to/terraform-provider-updown"
  }
  direct {}
}
```

## Usage

### Provider Configuration

```hcl
terraform {
  required_providers {
    updown = {
      source = "sergo-techhub/updown"
    }
  }
}

provider "updown" {
  # API key can also be set via UPDOWN_API_KEY environment variable
  api_key = "<YOUR_UPDOWN_API_KEY>"
}
```

### Basic HTTP/HTTPS Check

```hcl
resource "updown_check" "website" {
  url          = "https://example.com"
  alias        = "Example Website"
  period       = 60
  apdex_t      = 1.0
  enabled      = true
  published    = false
  string_match = "Welcome"

  disabled_locations = ["mia", "syd"]

  custom_headers = {
    "X-Custom-Header" = "value"
  }
}
```

### ICMP Ping Check

```hcl
resource "updown_check" "ping_server" {
  url   = "192.168.1.1"
  type  = "icmp"
  alias = "Server Ping"
}
```

### TCP Port Check

```hcl
resource "updown_check" "postgres" {
  url   = "tcp://db.example.com:5432"
  type  = "tcp"
  alias = "PostgreSQL Database"
}
```

### HTTPS with Custom HTTP Method

```hcl
resource "updown_check" "api_health" {
  url       = "https://api.example.com/health"
  alias     = "API Health Check"
  http_verb = "GET/HEAD"  # Use HEAD requests first, fallback to GET

  custom_headers = {
    "Accept" = "application/json"
  }
}
```

### Recipients and Alerts

```hcl
resource "updown_recipient" "email_alert" {
  type  = "email"
  value = "alerts@example.com"
}

resource "updown_recipient" "slack_alert" {
  type  = "slack_compatible"
  value = "https://hooks.slack.com/services/..."
}

resource "updown_check" "monitored_site" {
  url   = "https://example.com"
  alias = "Monitored Site"

  recipients = [
    updown_recipient.email_alert.id,
    updown_recipient.slack_alert.id,
  ]
}
```

### Get Monitoring Node IPs

```hcl
data "updown_nodes" "all" {}

output "monitoring_nodes_ipv4" {
  value = data.updown_nodes.all.ipv4
}

output "monitoring_nodes_ipv6" {
  value = data.updown_nodes.all.ipv6
}
```

### Status Pages

```hcl
# Public status page
resource "updown_status_page" "public" {
  name        = "My Services Status"
  description = "Public status page for all monitored services"
  visibility  = "public"

  checks = [
    updown_check.website.id,
    updown_check.api.id,
  ]
}

# Protected status page with access key
resource "updown_status_page" "protected" {
  name        = "Internal Services"
  description = "Protected status page for internal services"
  visibility  = "protected"
  access_key  = "my-secret-access-key"

  checks = [
    updown_check.internal_api.id,
  ]
}

# Private status page (only visible when logged in)
resource "updown_status_page" "private" {
  name       = "Private Infrastructure"
  visibility = "private"

  checks = [
    updown_check.database.id,
    updown_check.cache.id,
  ]
}
```

## Resource Reference

### updown_check

| Attribute | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `url` | string | Yes | - | The URL to monitor |
| `type` | string | No | _(inferred)_ | Check type: `http`, `https`, `icmp`, `tcp`, `tcps` |
| `alias` | string | No | - | Human-readable name |
| `period` | number | No | `60` | Check interval in seconds (15, 30, 60, 120, 300, 600, 1800, 3600) |
| `apdex_t` | number | No | `0.5` | APDEX threshold in seconds |
| `enabled` | bool | No | `true` | Whether the check is enabled |
| `published` | bool | No | `false` | Whether to show on public status page |
| `string_match` | string | No | - | String to search for in response |
| `mute_until` | string | No | - | Mute notifications until time, `recovery`, or `forever` |
| `http_verb` | string | No | `GET` | HTTP method for http/https checks: `GET`, `GET/HEAD`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS` |
| `http_body` | string | No | - | Request body for POST/PUT/PATCH |
| `disabled_locations` | set(string) | No | - | Locations to exclude from monitoring (max 8) |
| `recipients` | set(string) | No | - | Recipient IDs for alerts |
| `custom_headers` | map(string) | No | - | Custom HTTP headers |

### updown_recipient

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `type` | string | Yes | Recipient type: `email`, `webhook`, `slack_compatible` (Note: `sms` and `msteams` require web UI setup) |
| `value` | string | Yes | Email address, phone number, or webhook URL |

### updown_status_page

| Attribute | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `checks` | list(string) | Yes | - | List of check tokens to display (order is preserved) |
| `name` | string | No | - | Name of the status page |
| `description` | string | No | - | Description text (supports newlines and links) |
| `visibility` | string | No | `public` | Page visibility: `public`, `protected`, or `private` |
| `access_key` | string | No | _(auto-generated)_ | Access key for protected pages |
| `url` | string | Read-only | - | The URL of the status page |

## API Reference

For the complete updown.io API documentation, visit: https://updown.io/api

## License

This project is licensed under the Apache 2.0 License.
