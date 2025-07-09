package grafana_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/grafana/terraform-provider-grafana/v3/internal/testutils"
)
func TestAccDataSourceCurrentOrganization_basic(t *testing.T) {
	grafanaUrl := os.Getenv("GRAFANA_URL")

	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testutils.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCurrentOrganizationDatasourceConfig(grafanaUrl),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.grafana_current_organization.test_current_org", "name", "test-current-org"),
					resource.TestCheckResourceAttr("data.grafana_current_organization.test_current_org", "admins.0", "admin@localhost"),
					resource.TestCheckResourceAttr("data.grafana_current_organization.test_current_org", "editors.0", "editor-01@example.com"),
					resource.TestCheckResourceAttr("data.grafana_current_organization.test_current_org", "viewers.0", "viewer-01@example.com"),
				),
			},
		},
	})
}

func testCurrentOrganizationDatasourceConfig(grafanaUrl string) string {
	return fmt.Sprintf(`
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
	seconds_to_live    = 3600
}

provider "grafana" {
	alias = "test_current_org"

	url  = "%[1]s"
	auth = grafana_service_account_token.sa_token.key
}

data "grafana_current_organization" "test_current_org" {
	provider = grafana.test_current_org
	depends_on = [grafana_service_account_token.sa_token]
}
`, grafanaUrl)
}
