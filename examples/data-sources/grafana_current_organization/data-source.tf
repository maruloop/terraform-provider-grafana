provider "grafanasa" {}

data "grafana_current_organization" "test_current_org" {
  provider = grafanasa
}
