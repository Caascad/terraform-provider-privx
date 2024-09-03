package provider

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &WebproxyDataSource{}

func NewWebproxyDataSource() datasource.DataSource {
	return &WebproxyDataSource{}
}

// WebproxyDataSource defines the DataSource implementation.
type WebproxyDataSource struct {
	client *userstore.UserStore
}

// Carrier contains PrivX webproxy information.
type WebproxyDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	RoutingPrefix   types.String `tfsdk:"routing_prefix"`
	Name            types.String `tfsdk:"name"`
	Permissions     types.List   `tfsdk:"permissions"`
	Secret          types.String `tfsdk:"secret"`
	WebProxyAddress types.String `tfsdk:"web_proxy_address"`
	ExtenderAddress types.List   `tfsdk:"extender_address"`
	Subnets         types.List   `tfsdk:"subnets"`
	Registered      types.Bool   `tfsdk:"registered"`
	GroupId         types.String `tfsdk:"group_id"`
}

func (r *WebproxyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webproxy"
}

func (r *WebproxyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Carrier DataSource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Carrier ID",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Carrier enabled",
				Computed:            true,
			},
			"routing_prefix": schema.StringAttribute{
				MarkdownDescription: "Routing Prefix",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Carrier name",
				Computed:            true,
			},
			"permissions": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Carrier permissions",
				Computed:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "Carrier secret",
				Sensitive:           true,
				Computed:            true,
			},
			"web_proxy_address": schema.StringAttribute{
				MarkdownDescription: "Web Proxy address",
				Computed:            true,
			},
			"extender_address": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Extender addresses",
				Computed:            true,
			},
			"group_id": schema.StringAttribute{
				MarkdownDescription: "Group ID",
				Required:            true,
			},
			"subnets": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Subnets",
				Computed:            true,
			},
			"registered": schema.BoolAttribute{
				MarkdownDescription: "Carrier registered",
				Computed:            true,
			},
		},
	}
}

func (r *WebproxyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	r.client = userstore.New(*connector)
}

func (r *WebproxyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *WebproxyDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	trustedClientList, err := r.client.TrustedClients()

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read trustedClient list, got error: %s", err))
		return
	}

	var Client userstore.TrustedClient
	found_client := false

	for _, client := range trustedClientList {
		if client.GroupId == data.GroupId.ValueString() && client.Type == "ICAP" {
			found_client = true
			Client = client
			break // Arrêtez la boucle dès que vous trouvez la bonne valeur
		}
	}
	if !found_client {
		resp.Diagnostics.AddError("Proxy Error", fmt.Sprintf("Unable to find associated WebProxy, got error: %s", data.GroupId.ValueString()))
		return
	}

	data.ID = types.StringValue(Client.ID)
	data.Name = types.StringValue(Client.Name)
	data.Secret = types.StringValue(Client.Secret)
	data.Registered = types.BoolValue(Client.Registered)
	data.Enabled = types.BoolValue(Client.Enabled)
	data.WebProxyAddress = types.StringValue(Client.WebProxyAddress)
	data.RoutingPrefix = types.StringValue(Client.RoutingPrefix)

	subnets, diags := types.ListValueFrom(ctx, data.Subnets.ElementType(ctx), Client.Subnets)
	if diags.HasError() {
		return
	}
	data.Subnets = subnets

	permissions, diags := types.ListValueFrom(ctx, data.Permissions.ElementType(ctx), Client.Permissions)
	if diags.HasError() {
		return
	}
	data.Permissions = permissions

	extenderAddress, diags := types.ListValueFrom(ctx, data.ExtenderAddress.ElementType(ctx), Client.ExtenderAddress)
	if diags.HasError() {
		return
	}
	data.ExtenderAddress = extenderAddress

	tflog.Debug(ctx, "Storing webproxy type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", data),
	})
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
