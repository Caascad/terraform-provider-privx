package privx

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	errorSourceCreate = "error creating PrivX Source (%s): %s"
	errorSourceUpdate = "error updating PrivX Source (%s): %s"
	errorSourceDelete = "error deleting PrivX Source (%s): %s"
	errorSourceRead   = "error reading PrivX Source (%v): %s"
)

type Oidc_connection struct {
	oidc_enabled             bool
	oidc_issuer              string
	oidc_button_title        string
	oidc_client_id           string
	oidc_client_secret       string
	oidc_tags_attribute_name string
}

func resourcePrivXSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivxSourceCreate,
		ReadContext:   resourcePrivxSourceRead,
		UpdateContext: resourcePrivxSourceUpdate,
		DeleteContext: resourcePrivxSourceDelete,
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
			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"username_pattern": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"oidc_client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oidc_client_secret": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"oidc_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"oidc_issuer": {
				Type:     schema.TypeString,
				Required: true,
			},
			"oidc_tags_attribute_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oidc_button_title": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourcePrivxSourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var new_Source_id string

	var Source = rolestore.Source{
		Name:            d.Get("name").(string),
		Comment:         d.Get("comment").(string),
		Enabled:         d.Get("enabled").(bool),
		TTL:             d.Get("ttl").(int),
		UsernamePattern: flattenSimpleSlice(d.Get("username_pattern").([]interface{})),
		Tags:            flattenSimpleSlice(d.Get("tags").([]interface{})),
		Connection:      create_oidc_connection(d),
	}

	new_Source_id, err := createRoleClient(ctx, meta.(privx_API_client_connector).Connector).CreateSource(Source)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceCreate, new_Source_id, err))
	}

	d.SetId(new_Source_id)

	return resourcePrivxSourceRead(ctx, d, meta)
}

func resourcePrivxSourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/*Get Source List*/
	Sources_list, err := createRoleClient(ctx, meta.(privx_API_client_connector).Connector).Sources()

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	ind := findSourceIndex(Sources_list, d.Get("id").(string))
	if ind == -1 {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	Source := Sources_list[ind]

	if err := d.Set("name", Source.Name); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}

	if err := d.Set("comment", Source.Comment); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}

	if err := d.Set("ttl", Source.TTL); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("enabled", Source.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("username_pattern", Source.UsernamePattern); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("tags", Source.Tags); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}

	if err := d.Set("oidc_client_id", Source.Connection.OIDCClientID); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("oidc_enabled", Source.Connection.OIDCEnabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("oidc_issuer", Source.Connection.OIDCIssuer); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("oidc_tags_attribute_name", Source.Connection.OIDCTagsAttributeName); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("oidc_button_title", Source.Connection.OIDCButtonTitle); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}

	return nil
}

func resourcePrivxSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("name") || d.HasChange("comment") || d.HasChange("enabled") || d.HasChange("tags") || d.HasChange("username_pattern") || d.HasChange("ttl") {
		var Source = rolestore.Source{
			ID:              d.Get("id").(string),
			Name:            d.Get("name").(string),
			Comment:         d.Get("comment").(string),
			Enabled:         d.Get("enabled").(bool),
			TTL:             d.Get("ttl").(int),
			UsernamePattern: flattenSimpleSlice(d.Get("username_pattern").([]interface{})),
			Tags:            flattenSimpleSlice(d.Get("tags").([]interface{})),
			Connection:      create_oidc_connection(d),
		}
		err := createRoleClient(ctx, meta.(privx_API_client_connector).Connector).UpdateSource(d.Get("id").(string), &Source)
		if err != nil {
			return diag.FromErr(fmt.Errorf(errorSourceUpdate, d.Get("id").(string), err))
		}
	}

	return nil
}

func resourcePrivxSourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := createRoleClient(ctx, meta.(privx_API_client_connector).Connector).DeleteSource(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceDelete, d.Id(), err))
	}

	d.SetId("")

	return nil
}

func findSourceIndex(mySlice []rolestore.Source, id string) int {
	for i, item := range mySlice {
		if item.ID == id {
			return i
		}
	}
	return -1
}

func create_oidc_connection(d *schema.ResourceData) rolestore.Connection {
	connection := rolestore.Connection{
		Type:                  "OIDC",
		OIDCEnabled:           d.Get("oidc_enabled").(bool),
		OIDCIssuer:            d.Get("oidc_issuer").(string),
		OIDCClientID:          d.Get("oidc_client_id").(string),
		OIDCButtonTitle:       d.Get("oidc_button_title").(string),
		OIDCClientSecret:      d.Get("oidc_client_secret").(string),
		OIDCTagsAttributeName: d.Get("oidc_tags_attribute_name").(string),
	}
	return connection
}
