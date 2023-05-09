package privx

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/SSHcom/privx-sdk-go/api/hoststore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	errorHostCreate = "error creating PrivX Host (%s): %s"
	errorHostUpdate = "error updating PrivX Host (%s): %s"
	errorHostDelete = "error deleting PrivX Host (%s): %s"
	errorHostRead   = "error reading PrivX Host (%s): %s"
)

func resourcePrivXHost() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivxHostCreate,
		ReadContext:   resourcePrivxHostRead,
		UpdateContext: resourcePrivxHostUpdate,
		DeleteContext: resourcePrivxHostDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePrivXHostImportState,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"source_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"contact_adress": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloud_provider": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cloud_provider_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"addresses": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func convertToAddressList(address_list []interface{}) []hoststore.Address {
	listAddress := make([]hoststore.Address, 0, len(address_list))
	for _, address := range address_list {
		listAddress = append(listAddress, hoststore.Address(address.(string)))
	}
	return listAddress
}

func flattenAdressList(address_list []hoststore.Address) []string {
	listAddress := make([]string, 0, len(address_list))
	for _, address := range address_list {
		listAddress = append(listAddress, string(address))
	}
	return listAddress
}

func resourcePrivxHostCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var new_host_id string

	var host = hoststore.Host{
		ID:                  d.Get("id").(string),
		AccessGroupID:       d.Get("access_group_id").(string),
		ExternalID:          d.Get("external_id").(string),
		InstanceID:          d.Get("instance_id").(string),
		SourceID:            d.Get("source_id").(string),
		Name:                d.Get("name").(string),
		ContactAdress:       d.Get("contact_adress").(string),
		CloudProvider:       d.Get("cloud_provider").(string),
		CloudProviderRegion: d.Get("cloud_provider_region").(string),
		Tags:                flattenSimpleSlice(d.Get("tags").([]interface{})),
		Addresses:           convertToAddressList(d.Get("addresses").([]interface{})),
		// DistinguishedName   string         `json:"distinguished_name,omitempty"`
		// Organization        string         `json:"organization,omitempty"`
		// OrganizationUnit    string         `json:"organizational_unit,omitempty"`
		// Zone                string         `json:"zone,omitempty"`
		// HostType            string         `json:"host_type,omitempty"`
		// HostClassification  string         `json:"host_classification,omitempty"`
		// Comment             string         `json:"comment,omitempty"`
		// Disabled            string         `json:"disabled,omitempty"`
		// Deployable          bool           `json:"deployable,omitempty"`
		// Tofu                bool           `json:"tofu,omitempty"`
		// StandAlone          bool           `json:"stand_alone_host,omitempty"`
		// Audit               bool           `json:"audit_enabled,omitempty"`
		// Scope               []string       `json:"scope,omitempty"`
		// Services            []Service      `json:"services,omitempty"`
		// Principals          []Principal    `json:"principals,omitempty"`
		// PublicKeys          []SSHPublicKey `json:"ssh_host_public_keys,omitempty"`
		// Status              []Status       `json:"status,omitempty"`
	}

	new_host_id, err := createHostClient(ctx, meta.(privx_API_client_connector).Connector).CreateHost(host)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorHostCreate, new_host_id, err))
	}

	d.SetId(new_host_id)

	return resourcePrivxHostRead(ctx, d, meta)
}

func resourcePrivxHostRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	/*Get Extender List*/
	host, err := createHostClient(ctx, meta.(privx_API_client_connector).Connector).Host(d.Get("id").(string))

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorHostRead, d.Get("id").(string), err))
	}

	if err := d.Set("access_group_id", host.AccessGroupID); err != nil {
		return diag.FromErr(fmt.Errorf(errorHostRead, d.Get("id").(string), err))
	}
	if err := d.Set("external_id", host.ExternalID); err != nil {
		return diag.FromErr(fmt.Errorf(errorHostRead, d.Get("id").(string), err))
	}
	if err := d.Set("instance_id", host.InstanceID); err != nil {
		return diag.FromErr(fmt.Errorf(errorHostRead, d.Get("id").(string), err))
	}
	if err := d.Set("source_id", host.SourceID); err != nil {
		return diag.FromErr(fmt.Errorf(errorHostRead, d.Get("id").(string), err))
	}
	if err := d.Set("name", host.Name); err != nil {
		return diag.FromErr(fmt.Errorf(errorHostRead, d.Get("id").(string), err))
	}
	if err := d.Set("contact_adress", host.ContactAdress); err != nil {
		return diag.FromErr(fmt.Errorf(errorHostRead, d.Get("id").(string), err))
	}
	if err := d.Set("cloud_provider", host.CloudProvider); err != nil {
		return diag.FromErr(fmt.Errorf(errorHostRead, d.Get("id").(string), err))
	}
	if err := d.Set("cloud_provider_region", host.CloudProvider); err != nil {
		return diag.FromErr(fmt.Errorf(errorHostRead, d.Get("id").(string), err))
	}
	if err := d.Set("tags", host.Tags); err != nil {
		return diag.FromErr(fmt.Errorf(errorHostRead, d.Get("id").(string), err))
	}
	if err := d.Set("addresses", flattenAdressList(host.Addresses)); err != nil {
		return diag.FromErr(fmt.Errorf(errorHostRead, d.Get("id").(string), err))
	}
	return nil
}

func resourcePrivxHostUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var host = hoststore.Host{
		ID:                  d.Get("id").(string),
		AccessGroupID:       d.Get("access_group_id").(string),
		ExternalID:          d.Get("external_id").(string),
		InstanceID:          d.Get("instance_id").(string),
		SourceID:            d.Get("source_id").(string),
		Name:                d.Get("name").(string),
		ContactAdress:       d.Get("contact_adress").(string),
		CloudProvider:       d.Get("cloud_provider").(string),
		CloudProviderRegion: d.Get("cloud_provider_region").(string),
		Tags:                flattenSimpleSlice(d.Get("tags").([]interface{})),
		Addresses:           convertToAddressList(d.Get("addresses").([]interface{})),
		// DistinguishedName   string         `json:"distinguished_name,omitempty"`
		// Organization        string         `json:"organization,omitempty"`
		// OrganizationUnit    string         `json:"organizational_unit,omitempty"`
		// Zone                string         `json:"zone,omitempty"`
		// HostType            string         `json:"host_type,omitempty"`
		// HostClassification  string         `json:"host_classification,omitempty"`
		// Comment             string         `json:"comment,omitempty"`
		// Disabled            string         `json:"disabled,omitempty"`
		// Deployable          bool           `json:"deployable,omitempty"`
		// Tofu                bool           `json:"tofu,omitempty"`
		// StandAlone          bool           `json:"stand_alone_host,omitempty"`
		// Audit               bool           `json:"audit_enabled,omitempty"`
		// Scope               []string       `json:"scope,omitempty"`
		// Services            []Service      `json:"services,omitempty"`
		// Principals          []Principal    `json:"principals,omitempty"`
		// PublicKeys          []SSHPublicKey `json:"ssh_host_public_keys,omitempty"`
		// Status              []Status       `json:"status,omitempty"`
	}

	err := createHostClient(ctx, meta.(privx_API_client_connector).Connector).UpdateHost(host.ID, &host)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorHostUpdate, host.ID, err))
	}
	return nil
}

func resourcePrivxHostDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := createHostClient(ctx, meta.(privx_API_client_connector).Connector).DeleteHost(d.Get("id").(string))

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorHostUpdate, d.Get("id").(string), err))
	}
	d.SetId("")
	return nil
}

func resourcePrivXHostImportState(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	conn := createHostClient(ctx, meta.(privx_API_client_connector).Connector)

	parts := strings.SplitN(d.Id(), "-", -1)
	if len(parts) != 5 {
		return nil, errors.New("import format error: invalid privx host ID")
	}

	host, err := conn.Host(d.Id())
	if err != nil {
		return nil, fmt.Errorf("couldn't import host %s, %v", d.Id(), err)
	}

	if err := d.Set("access_group_id", host.AccessGroupID); err != nil {
		return nil, fmt.Errorf(errorHostRead, d.Get("id").(string), err)
	}
	if err := d.Set("external_id", host.ExternalID); err != nil {
		return nil, fmt.Errorf(errorHostRead, d.Get("id").(string), err)
	}
	if err := d.Set("instance_id", host.InstanceID); err != nil {
		return nil, fmt.Errorf(errorHostRead, d.Get("id").(string), err)
	}
	if err := d.Set("source_id", host.SourceID); err != nil {
		return nil, fmt.Errorf(errorHostRead, d.Get("id").(string), err)
	}
	if err := d.Set("name", host.Name); err != nil {
		return nil, fmt.Errorf(errorHostRead, d.Get("id").(string), err)
	}
	if err := d.Set("contact_adress", host.ContactAdress); err != nil {
		return nil, fmt.Errorf(errorHostRead, d.Get("id").(string), err)
	}
	if err := d.Set("cloud_provider", host.CloudProvider); err != nil {
		return nil, fmt.Errorf(errorHostRead, d.Get("id").(string), err)
	}
	if err := d.Set("cloud_provider_region", host.CloudProvider); err != nil {
		return nil, fmt.Errorf(errorHostRead, d.Get("id").(string), err)
	}
	if err := d.Set("tags", host.Tags); err != nil {
		return nil, fmt.Errorf(errorHostRead, d.Get("id").(string), err)
	}
	if err := d.Set("addresses", flattenAdressList(host.Addresses)); err != nil {
		return nil, fmt.Errorf(errorHostRead, d.Get("id").(string), err)
	}

	return []*schema.ResourceData{d}, nil
}
