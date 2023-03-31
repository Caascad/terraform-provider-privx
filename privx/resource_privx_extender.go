package privx

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	errorExtenderCreate = "error creating PrivX Extender (%s): %s"
	errorExtenderUpdate = "error updating PrivX Extender (%s): %s"
	errorExtenderDelete = "error deleting PrivX Extender (%s): %s"
	errorExtenderRead   = "error reading PrivX Extender (%s): %s"
)

func resourcePrivXExtender() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivxExtenderCreate,
		ReadContext:   resourcePrivxExtenderRead,
		UpdateContext: resourcePrivxExtenderUpdate,
		DeleteContext: resourcePrivxExtenderDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"registered": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"access_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"web_proxy_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"web_proxy_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"routing_prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id": {
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
			"extender_address": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"subnets": {
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

func resourcePrivxExtenderCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var new_extender_id string

	var clientType userstore.ClientType = userstore.ClientType(d.Get("type").(string))

	var trusted_client = userstore.TrustedClient{
		Enabled:         d.Get("enabled").(bool),
		Name:            d.Get("name").(string),
		Secret:          d.Get("secret").(string),
		Type:            clientType,
		WebProxyAddress: d.Get("web_proxy_address").(string),
		Permissions:     flattenSimpleSlice(d.Get("permissions").([]interface{})),
		ExtenderAddress: flattenSimpleSlice(d.Get("extender_address").([]interface{})),
	}

	new_extender_id, err := createUserStoreClient(ctx, meta.(privx_API_client_connector).Connector).CreateTrustedClient(trusted_client)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorExtenderCreate, new_extender_id, err))
	}

	d.SetId(new_extender_id)

	return resourcePrivxExtenderRead(ctx, d, meta) //Extender API read gives less attributes than needed for extender creation
}

func resourcePrivxExtenderRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/*Get Extender List*/
	extenders_list, err := createUserStoreClient(ctx, meta.(privx_API_client_connector).Connector).TrustedClients()

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorExtenderRead, d.Id(), err))
	}
	extender := extenders_list[findExtenderIndex(extenders_list, d.Get("id").(string))]

	if err := d.Set("secret", extender.Secret); err != nil {
		return diag.FromErr(fmt.Errorf(errorExtenderRead, d.Id(), err))
	}

	if err := d.Set("name", extender.Name); err != nil {
		return diag.FromErr(fmt.Errorf(errorExtenderRead, d.Id(), err))
	}

	if err := d.Set("enabled", extender.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorExtenderRead, d.Id(), err))
	}

	if err := d.Set("registered", extender.Registered); err != nil {
		return diag.FromErr(fmt.Errorf(errorExtenderRead, d.Id(), err))
	}

	if err := d.Set("type", extender.Type); err != nil {
		return diag.FromErr(fmt.Errorf(errorExtenderRead, d.Id(), err))
	}

	if err := d.Set("extender_address", extender.ExtenderAddress); err != nil {
		return diag.FromErr(fmt.Errorf(errorExtenderRead, d.Id(), err))
	}

	if err := d.Set("permissions", extender.Permissions); err != nil {
		return diag.FromErr(fmt.Errorf(errorExtenderRead, d.Id(), err))
	}

	if err := d.Set("web_proxy_address", extender.WebProxyAddress); err != nil {
		return diag.FromErr(fmt.Errorf(errorExtenderRead, d.Id(), err))
	}
	return nil
}

func resourcePrivxExtenderUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	if d.HasChange("name") || d.HasChange("extender_address") || d.HasChange("permissions") { //Routing_prefix is not handled by SDK.
		var clientType userstore.ClientType = userstore.ClientType(d.Get("type").(string))
		var trusted_client = userstore.TrustedClient{
			Enabled:         d.Get("enabled").(bool),
			Name:            d.Get("name").(string),
			Secret:          d.Get("secret").(string),
			Type:            clientType,
			WebProxyAddress: d.Get("web_proxy_address").(string),
			Permissions:     flattenSimpleSlice(d.Get("permissions").([]interface{})),
			ExtenderAddress: flattenSimpleSlice(d.Get("extender_address").([]interface{})),
		}
		err := createUserStoreClient(ctx, meta.(privx_API_client_connector).Connector).UpdateTrustedClient(d.Get("id").(string), &trusted_client)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorExtenderUpdate, d.Get("id").(string), err))
		}
	}

	return nil
}

func resourcePrivxExtenderDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := createUserStoreClient(ctx, meta.(privx_API_client_connector).Connector).DeleteTrustedClient(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorExtenderDelete, d.Id(), err))
	}

	d.SetId("")

	return nil
}

func findExtenderIndex(mySlice []userstore.TrustedClient, id string) int {
	for i, item := range mySlice {
		if item.ID == id {
			return i
		}
	}
	return -1
}
