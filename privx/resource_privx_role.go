package privx

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	errorRoleCreate = "error creating PrivX Role (%s): %s"
	errorRoleUpdate = "error updating PrivX Role (%s): %s"
	errorRoleDelete = "error deleting PrivX Role (%s): %s"
	errorRoleRead   = "error reading PrivX Role (%s): %s"
)

func resourcePrivXRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivxRoleCreate,
		ReadContext:   resourcePrivxRoleRead,
		UpdateContext: resourcePrivxRoleUpdate,
		DeleteContext: resourcePrivxRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
				Required: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permissions": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"source_rules": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourcePrivxRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var new_role_id string

	var role = rolestore.Role{
		Name:        d.Get("name").(string),
		Permissions: flattenSimpleSlice(d.Get("permissions").([]interface{})),
		SourceRule:  d.Get("source_rules").(rolestore.SourceRule),
	}

	new_role_id, err := createRoleClient(ctx, meta.(privx_API_client_connector).Connector).CreateRole(role)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorRoleCreate, new_role_id, err))
	}

	d.SetId(new_role_id)

	return resourcePrivxRoleRead(ctx, d, meta) //Role API read gives less attributes than needed for extender creation
}

func resourcePrivxRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/*Get Role List*/
	roles_list, err := createRoleClient(ctx, meta.(privx_API_client_connector).Connector).Roles()

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorRoleRead, d.Id(), err))
	}
	role := roles_list[findRoleIndex(roles_list, d.Get("id").(string))]

	if err := d.Set("name", role.Name); err != nil {
		return diag.FromErr(fmt.Errorf(errorRoleRead, d.Id(), err))
	}

	return nil
}

func resourcePrivxRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange("name") || d.HasChange("source_rules") || d.HasChange("") { //Routing_prefix is not handled by SDK.
		var role = rolestore.Role{
			Name:       d.Get("name").(string),
			SourceRule: d.Get("source_rules").(rolestore.SourceRule),
			Comment:    d.Get("comment").(string),
		}
		err := createRoleClient(ctx, meta.(privx_API_client_connector).Connector).UpdateRole(d.Get("id").(string), &role)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorRoleUpdate, d.Get("id").(string), err))
		}
	}

	return nil
}

func resourcePrivxRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := createUserStoreClient(ctx, meta.(privx_API_client_connector).Connector).DeleteTrustedClient(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorRoleDelete, d.Id(), err))
	}

	d.SetId("")

	return nil
}

func findRoleIndex(mySlice []rolestore.Role, id string) int {
	for i, item := range mySlice {
		if item.ID == id {
			return i
		}
	}
	return -1
}
