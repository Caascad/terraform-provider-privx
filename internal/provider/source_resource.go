package provider

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SourceResource{}
var _ resource.ResourceWithImportState = &SourceResource{}

func NewSourceResource() resource.Resource {
	return &SourceResource{}
}

// SourceResource defines the resource implementation.
type SourceResource struct {
	client *rolestore.RoleStore
}

type (
	OIDCConnectionModel struct {
		Address           types.String `tfsdk:"address"`
		Enabled           types.Bool   `tfsdk:"enabled"`
		Issuer            types.String `tfsdk:"issuer"`
		ButtonTitle       types.String `tfsdk:"button_title"`
		ClientID          types.String `tfsdk:"client_id"`
		ClientSecret      types.String `tfsdk:"client_secret"`
		TagsAttributeName types.String `tfsdk:"tags_attribute_name"`
		ScopesSecret      types.List   `tfsdk:"additional_scopes_secret"`
	}

	EUMModel struct {
		SourceID           types.String `tfsdk:"source_id"`
		SourceSeaerchField types.String `tfsdk:"source_search_field"`
	}

	// SourceResourceModel describes the resource data model.
	SourceResourceModel struct {
		ID                  types.String         `tfsdk:"id"`
		Enabled             types.Bool           `tfsdk:"enabled"`
		TTL                 types.Int64          `tfsdk:"ttl"`
		Name                types.String         `tfsdk:"name"`
		Comment             types.String         `tfsdk:"comment"`
		Tags                types.List           `tfsdk:"tags"`
		UsernamePattern     types.List           `tfsdk:"username_pattern"`
		ExternalUserMapping []*EUMModel          `tfsdk:"external_user_mapping"`
		OIDCConnection      *OIDCConnectionModel `tfsdk:"oidc_connection"`
	}
)

func (r *SourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_source"
}

func (r *SourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Source resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Source ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Source name",
				Optional:            true,
			},
			"comment": schema.StringAttribute{
				MarkdownDescription: "Source comment",
				Optional:            true,
			},
			"ttl": schema.Int64Attribute{
				MarkdownDescription: "Source ttl",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(900),
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Source enabled",
				Optional:            true,
			},
			"tags": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Source tags",
				Optional:            true,
			},
			"username_pattern": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Source external user pattern",
				Optional:            true,
			},
			"external_user_mapping": schema.ListAttribute{
				MarkdownDescription: "Source external user mapping",
				Optional:            true,
				ElementType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"source_id":           types.StringType,
						"source_search_field": types.StringType,
					}},
			},
			"oidc_connection": schema.SingleNestedAttribute{
				MarkdownDescription: "OIDC connection",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"address": schema.StringAttribute{
						MarkdownDescription: "oidc connection address",
						Optional:            true,
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "oidc connection enabled",
						Optional:            true,
					},
					"issuer": schema.StringAttribute{
						MarkdownDescription: "oidc connection issuer",
						Optional:            true,
					},
					"button_title": schema.StringAttribute{
						MarkdownDescription: "oidc connection title",
						Optional:            true,
					},
					"client_id": schema.StringAttribute{
						MarkdownDescription: "oidc connection client ID",
						Optional:            true,
					},
					"client_secret": schema.StringAttribute{
						MarkdownDescription: "oidc connection client Secret",
						Computed:            true,
						Optional:            true,
						Sensitive:           true,
					},
					"tags_attribute_name": schema.StringAttribute{
						MarkdownDescription: "oidc connection tags attribute name",
						Optional:            true,
					},
					"additional_scopes_secret": schema.ListAttribute{
						ElementType:         types.StringType,
						MarkdownDescription: "oidc additional scopes secret",
						Optional:            true,
					},
				},
			},
		},
	}
}

func (r *SourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = rolestore.New(*connector)
}

func (r *SourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data SourceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tagsPayload := make([]string, len(data.Tags.Elements()))
	resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tagsPayload, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userNamePatternPayload := make([]string, len(data.UsernamePattern.Elements()))
	resp.Diagnostics.Append(data.UsernamePattern.ElementsAs(ctx, &userNamePatternPayload, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var externalUserMappingPayload []rolestore.EUM
	for _, eum := range data.ExternalUserMapping {
		externalUserMappingPayload = append(externalUserMappingPayload,
			rolestore.EUM{SourceID: eum.SourceID.ValueString(), SourceSeaerchField: eum.SourceSeaerchField.ValueString()},
		)
	}

	OIDCAdditionalScopesSecretPayload := make([]string, len(data.OIDCConnection.ScopesSecret.Elements()))
	resp.Diagnostics.Append(data.OIDCConnection.ScopesSecret.ElementsAs(ctx, &OIDCAdditionalScopesSecretPayload, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connectionPayload := rolestore.Connection{
		Type:                  "OIDC",
		Address:               data.OIDCConnection.Address.ValueString(),
		OIDCEnabled:           data.OIDCConnection.Enabled.ValueBool(),
		OIDCIssuer:            data.OIDCConnection.Issuer.ValueString(),
		OIDCButtonTitle:       data.OIDCConnection.ButtonTitle.ValueString(),
		OIDCClientID:          data.OIDCConnection.ClientID.ValueString(),
		OIDCClientSecret:      data.OIDCConnection.ClientSecret.ValueString(),
		OIDCTagsAttributeName: data.OIDCConnection.TagsAttributeName.ValueString(),
		OIDCScopesSecret:      OIDCAdditionalScopesSecretPayload,
	}

	source := rolestore.Source{
		Name:                data.Name.ValueString(),
		Comment:             data.Comment.ValueString(),
		TTL:                 int(data.TTL.ValueInt64()),
		Enabled:             data.Enabled.ValueBool(),
		Tags:                tagsPayload,
		UsernamePattern:     userNamePatternPayload,
		ExternalUserMapping: externalUserMappingPayload,
		Connection:          connectionPayload,
	}

	sourceID, err := r.client.CreateSource(source)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while attempting to create the resource. "+
				"Please retry the operation or report this issue to the provider developers.\n\n"+
				err.Error(),
		)
		return
	}

	// Convert from the API data model to the Terraform data model
	// and set any unknown attribute values.
	data.ID = types.StringValue(sourceID)

	tflog.Info(ctx, fmt.Sprintf("data stored: %+v", data))

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SourceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the source object from PrivX API
	source, err := r.client.Source(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read source, got error: %s", err))
		return
	}

	tags, diags := types.ListValueFrom(ctx, data.Tags.ElementType(ctx), source.Tags)
	if diags.HasError() {
		return
	}

	usernamePattern, diags := types.ListValueFrom(ctx, data.UsernamePattern.ElementType(ctx), source.UsernamePattern)
	if diags.HasError() {
		return
	}

	var eum []*EUMModel
	for _, v := range source.ExternalUserMapping {
		eum = append(eum, &EUMModel{
			types.StringValue(v.SourceID),
			types.StringValue(v.SourceSeaerchField)})
	}

	scopesSecret, diags := types.ListValueFrom(ctx, data.OIDCConnection.ScopesSecret.ElementType(ctx), source.Connection.OIDCScopesSecret)
	if diags.HasError() {
		return
	}

	connection := &OIDCConnectionModel{
		Address:           types.StringValue(source.Connection.Address),
		Enabled:           types.BoolValue(source.Connection.OIDCEnabled),
		ButtonTitle:       types.StringValue(source.Connection.OIDCButtonTitle),
		Issuer:            types.StringValue(source.Connection.OIDCIssuer),
		ClientID:          types.StringValue(source.Connection.OIDCClientID),
		ClientSecret:      data.OIDCConnection.ClientSecret, // Do not update client_secret. We keep the state value since PrivX returns "*****" as password.
		TagsAttributeName: types.StringValue(source.Connection.OIDCTagsAttributeName),
		ScopesSecret:      scopesSecret,
	}

	data.Name = types.StringValue(source.Name)
	data.Enabled = types.BoolValue(source.Enabled)
	data.TTL = types.Int64Value(int64(source.TTL))
	data.Comment = types.StringValue(source.Comment)
	data.Tags = tags
	data.UsernamePattern = usernamePattern
	data.ExternalUserMapping = eum
	data.OIDCConnection = connection

	tflog.Info(ctx, fmt.Sprintf("data stored: %+v", data))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SourceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tagsPayload := make([]string, len(data.Tags.Elements()))
	resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tagsPayload, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	userNamePatternPayload := make([]string, len(data.UsernamePattern.Elements()))
	resp.Diagnostics.Append(data.UsernamePattern.ElementsAs(ctx, &userNamePatternPayload, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var externalUserMappingPayload []rolestore.EUM
	for _, eum := range data.ExternalUserMapping {
		externalUserMappingPayload = append(
			externalUserMappingPayload,
			rolestore.EUM{SourceID: eum.SourceID.ValueString(), SourceSeaerchField: eum.SourceSeaerchField.ValueString()},
		)
	}

	OIDCAdditionalScopesSecretPayload := make([]string, len(data.OIDCConnection.ScopesSecret.Elements()))
	resp.Diagnostics.Append(data.OIDCConnection.ScopesSecret.ElementsAs(ctx, &OIDCAdditionalScopesSecretPayload, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	connectionPayload := rolestore.Connection{
		Type:                  "OIDC",
		Address:               data.OIDCConnection.Address.ValueString(),
		OIDCEnabled:           data.OIDCConnection.Enabled.ValueBool(),
		OIDCIssuer:            data.OIDCConnection.Issuer.ValueString(),
		OIDCButtonTitle:       data.OIDCConnection.ButtonTitle.ValueString(),
		OIDCClientID:          data.OIDCConnection.ClientID.ValueString(),
		OIDCClientSecret:      data.OIDCConnection.ClientSecret.ValueString(),
		OIDCTagsAttributeName: data.OIDCConnection.TagsAttributeName.ValueString(),
		OIDCScopesSecret:      OIDCAdditionalScopesSecretPayload,
	}

	source := rolestore.Source{
		ID:                  data.ID.ValueString(),
		Enabled:             data.Enabled.ValueBool(),
		TTL:                 int(data.TTL.ValueInt64()),
		Name:                data.Name.ValueString(),
		Comment:             data.Comment.ValueString(),
		Tags:                tagsPayload,
		UsernamePattern:     userNamePatternPayload,
		ExternalUserMapping: externalUserMappingPayload,
		Connection:          connectionPayload,
	}

	// Update source with the API
	err := r.client.UpdateSource(data.ID.ValueString(), &source)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update source, got error: %s", err))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("data stored: %+v", data))

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SourceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSource(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete source, got error: %s", err))
		return
	}
}

func (r *SourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
