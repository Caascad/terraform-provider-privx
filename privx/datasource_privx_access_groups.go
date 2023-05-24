package privx

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/authorizer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePrivxAccessGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: listAccessGroups,
		Schema: map[string]*schema.Schema{
			"access_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
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
							Optional: true,
						},
						"ca_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"updated_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"author": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func listAccessGroups(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	access_groups_list, err := authorizer.New(meta.(privx_API_client_connector).Connector).AccessGroups(0, 1000, "id", "ASC")

	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting API Client: %s", err))
	}
	if err := d.Set("access_groups", flattenAccessGroups(access_groups_list)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `name`: %s", err))
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenAccessGroups(access_groups_list []authorizer.AccessGroup) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(access_groups_list))

	for _, access_group := range access_groups_list {
		results = append(results, map[string]interface{}{
			"id":         access_group.ID,
			"name":       access_group.Name,
			"comment":    access_group.Comment,
			"ca_id":      access_group.CAID,
			"created":    access_group.Created,
			"updated":    access_group.Updated,
			"updated_by": access_group.UpdatedBy,
			"author":     access_group.Author,
		})
	}

	return results
}
