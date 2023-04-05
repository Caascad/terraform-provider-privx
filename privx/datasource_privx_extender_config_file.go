package privx

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourcePrivxExtenderConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: extenderConfig,

		Schema: map[string]*schema.Schema{
			"trusted_client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"session_id": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"extender_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func extenderConfig(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	connector := meta.(privx_API_client_connector).Connector
	authorizer_client := createAuthorizerClient(ctx, connector)
	extender_session_id, err := authorizer_client.ExtenderConfigDownloadHandle(d.Get("trusted_client_id").(string))

	if err != nil {
		return diag.FromErr(fmt.Errorf("error generating ExtenderConfigDownloadHandle : %s", err))
	}

	if err := d.Set("session_id", extender_session_id.SessionID); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `session_id`: %s", err))
	}

	extender_config, err := GetExtenderConfig(connector, d.Get("trusted_client_id").(string), extender_session_id.SessionID)

	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting extender_config : %s", err))
	}

	if err := d.Set("extender_config", extender_config); err != nil {
		return diag.FromErr(fmt.Errorf("error setting `session_id`: %s", err))
	}

	d.SetId(d.Get("trusted_client_id").(string))

	return nil
}

func GetExtenderConfig(restapi_connector restapi.Connector, trusted_client_id, session_id string) (string, error) {
	curl := restapi_connector.URL(fmt.Sprintf("/authorizer/api/v1/extender/conf/%s/%s", trusted_client_id, session_id))

	resp, err := curl.Fetch()

	return string(resp), err
}
