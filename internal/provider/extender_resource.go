package provider

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/userstore"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ExtenderResource{}
var _ resource.ResourceWithImportState = &ExtenderResource{}

func NewExtenderResource() resource.Resource {
	return &ExtenderResource{}
}

// ExtenderResource defines the resource implementation.
type ExtenderResource struct {
	client *userstore.UserStore
}

// Extender contains PrivX extender information.
type ExtenderResourceModel struct {
	ID              types.String `tfsdk:"id"`
	Secret          types.String `tfsdk:"secret"`
	Name            types.String `tfsdk:"name"`
	WebProxyAddress types.String `tfsdk:"web_proxy_address"`
	Registered      types.Bool   `tfsdk:"registered"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	Permissions     types.List   `tfsdk:"permissions"`
	ExtenderAddress types.List   `tfsdk:"extender_address"`
}

func (r *ExtenderResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extender"
}

func (r *ExtenderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Extender resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Extender ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Extender name",
				Optional:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "Extender secret",
				Sensitive:           true,
				Computed:            true,
			},
			"registered": schema.BoolAttribute{
				MarkdownDescription: "Extender registered",
				Computed:            true,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Extender enabled",
				Optional:            true,
				Computed:            true,
			},
			"permissions": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Extender permissions",
				Optional:            true,
			},
			"extender_address": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Extender addresses",
				Optional:            true,
			},
		},
	}
}

func (r *ExtenderResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ExtenderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ExtenderResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Loaded extender type data", map[string]interface{}{
		"data": fmt.Sprintf("%+v", data),
	})

	var permissionsPayload []string
	if len(data.Permissions.Elements()) > 0 {
		resp.Diagnostics.Append(data.Permissions.ElementsAs(ctx, &permissionsPayload, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var extenderAddressPayload []string
	if len(data.Permissions.Elements()) > 0 {
		resp.Diagnostics.Append(data.ExtenderAddress.ElementsAs(ctx, &extenderAddressPayload, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	extender := userstore.TrustedClient{
		Type:            userstore.ClientExtender,
		Name:            data.Name.ValueString(),
		Registered:      data.Registered.ValueBool(),
		Enabled:         data.Enabled.ValueBool(),
		Permissions:     permissionsPayload,
		ExtenderAddress: extenderAddressPayload,
	}

	tflog.Debug(ctx, fmt.Sprintf("userstore.Extender model used: %+v", extender))

	extenderID, err := r.client.CreateTrustedClient(extender)

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
	data.ID = types.StringValue(extenderID)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Debug(ctx, "created extender resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExtenderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ExtenderResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

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

func (r *ExtenderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ExtenderResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var permissionsPayload []string
	if len(data.Permissions.Elements()) > 0 {
		resp.Diagnostics.Append(data.Permissions.ElementsAs(ctx, &permissionsPayload, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var extenderAddressPayload []string
	if len(data.ExtenderAddress.Elements()) > 0 {
		resp.Diagnostics.Append(data.ExtenderAddress.ElementsAs(ctx, &extenderAddressPayload, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	extender := userstore.TrustedClient{
		Name:            data.Name.ValueString(),
		Secret:          data.Secret.ValueString(),
		Permissions:     permissionsPayload,
		ExtenderAddress: extenderAddressPayload,
		Registered:      data.Registered.ValueBool(),
		Enabled:         data.Enabled.ValueBool(),
	}

	tflog.Debug(ctx, fmt.Sprintf("userstore.Extender model used: %+v", extender))

	err := r.client.UpdateTrustedClient(
		data.ID.ValueString(),
		&extender)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update extender, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ExtenderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ExtenderResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteTrustedClient(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete extender, got error: %s", err))
		return
	}
}

func (r *ExtenderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
