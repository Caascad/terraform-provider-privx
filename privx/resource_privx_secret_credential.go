package privx

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	errorSecretCreate = "error creating PrivX Secret (%s): %s"
	errorSecretUpdate = "error updating PrivX Secret (%s): %s"
	errorSecretDelete = "error deleting PrivX Secret (%s): %s"
	errorSecretRead   = "error reading PrivX Secret (%s): %s"
	_privx_schema     = "credentials"
)

type SecretCredentialsData struct {
	User         string `json:"user"`
	Pass         string `json:"pass"`
	Comment      string `json:"comment"`
	Privx_schema string `json:"_privx_schema"`
}

// tVaultReq t vault request definition
type tVaultReq struct {
	Name       string                `json:"name"`
	Data       SecretCredentialsData `json:"data"`
	AllowRead  []rolestore.RoleRef   `json:"read_roles,omitempty"`
	AllowWrite []rolestore.RoleRef   `json:"write_roles,omitempty"`
	OwnerID    string                `json:"owner_id,omitempty"`
}

func resourcePrivXSecretCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePrivxSecretCreate,
		ReadContext:   resourcePrivxSecretRead,
		UpdateContext: resourcePrivxSecretUpdate,
		DeleteContext: resourcePrivxSecretDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"pass": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"comment": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
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

func resourcePrivxSecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var secretData = SecretCredentialsData{
		User:         d.Get("user").(string),
		Pass:         d.Get("pass").(string),
		Comment:      d.Get("comment").(string),
		Privx_schema: _privx_schema,
	}

	curl := meta.(privx_API_client_connector).Connector.URL(fmt.Sprintf("/vault/api/v1/secrets"))

	vaultReq := tVaultReq{
		Data:       secretData,
		Name:       d.Get("name").(string),
		AllowRead:  createSecretListRoleRef(flattenSimpleSlice(d.Get("read_roles").([]interface{}))),
		AllowWrite: createSecretListRoleRef(flattenSimpleSlice(d.Get("write_roles").([]interface{}))),
		OwnerID:    d.Get("owner_id").(string),
	}
	_, err := curl.Post(vaultReq)

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorSecretCreate, d.Get("name").(string), err))
	}

	d.SetId(d.Get("name").(string))

	return nil
}

func check_value_string(value interface{}) string {
	if s, ok := value.(string); ok {
		return s
	}
	return "" // Default value
}

func check_value_SecretCredentialsData(value interface{}) SecretCredentialsData {
	valuem, _ := value.(map[string]interface{})
	secretData := SecretCredentialsData{
		User:    valuem["user"].(string),
		Pass:    valuem["pass"].(string),
		Comment: valuem["comment"].(string),
	}
	return secretData
}

func check_value_roleRef(value interface{}) rolestore.RoleRef {
	valuem, _ := value.(map[string]interface{})
	role := rolestore.RoleRef{
		ID:   valuem["id"].(string),
		Name: valuem["name"].(string),
	}
	return role
}

func check_value_sliceRoleRef(value interface{}) []rolestore.RoleRef {
	valuem, ok := value.(map[string]interface{})
	list_Role := make([]rolestore.RoleRef, 0, len(valuem))
	if ok {
		for _, role := range valuem {
			list_Role = append(list_Role, check_value_roleRef(role))
		}
	}
	return list_Role
}

func resourcePrivxSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	/*Get Secret List*/
	curl := meta.(privx_API_client_connector).Connector.URL(fmt.Sprintf("/vault/api/v1/secrets/%s", d.Get("name").(string)))
	var secret interface{}
	_, err := curl.Get(&secret)
	secretm, ok := secret.(map[string]interface{})
	if !ok {
		return diag.FromErr(fmt.Errorf(errorSecretRead, d.Get("name"), "convertion error"))
	} else {
		secret = tVaultReq{
			Name:       secretm["name"].(string),
			OwnerID:    check_value_string(secretm["owner_id"]),
			Data:       check_value_SecretCredentialsData(secretm["data"]),
			AllowRead:  check_value_sliceRoleRef(secretm["read_roles"]),
			AllowWrite: check_value_sliceRoleRef(secretm["write_roles"]),
		}
	}

	if err != nil {
		return diag.FromErr(fmt.Errorf(errorSecretRead, d.Get("name"), err))
	}

	d.SetId(d.Get("name").(string))

	return nil
}

func resourcePrivxSecretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("user") || d.HasChange("pass") || d.HasChange("comment") || d.HasChange("read_roles") || d.HasChange("write_roles") || d.HasChange("owner_id") {

		var secretData = SecretCredentialsData{
			User:         d.Get("user").(string),
			Pass:         d.Get("pass").(string),
			Comment:      d.Get("comment").(string),
			Privx_schema: _privx_schema,
		}

		vaultReq := &tVaultReq{
			Data:       secretData,
			Name:       d.Get("name").(string),
			AllowRead:  createSecretListRoleRef(flattenSimpleSlice(d.Get("read_roles").([]interface{}))),
			AllowWrite: createSecretListRoleRef(flattenSimpleSlice(d.Get("write_roles").([]interface{}))),
			OwnerID:    d.Get("owner_id").(string),
		}
		curl := meta.(privx_API_client_connector).Connector.URL(fmt.Sprintf("/vault/api/v1/secrets/%s", vaultReq.Name))
		_, err := curl.Put(vaultReq)

		if err != nil {
			return diag.FromErr(fmt.Errorf(errorSecretUpdate, d.Get("name").(string), err))
		}
	}

	return nil
}

func resourcePrivxSecretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	curl := meta.(privx_API_client_connector).Connector.URL(fmt.Sprintf("/vault/api/v1/secrets/%s", d.Get("name").(string)))
	curl.Delete()

	d.SetId("")
	return nil

}

func createSecretListRoleRef(list_role_id []string) []rolestore.RoleRef {
	results := make([]rolestore.RoleRef, 0, len(list_role_id))
	for _, role_id := range list_role_id {
		results = append(results, rolestore.RoleRef{ID: role_id})
	}
	return results
}
