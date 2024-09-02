package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &RoleDataSource{}

func NewRoleDataSource() datasource.DataSource {
	return &RoleDataSource{}
}

// RoleDataSource defines the data source implementation.
type RoleDataSource struct {
	client *rolestore.RoleStore
}

// RoleDataSourceModel describes the data source data model.
type RoleDataSourceModel RoleResourceModel

func (d *RoleDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

func (d *RoleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Role data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Role ID",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the role",
				Required:            true,
			},
			"access_group_id": schema.StringAttribute{
				MarkdownDescription: "Scopes host and connection permissions to an access group. (Defaults to Default access group)",
				Computed:            true,
			},
			"comment": schema.StringAttribute{
				MarkdownDescription: "A comment describing the object",
				Computed:            true,
			},
			"permissions": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Role permissions",
				Computed:            true,
			},
			"principal_public_key_strings": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "List of role's principal public keys",
				Computed:            true,
			},
			"permit_agent": schema.BoolAttribute{
				MarkdownDescription: "Role permit agent",
				Computed:            true,
			},
			"source_rules": schema.StringAttribute{
				MarkdownDescription: `A source rule(s) definition. Can be a single rule or a rule group, in which case either "single" or "group" attributes are requrired. Defined in JSON`,
				Computed:            true,
			},
		},
	}
}

func (d *RoleDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = rolestore.New(*connector)
}

func (d *RoleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RoleDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// resolve role from role name
	role_names_as_list := []string{data.Name.ValueString()}
	roles, err := d.client.ResolveRoles(role_names_as_list)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to resolves roles, got error: %s", err))
		return
	}

	if len(roles) == 0 {
		resp.Diagnostics.AddError("ResolveRoles Error", fmt.Sprintf("Could not retrieve a role from name: %s", data.Name.ValueString()))
		return
	}

	// retrieve role from id
	role, err := d.client.Role(roles[0].ID)

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
	data.SourceRule = types.StringValue(string(sourceRuleData))

	tflog.Debug(ctx, "Storing role type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", data),
	})
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
