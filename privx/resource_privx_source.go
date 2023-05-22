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
			"oidc_connection": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"oidc_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"oidc_issuer": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"oidc_button_title": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"oidc_client_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"oidc_client_secret": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"oidc_tags_attribute_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
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
		Connection:      create_oidc_connection(d.Get("oidc_connection").([]interface{})),
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
	Source := Sources_list[findSourceIndex(Sources_list, d.Get("id").(string))]

	if err := d.Set("name", Source.Name); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}

	if err := d.Set("comment", Source.Comment); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}

	if err := d.Set("enabled", Source.Enabled); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}

	if err := d.Set("oidc_connection", flatten_oidc_connection(Source.Connection)); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, Source.Connection, err))
	}

	return nil
}

func resourcePrivxSourceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("name") || d.HasChange("comment") || d.HasChange("enabled") || d.HasChange("tags") || d.HasChange("username_pattern") || d.HasChange("ttl") || d.HasChange("oidc_connection") {
		var Source = rolestore.Source{
			ID:              d.Get("id").(string),
			Name:            d.Get("name").(string),
			Comment:         d.Get("comment").(string),
			Enabled:         d.Get("enabled").(bool),
			TTL:             d.Get("ttl").(int),
			UsernamePattern: flattenSimpleSlice(d.Get("username_pattern").([]interface{})),
			Tags:            flattenSimpleSlice(d.Get("tags").([]interface{})),
			Connection:      create_oidc_connection(d.Get("oidc_connection").([]interface{})),
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

func create_oidc_connection(oidc_connection []interface{}) rolestore.Connection {
	results := rolestore.Connection{
		Type:                  "OIDC",
		OIDCEnabled:           oidc_connection[0].(map[string]interface{})["oidc_enabled"].(bool),
		OIDCIssuer:            oidc_connection[0].(map[string]interface{})["oidc_issuer"].(string),
		OIDCClientID:          oidc_connection[0].(map[string]interface{})["oidc_client_id"].(string),
		OIDCButtonTitle:       oidc_connection[0].(map[string]interface{})["oidc_button_title"].(string),
		OIDCClientSecret:      oidc_connection[0].(map[string]interface{})["oidc_client_secret"].(string),
		OIDCTagsAttributeName: oidc_connection[0].(map[string]interface{})["oidc_tags_attribute_name"].(string),
	}
	return results
}

func flatten_oidc_connection(connection rolestore.Connection) []Oidc_connection {
	var result []Oidc_connection
	flattened_connection := Oidc_connection{
		oidc_enabled:             connection.OIDCEnabled,
		oidc_button_title:        connection.OIDCButtonTitle,
		oidc_client_id:           connection.OIDCClientID,
		oidc_issuer:              connection.OIDCIssuer,
		oidc_client_secret:       connection.OIDCClientSecret,
		oidc_tags_attribute_name: connection.OIDCTagsAttributeName,
	}
	result = append(result, flattened_connection)

	return result
}
