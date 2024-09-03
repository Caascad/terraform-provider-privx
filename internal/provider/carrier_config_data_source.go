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
var _ datasource.DataSource = &CarrierConfigDataSource{}

func NewCarrierConfigDataSource() datasource.DataSource {
	return &CarrierConfigDataSource{}
}

// CarrierConfigDataSource defines the DataSource implementation.
type CarrierConfigDataSource struct {
	client    *authorizer.Client
	connector *restapi.Connector
}

// CarrierConfig contains PrivX CarrierConfig information.
type CarrierConfigDataSourceModel struct {
	TrustedClientID types.String `tfsdk:"trusted_client_id"`
	CarrierConfig   types.String `tfsdk:"carrier_config"`
}

func (r *CarrierConfigDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_carrier_config"
}

func (r *CarrierConfigDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "CarrierConfig DataSource",
		Attributes: map[string]schema.Attribute{
			"trusted_client_id": schema.StringAttribute{
				MarkdownDescription: "CarrierConfig ID",
				Required:            true,
			},
			"carrier_config": schema.StringAttribute{
				MarkdownDescription: "Carrier config",
				Computed:            true,
			},
		},
	}
}

func (r *CarrierConfigDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *CarrierConfigDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *CarrierConfigDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	DownloadHandle, err := r.client.CarrierConfigDownloadHandle(data.TrustedClientID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to get carrier download sessionid, got error: %s", err))
		return
	}
	carrier_config, err := GetCarrierConfig(*r.connector, data.TrustedClientID.ValueString(), DownloadHandle.SessionID)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Cannot get carrier config: %s", err))
		return
	}

	data.CarrierConfig = types.StringValue(carrier_config)

	tflog.Debug(ctx, "Storing CarrierConfig type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", data),
	})
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func GetCarrierConfig(restapi_connector restapi.Connector, trusted_client_id, session_id string) (string, error) {
	curl := restapi_connector.URL(fmt.Sprintf("/authorizer/api/v1/carrier/conf/%s/%s", trusted_client_id, session_id))
	resp, err := curl.Fetch()
	if err != nil {
		return "", err
	}
	return string(resp), nil
}
