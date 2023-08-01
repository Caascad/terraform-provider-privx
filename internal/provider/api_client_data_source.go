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
var _ datasource.DataSource = &APIClientDataSource{}

func NewAPIClientDataSource() datasource.DataSource {
	return &APIClientDataSource{}
}

// APIClientDataSource defines the data source implementation.
type APIClientDataSource struct {
	client *userstore.UserStore
}

// APIClientDataSourceModel describes the data source data model.
type APIClientDataSourceModel struct {
	ID                types.String   `tfsdk:"id"`
	Secret            types.String   `tfsdk:"secret"`
	Name              types.String   `tfsdk:"name"`
	Created           types.String   `tfsdk:"created"`
	Author            types.String   `tfsdk:"author"`
	Roles             []RoleRefModel `tfsdk:"roles"`
	OAuthClientID     types.String   `tfsdk:"oauth_client_id"`
	OAuthClientSecret types.String   `tfsdk:"oauth_client_secret"`
}

func (d *APIClientDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_api_client"
}

func (d *APIClientDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "ID of the API client",
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "name of the API client",
				Computed:            true,
			},
			"secret": schema.StringAttribute{
				MarkdownDescription: "secret of the API client",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "When the object was created",
				Computed:            true,
			},
			"author": schema.StringAttribute{
				MarkdownDescription: "ID of the user who originally authored the object",
				Computed:            true,
			},
			"oauth_client_id": schema.StringAttribute{
				MarkdownDescription: "ID for OAuth2 client, used for authentication",
				Computed:            true,
			},
			"oauth_client_secret": schema.StringAttribute{
				MarkdownDescription: "Secret for OAuth2 client, used for authentication",
				Computed:            true,
			},

			"roles": schema.SetNestedAttribute{
				MarkdownDescription: "List of roles possessed by the API client",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Role ID",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Role name, ignored by server in requests.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *APIClientDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = userstore.New(*connector)
}

func (d *APIClientDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data APIClientDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiClient, err := d.client.APIClient(data.ID.ValueString())
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
	data.Created = types.StringValue(apiClient.Created)
	data.Author = types.StringValue(apiClient.Author)
	data.OAuthClientID = types.StringValue(apiClient.AuthClientID)
	data.OAuthClientSecret = types.StringValue(apiClient.AuthClientSecret)

	tflog.Trace(ctx, "read API Client data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
