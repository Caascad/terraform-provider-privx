package provider

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/authorizer"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &AccessGroupDataSource{}

func NewAccessGroupDataSource() datasource.DataSource {
	return &AccessGroupDataSource{}
}

// AccessGroupDataSource defines the data source implementation.
type AccessGroupDataSource struct {
	client *authorizer.Client
}

// AccessGroupDataSourceModel describes the data source data model.
type AccessGroupDataSourceModel struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	Comment   types.String `tfsdk:"comment"`
	CAID      types.String `tfsdk:"ca_id"`
	Author    types.String `tfsdk:"author"`
	Created   types.String `tfsdk:"created"`
	Updated   types.String `tfsdk:"updated"`
	UpdatedBy types.String `tfsdk:"updated_by"`
	Default   types.Bool   `tfsdk:"default"`
}

func (d *AccessGroupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_group"
}

func (d *AccessGroupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Role data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "UUID",
				Computed:            true,
			},
			"comment": schema.StringAttribute{
				MarkdownDescription: "optional human readable description",
				Optional:            true,
			},
			"ca_id": schema.StringAttribute{
				MarkdownDescription: "UUID of access group's CA",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "When the object was created",
				Computed:            true,
			},
			"updated": schema.StringAttribute{
				MarkdownDescription: "When the object was created",
				Computed:            true,
			},
			"updated_by": schema.StringAttribute{
				MarkdownDescription: "ID of the user who updated the object",
				Computed:            true,
			},
			"author": schema.StringAttribute{
				MarkdownDescription: "ID of the user who originally authored the object",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "unique human reabable name for access group",
				Optional:            true,
			},
			"default": schema.BoolAttribute{
				MarkdownDescription: "Is default access group",
				Optional:            true,
			},
		},
	}
}

func (d *AccessGroupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	tflog.Debug(ctx, "Creating rolestore", map[string]interface{}{
		"connector : ": fmt.Sprintf("%+v", *connector),
	})

	d.client = authorizer.New(*connector)
}

func (d AccessGroupDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.Conflicting(
			path.MatchRoot("name"),
			path.MatchRoot("default"),
		),
	}
}

func (d *AccessGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AccessGroupDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if data.Name.IsNull() && data.Default.IsNull() {
		resp.Diagnostics.AddError("Configuration Error", "Name and Default cannot be null at the same time")
		return
	}
	searchResult, err := d.client.AccessGroups(0, 1000, "id", "ASC")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read access group, got error: %s", err))
		return
	}
	var accessGroup authorizer.AccessGroup
	for _, result := range searchResult {
		if !data.Name.IsNull() {
			if result.Name == data.Name.ValueString() {
				accessGroup = result
				break
			}
		}
		if !data.Default.IsNull() {
			if result.Default == data.Default.ValueBool() {
				accessGroup = result
				break
			}
		}
	}

	data.ID = types.StringValue(accessGroup.ID)
	data.Name = types.StringValue(accessGroup.Name)
	data.Comment = types.StringValue(accessGroup.Comment)
	data.CAID = types.StringValue(accessGroup.CAID)
	data.Author = types.StringValue(accessGroup.Author)
	data.Created = types.StringValue(accessGroup.Created)
	data.Updated = types.StringValue(accessGroup.Updated)
	data.UpdatedBy = types.StringValue(accessGroup.UpdatedBy)
	data.Default = types.BoolValue(accessGroup.Default)

	tflog.Debug(ctx, "Storing access group type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", data),
	})
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
