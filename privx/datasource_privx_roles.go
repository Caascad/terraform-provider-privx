package privx

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePrivxRoles() *schema.Resource {
	return &schema.Resource{
		ReadContext: listRoles,
		Schema: map[string]*schema.Schema{
			"roles": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permissions": {
							Type:     schema.TypeList,
							Computed: true,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func listRoles(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	roleslist, err := createRoleClient(ctx, meta.(privx_API_client_connector).Connector).Roles()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting API Client: %v", roleslist))
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting API Client: %s", err))
	}
	if err := d.Set("roles", flattenRoles(roleslist)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `roles`: %s", err))
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenRoles(roles_list []rolestore.Role) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(roles_list))

	for _, role := range roles_list {
		results = append(results, map[string]interface{}{
			"id":              role.ID,
			"access_group_id": role.AccessGroupID,
			"name":            role.Name,
			"permissions":     role.Permissions,
			"comment":         role.Comment,
		})
	}
	return results

}
