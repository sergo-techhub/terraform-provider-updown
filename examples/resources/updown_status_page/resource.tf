resource "updown_status_page" "public" {
  name        = "My Services Status"
  description = "Public status page for all monitored services"
  visibility  = "public"

  checks = [
    updown_check.website.id,
    updown_check.api.id,
  ]
}

resource "updown_status_page" "protected" {
  name        = "Internal Services"
  description = "Protected status page for internal services"
  visibility  = "protected"
  access_key  = "my-secret-access-key"

  checks = [
    updown_check.internal_api.id,
  ]
}

resource "updown_status_page" "private" {
  name       = "Private Infrastructure"
  visibility = "private"

  checks = [
    updown_check.database.id,
    updown_check.cache.id,
  ]
}
