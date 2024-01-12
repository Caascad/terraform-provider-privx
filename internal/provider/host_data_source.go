package provider

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/hoststore"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &HostDataSource{}

func NewHostDataSource() datasource.DataSource {
	return &HostDataSource{}
}

// HostDataSource defines the data source implementation.
type HostDataSource struct {
	client *hoststore.HostStore
}

// HostDataSourceModel describes the data source data model.
type (
	StatusModel struct {
		K types.String `tfsdk:"k"`
		V types.String `tfsdk:"v"`
	}
	ApplicationDataSourceModel struct {
		Name types.String `tfsdk:"name"`
		/* FIXME: Not implemented in privx-sdk-go v1.29.0
		Application      types.String `tfsdk:"application"`
		Arguments        types.String `tfsdk:"arguments"`
		WorkingDirectory types.String `tfsdk:"working_directory"`
		*/
	}

	PrincipalDataSourceModel struct {
		ID             types.String                 `tfsdk:"principal"`
		Passphrase     types.String                 `tfsdk:"passphrase"`
		Source         types.String                 `tfsdk:"Source"`
		UseUserAccount types.Bool                   `tfsdk:"use_user_account"`
		Roles          []RoleRefModel               `tfsdk:"roles"`
		Applications   []ApplicationDataSourceModel `tfsdk:"applications"`

		/* FIXME: Not implemented in privx-sdk-go v1.29.0
		Rotate                 types.Bool               `tfsdk:"rotate"`
		UseForPasswordRotation types.Bool               `tfsdk:"use_for_password_rotation"`
		ServiceOptions         ServiceOptionsModel      `tfsdk:"service_options"`
		CommandRestrictions    CommandRestrictionsModel `tfsdk:"command_restrictions"`
		*/
	}

	HostDataSourceModel struct {
		ID                  types.String               `tfsdk:"id"`
		AccessGroupID       types.String               `tfsdk:"access_group_id"`
		ExternalID          types.String               `tfsdk:"external_id"`
		InstanceID          types.String               `tfsdk:"instance_id"`
		SourceID            types.String               `tfsdk:"source_id"`
		Name                types.String               `tfsdk:"common_name"`
		ContactAddress      types.String               `tfsdk:"contact_address"`
		CloudProvider       types.String               `tfsdk:"cloud_provider"`
		CloudProviderRegion types.String               `tfsdk:"cloud_provider_region"`
		DistinguishedName   types.String               `tfsdk:"distinguished_name"`
		Organization        types.String               `tfsdk:"organization"`
		OrganizationUnit    types.String               `tfsdk:"organizational_unit"`
		Zone                types.String               `tfsdk:"zone"`
		HostType            types.String               `tfsdk:"host_type"`
		HostClassification  types.String               `tfsdk:"host_classification"`
		Comment             types.String               `tfsdk:"comment"`
		Disabled            types.String               `tfsdk:"disabled"`
		Deployable          types.Bool                 `tfsdk:"deployable"`
		Tofu                types.Bool                 `tfsdk:"tofu"`
		StandAlone          types.Bool                 `tfsdk:"stand_alone_host"`
		Audit               types.Bool                 `tfsdk:"audit_enabled"`
		Scope               types.Set                  `tfsdk:"scope"`
		Tags                types.Set                  `tfsdk:"tags"`
		Addresses           types.Set                  `tfsdk:"addresses"`
		Services            []ServiceModel             `tfsdk:"services"`
		Principals          []PrincipalDataSourceModel `tfsdk:"principals"`
		PublicKeys          []SSHPublicKeyModel        `tfsdk:"ssh_host_public_keys"`

		Created   types.String  `tfsdk:"created"`
		Updated   types.String  `tfsdk:"updated"`
		UpdatedBy types.String  `tfsdk:"updated_by"`
		Status    []StatusModel `tfsdk:"status"`

		/* FIXME: Not implemented in privx-sdk-go v1.29.0
		CertificateTemplate     types.String          `tfsdk:"certificate_template"`
		HostCertificateRaw      types.String          `tfsdk:"host_certificate_raw"`
		PasswordRotationEnabled types.Bool            `tfsdk:"password_rotation_enabled"`
		PasswordRotation        PasswordRotationModel `tfsdk:"password_rotation"`
		*/
	}
)

func (d *HostDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (d *HostDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example data source",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Host ID",
				Computed:            true,
			},
			"access_group_id": schema.StringAttribute{
				MarkdownDescription: "Defines host's access group",
				Computed:            true,
			},
			"external_id": schema.StringAttribute{
				MarkdownDescription: "The equipment ID from the originating equipment store",
				Computed:            true,
			},
			"instance_id": schema.StringAttribute{
				MarkdownDescription: "The instance ID from the originating cloud service (searchable by keyword)",
				Computed:            true,
			},
			"source_id": schema.StringAttribute{
				MarkdownDescription: "A unique import-source identifier for the host entry, for example a hash for AWS account ID. (searchable by keyword)",
				Computed:            true,
			},
			"common_name": schema.StringAttribute{
				MarkdownDescription: "X.500 Common name (searchable by keyword)",
				Computed:            true,
			},
			"created": schema.StringAttribute{
				MarkdownDescription: "When the object was created",
				Computed:            true,
			},
			"updated": schema.StringAttribute{
				MarkdownDescription: "When the object was updated",
				Computed:            true,
			},
			"updated_by": schema.StringAttribute{
				MarkdownDescription: "Id of the user who updated the object",
				Computed:            true,
			},
			"contact_address": schema.StringAttribute{
				MarkdownDescription: "The host public address scanning script instructs the host store to use in service address-field.",
				Computed:            true,
			},
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "The cloud provider the host resides in",
				Computed:            true,
			},
			"cloud_provider_region": schema.StringAttribute{
				MarkdownDescription: "The cloud provider region the host resides in",
				Computed:            true,
			},
			"distinguished_name": schema.StringAttribute{
				MarkdownDescription: "LDAPv3 Disinguished name (searchable by keyword)",
				Computed:            true,
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "X.500 Organization (searchable by keyword)",
				Computed:            true,
			},
			"organizational_unit": schema.StringAttribute{
				MarkdownDescription: "X.500 Organizational unit (searchable by keyword)",
				Computed:            true,
			},
			"zone": schema.StringAttribute{
				MarkdownDescription: "Equipment zone (development, production, user acceptance testing, ..) (searchable by keyword)",
				Computed:            true,
			},
			"host_type": schema.StringAttribute{
				MarkdownDescription: "Equipment type (virtual, physical) (searchable by keyword)",
				Computed:            true,
			},
			"host_classification": schema.StringAttribute{
				MarkdownDescription: "Classification (Windows desktop, Windows server, AIX, Linux RH, ..) (searchable by keyword)",
				Computed:            true,
			},
			"comment": schema.StringAttribute{
				MarkdownDescription: "A comment describing the host",
				Computed:            true,
			},
			/* FIXME: Not implemented in privx-sdk-go v1.29.0
			"host_certificate_raw": schema.StringAttribute{
				MarkdownDescription: "Host certificate, used to verify that the target host is the correct one.",
				Computed:            true,
			},
			*/
			"disabled": schema.StringAttribute{
				MarkdownDescription: `disabled ("BY_ADMIN" | "BY_LISCENCE" | "false")`,
				Computed:            true,
			},
			"deployable": schema.BoolAttribute{
				MarkdownDescription: "Whether the host is writable through /deploy end point with deployment credentials",
				Computed:            true,
			},
			"tofu": schema.BoolAttribute{
				MarkdownDescription: "Whether the host key should be accepted and stored on first connection",
				Computed:            true,
			},
			"stand_alone_host": schema.BoolAttribute{
				MarkdownDescription: "Indicates it is a standalone host - bound to local host directory",
				Computed:            true,
			},
			"audit_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the host is set to be audited",
				Computed:            true,
			},
			"scope": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Under what compliance scopes the listed equipment falls under (searchable by keyword)",
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Host tags",
				Computed:            true,
			},
			"addresses": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Host addresses",
				Computed:            true,
			},
			/* FIXME: Not implemented in privx-sdk-go v1.29.0
			"certificate_template": schema.StringAttribute{
				MarkdownDescription: "Name of the certificate template used for certificate authentication for this host",
				Computed:            true,
			},
			*/
			"status": schema.SetNestedAttribute{
				MarkdownDescription: "Status",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"k": schema.StringAttribute{
							MarkdownDescription: "k",
							Computed:            true,
						},
						"v": schema.StringAttribute{
							MarkdownDescription: "v",
							Computed:            true,
						},
					},
				},
			},
			"ssh_host_public_keys": schema.SetNestedAttribute{
				MarkdownDescription: "Host public keys, used to verify the identity of the accessed host",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "Host public key, used to verify the identity of the accessed host",
							Computed:            true,
						},
					},
				},
			},
			"services": schema.SetNestedAttribute{
				MarkdownDescription: "Host services",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"service": schema.StringAttribute{
							MarkdownDescription: "Allowed protocol - SSH, RDP, VNC, HTTP, HTTPS (searchable)",
							Computed:            true,
						},
						"address": schema.StringAttribute{
							MarkdownDescription: "Service address, IPv4, IPv6 or FQDN",
							Computed:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "Service port",
							Computed:            true,
						},
						/*
							"use_for_password_rotation": schema.BoolAttribute{
								MarkdownDescription: "if service SSH, informs whether this service is used to rotate password",
								Computed:            true,
							},
						*/
					},
				},
			},
			"principals": schema.SetNestedAttribute{
				MarkdownDescription: "What principals (target server user names/ accounts) the host has",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"principal": schema.StringAttribute{
							MarkdownDescription: "The account name",
							Computed:            true,
						},
						/* FIXME: Not implemented in privx-sdk-go v1.29.0
						"rotate": schema.BoolAttribute{
							MarkdownDescription: "Rotate password of this account",
							Computed:            true,
						},
						"use_for_password_rotation": schema.BoolAttribute{
							MarkdownDescription: "marks account to be used as the account through which password rotation takes place, when flag use_main_account set in rotation_metadata",
							Computed:            true,
						},
						*/
						"use_user_account": schema.BoolAttribute{
							MarkdownDescription: "Use user account as host principal name",
							Computed:            true,
						},
						"passphrase": schema.StringAttribute{
							MarkdownDescription: "The account static passphrase or the initial rotating password value. If rotate selected, active in create, disabled/hidden in edit",
							Sensitive:           true,
							Computed:            true,
						},
						"source": schema.StringAttribute{
							MarkdownDescription: `Identifies the source of the principals object "UI" or "SCAN". Deploy is also treated as "UI"`,
							Computed:            true,
						},
						"roles": schema.SetNestedAttribute{
							MarkdownDescription: "An array of roles entitled to access this principal on the host",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: "Role UUID",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										MarkdownDescription: "Role UUID",
										Computed:            true,
									},
								},
							},
						},
						"applications": schema.SetNestedAttribute{
							MarkdownDescription: "An array of application the principal may launch on the target host",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed: true,
									},
									/* FIXME: Not implemented in privx-sdk-go v1.29.0
									"application": schema.StringAttribute{
										Computed: true,
									},
									"arguments": schema.StringAttribute{
										Computed: true,
									},
									"working_directory": schema.StringAttribute{
										Computed: true,
									},
									*/
								},
							},
						},
						/* FIXME: Not implemented in privx-sdk-go v1.29.0
						"service_options": schema.SingleNestedAttribute{
							MarkdownDescription: "Object for service options",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"ssh": schema.SingleNestedAttribute{
									MarkdownDescription: "SSH service options",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"shell": schema.BoolAttribute{
											MarkdownDescription: "Shell channel",
											Computed:            true,
										},
										"file_transfer": schema.BoolAttribute{
											MarkdownDescription: "File transfer channel",
											Computed:            true,
										},
										"exec": schema.BoolAttribute{
											MarkdownDescription: "exec_channel",
											Computed:            true,
										},
										"tunnels": schema.BoolAttribute{
											MarkdownDescription: "tunnels",
											Computed:            true,
										},
										"x11": schema.BoolAttribute{
											MarkdownDescription: "x11",
											Computed:            true,
										},
										"other": schema.BoolAttribute{
											MarkdownDescription: "other options",
											Computed:            true,
										},
									},
								},
								"rdp": schema.SingleNestedAttribute{
									MarkdownDescription: "SSH service options",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"file_transfer": schema.BoolAttribute{
											MarkdownDescription: "File transfer channel",
											Computed:            true,
										},
										"audio": schema.BoolAttribute{
											MarkdownDescription: "audio",
											Computed:            true,
										},
										"clipboard": schema.BoolAttribute{
											MarkdownDescription: "clipboard",
											Computed:            true,
										},
										"web": schema.BoolAttribute{
											MarkdownDescription: "WEB service options",
											Computed:            true,
										},
									},
								},
								"web": schema.SingleNestedAttribute{
									MarkdownDescription: "SSH service options",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"file_transfer": schema.BoolAttribute{
											MarkdownDescription: "File transfer channel",
											Computed:            true,
										},
										"audio": schema.BoolAttribute{
											MarkdownDescription: "audio",
											Computed:            true,
										},
										"clipboard": schema.BoolAttribute{
											MarkdownDescription: "clipboard",
											Computed:            true,
										},
									},
								},
							},
						},
						"command_restrictions": schema.SingleNestedAttribute{
							MarkdownDescription: "Host services",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									MarkdownDescription: "Are command restrictions enabled",
									Computed:            true,
								},
								"default_whitelist": schema.SingleNestedAttribute{
									MarkdownDescription: "Default whitelist handle, required if command restrictions are enabled",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											MarkdownDescription: "Whitelist ID",
											Computed:            true,
										},
										"name": schema.StringAttribute{
											MarkdownDescription: "Whitelist name",
											Computed:            true,
										},
										"deleted": schema.BoolAttribute{
											MarkdownDescription: "Has whitelist been deleted, ignored in requests",
											Computed:            true,
										},
									},
								},
								"rshell_variant": schema.StringAttribute{
									MarkdownDescription: "Restricted shell variant, required if command restrictions are enabled",
									Computed:            true,
								},
								"banner": schema.StringAttribute{
									MarkdownDescription: "Computed banner displayed in SSH terminal",
									Computed:            true,
								},
								"allow_no_match": schema.BoolAttribute{
									MarkdownDescription: "If true then commands that do not match any whitelist pattern are allowed to execute",
									Computed:            true,
								},
								"audit_match": schema.BoolAttribute{
									MarkdownDescription: "If true then an audit event is generated for every allowed command",
									Computed:            true,
								},
								"audit_no_match": schema.BoolAttribute{
									MarkdownDescription: "If true then an audit event is generated for every disallowed command",
									Computed:            true,
								},
								"whitelists": schema.SetNestedAttribute{
									Computed: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"whitelist": schema.SingleNestedAttribute{
												Computed: true,
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														MarkdownDescription: "Whitelist ID",
														Computed:            true,
													},
													"name": schema.StringAttribute{
														MarkdownDescription: "Whitelist name",
														Computed:            true,
													},

													"deleted": schema.BoolAttribute{
														MarkdownDescription: "Has whitelist been deleted, ignored in requests",
														Computed:            true,
													},
												},
											},
											"roles": schema.SetNestedAttribute{
												MarkdownDescription: "List of roles granting access to the whitelist",
												Computed:            true,

												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															MarkdownDescription: "Role ID",
															Computed:            true,
														},
														"name": schema.StringAttribute{
															MarkdownDescription: "Role Name",
															Computed:            true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
						*/
					},
				},
			},
			/* FIXME: Not implemented in privx-sdk-go v1.29.0
			"password_rotation_enabled": schema.BoolAttribute{
				MarkdownDescription: "set, if there are accounts, in which passwords need to be rotated",
				Computed:            true,
			},
			*/
			/* FIXME: Not implemented in privx-sdk-go v1.29.0
			"password_rotation": schema.SingleNestedAttribute{
				MarkdownDescription: "password rotation settings for host",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"use_main_account": schema.BoolAttribute{
						MarkdownDescription: "rotate passwords of all accounts in host through one account",
						Computed:            true,
					},
					"operating_system": schema.StringAttribute{
						MarkdownDescription: "Bash for Linux, Powershell for windows for shell access (LINUX | WINDOWS)",
						Computed:            true,
					},
					"winrm_address": schema.StringAttribute{
						MarkdownDescription: "IPv4 address or FQDN to use for winrm connection",
						Computed:            true,
					},
					"winrm_port": schema.Int64Attribute{
						MarkdownDescription: "port to use for password rotation with winrm, zero for winrm default",
						Computed:            true,
					},
					"protocol": schema.StringAttribute{
						MarkdownDescription: "protocol (SSH | WINRM)",
						Computed:            true,
					},
					"password_policy_id": schema.StringAttribute{
						MarkdownDescription: "password policy to be applied",
						Computed:            true,
					},
					"script_template_id": schema.StringAttribute{
						MarkdownDescription: "script template to be run in host",
						Computed:            true,
					},
				},
			},
			*/
		},
	}
}

func (d *HostDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	tflog.Debug(ctx, "Creating HostStore client", map[string]interface{}{
		"connector : ": fmt.Sprintf("%+v", *connector),
	})

	d.client = hoststore.New(*connector)
}

func (d *HostDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data HostDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	host, err := d.client.Host(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host, got error: %s", err))
		return
	}

	data.AccessGroupID = types.StringValue(host.AccessGroupID)
	data.ExternalID = types.StringValue(host.ExternalID)
	data.InstanceID = types.StringValue(host.InstanceID)
	data.SourceID = types.StringValue(host.SourceID)
	data.Name = types.StringValue(host.Name)
	data.Created = types.StringValue(host.Created)
	data.Updated = types.StringValue(host.Updated)
	data.UpdatedBy = types.StringValue(host.UpdatedBy)
	data.ContactAddress = types.StringValue(host.ContactAdress)
	data.CloudProvider = types.StringValue(host.CloudProvider)
	data.CloudProviderRegion = types.StringValue(host.CloudProviderRegion)
	data.DistinguishedName = types.StringValue(host.DistinguishedName)
	data.Organization = types.StringValue(host.Organization)
	data.OrganizationUnit = types.StringValue(host.OrganizationUnit)
	data.Zone = types.StringValue(host.Zone)
	data.HostType = types.StringValue(host.HostType)
	data.HostClassification = types.StringValue(host.HostClassification)
	data.Comment = types.StringValue(host.Comment)
	data.Disabled = types.StringValue(host.Disabled)
	data.Deployable = types.BoolValue(host.Deployable)
	data.Tofu = types.BoolValue(host.Tofu)
	data.StandAlone = types.BoolValue(host.Tofu)
	data.Audit = types.BoolValue(host.Audit)

	scope, diags := types.SetValueFrom(ctx, data.Scope.ElementType(ctx), host.Scope)
	if diags.HasError() {
		return
	}
	data.Scope = scope

	tags, diags := types.SetValueFrom(ctx, data.Tags.ElementType(ctx), host.Tags)
	if diags.HasError() {
		return
	}
	data.Tags = tags

	addresses, diags := types.SetValueFrom(ctx, data.Addresses.ElementType(ctx), host.Addresses)
	if diags.HasError() {
		return
	}
	data.Addresses = addresses

	var status []StatusModel
	for _, st := range host.Status {
		status = append(status, StatusModel{
			K: types.StringValue(st.K),
			V: types.StringValue(st.V),
		})
	}
	data.Status = status

	var services []ServiceModel
	for _, s := range host.Services {
		services = append(services, ServiceModel{
			Scheme:  types.StringValue(string(s.Scheme)),
			Address: types.StringValue(string(s.Address)),
			Port:    types.Int64Value(int64(s.Port)),
			// UseForPasswordRotation: types.StringValue(s.UseForPasswordRotation), // FIXME: Not implemented in privx-sdk-go v1.29.0
		})
	}
	data.Services = services

	var principals []PrincipalDataSourceModel
	for _, p := range host.Principals {
		var roles []RoleRefModel
		for _, r := range p.Roles {
			roles = append(roles, RoleRefModel{
				ID: types.StringValue(r.ID),
			})
		}
		var applications []ApplicationDataSourceModel
		for _, a := range p.Applications {
			applications = append(applications, ApplicationDataSourceModel{
				Name: types.StringValue(a),
			})
		}
		principals = append(principals, PrincipalDataSourceModel{
			Passphrase:     types.StringValue(p.Passphrase),
			Source:         types.StringValue(string(p.Source)),
			UseUserAccount: types.BoolValue(p.UseUserAccount),
			// Rotate:     types.BoolValue(p.Rotate), // FIXME: Not implemented in privx-sdk-go v1.29.0
			// UseForPasswordRotation:     types.BoolValue(p.UseForPasswordRotation), // FIXME: Not implemented in privx-sdk-go v1.29.0
			// ServiceOptions: serviceOptions, // FIXME: Not implemented in privx-sdk-go v1.29.0
			// CommandRestrictions: commandRestrictions, // FIXME: Not implemented in privx-sdk-go v1.29.0
			Roles:        roles,
			Applications: applications,
		})
	}
	data.Principals = principals

	var publickeys []SSHPublicKeyModel
	for _, pb := range host.PublicKeys {
		publickeys = append(publickeys, SSHPublicKeyModel{
			Key: types.StringValue(pb.Key),
		})
	}
	data.PublicKeys = publickeys

	tflog.Debug(ctx, "Storing host type into the state", map[string]interface{}{
		"createNewState": fmt.Sprintf("%+v", data),
	})
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
