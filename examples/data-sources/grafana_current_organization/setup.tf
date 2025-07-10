resource "grafana_organization" "test_current_org" {
	name = "test-current-org"
	editors = [
		"editor-01@example.com",
		"editor-02@example.com",
	]
	viewers = [
		"viewer-01@example.com",
		"viewer-02@example.com",
	]
}

resource "grafana_service_account" "test_current_org_sa" {
	name   = "sa-current-org-test"
	role   = "Admin"
	org_id = grafana_organization.test_current_org.org_id
}

resource "grafana_service_account_token" "sa_token" {
	name               = "sa-token"
	service_account_id = grafana_service_account.test_current_org_sa.id
}
