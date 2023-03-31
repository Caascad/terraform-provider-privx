package privx

import (
	"context"
	"fmt"
	"reflect"

	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePrivxExtender() *schema.Resource {
	return &schema.Resource{
		ReadContext: extenderClients,
		Schema: map[string]*schema.Schema{
			"extenders": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"secret": {
							Type:      schema.TypeString,
							Required:  true,
							Sensitive: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"registered": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"web_proxy_address": {
							Type:     schema.TypeString,
							Optional: true,
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
						"extender_address": {
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

func extenderClients(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	extenders_list, err := createUserStoreClient(ctx, meta.(privx_API_client_connector).Connector).TrustedClients()

	if err != nil {
		return diag.FromErr(fmt.Errorf("error setting API Client: %s", err))
	}

	if err := d.Set("extenders", flattenExtenders(extenders_list)); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `extenders`: %s", err))
	}

	d.SetId(resource.UniqueId())

	return nil
}

func flattenExtenders(extenders_list []userstore.TrustedClient) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(extenders_list))

	for _, extender := range extenders_list {
		results = append(results, map[string]interface{}{
			"id":                extender.ID,
			"secret":            extender.Secret,
			"name":              extender.Name,
			"registered":        extender.Registered,
			"enabled":           extender.Enabled,
			"type":              extender.Type,
			"extender_address":  extender.ExtenderAddress,
			"permissions":       extender.Permissions,
			"web_proxy_address": extender.WebProxyAddress,
		})
	}

	return results
}

func flattenComplexPrivXObj(list_objects []map[string]interface{}) []map[string]interface{} {
	results := make([]map[string]interface{}, 0, len(list_objects))
	for _, obj := range list_objects {
		obj_value := reflect.ValueOf(obj)
		var obj_struct map[string]interface{}
		for i := 0; i < obj_value.NumField(); i++ {
			obj_struct[obj_value.Field(i).Type().Name()] = obj_value.Field(i).Interface()
		}
		results = append(results, obj_struct)
	}
	return results
}
