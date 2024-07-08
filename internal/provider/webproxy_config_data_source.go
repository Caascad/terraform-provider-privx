package provider

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/authorizer"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &WebproxyConfigDataSource{}

func NewWebproxyConfigDataSource() datasource.DataSource {
	return &WebproxyConfigDataSource{}
}

// WebproxyConfigDataSource defines the DataSource implementation.
type WebproxyConfigDataSource struct {
	client    *authorizer.Client
	connector *restapi.Connector
}

// WebproxyConfig contains PrivX WebproxyConfig information.
type WebproxyConfigDataSourceModel struct {
	TrustedClientID types.String `tfsdk:"trusted_client_id"`
	WebproxyConfig  types.String `tfsdk:"webproxy_config"`
}

func (r *WebproxyConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webproxy_config"
}

func (r *WebproxyConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "WebproxyConfig DataSource",
		Attributes: map[string]schema.Attribute{
			"trusted_client_id": schema.StringAttribute{
				MarkdownDescription: "WebproxyConfig ID",
				Required:            true,
			},
			"webproxy_config": schema.StringAttribute{
				MarkdownDescription: "Web Proxy config",
				Computed:            true,
			},
		},
	}
}

func (r *WebproxyConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	connector, ok := req.ProviderData.(*restapi.Connector)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *restapi.Connector, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	tflog.Debug(ctx, "Creating userstore", map[string]interface{}{
		"connector : ": fmt.Sprintf("%+v", *connector),
	})

	r.connector = connector

	r.client = authorizer.New(*connector)
}

func (r *WebproxyConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *WebproxyConfigDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	DownloadHandle, err := r.client.WebProxySessionDownloadHandle(data.TrustedClientID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get web proxy download sessionid, got error: %s", err))
		return
	}
	webproxy_config, err := GetWebproxyConfig(*r.connector, data.TrustedClientID.ValueString(), DownloadHandle.SessionID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Cannot get web proxy config: %s", err))
		return
	}

	data.WebproxyConfig = types.StringValue(webproxy_config)

	tflog.Debug(ctx, "Storing WebproxyConfig type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", data),
	})
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func GetWebproxyConfig(restapi_connector restapi.Connector, trusted_client_id, session_id string) (string, error) {
	curl := restapi_connector.URL(fmt.Sprintf("/authorizer/api/v1/icap/conf/%s/%s", trusted_client_id, session_id))
	resp, err := curl.Fetch()
	if err != nil {
		return "", err
	}
	return string(resp), nil
}
