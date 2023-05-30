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
			"source_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"match": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"rules": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"source": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"search_string": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
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
		SourceRule:    CreateSourceRule(d.Get("source_rules").([]interface{})),
	}

	new_role_id, err := createRoleClient(ctx, meta.(privx_API_client_connector).Connector).CreateRole(role)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorRoleCreate, new_role_id, err))
	}

	d.SetId(new_role_id)

	return resourcePrivxRoleRead(ctx, d, meta)
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

	if d.HasChange("name") || d.HasChange("comment") || d.HasChange("permissions") || d.HasChange("access_group_id") || d.HasChange("permissions") || d.HasChange("source_rules") {
		var role = rolestore.Role{
			ID:            d.Get("id").(string),
			Name:          d.Get("name").(string),
			Comment:       d.Get("comment").(string),
			PermitAgent:   d.Get("permit_agent").(bool),
			AccessGroupID: d.Get("access_group_id").(string),
			Permissions:   flattenSimpleSlice(d.Get("permissions").([]interface{})),
			SourceRule:    CreateSourceRule(d.Get("source_rules").([]interface{})),
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

func CreateSourceRule(source_rule []interface{}) rolestore.SourceRule {
	if len(source_rule) < 1 {
		return rolestore.SourceRuleNone()
	}
	var obj_source_rule = rolestore.SourceRule{
		Type:  "GROUP", /*Force group as API seems to only generate arrays*/
		Match: source_rule[0].(map[string]interface{})["match"].(string),
		Rules: CreateNestedSourceRule(source_rule[0].(map[string]interface{})["rules"].([]interface{})),
	}
	return obj_source_rule
}

func CreateNestedSourceRule(source_rules_slice []interface{}) []rolestore.SourceRule {
	var nested_source_rules = make([]rolestore.SourceRule, len(source_rules_slice))
	for cpt, rule := range source_rules_slice {
		var nested_rule = rolestore.SourceRule{
			Type:    "RULE",
			Source:  rule.(map[string]interface{})["source"].(string),
			Pattern: rule.(map[string]interface{})["search_string"].(string),
		}
		nested_source_rules[cpt] = nested_rule
	}
	return nested_source_rules
}
