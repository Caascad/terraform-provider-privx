package provider

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
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
var _ resource.Resource = &APIClientResource{}
var _ resource.ResourceWithImportState = &APIClientResource{}

func NewAPIClientResource() resource.Resource {
	return &APIClientResource{}
}

// APIClientResource defines the resource implementation.
type APIClientResource struct {
	client *userstore.UserStore
}

// APIClientModel describes the resource data model.
type APIClientModel struct {
	ID                types.String   `tfsdk:"id"`
	Name              types.String   `tfsdk:"name"`
	Secret            types.String   `tfsdk:"secret"`
	OauthClientId     types.String   `tfsdk:"oauth_client_id"`
	OauthClientSecret types.String   `tfsdk:"oauth_client_secret"`
	Roles             []RoleRefModel `tfsdk:"roles"`
}

func (r *APIClientResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_client"
}

func (r *APIClientResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "API client resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "ID of the API client",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "name of the API client",
				Required:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "secret of the API client",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"oauth_client_id": schema.StringAttribute{
				MarkdownDescription: "oauth_client_id of the API client",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"oauth_client_secret": schema.StringAttribute{
				MarkdownDescription: "oauth_client_secret of the API client",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"roles": schema.SetNestedAttribute{
				MarkdownDescription: "List of roles possessed by the API client",
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

func (r *APIClientResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *APIClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *APIClientModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var rolesPayload []string
	for _, roleRef := range data.Roles {
		rolesPayload = append(rolesPayload, roleRef.ID.ValueString())
	}

	id, err := r.client.CreateAPIClient(data.Name.ValueString(), rolesPayload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create API client Resource",
			"An unexpected error occurred while attempting to create the resource.\n"+
				err.Error(),
		)
		return
	}

	api_client, err := r.client.APIClient(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read api_client, got error: %s", err))
		return
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.ID = types.StringValue(id)
	data.Secret = types.StringValue(api_client.Secret)
	data.OauthClientId = types.StringValue(api_client.AuthClientID)
	data.OauthClientSecret = types.StringValue(api_client.AuthClientSecret)
	var roles []RoleRefModel
	for _, role := range api_client.Roles {
		roles = append(roles,
			RoleRefModel{ID: types.StringValue(role.ID),
				Name: types.StringValue(role.Name),
			})
	}
	data.Roles = roles

	ctx = tflog.SetField(ctx, "API client name", data.Name.ValueString())
	ctx = tflog.SetField(ctx, "API client roles", data.Roles)
	tflog.Debug(ctx, "Created API client")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *APIClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *APIClientModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := r.client.APIClient(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read API client, got error: %s", err))
		return
	}

	var roles []RoleRefModel
	for _, role := range apiClient.Roles {
		roles = append(roles,
			RoleRefModel{ID: types.StringValue(role.ID),
				Name: types.StringValue(role.Name),
			})
	}
	data.Roles = roles
	data.Name = types.StringValue(apiClient.Name)
	data.Secret = types.StringValue(apiClient.Secret)
	data.OauthClientId = types.StringValue(apiClient.AuthClientID)
	data.OauthClientSecret = types.StringValue(apiClient.AuthClientSecret)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *APIClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *APIClientModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var rolesPayload []rolestore.RoleRef
	for _, roleRef := range data.Roles {
		rolesPayload = append(rolesPayload,
			rolestore.RoleRef{ID: roleRef.ID.ValueString(), Name: roleRef.Name.ValueString()})
	}

	apiClientPayload := userstore.APIClient{
		ID:               data.ID.ValueString(),
		Name:             data.Name.ValueString(),
		Secret:           data.Secret.ValueString(),
		AuthClientID:     data.OauthClientId.ValueString(),
		AuthClientSecret: data.OauthClientSecret.ValueString(),
		Roles:            rolesPayload,
	}

	if err := r.client.UpdateAPIClient(data.ID.ValueString(), &apiClientPayload); err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create API client Resource",
			"An unexpected error occurred while attempting to create the resource.\n"+
				err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *APIClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *APIClientModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteAPIClient(data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete APIClient, got error: %s", err))
		return
	}
}

func (r *APIClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
