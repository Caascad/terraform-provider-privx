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
var _ datasource.DataSource = &ExtenderDataSource{}

func NewExtenderDataSource() datasource.DataSource {
	return &ExtenderDataSource{}
}

// ExtenderDataSource defines the DataSource implementation.
type ExtenderDataSource struct {
	client *userstore.UserStore
}

// Extender contains PrivX extender information.
type ExtenderDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	RoutingPrefix   types.String `tfsdk:"routing_prefix"`
	Name            types.String `tfsdk:"name"`
	Permissions     types.List   `tfsdk:"permissions"`
	Secret          types.String `tfsdk:"secret"`
	WebProxyAddress types.String `tfsdk:"web_proxy_address"`
	WebProxyPort    types.Int64  `tfsdk:"web_proxy_port"`
	ExtenderAddress types.List   `tfsdk:"extender_address"`
	Subnets         types.List   `tfsdk:"subnets"`
	Registered      types.Bool   `tfsdk:"registered"`
}

func (r *ExtenderDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extender"
}

func (r *ExtenderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Extender DataSource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Extender ID",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Extender enabled",
				Computed:            true,
			},
			"routing_prefix": schema.StringAttribute{
				MarkdownDescription: "Routing Prefix",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Extender name",
				Required:            true,
			},
			"permissions": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Extender permissions",
				Computed:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "Extender secret",
				Sensitive:           true,
				Computed:            true,
			},
			"web_proxy_address": schema.StringAttribute{
				MarkdownDescription: "Web Proxy address",
				Computed:            true,
			},
			"web_proxy_port": schema.Int64Attribute{
				MarkdownDescription: "Web Proxy address",
				Computed:            true,
			},
			"extender_address": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Extender addresses",
				Computed:            true,
			},
			"subnets": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Subnets",
				Computed:            true,
			},
			"registered": schema.BoolAttribute{
				MarkdownDescription: "Extender registered",
				Computed:            true,
			},
		},
	}
}

func (r *ExtenderDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *ExtenderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *ExtenderDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	extender, err := r.client.TrustedClient(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read extender, got error: %s", err))
		return
	}

	data.Name = types.StringValue(extender.Name)
	data.Secret = types.StringValue(extender.Secret)
	data.Registered = types.BoolValue(extender.Registered)
	data.Enabled = types.BoolValue(extender.Enabled)
	data.RoutingPrefix = types.StringValue(extender.RoutingPrefix)

	subnets, diags := types.ListValueFrom(ctx, data.Subnets.ElementType(ctx), extender.Subnets)
	if diags.HasError() {
		return
	}
	data.Subnets = subnets

	permissions, diags := types.ListValueFrom(ctx, data.Permissions.ElementType(ctx), extender.Permissions)
	if diags.HasError() {
		return
	}
	data.Permissions = permissions

	extenderAddress, diags := types.ListValueFrom(ctx, data.ExtenderAddress.ElementType(ctx), extender.ExtenderAddress)
	if diags.HasError() {
		return
	}
	data.ExtenderAddress = extenderAddress

	tflog.Debug(ctx, "Storing extender type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", data),
	})
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
