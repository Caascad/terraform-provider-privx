package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"terraform-provider-privx/internal/utils"
	"time"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &RoleResource{}
var _ resource.ResourceWithImportState = &RoleResource{}

func NewRoleResource() resource.Resource {
	return &RoleResource{}
}

// RoleResource defines the resource implementation.
type RoleResource struct {
	client *rolestore.RoleStore
}

// Role contains PrivX role information.
type RoleResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Comment       types.String `tfsdk:"comment"`
	AccessGroupID types.String `tfsdk:"access_group_id"`
	Permissions   types.Set    `tfsdk:"permissions"`
	PublicKey     types.Set    `tfsdk:"principal_public_key_strings"`
	PermitAgent   types.Bool   `tfsdk:"permit_agent"`
	SourceRule    types.String `tfsdk:"source_rules"`
}

func (r *RoleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (r *RoleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Role resource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Role ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the role",
				Required:            true,
			},
			"access_group_id": schema.StringAttribute{
				MarkdownDescription: "Scopes host and connection permissions to an access group. (Defaults to Default access group)",
				Required:            true,
			},
			"comment": schema.StringAttribute{
				MarkdownDescription: "A comment describing the object",
				Optional:            true,
			},
			"permissions": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Role permissions",
				Optional:            true,
				Validators: []validator.Set{
					setvalidator.ValueStringsAre(
						stringvalidator.OneOf(
							"access-groups-manage",
							"api-clients-manage",
							"authorized-keys-manage",
							"certificates-view",
							"connections-authorize",
							"connections-manage",
							"connections-manual",
							"connections-playback",
							"connections-terminate",
							"connections-trail",
							"connections-view",
							"hosts-manage",
							"hosts-view",
							"idp-clients-manage",
							"idp-clients-view",
							"licenses-manage",
							"logs-manage",
							"logs-view",
							"network-targets-manage",
							"network-targets-view",
							"requests-view",
							"role-target-resources-manage",
							"role-target-resources-view",
							"roles-manage",
							"roles-view",
							"settings-manage",
							"settings-view",
							"sources-data-push",
							"sources-manage",
							"sources-view",
							"ueba-manage",
							"ueba-view",
							"users-manage",
							"users-view",
							"vault-add",
							"vault-manage",
							"webauthn-credentials-manage",
							"workflows-manage",
							"workflows-requests",
							"workflows-requests-on-behalf",
							"workflows-view",
						),
					),
				},
			},
			"principal_public_key_strings": schema.SetAttribute{
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "List of role's principal public keys",
				Computed:            true,
			},
			"permit_agent": schema.BoolAttribute{
				MarkdownDescription: "Role permit agent",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"source_rules": schema.StringAttribute{
				MarkdownDescription: `A source rule(s) definition. Can be a single rule or a rule group, in which case either "single" or "group" attributes are requrired. Defined in JSON`,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(`{"type":"GROUP","match":"ANY","rules":[]}`),
			},
		},
	}
}

func (r *RoleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	tflog.Debug(ctx, "Creating rolestore", map[string]interface{}{
		"connector : ": fmt.Sprintf("%+v", *connector),
	})

	r.client = rolestore.New(*connector)
}

func (r *RoleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RoleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Loaded role type data", map[string]interface{}{
		"data": fmt.Sprintf("%+v", data),
	})

	if !data.AccessGroupID.IsUnknown() && data.AccessGroupID.ValueString() == "" {
		resp.Diagnostics.AddError("Attribute error", "access_group_id cannot be set to empty string")
		return
	}

	var permissionsPayload []string
	if len(data.Permissions.Elements()) > 0 {
		resp.Diagnostics.Append(data.Permissions.ElementsAs(ctx, &permissionsPayload, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var sourceRule rolestore.SourceRule
	err := json.Unmarshal([]byte(data.SourceRule.ValueString()), &sourceRule)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"Cannot unmarshal sourceRule to json.\n"+
				err.Error(),
		)
		return
	}

	role := rolestore.Role{
		Name:          data.Name.ValueString(),
		Comment:       data.Comment.ValueString(),
		AccessGroupID: data.AccessGroupID.ValueString(),
		Permissions:   permissionsPayload,
		PermitAgent:   data.PermitAgent.ValueBool(),
		SourceRule:    sourceRule,
	}

	tflog.Debug(ctx, fmt.Sprintf("rolestore.Role model used: %+v", role))

	roleID, err := r.client.CreateRole(role)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create the role, got error: %s", err))
		return
	}

	// generate the principal key for the role (using the role ID)
	principalKeyID, err := r.client.GeneratePrincipalKey(roleID)

	// Get role public key into state.
	// PrivX takes some time to generate them.
	publicKeyData := []string{}
	timeout := 12 * time.Second
	startTime := time.Now()
	for {
		principalKeyRead, err := r.client.PrincipalKey(roleID, principalKeyID)

		if err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read the principal key, got error: %s", err))
			return
		}

		if strings.Contains(principalKeyRead.PublicKey, "ssh-rsa ") {
			publicKeyData = append(publicKeyData, principalKeyRead.PublicKey)
			break
		}
		if time.Since(startTime) > timeout {
			tflog.Debug(ctx, "timeout reached")
			break
		}
		time.Sleep(time.Second)
		tflog.Debug(ctx, fmt.Sprintf("Waiting for public keys to be generated (%s timeout)", timeout))
	}
	publicKey, diags := types.SetValueFrom(ctx, data.PublicKey.ElementType(ctx), publicKeyData)
	if diags.HasError() {
		return
	}
	data.PublicKey = publicKey
	data.ID = types.StringValue(roleID)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while attempting to create the resource.\n"+
				err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "created role resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RoleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *RoleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	role, err := r.client.Role(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read role, got error: %s", err))
		return
	}

	data.Name = types.StringValue(role.Name)
	data.Comment = types.StringValue(role.Comment)
	data.AccessGroupID = types.StringValue(role.AccessGroupID)
	data.PermitAgent = types.BoolValue(role.PermitAgent)

	permissions, diags := types.SetValueFrom(ctx, data.Permissions.ElementType(ctx), role.Permissions)
	if diags.HasError() {
		return
	}
	data.Permissions = permissions

	publicKey, diags := types.SetValueFrom(ctx, data.PublicKey.ElementType(ctx), role.PublicKey)
	if diags.HasError() {
		return
	}
	data.PublicKey = publicKey

	sourceRuleData, err := json.Marshal(role.SourceRule)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Resource",
			"Cannot marshal SourceRule data to json.\n"+
				err.Error(),
		)
		return
	}

	equal, _ := utils.JSONBytesEqual(sourceRuleData, []byte(data.SourceRule.ValueString()))
	if !equal {
		data.SourceRule = types.StringValue(string(sourceRuleData))
	}

	tflog.Debug(ctx, "Storing role type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", data),
	})
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RoleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *RoleResourceModel

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

	var sourceRule rolestore.SourceRule
	err := json.Unmarshal([]byte(data.SourceRule.ValueString()), &sourceRule)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"Cannot unmarshal json.\n"+
				err.Error(),
		)
		return
	}

	var publicKeyPayload []string
	if len(data.PublicKey.Elements()) > 0 {
		resp.Diagnostics.Append(data.PublicKey.ElementsAs(ctx, &publicKeyPayload, false)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	role := rolestore.Role{
		ID:            data.ID.ValueString(),
		Name:          data.Name.ValueString(),
		Comment:       data.Comment.ValueString(),
		AccessGroupID: data.AccessGroupID.ValueString(),
		Permissions:   permissionsPayload,
		PermitAgent:   data.PermitAgent.ValueBool(),
		PublicKey:     publicKeyPayload,
		SourceRule:    sourceRule,
	}

	tflog.Debug(ctx, fmt.Sprintf("rolestore.Role model used: %+v", role))

	err = r.client.UpdateRole(
		data.ID.ValueString(),
		&role)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update role, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RoleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *RoleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteRole(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete role, got error: %s", err))
		return
	}
}

func (r *RoleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
