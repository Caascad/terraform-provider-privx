package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/vault"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SecretResource{}
var _ resource.ResourceWithImportState = &SecretResource{}

func NewSecretResource() resource.Resource {
	return &SecretResource{}
}

// SecretResource defines the resource implementation.
type SecretResource struct {
	client *vault.Vault
}

// SecretResourceModel describes the resource data model.
type (
	RoleRefModel struct {
		ID   types.String `tfsdk:"id"`
		Name types.String `tfsdk:"name"`
	}

	SecretResourceModel struct {
		Name       types.String   `tfsdk:"name"`
		Data       types.String   `tfsdk:"data"`
		ReadRoles  []RoleRefModel `tfsdk:"read_roles"`
		WriteRoles []RoleRefModel `tfsdk:"write_roles"`
	}
)

func (r *SecretResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (r *SecretResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Secret resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Secret's name",
				Required:            true,
			},
			"data": schema.StringAttribute{
				MarkdownDescription: "Secret to be stored",
				Optional:            true,
				Computed:            true,
				Sensitive:           true,
				Default:             stringdefault.StaticString("{}"),
			},
			"read_roles": schema.SetNestedAttribute{
				MarkdownDescription: "List of roles that can read secret.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Role ID",
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Role name, ignored by server in requests.",
							Optional:            true,
						},
					},
				},
			},
			"write_roles": schema.SetNestedAttribute{
				MarkdownDescription: "List of roles that can replace secret.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Role ID",
							Required:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Role name, ignored by server in requests.",
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

func (r *SecretResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	tflog.Debug(ctx, "Creating vault client", map[string]interface{}{
		"connector : ": fmt.Sprintf("%+v", *connector),
	})

	r.client = vault.New(*connector)
}

func (r *SecretResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SecretResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Loaded secret type data", map[string]interface{}{
		"data": fmt.Sprintf("%+v", data),
	})

	var readRolesPayload []string
	for _, roleRef := range data.ReadRoles {
		readRolesPayload = append(readRolesPayload, roleRef.ID.ValueString())
	}

	var writeRolesPayload []string
	for _, roleRef := range data.WriteRoles {
		writeRolesPayload = append(writeRolesPayload, roleRef.ID.ValueString())
	}

	var secretPayload interface{}
	if err := json.Unmarshal([]byte(data.Data.ValueString()), &secretPayload); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"Cannot unmarshal secret data json.\n"+
				err.Error(),
		)
		return
	}

	ctx = tflog.SetField(ctx, "secret name", data.Name.ValueString())
	ctx = tflog.SetField(ctx, "allowed read", readRolesPayload)
	ctx = tflog.SetField(ctx, "allowed write", writeRolesPayload)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "secret data")
	tflog.Debug(ctx, "Created secret")

	if err := r.client.CreateSecret(data.Name.ValueString(), readRolesPayload, writeRolesPayload, secretPayload); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while attempting to create the resource.\n"+
				err.Error(),
		)
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SecretResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SecretResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := r.client.Secret(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read secret : %s, got error: %s", data.Name.ValueString(), err))
		return
	}

	var allowedRead []RoleRefModel
	for _, v := range secret.AllowRead {
		allowedRead = append(allowedRead, RoleRefModel{types.StringValue(v.ID), types.StringValue(v.Name)})
	}
	data.ReadRoles = allowedRead

	var allowedWrite []RoleRefModel
	for _, v := range secret.AllowWrite {
		allowedWrite = append(allowedWrite, RoleRefModel{types.StringValue(v.ID), types.StringValue(v.Name)})
	}
	data.WriteRoles = allowedWrite

	secretData, err := json.Marshal(secret.Data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Resource",
			"Cannot marshal SourceRule data to json.\n"+
				err.Error(),
		)
		return
	}
	data.Data = types.StringValue(string(secretData))

	tflog.Debug(ctx, "Storing secret type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", data),
	})
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SecretResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SecretResourceModel
	var name_from_state string

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("name"), &name_from_state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var readRolesPayload []string
	for _, roleRef := range data.ReadRoles {
		readRolesPayload = append(readRolesPayload, roleRef.ID.ValueString())
	}

	var writeRolesPayload []string
	for _, roleRef := range data.WriteRoles {
		writeRolesPayload = append(writeRolesPayload, roleRef.ID.ValueString())
	}

	var secretPayload interface{}
	if err := json.Unmarshal([]byte(data.Data.ValueString()), &secretPayload); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"Cannot unmarshal secret data json.\n"+
				err.Error(),
		)
		return
	}

	//If secret change name trigger delete of old secret then recreate
	if data.Name.ValueString() != name_from_state {
		if err := r.client.CreateSecret(data.Name.ValueString(), readRolesPayload, writeRolesPayload, secretPayload); err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create Resource",
				"An unexpected error occurred while attempting to create the resource.\n"+
					err.Error(),
			)
			return
		}
		if err := r.client.DeleteSecret(name_from_state); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete secret, got error: %s", err))
			return
		}
	} else {
		if err := r.client.UpdateSecret(data.Name.ValueString(), readRolesPayload, writeRolesPayload, secretPayload); err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create Resource",
				"An unexpected error occurred while attempting to create the resource.\n"+
					err.Error(),
			)
			return
		}
	}

	ctx = tflog.SetField(ctx, "secret name", data.Name.ValueString())
	ctx = tflog.SetField(ctx, "allowed read", readRolesPayload)
	ctx = tflog.SetField(ctx, "allowed write", writeRolesPayload)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "secret data")
	tflog.Debug(ctx, "Updated secret")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SecretResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SecretResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteSecret(data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete secret, got error: %s", err))
		return
	}
}

func (r *SecretResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
