package privx

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	errorJsonSecretCreate = "error creating PrivX JsonSecret (%s): %s"
	errorJsonSecretUpdate = "error updating PrivX JsonSecret (%s): %s"
	errorJsonSecretDelete = "error deleting PrivX JsonSecret (%s): %s"
	errorJsonSecretRead   = "error reading PrivX JsonSecret (%s): %s"
)

type SecretJsonData struct {
	Data string `json:"-"`
}

// tVaultReq t vault request definition
type tVaultReqJson struct {
	Name       string              `json:"name"`
	Data       SecretJsonData      `json:"data"`
	AllowRead  []rolestore.RoleRef `json:"read_roles,omitempty"`
	AllowWrite []rolestore.RoleRef `json:"write_roles,omitempty"`
	OwnerID    string              `json:"owner_id,omitempty"`
}

func resourcePrivxJsonSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivxJsonSecretCreate,
		ReadContext:   resourcePrivxJsonSecretRead,
		UpdateContext: resourcePrivxJsonSecretUpdate,
		DeleteContext: resourcePrivxJsonSecretDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"read_roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"write_roles": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourcePrivxJsonSecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var secretData = SecretJsonData{
		Data: d.Get("data").(string),
	}

	curl := meta.(privx_API_client_connector).Connector.URL(fmt.Sprintf("/vault/api/v1/secrets"))

	vaultReq := tVaultReqJson{
		Data:       secretData,
		Name:       d.Get("name").(string),
		AllowRead:  createJsonSecretListRoleRef(flattenSimpleSlice(d.Get("read_roles").([]interface{}))),
		AllowWrite: createJsonSecretListRoleRef(flattenSimpleSlice(d.Get("write_roles").([]interface{}))),
		OwnerID:    "",
	}
	_, err := curl.Post(vaultReq)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorJsonSecretCreate, d.Get("name").(string), err))
	}

	d.SetId(d.Get("name").(string))

	return nil
}

func check_value_JsonSecretData(value interface{}) SecretJsonData {
	valuem, _ := value.(map[string]interface{})
	if valuem["data"] != nil {
		secretData := SecretJsonData{
			Data: valuem["data"].(string),
		}
		return secretData
	} else {
		secretData := SecretJsonData{
			Data: "",
		}
		return secretData
	}
}

func resourcePrivxJsonSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	/*Get JsonSecret List*/
	curl := meta.(privx_API_client_connector).Connector.URL(fmt.Sprintf("/vault/api/v1/secrets/%s", d.Get("name").(string)))
	var secret interface{}
	_, err := curl.Get(&secret)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorJsonSecretRead, d.Get("name"), err))
	}

	secretm, _ := secret.(map[string]interface{})
	secret = tVaultReqJson{
		Name:       secretm["name"].(string),
		OwnerID:    check_value_string(secretm["owner_id"]),
		Data:       check_value_JsonSecretData(secretm["data"]),
		AllowRead:  check_value_sliceRoleRef(secretm["read_roles"]),
		AllowWrite: check_value_sliceRoleRef(secretm["write_roles"]),
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorJsonSecretRead, d.Get("name"), err))
	}

	d.SetId(d.Get("name").(string))

	return nil
}

func resourcePrivxJsonSecretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("data") || d.HasChange("read_roles") || d.HasChange("write_roles") {

		var secretData = SecretJsonData{
			Data: d.Get("data").(string),
		}

		vaultReq := &tVaultReqJson{
			Data:       secretData,
			Name:       d.Get("name").(string),
			AllowRead:  createJsonSecretListRoleRef(flattenSimpleSlice(d.Get("read_roles").([]interface{}))),
			AllowWrite: createJsonSecretListRoleRef(flattenSimpleSlice(d.Get("write_roles").([]interface{}))),
			OwnerID:    "",
		}
		curl := meta.(privx_API_client_connector).Connector.URL(fmt.Sprintf("/vault/api/v1/secrets/%s", vaultReq.Name))
		_, err := curl.Put(vaultReq)

		if err != nil {
			return diag.FromErr(fmt.Errorf(errorJsonSecretUpdate, d.Get("name").(string), err))
		}
	}

	return nil
}

func resourcePrivxJsonSecretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	curl := meta.(privx_API_client_connector).Connector.URL(fmt.Sprintf("/vault/api/v1/secrets/%s", d.Get("name").(string)))
	curl.Delete()

	d.SetId("")
	return nil

}

func createJsonSecretListRoleRef(list_role_id []string) []rolestore.RoleRef {
	results := make([]rolestore.RoleRef, 0, len(list_role_id))
	for _, role_id := range list_role_id {
		results = append(results, rolestore.RoleRef{ID: role_id})
	}
	return results
}
