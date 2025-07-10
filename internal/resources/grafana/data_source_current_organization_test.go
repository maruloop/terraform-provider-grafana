package grafana_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/grafana/terraform-provider-grafana/v3/internal/testutils"
)
func TestAccDataSourceCurrentOrganization_basic(t *testing.T) {
	testutils.CheckOSSTestsEnabled(t, ">=9.1.0")
	orgId := orgScopedTest(t)
	resource.Test(t, resource.TestCase{
		ProtoV5ProviderFactories: testutils.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testCurrentOrganizationDatasourceConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.grafana_current_organization.test_current_org", "org_id", fmt.Sprintf("%d", orgId)),
					resource.TestCheckResourceAttr("data.grafana_current_organization.test_current_org", "admins.0", "admin@localhost"),
					resource.TestCheckResourceAttr("data.grafana_current_organization.test_current_org", "editors.0", "editor-01@example.com"),
					resource.TestCheckResourceAttr("data.grafana_current_organization.test_current_org", "viewers.0", "viewer-01@example.com"),
				),
			},
		},
	})
}

func testCurrentOrganizationDatasourceConfig() string {
	return `
data "grafana_current_organization" "test_current_org" {
}
`
}
