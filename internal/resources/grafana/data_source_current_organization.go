package grafana

import (
	"fmt"
	"context"
	"strconv"
	"strings"

	"github.com/grafana/terraform-provider-grafana/v3/internal/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)
func datasourceCurrentOrganization() *common.DataSource {
	s := &schema.Resource{
		Description: "Data source to fetch the current organization info of the scoped Grafana provider.",
		ReadContext: dataSourceCurrentOrganizationRead,
		Schema: map[string]*schema.Schema{
			"org_id": orgIDAttribute(),
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the current organization.",
			},
			"admins": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of admin user logins in the current organization.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

	return common.NewLegacySDKDataSource(common.CategoryGrafanaOSS, "grafana_current_organization", s)
}

func dataSourceCurrentOrganizationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, _ := OAPIClientFromNewOrgResource(meta, d)

	resp, err := client.Org.GetCurrentOrg(nil)
	if err != nil {
		return diag.Errorf("failed to get current organization: %s", err)
	}
	fmt.Printf("🔍 Auth'd Org ID: %d, Name: %s\n", resp.Payload.ID, resp.Payload.Name)
	currentOrg := resp.GetPayload()

	orgUsers, err := client.Org.GetOrgUsersForCurrentOrg(nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var admins []string
	for _, user := range orgUsers.Payload {
		if strings.ToLower(user.Role) == "admin" {
			admins = append(admins, user.Email)
		}
	}

	if err := d.Set("admins", admins); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", currentOrg.Name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(currentOrg.ID, 10))
	return nil
}
