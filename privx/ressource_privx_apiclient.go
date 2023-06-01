package privx

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	errorApiClientCreate = "error creating PrivX ApiClient (%s): %s"
	errorApiClientUpdate = "error updating PrivX ApiClient (%s): %s"
	errorApiClientDelete = "error deleting PrivX ApiClient (%s): %s"
	errorApiClientRead   = "error reading PrivX ApiClient (%v): %s"
)

type createdData struct {
	ID string `json:"id"`
}

func resourcePrivXApiClient() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivxApiClientCreate,
		ReadContext:   resourcePrivxApiClientRead,
		UpdateContext: resourcePrivxApiClientUpdate,
		DeleteContext: resourcePrivxApiClientDelete,
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
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"authclientid": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"authclientsecret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created": {
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

func resourcePrivxApiClientCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var userStoreAPIClient = userstore.APIClient{
		Name:  d.Get("name").(string),
		Roles: createSecretListRoleRef(flattenSimpleSlice(d.Get("roles").([]interface{}))),
	}
	var CreatedData createdData
	curl := meta.(privx_API_client_connector).Connector.URL(fmt.Sprintf("/local-user-store/api/v1/api-clients"))
	_, err := curl.Post(userStoreAPIClient, &CreatedData)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorApiClientCreate, d.Get("name").(string), err))
	}
	d.SetId(CreatedData.ID)
	return resourcePrivxApiClientRead(ctx, d, meta)
}

func resourcePrivxApiClientRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var userStoreAPIClient userstore.APIClient
	curl := meta.(privx_API_client_connector).Connector.URL(fmt.Sprintf("/local-user-store/api/v1/api-clients/%s", d.Get("id").(string)))
	_, err := curl.Get(&userStoreAPIClient)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorApiClientRead, d.Id(), err))
	}

	if err := d.Set("name", userStoreAPIClient.Name); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("id", userStoreAPIClient.ID); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("name", userStoreAPIClient.Name); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("secret", userStoreAPIClient.Secret); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("authclientid", userStoreAPIClient.AuthClientID); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("authclientsecret", userStoreAPIClient.AuthClientSecret); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("roles", flattenSimpleRoleList(userStoreAPIClient.Roles)); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("created", userStoreAPIClient.Created); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}
	if err := d.Set("author", userStoreAPIClient.Author); err != nil {
		return diag.FromErr(fmt.Errorf(errorSourceRead, d.Id(), err))
	}

	return nil
}

func resourcePrivxApiClientUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var userStoreAPIClient = userstore.APIClient{
		ID:               d.Get("id").(string),
		Name:             d.Get("name").(string),
		Secret:           d.Get("secret").(string),
		AuthClientID:     d.Get("authclientid").(string),
		AuthClientSecret: d.Get("authclientsecret").(string),
		Roles:            createSecretListRoleRef(flattenSimpleSlice(d.Get("roles").([]interface{}))),
		Created:          d.Get("created").(string),
		Author:           d.Get("author").(string),
	}

	curl := meta.(privx_API_client_connector).Connector.URL(fmt.Sprintf("/local-user-store/api/v1/api-clients/%s", d.Get("id").(string)))
	_, err := curl.Put(userStoreAPIClient)
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorApiClientUpdate, d.Id(), err))
	}
	return resourcePrivxApiClientRead(ctx, d, meta)

}

func resourcePrivxApiClientDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	curl := meta.(privx_API_client_connector).Connector.URL(fmt.Sprintf("/local-user-store/api/v1/api-clients/%s", d.Get("id").(string)))
	_, err := curl.Delete()
	if err != nil {
		return diag.FromErr(fmt.Errorf(errorApiClientDelete, d.Id(), err))
	}

	d.SetId("")

	return nil
}

func flattenSimpleRoleList(roles_list []rolestore.RoleRef) []string {
	results := make([]string, 0, len(roles_list))
	for _, role := range roles_list {
		results = append(results, role.ID)
	}
	return results
}
