package provider

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &CarrierResource{}
var _ resource.ResourceWithImportState = &CarrierResource{}

func NewCarrierResource() resource.Resource {
	return &CarrierResource{}
}

// CarrierResource defines the resource implementation.
type CarrierResource struct {
	client *userstore.UserStore
}

// Carrier contains PrivX carrier information.
type CarrierResourceModel struct {
	ID                            types.String `tfsdk:"id"`
	Type                          types.String `tfsdk:"type"`
	Enabled                       types.Bool   `tfsdk:"enabled"`
	RoutingPrefix                 types.String `tfsdk:"routing_prefix"`
	Name                          types.String `tfsdk:"name"`
	Permissions                   types.List   `tfsdk:"permissions"`
	WebProxyAddress               types.String `tfsdk:"web_proxy_address"`
	WebProxyPort                  types.Int64  `tfsdk:"web_proxy_port"`
	WebProxyExtenderRoutePatterns types.List   `tfsdk:"web_proxy_extender_route_patterns"`
	ExtenderAddress               types.List   `tfsdk:"extender_address"`
	Subnets                       types.List   `tfsdk:"subnets"`
	Registered                    types.Bool   `tfsdk:"registered"`
	AccessGroupId                 types.String `tfsdk:"access_group_id"`
	GroupID                       types.String `tfsdk:"group_id"`
}

func (r *CarrierResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_carrier"
}

func (r *CarrierResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Carrier resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Carrier ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Trusted client Type",
				Computed:            true,
				Default:             stringdefault.StaticString("CARRIER"),
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Carrier enabled",
				Required:            true,
			},
			"routing_prefix": schema.StringAttribute{
				MarkdownDescription: "Routing Prefix",
				Optional:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Carrier name",
				Required:            true,
			},
			"permissions": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Carrier permissions",
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"web_proxy_address": schema.StringAttribute{
				MarkdownDescription: "Web Proxy address",
				Required:            true,
			},
			"web_proxy_port": schema.Int64Attribute{
				MarkdownDescription: "Web Proxy address",
				Optional:            true,
			},
			"web_proxy_extender_route_patterns": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Web Proxy Extender Route Patterns",
				Optional:            true,
			},
			"extender_address": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Extender addresses",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"subnets": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Subnets",
				Optional:            true,
			},
			"access_group_id": schema.StringAttribute{
				MarkdownDescription: "Access Group ID",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"group_id": schema.StringAttribute{
				MarkdownDescription: "Group ID",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"registered": schema.BoolAttribute{
				MarkdownDescription: "Carrier registered",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *CarrierResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	connector, ok := req.ProviderData.(*restapi.Connector)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *restapi.Connector, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	tflog.Debug(ctx, "Creating userstore", map[string]interface{}{
		"connector : ": fmt.Sprintf("%+v", *connector),
	})

	r.client = userstore.New(*connector)
}

func (r *CarrierResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CarrierResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Loaded Carrier type data", map[string]interface{}{
		"data": fmt.Sprintf("%+v", data),
	})

	permissionPayload := []string{"privx-carrier"}

	var extenderAddressPayload []string
	if len(data.ExtenderAddress.Elements()) > 0 {
		resp.Diagnostics.Append(data.ExtenderAddress.ElementsAs(ctx, &extenderAddressPayload, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var subnetsPayload []string
	if len(data.Subnets.Elements()) > 0 {
		resp.Diagnostics.Append(data.Subnets.ElementsAs(ctx, &subnetsPayload, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var WebProxyExtenderRoutPattersPayload []string
	if len(data.WebProxyExtenderRoutePatterns.Elements()) > 0 {
		resp.Diagnostics.Append(data.WebProxyExtenderRoutePatterns.ElementsAs(ctx, &WebProxyExtenderRoutPattersPayload, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	carrier := userstore.TrustedClient{
		Type:                          userstore.ClientType(data.Type.ValueString()),
		Name:                          data.Name.ValueString(),
		Enabled:                       data.Enabled.ValueBool(),
		Permissions:                   permissionPayload,
		AccessGroupId:                 data.AccessGroupId.ValueString(),
		GroupId:                       data.GroupID.ValueString(),
		ExtenderAddress:               extenderAddressPayload,
		WebProxyAddress:               data.WebProxyAddress.ValueString(),
		WebProxyExtenderRoutePatterns: WebProxyExtenderRoutPattersPayload,
		Subnets:                       subnetsPayload,
		RoutingPrefix:                 data.RoutingPrefix.ValueString(),
	}

	tflog.Debug(ctx, fmt.Sprintf("userstore.TrustedClient model used: %+v", carrier))

	carrierID, err := r.client.CreateTrustedClient(carrier)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while attempting to create the resource.\n"+
				err.Error(),
		)
		return
	}

	// Convert from the API data model to the Terraform data model
	// and set any unknown attribute values.
	data.ID = types.StringValue(carrierID)

	carrierRead, err := r.client.TrustedClient(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Resource",
			"An unexpected error occurred while attempting to read the resource.\n"+
				err.Error(),
		)
		return
	}
	data.Registered = types.BoolValue(carrierRead.Registered)
	data.AccessGroupId = types.StringValue(carrierRead.AccessGroupId)
	data.GroupID = types.StringValue(carrierRead.GroupId)
	permissions, diags := types.ListValueFrom(ctx, data.Permissions.ElementType(ctx), carrierRead.Permissions)
	if diags.HasError() {
		return
	}
	data.Permissions = permissions

	extender_address, diags := types.ListValueFrom(ctx, data.ExtenderAddress.ElementType(ctx), carrierRead.ExtenderAddress)
	if diags.HasError() {
		return
	}
	data.ExtenderAddress = extender_address

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Debug(ctx, "created carrier resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CarrierResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CarrierResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	carrier, err := r.client.TrustedClient(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read carrier, got error: %s", err))
		return
	}

	data.Name = types.StringValue(carrier.Name)
	data.Registered = types.BoolValue(carrier.Registered)
	data.Enabled = types.BoolValue(carrier.Enabled)
	data.RoutingPrefix = types.StringValue(carrier.RoutingPrefix)

	subnets, diags := types.ListValueFrom(ctx, data.Subnets.ElementType(ctx), carrier.Subnets)
	if diags.HasError() {
		return
	}
	data.Subnets = subnets

	permissions, diags := types.ListValueFrom(ctx, data.Permissions.ElementType(ctx), carrier.Permissions)
	if diags.HasError() {
		return
	}
	data.Permissions = permissions

	extenderAddress, diags := types.ListValueFrom(ctx, data.ExtenderAddress.ElementType(ctx), carrier.ExtenderAddress)
	if diags.HasError() {
		return
	}
	data.ExtenderAddress = extenderAddress

	tflog.Debug(ctx, "Storing carrier type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", data),
	})
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CarrierResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CarrierResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var extenderAddressPayload []string
	if len(data.ExtenderAddress.Elements()) > 0 {
		resp.Diagnostics.Append(data.ExtenderAddress.ElementsAs(ctx, &extenderAddressPayload, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var subnetsPayload []string
	if len(data.Subnets.Elements()) > 0 {
		resp.Diagnostics.Append(data.Subnets.ElementsAs(ctx, &subnetsPayload, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var WebProxyExtenderRoutPattersPayload []string
	if len(data.WebProxyExtenderRoutePatterns.Elements()) > 0 {
		resp.Diagnostics.Append(data.WebProxyExtenderRoutePatterns.ElementsAs(ctx, &WebProxyExtenderRoutPattersPayload, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	permissionPayload := []string{"privx-carrier"}

	carrier := userstore.TrustedClient{
		Type:                          userstore.ClientType(data.Type.ValueString()),
		Name:                          data.Name.ValueString(),
		Enabled:                       data.Enabled.ValueBool(),
		Permissions:                   permissionPayload,
		AccessGroupId:                 data.AccessGroupId.ValueString(),
		GroupId:                       data.GroupID.ValueString(),
		ExtenderAddress:               extenderAddressPayload,
		WebProxyAddress:               data.WebProxyAddress.ValueString(),
		WebProxyExtenderRoutePatterns: WebProxyExtenderRoutPattersPayload,
		Subnets:                       subnetsPayload,
		RoutingPrefix:                 data.RoutingPrefix.ValueString(),
	}

	tflog.Debug(ctx, fmt.Sprintf("userstore.TrustedClient model used: %+v", carrier))

	err := r.client.UpdateTrustedClient(
		data.ID.ValueString(),
		&carrier)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update carrier, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CarrierResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CarrierResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTrustedClient(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete carrier, got error: %s", err))
		return
	}
}

func (r *CarrierResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
