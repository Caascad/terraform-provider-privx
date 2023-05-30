package privx

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/authorizer"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	errorAccessGroupCreate = "error creating PrivX AccessGroup (%s): %s"
	errorAccessGroupUpdate = "error updating PrivX AccessGroup (%s): %s"
	errorAccessGroupDelete = "error deleting PrivX AccessGroup (%s): %s"
	errorAccessGroupRead   = "error reading PrivX AccessGroup (%s): %s"
)

func resourcePrivXAccessGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivxAccessGroupCreate,
		ReadContext:   resourcePrivxAccessGroupRead,
		UpdateContext: resourcePrivxAccessGroupUpdate,
		DeleteContext: resourcePrivxAccessGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePrivXAccessGroupImportState,
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
	}
}

func resourcePrivxAccessGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var new_access_group_id string

	var access_group = authorizer.AccessGroup{
		Name:    d.Get("name").(string),
		Comment: d.Get("comment").(string),
	}

	new_access_group_id, err := createAuthorizerClient(ctx, meta.(privx_API_client_connector).Connector).CreateAccessGroup(&access_group)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorAccessGroupCreate, new_access_group_id, err))
	}

	d.SetId(new_access_group_id)

	return resourcePrivxAccessGroupRead(ctx, d, meta)
}

func resourcePrivxAccessGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/*Get AccessGroup List*/
	access_groups_list, err := createAuthorizerClient(ctx, meta.(privx_API_client_connector).Connector).AccessGroups(0, 1000, "id", "ASC")

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorAccessGroupRead, d.Id(), err))
	}

	ind := findAccessGroupIndex(access_groups_list, d.Get("id").(string))
	if ind == -1 {
		return diag.FromErr(fmt.Errorf(errorAccessGroupRead, d.Id(), fmt.Errorf("%v", nil)))
	}

	access_group := access_groups_list[findAccessGroupIndex(access_groups_list, d.Get("id").(string))]

	if err := d.Set("name", access_group.Name); err != nil {
		return diag.FromErr(fmt.Errorf(errorAccessGroupRead, d.Id(), err))
	}

	if err := d.Set("comment", access_group.Comment); err != nil {
		return diag.FromErr(fmt.Errorf(errorAccessGroupRead, d.Id(), err))
	}
	if err := d.Set("ca_id", access_group.CAID); err != nil {
		return diag.FromErr(fmt.Errorf(errorAccessGroupRead, d.Id(), err))
	}
	if err := d.Set("created", access_group.Created); err != nil {
		return diag.FromErr(fmt.Errorf(errorAccessGroupRead, d.Id(), err))
	}
	if err := d.Set("updated", access_group.Updated); err != nil {
		return diag.FromErr(fmt.Errorf(errorAccessGroupRead, d.Id(), err))
	}
	if err := d.Set("updated_by", access_group.UpdatedBy); err != nil {
		return diag.FromErr(fmt.Errorf(errorAccessGroupRead, d.Id(), err))
	}
	if err := d.Set("author", access_group.Author); err != nil {
		return diag.FromErr(fmt.Errorf(errorAccessGroupRead, d.Id(), err))
	}

	return nil
}

func resourcePrivxAccessGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange("name") || d.HasChange("comment") {
		var access_group = authorizer.AccessGroup{
			Name:      d.Get("name").(string),
			Comment:   d.Get("comment").(string),
			CAID:      d.Get("ca_id").(string),
			Author:    d.Get("author").(string),
			Created:   d.Get("created").(string),
			Updated:   d.Get("updated").(string),
			UpdatedBy: d.Get("updated_by").(string),
		}
		err := createAuthorizerClient(ctx, meta.(privx_API_client_connector).Connector).UpdateAccessGroup(d.Get("id").(string), &access_group)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorAccessGroupUpdate, d.Get("id").(string), err))
		}
	}

	return nil
}

func resourcePrivxAccessGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := createAuthorizerClient(ctx, meta.(privx_API_client_connector).Connector).DeleteAccessGroup(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorAccessGroupDelete, d.Id(), err))
	}

	d.SetId("")

	return nil
}

func resourcePrivXAccessGroupImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := createAuthorizerClient(ctx, meta.(privx_API_client_connector).Connector)

	parts := strings.SplitN(d.Id(), "-", -1)
	if len(parts) != 5 {
		return nil, errors.New("import format error: invalid privx access group ID")
	}

	access_group, err := conn.AccessGroup(d.Id())
	if err != nil {
		return nil, fmt.Errorf("couldn't import access group %s, %v", d.Id(), err)
	}

	if err := d.Set("name", access_group.Name); err != nil {
		return nil, fmt.Errorf(errorAccessGroupRead, d.Id(), err)
	}

	if err := d.Set("comment", access_group.Comment); err != nil {
		return nil, fmt.Errorf(errorAccessGroupRead, d.Id(), err)
	}
	if err := d.Set("ca_id", access_group.CAID); err != nil {
		return nil, fmt.Errorf(errorAccessGroupRead, d.Id(), err)
	}
	if err := d.Set("created", access_group.Created); err != nil {
		return nil, fmt.Errorf(errorAccessGroupRead, d.Id(), err)
	}
	if err := d.Set("updated", access_group.Updated); err != nil {
		return nil, fmt.Errorf(errorAccessGroupRead, d.Id(), err)
	}
	if err := d.Set("updated_by", access_group.UpdatedBy); err != nil {
		return nil, fmt.Errorf(errorAccessGroupRead, d.Id(), err)
	}
	if err := d.Set("author", access_group.Author); err != nil {
		return nil, fmt.Errorf(errorAccessGroupRead, d.Id(), err)
	}

	return []*schema.ResourceData{d}, nil
}

func findAccessGroupIndex(mySlice []authorizer.AccessGroup, id string) int {
	for i, item := range mySlice {
		if item.ID == id {
			return i
		}
	}
	return -1
}
