package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/vault"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SecretDataSource{}

func NewSecretDataSource() datasource.DataSource {
	return &SecretDataSource{}
}

// SecretDataSource defines the data source implementation.
type SecretDataSource struct {
	client *vault.Vault
}

// SecretDataSourceModel describes the data source data model.
type SecretDataSourceModel struct {
	Name       types.String   `tfsdk:"name"`
	Data       types.String   `tfsdk:"data"`
	Author     types.String   `tfsdk:"author"`
	UpdatedBy  types.String   `tfsdk:"updated_by"`
	Created    types.String   `tfsdk:"created"`
	Updated    types.String   `tfsdk:"updated"`
	AllowRead  []RoleRefModel `tfsdk:"read_roles"`
	AllowWrite []RoleRefModel `tfsdk:"write_roles"`
}

func (d *SecretDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_secret"
}

func (d *SecretDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				MarkdownDescription: "Secret's name",
				Required:            true,
			},
			"data": schema.StringAttribute{
				MarkdownDescription: "Secret to be stored",
				Computed:            true,
				Sensitive:           true,
			},
			"author": schema.StringAttribute{
				MarkdownDescription: "ID of secret's author",
				Computed:            true,
			},
			"updated_by": schema.StringAttribute{
				MarkdownDescription: "ID of last user to update secret",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "Creation time",
				Computed:            true,
			},
			"updated": schema.StringAttribute{
				MarkdownDescription: "Update time",
				Computed:            true,
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

func (d *SecretDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	connector, ok := req.ProviderData.(*restapi.Connector)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *restapi.Connector, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	tflog.Debug(ctx, "Creating vault client", map[string]interface{}{
		"connector : ": fmt.Sprintf("%+v", *connector),
	})

	d.client = vault.New(*connector)
}

func (d *SecretDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SecretDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	secret, err := d.client.Secret(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read secret, got error: %s", err))
		return
	}

	var allowedRead []RoleRefModel
	for _, v := range secret.AllowRead {
		allowedRead = append(allowedRead, RoleRefModel{types.StringValue(v.ID), types.StringValue(v.Name)})
	}
	data.AllowRead = allowedRead

	var allowedWrite []RoleRefModel
	for _, v := range secret.AllowWrite {
		allowedWrite = append(allowedWrite, RoleRefModel{types.StringValue(v.ID), types.StringValue(v.Name)})
	}
	data.AllowWrite = allowedWrite

	secretData, err := json.Marshal(secret.Data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Data Source",
			"Cannot marshal SourceRule data to json.\n"+
				err.Error(),
		)
		return
	}
	data.Data = types.StringValue(string(secretData))

	data.Author = types.StringValue(secret.Author)
	data.Created = types.StringValue(secret.Created)
	data.Updated = types.StringValue(secret.Updated)
	data.UpdatedBy = types.StringValue(secret.Editor)

	tflog.Debug(ctx, "Storing secret type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", data),
	})
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
