package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/settings"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &SettingsDataSource{}

func NewSettingsDataSource() datasource.DataSource {
	return &SettingsDataSource{}
}

// SettingsDataSource defines the DataSource implementation.
type SettingsDataSource struct {
	client *settings.Settings
}

// settings contains PrivX settings information.
type SettingsDataSourceModel struct {
	Scope              types.String `tfsdk:"scope"`
	Merge              types.Bool   `tfsdk:"merge"`
	Section            types.String `tfsdk:"section"`
}

func (r *SettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_settings"
}

func (r *SettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Settings DataSource",
		Attributes: map[string]schema.Attribute{
			"scope": schema.StringAttribute{
				MarkdownDescription: "Scope",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf("AUTH","AUTHORIZER","CONNECTION-MANAGER","HOST-STORE","TRAIL-INDEX","KEYVAULT","MONITOR-SERVICE","RDP-PROXY","RDP-MITM","SSH-PROXY","SSH-MITM","ROLE-STORE","USER-STORE","WORKFLOW-ENGINE","SETTINGS","SEARCH","VAULT","LICENSE-MANAGER","SECRETS-MANAGER","NETWORK-ACCESS-MANAGER","EXTENDER-SERVICE","DB-PROXY","PRIVX-CARRIER","PRIVX-WEB-PROXY","PRIVX-EXTENDER","GLOBAL"),
				},

			},
			"merge": schema.BoolAttribute{
				MarkdownDescription: "Merge",
				Optional: true,
			},
			"section": schema.StringAttribute{
				MarkdownDescription: "Section",
				Computed:            true,
			},
		},
	}
}

func (r *SettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	connector, ok := req.ProviderData.(*restapi.Connector)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("Expected *restapi.Connector, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	tflog.Debug(ctx, "Creating settings", map[string]interface{}{
		"connector : ": fmt.Sprintf("%+v", *connector),
	})

	r.client = settings.New(*connector)
}

func (r *SettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *SettingsDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}


	settings, err := r.client.ScopeSettings(data.Scope.ValueString(), "false")
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read settings, got error: %s", err))
		return
	}


	toto, err:= json.Marshal(settings)

	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read settings, got error: %s", toto))

	// data.Scope = types.StringValue(settings.Scope)
	// data.Section = types.StringValue(toto.Section)

	tflog.Debug(ctx, "Storing settings type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", data),
	})
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
