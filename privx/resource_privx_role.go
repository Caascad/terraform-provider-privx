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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"permit_agent": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"access_group_id": {
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
		},
	}
}

func resourcePrivxRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var new_role_id string

	var role = rolestore.Role{
		Name:          d.Get("name").(string),
		Comment:       d.Get("comment").(string),
		PermitAgent:   d.Get("permit_agent").(bool),
		AccessGroupID: d.Get("access_group_id").(string),
		Permissions:   flattenSimpleSlice(d.Get("permissions").([]interface{})),
		SourceRule:    rolestore.SourceRuleNone(), // Creates an empty mapping, till we know if we need this.
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

	if err := d.Set("comment", role.Comment); err != nil {
		return diag.FromErr(fmt.Errorf(errorRoleRead, d.Id(), err))
	}

	if err := d.Set("permit_agent", role.PermitAgent); err != nil {
		return diag.FromErr(fmt.Errorf(errorRoleRead, d.Id(), err))
	}

	if err := d.Set("access_group_id", role.AccessGroupID); err != nil {
		return diag.FromErr(fmt.Errorf(errorRoleRead, d.Id(), err))
	}

	if err := d.Set("permissions", role.Permissions); err != nil {
		return diag.FromErr(fmt.Errorf(errorRoleRead, d.Id(), err))
	}

	return nil
}

func resourcePrivxRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange("name") || d.HasChange("comment") || d.HasChange("permissions") || d.HasChange("access_group_id") || d.HasChange("permissions") {
		var role = rolestore.Role{
			ID:            d.Get("id").(string),
			Name:          d.Get("name").(string),
			Comment:       d.Get("comment").(string),
			PermitAgent:   d.Get("permit_agent").(bool),
			AccessGroupID: d.Get("access_group_id").(string),
			Permissions:   flattenSimpleSlice(d.Get("permissions").([]interface{})),
			SourceRule:    rolestore.SourceRuleNone(),
		}
		err := createRoleClient(ctx, meta.(privx_API_client_connector).Connector).UpdateRole(d.Get("id").(string), &role)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorRoleUpdate, d.Get("id").(string), err))
		}
	}

	return nil
}

func resourcePrivxRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := createRoleClient(ctx, meta.(privx_API_client_connector).Connector).DeleteRole(d.Get("id").(string))
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

/* Here lies the definition of a source_rules of type single for potential later use.

"source_rules": {
	Type:     schema.TypeList,
	Required: true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source": {
				Type:     schema.TypeString,
				Required: true,
			},
			"search_string": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	},
},

*/
