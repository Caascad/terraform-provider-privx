package provider

import (
	"context"
	"fmt"

	"github.com/SSHcom/privx-sdk-go/api/hoststore"
	"github.com/SSHcom/privx-sdk-go/api/rolestore"
	"github.com/SSHcom/privx-sdk-go/restapi"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &HostResource{}
var _ resource.ResourceWithImportState = &HostResource{}

type Address types.String

func NewHostResource() resource.Resource {
	return &HostResource{}
}

type (
	// HostResource defines the resource implementation.
	HostResource struct {
		client *hoststore.HostStore
	}

	ServiceModel struct {
		Scheme  types.String `tfsdk:"service"`
		Address types.String `tfsdk:"address"`
		Port    types.Int64  `tfsdk:"port"`

		/* FIXME: Not implemented in privx-sdk-go v1.29.0
		UseForPasswordRotation types.Bool   `tfsdk:"use_for_password_rotation"`
		*/
	}

	/* FIXME: Not implemented in privx-sdk-go v1.29.0
	ApplicationModel struct {
		Name types.String `tfsdk:"name"`

		Application      types.String `tfsdk:"application"`
		Arguments        types.String `tfsdk:"arguments"`
		WorkingDirectory types.String `tfsdk:"working_directory"`
	}
	*/

	SSHServiceModel struct {
		Shell        types.Bool `tfsdk:"shell"`
		FileTransfer types.Bool `tfsdk:"file_transfer"`
		Exec         types.Bool `tfsdk:"exec"`
		Tunnels      types.Bool `tfsdk:"tunnels"`
		X11          types.Bool `tfsdk:"x11"`
		Other        types.Bool `tfsdk:"other"`
	}

	RDPServiceModel struct {
		FileTransfer types.Bool `tfsdk:"file_transfer"`
		Audio        types.Bool `tfsdk:"audio"`
		Clipboard    types.Bool `tfsdk:"clipboard"`
	}

	ServiceOptionsModel struct {
		SSH SSHServiceModel `tfsdk:"ssh"`
		RDP RDPServiceModel `tfsdk:"rdp"`
		Web RDPServiceModel `tfsdk:"Web"` //RDP and Web models are the same
	}

	WhitelistModel struct {
		ID      types.String `tfsdk:"id"`
		Name    types.String `tfsdk:"name"`
		Deleted types.Bool   `tfsdk:"deleted"`
	}

	/* FIXME: Not implemented in privx-sdk-go v1.29.0
	CommandRestrictionsModel struct {
		RShellVariant    types.String     `tfsdk:"rshell_variant"`
		Banner           types.String     `tfsdk:"banner"`
		Enabled          types.Bool       `tfsdk:"enabled"`
		AllowNoMatch     types.Bool       `tfsdk:"allow_no_match"`
		AuditMatch       types.Bool       `tfsdk:"audit_match"`
		AuditNoMatch     types.Bool       `tfsdk:"audit_no_match"`
		DefalutWhitelist WhitelistModel   `tfsdk:"default_whitelist"`
		Whitelists       []WhitelistModel `tfsdk:"whitelists"`
	}
	*/

	/* FIXME: Not implemented in privx-sdk-go v1.29.0
	PasswordRotationModel struct {
		OperatingSystem  types.String `tfsdk:"operating_system"`
		WINRMAddress     types.String `tfsdk:"winrm_address"`
		WINRMPort        types.String `tfsdk:"winrm_port"`
		Protocol         types.String `tfsdk:"protocol"`
		PasswordPolicyID types.String `tfsdk:"password_policy_id"`
		ScriptTemplateID types.String `tfsdk:"script_template_id"`
		UseMainAccount   types.Bool   `tfsdk:"use_main_account"`
	}
	*/

	RoleRefResourceModel struct {
		ID types.String `tfsdk:"id"`
	}

	// Principal of the target host.
	PrincipalModel struct {
		ID             types.String           `tfsdk:"principal"`
		Passphrase     types.String           `tfsdk:"passphrase"`
		UseUserAccount types.Bool             `tfsdk:"use_user_account"`
		Roles          []RoleRefResourceModel `tfsdk:"roles"`

		/* FIXME: Not implemented in privx-sdk-go v1.29.0
		Applications   []ApplicationModel     `tfsdk:"applications"`
		Rotate                 types.Bool               `tfsdk:"rotate"`
		UseForPasswordRotation types.Bool               `tfsdk:"use_for_password_rotation"`
		ServiceOptions         ServiceOptionsModel      `tfsdk:"service_options"`
		CommandRestrictions    CommandRestrictionsModel `tfsdk:"command_restrictions"`
		*/
	}

	SSHPublicKeyModel struct {
		Key types.String `tfsdk:"key"`
	}

	HostResourceModel struct {
		ID                  types.String        `tfsdk:"id"`
		AccessGroupID       types.String        `tfsdk:"access_group_id"`
		ExternalID          types.String        `tfsdk:"external_id"`
		InstanceID          types.String        `tfsdk:"instance_id"`
		Name                types.String        `tfsdk:"common_name"`
		ContactAddress      types.String        `tfsdk:"contact_address"`
		CloudProvider       types.String        `tfsdk:"cloud_provider"`
		CloudProviderRegion types.String        `tfsdk:"cloud_provider_region"`
		DistinguishedName   types.String        `tfsdk:"distinguished_name"`
		Organization        types.String        `tfsdk:"organization"`
		OrganizationUnit    types.String        `tfsdk:"organizational_unit"`
		Zone                types.String        `tfsdk:"zone"`
		HostType            types.String        `tfsdk:"host_type"`
		HostClassification  types.String        `tfsdk:"host_classification"`
		Comment             types.String        `tfsdk:"comment"`
		Tofu                types.Bool          `tfsdk:"tofu"`
		StandAlone          types.Bool          `tfsdk:"stand_alone_host"`
		Audit               types.Bool          `tfsdk:"audit_enabled"`
		Scope               types.Set           `tfsdk:"scope"`
		Tags                types.Set           `tfsdk:"tags"`
		Addresses           types.Set           `tfsdk:"addresses"`
		Services            []ServiceModel      `tfsdk:"services"`
		Principals          []PrincipalModel    `tfsdk:"principals"`
		PublicKeys          []SSHPublicKeyModel `tfsdk:"ssh_host_public_keys"`

		/* FIXME: Not implemented in privx-sdk-go v1.29.0
		CertificateTemplate     types.String          `tfsdk:"certificate_template"`
		HostCertificateRaw      types.String          `tfsdk:"host_certificate_raw"`
		PasswordRotationEnabled types.Bool            `tfsdk:"password_rotation_enabled"`
		PasswordRotation        PasswordRotationModel `tfsdk:"password_rotation"`
		*/

		/* Set by privx, not needed in resource
		SourceID            types.String        `tfsdk:"source_id"`
		Deployable          types.Bool          `tfsdk:"deployable"`
		*/
	}
)

func (r *HostResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

func (r *HostResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Host resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Host ID",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_group_id": schema.StringAttribute{
				MarkdownDescription: "Defines host's access group",
				Optional:            true,
			},
			"external_id": schema.StringAttribute{
				MarkdownDescription: "The equipment ID from the originating equipment store",
				Optional:            true,
			},
			"instance_id": schema.StringAttribute{
				MarkdownDescription: "The instance ID from the originating cloud service (searchable by keyword)",
				Optional:            true,
			},
			"common_name": schema.StringAttribute{
				MarkdownDescription: "X.500 Common name (searchable by keyword)",
				Optional:            true,
			},
			"contact_address": schema.StringAttribute{
				MarkdownDescription: "The host public address scanning script instructs the host store to use in service address-field.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"cloud_provider": schema.StringAttribute{
				MarkdownDescription: "The cloud provider the host resides in",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"cloud_provider_region": schema.StringAttribute{
				MarkdownDescription: "The cloud provider region the host resides in",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"distinguished_name": schema.StringAttribute{
				MarkdownDescription: "LDAPv3 Disinguished name (searchable by keyword)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"organization": schema.StringAttribute{
				MarkdownDescription: "X.500 Organization (searchable by keyword)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"organizational_unit": schema.StringAttribute{
				MarkdownDescription: "X.500 Organizational unit (searchable by keyword)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"zone": schema.StringAttribute{
				MarkdownDescription: "Equipment zone (development, production, user acceptance testing, ..) (searchable by keyword)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"host_type": schema.StringAttribute{
				MarkdownDescription: "Equipment type (virtual, physical) (searchable by keyword)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"host_classification": schema.StringAttribute{
				MarkdownDescription: "Classification (Windows desktop, Windows server, AIX, Linux RH, ..) (searchable by keyword)",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
			},
			"comment": schema.StringAttribute{
				MarkdownDescription: "A comment describing the host",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString(""),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			/* FIXME: Not implemented in privx-sdk-go v1.29.0
			"host_certificate_raw": schema.StringAttribute{
				MarkdownDescription: "Host certificate, used to verify that the target host is the correct one.",
				Optional:            true,
			},
			*/
			"tofu": schema.BoolAttribute{
				MarkdownDescription: "Whether the host key should be accepted and stored on first connection",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"stand_alone_host": schema.BoolAttribute{
				MarkdownDescription: "Indicates it is a standalone host - bound to local host directory",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"audit_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether the host is set to be audited",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"scope": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Under what compliance scopes the listed equipment falls under (searchable by keyword)",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})), // PrivX API uses empty lists instead of null values. We have to set empty sets
			},
			"tags": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Host tags",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			"addresses": schema.SetAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Host addresses",
				Optional:            true,
				Computed:            true,
				Default:             setdefault.StaticValue(types.SetValueMust(types.StringType, []attr.Value{})),
			},
			/* FIXME: Not implemented in privx-sdk-go v1.29.0
			"certificate_template": schema.StringAttribute{
				MarkdownDescription: "Name of the certificate template used for certificate authentication for this host",
				Optional:            true,
			},
			*/
			"ssh_host_public_keys": schema.SetNestedAttribute{
				MarkdownDescription: "Host public keys, used to verify the identity of the accessed host",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"key": schema.StringAttribute{
							MarkdownDescription: "Host public key, used to verify the identity of the accessed host",
							Optional:            true,
						},
					},
				},
			},
			"services": schema.SetNestedAttribute{
				MarkdownDescription: "Host services",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"service": schema.StringAttribute{
							MarkdownDescription: "Allowed protocol - SSH, RDP, VNC, HTTP, HTTPS (searchable)",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.OneOf("SSH", "RDP", "VNC", "HTTP", "HTTPS"),
							},
						},
						"address": schema.StringAttribute{
							MarkdownDescription: "Service address, IPv4, IPv6 or FQDN",
							Optional:            true,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "Service port",
							Optional:            true,
						},
						/*
							"use_for_password_rotation": schema.BoolAttribute{
								MarkdownDescription: "if service SSH, informs whether this service is used to rotate password",
								Optional:            true,
							},
						*/
					},
				},
			},
			"principals": schema.SetNestedAttribute{
				MarkdownDescription: "What principals (target server user names/ accounts) the host has",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"principal": schema.StringAttribute{
							MarkdownDescription: "The account name",
							Required:            true,
						},
						/* FIXME: Not implemented in privx-sdk-go v1.29.0
						"rotate": schema.BoolAttribute{
							MarkdownDescription: "Rotate password of this account",
							Optional:            true,
						},
						"use_for_password_rotation": schema.BoolAttribute{
							MarkdownDescription: "marks account to be used as the account through which password rotation takes place, when flag use_main_account set in rotation_metadata",
							Optional:            true,
						},
						*/
						"use_user_account": schema.BoolAttribute{
							MarkdownDescription: "Use user account as host principal name",
							Optional:            true,
						},
						"passphrase": schema.StringAttribute{
							MarkdownDescription: "The account static passphrase or the initial rotating password value. If rotate selected, active in create, disabled/hidden in edit",
							Optional:            true,
							Computed:            true,
							Sensitive:           true,
							Default:             stringdefault.StaticString(""),
						},
						"roles": schema.SetNestedAttribute{
							MarkdownDescription: "An array of roles entitled to access this principal on the host",
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: "Role UUID",
										Optional:            true,
									},
								},
							},
						},
						/* FIXME: Not implemented in privx-sdk-go v1.29.0
						"applications": schema.SetNestedAttribute{
							MarkdownDescription: "An array of application the principal may launch on the target host",
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Optional: true,
									},
									"application": schema.StringAttribute{
										Optional: true,
									},
									"arguments": schema.StringAttribute{
										Optional: true,
									},
									"working_directory": schema.StringAttribute{
										Optional: true,
									},
								},
							},
						},
						*/
						/* FIXME: Not implemented in privx-sdk-go v1.29.0
						"service_options": schema.SingleNestedAttribute{
							MarkdownDescription: "Object for service options",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"ssh": schema.SingleNestedAttribute{
									MarkdownDescription: "SSH service options",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"shell": schema.BoolAttribute{
											MarkdownDescription: "Shell channel",
											Optional:            true,
										},
										"file_transfer": schema.BoolAttribute{
											MarkdownDescription: "File transfer channel",
											Optional:            true,
										},
										"exec": schema.BoolAttribute{
											MarkdownDescription: "exec_channel",
											Optional:            true,
										},
										"tunnels": schema.BoolAttribute{
											MarkdownDescription: "tunnels",
											Optional:            true,
										},
										"x11": schema.BoolAttribute{
											MarkdownDescription: "x11",
											Optional:            true,
										},
										"other": schema.BoolAttribute{
											MarkdownDescription: "other options",
											Optional:            true,
										},
									},
								},
								"rdp": schema.SingleNestedAttribute{
									MarkdownDescription: "SSH service options",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"file_transfer": schema.BoolAttribute{
											MarkdownDescription: "File transfer channel",
											Optional:            true,
										},
										"audio": schema.BoolAttribute{
											MarkdownDescription: "audio",
											Optional:            true,
										},
										"clipboard": schema.BoolAttribute{
											MarkdownDescription: "clipboard",
											Optional:            true,
										},
										"web": schema.BoolAttribute{
											MarkdownDescription: "WEB service options",
											Optional:            true,
										},
									},
								},
								"web": schema.SingleNestedAttribute{
									MarkdownDescription: "SSH service options",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"file_transfer": schema.BoolAttribute{
											MarkdownDescription: "File transfer channel",
											Optional:            true,
										},
										"audio": schema.BoolAttribute{
											MarkdownDescription: "audio",
											Optional:            true,
										},
										"clipboard": schema.BoolAttribute{
											MarkdownDescription: "clipboard",
											Optional:            true,
										},
									},
								},
							},
						},
						"command_restrictions": schema.SingleNestedAttribute{
							MarkdownDescription: "Host services",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"enabled": schema.BoolAttribute{
									MarkdownDescription: "Are command restrictions enabled",
									Optional:            true,
								},
								"default_whitelist": schema.SingleNestedAttribute{
									MarkdownDescription: "Default whitelist handle, required if command restrictions are enabled",
									Optional:            true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											MarkdownDescription: "Whitelist ID",
											Required:            true,
										},
										"name": schema.StringAttribute{
											MarkdownDescription: "Whitelist name",
											Optional:            true,
										},
										"deleted": schema.BoolAttribute{
											MarkdownDescription: "Has whitelist been deleted, ignored in requests",
											Optional:            true,
										},
									},
								},
								"rshell_variant": schema.StringAttribute{
									MarkdownDescription: "Restricted shell variant, required if command restrictions are enabled",
									Optional:            true,
									Validators: []validator.String{
										stringvalidator.OneOf("bash", "posix"),
									},
								},
								"banner": schema.StringAttribute{
									MarkdownDescription: "Optional banner displayed in SSH terminal",
									Optional:            true,
								},
								"allow_no_match": schema.BoolAttribute{
									MarkdownDescription: "If true then commands that do not match any whitelist pattern are allowed to execute",
									Optional:            true,
								},
								"audit_match": schema.BoolAttribute{
									MarkdownDescription: "If true then an audit event is generated for every allowed command",
									Optional:            true,
								},
								"audit_no_match": schema.BoolAttribute{
									MarkdownDescription: "If true then an audit event is generated for every disallowed command",
									Optional:            true,
								},
								"whitelists": schema.SetNestedAttribute{
									Optional: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"whitelist": schema.SingleNestedAttribute{
												Optional: true,
												Attributes: map[string]schema.Attribute{
													"id": schema.StringAttribute{
														MarkdownDescription: "Whitelist ID",
														Required:            true,
													},
													"name": schema.StringAttribute{
														MarkdownDescription: "Whitelist name",
														Optional:            true,
													},

													"deleted": schema.BoolAttribute{
														MarkdownDescription: "Has whitelist been deleted, ignored in requests",
														Optional:            true,
													},
												},
											},
											"roles": schema.SetNestedAttribute{
												MarkdownDescription: "List of roles granting access to the whitelist",
												Optional:            true,

												NestedObject: schema.NestedAttributeObject{
													Attributes: map[string]schema.Attribute{
														"id": schema.StringAttribute{
															MarkdownDescription: "Role ID",
															Required:            true,
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
				Optional:            true,
			},
			*/
			/* FIXME: Not implemented in privx-sdk-go v1.29.0
			"password_rotation": schema.SingleNestedAttribute{
				MarkdownDescription: "password rotation settings for host",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"use_main_account": schema.BoolAttribute{
						MarkdownDescription: "rotate passwords of all accounts in host through one account",
						Required:            true,
					},
					"operating_system": schema.StringAttribute{
						MarkdownDescription: "Bash for Linux, Powershell for windows for shell access (LINUX | WINDOWS)",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("LINUX", "WINDOWS"),
						},
					},
					"winrm_address": schema.StringAttribute{
						MarkdownDescription: "IPv4 address or FQDN to use for winrm connection",
						Optional:            true,
					},
					"winrm_port": schema.Int64Attribute{
						MarkdownDescription: "port to use for password rotation with winrm, zero for winrm default",
						Optional:            true,
					},
					"protocol": schema.StringAttribute{
						MarkdownDescription: "protocol (SSH | WINRM)",
						Required:            true,
						Validators: []validator.String{
							stringvalidator.OneOf("SSH", "WINRM"),
						},
					},
					"password_policy_id": schema.StringAttribute{
						MarkdownDescription: "password policy to be applied",
						Required:            true,
					},
					"script_template_id": schema.StringAttribute{
						MarkdownDescription: "script template to be run in host",
						Required:            true,
					},
				},
			},
			*/
		},
	}
}

func (r *HostResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
	tflog.Debug(ctx, "Creating hoststore", map[string]interface{}{
		"connector : ": fmt.Sprintf("%+v", *connector),
	})

	r.client = hoststore.New(*connector)
}

func (r *HostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data HostResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Loaded host type data", map[string]interface{}{
		"data": fmt.Sprintf("%+v", data),
	})

	var scopePayload []string
	resp.Diagnostics.Append(data.Scope.ElementsAs(ctx, &scopePayload, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var tagsPayload []string
	resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tagsPayload, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var addressesPayload []hoststore.Address
	resp.Diagnostics.Append(data.Addresses.ElementsAs(ctx, &addressesPayload, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var servicesPayload []hoststore.Service
	for _, service := range data.Services {
		servicesPayload = append(servicesPayload,
			hoststore.Service{
				Scheme:  hoststore.Scheme(service.Scheme.ValueString()),
				Address: hoststore.Address(service.Address.ValueString()),
				Port:    int(service.Port.ValueInt64()),
				// UseForPasswordRotation: service.UseForPasswordRotation.ValueBool(), // FIXME: Not implemented in privx-sdk-go v1.29.0
			})
	}

	var principalsPayload []hoststore.Principal
	for _, principal := range data.Principals {
		var rolesPayload []rolestore.RoleRef
		for _, role := range principal.Roles {
			rolesPayload = append(rolesPayload,
				rolestore.RoleRef{
					ID: role.ID.ValueString(),
				})
		}
		/* FIXME: object application not implemented, principal only takes []string.
		var applicationsPayload []hoststore.Application
		for _, application := range principal.Applications {
			applicationsPayload = append(applicationsPayload,
			hoststore.Application {
				Name: application.Name.ValueString(),
				Application: application.Applictaion.ValueString(),
				Arguments: application.Arguments.ValueString(),
				WorkingDirectory: application.WorkingDirectory.ValueString(),
			})
		}
		*/

		principalsPayload = append(principalsPayload,
			hoststore.Principal{
				ID:             principal.ID.ValueString(),
				Source:         "terraform",
				UseUserAccount: principal.UseUserAccount.ValueBool(),
				Passphrase:     principal.Passphrase.ValueString(),
				Roles:          rolesPayload,
			})
	}

	var publicKeysPayload []hoststore.SSHPublicKey
	for _, SSHKey := range data.PublicKeys {
		publicKeysPayload = append(publicKeysPayload,
			hoststore.SSHPublicKey{
				Key: SSHKey.Key.ValueString(),
			})
	}

	host := hoststore.Host{
		AccessGroupID:       data.AccessGroupID.ValueString(),
		ExternalID:          data.ExternalID.ValueString(),
		InstanceID:          data.InstanceID.ValueString(),
		Name:                data.Name.ValueString(),
		ContactAdress:       data.ContactAddress.ValueString(),
		CloudProvider:       data.CloudProvider.ValueString(),
		CloudProviderRegion: data.CloudProviderRegion.ValueString(),
		DistinguishedName:   data.DistinguishedName.ValueString(),
		Organization:        data.Organization.ValueString(),
		OrganizationUnit:    data.OrganizationUnit.ValueString(),
		Zone:                data.Zone.ValueString(),
		HostType:            data.HostType.ValueString(),
		HostClassification:  data.HostClassification.ValueString(),
		Comment:             data.Comment.ValueString(),
		Tofu:                data.Tofu.ValueBool(),
		StandAlone:          data.StandAlone.ValueBool(),
		Audit:               data.Audit.ValueBool(),
		Scope:               scopePayload,
		Tags:                tagsPayload,
		Addresses:           addressesPayload,
		Services:            servicesPayload,
		Principals:          principalsPayload,
		PublicKeys:          publicKeysPayload,
	}

	tflog.Debug(ctx, fmt.Sprintf("hoststore.Host model used: %+v", host))

	hostID, err := r.client.CreateHost(host)

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Resource",
			"An unexpected error occurred while attempting to create the resource.\n"+
				err.Error(),
		)
		return
	}

	// Convert from the API data model to the Terraform data model
	// and set any unknown attribute values.
	data.ID = types.StringValue(hostID)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Debug(ctx, "created host resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *HostResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	host, err := r.client.Host(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read host, got error: %s", err))
		return
	}

	data.AccessGroupID = types.StringValue(host.AccessGroupID)
	data.ExternalID = types.StringValue(host.ExternalID)
	data.InstanceID = types.StringValue(host.InstanceID)
	data.Name = types.StringValue(host.Name)
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

	var principals []PrincipalModel
	for _, p := range host.Principals {
		var roles []RoleRefResourceModel
		for _, r := range p.Roles {
			roles = append(roles, RoleRefResourceModel{
				ID: types.StringValue(r.ID),
			})
		}
		var passphrase string
		for _, dp := range data.Principals {
			if dp.ID.ValueString() == p.ID {
				passphrase = dp.Passphrase.ValueString()
			}
		}
		principals = append(principals, PrincipalModel{
			ID:             types.StringValue(p.ID),
			Passphrase:     types.StringValue(passphrase),
			UseUserAccount: types.BoolValue(p.UseUserAccount),
			// Rotate:     types.BoolValue(p.Rotate), // FIXME: Not implemented in privx-sdk-go v1.29.0
			// UseForPasswordRotation:     types.BoolValue(p.UseForPasswordRotation), // FIXME: Not implemented in privx-sdk-go v1.29.0
			// ServiceOptions: serviceOptions, // FIXME: Not implemented in privx-sdk-go v1.29.0
			// CommandRestrictions: commandRestrictions, // FIXME: Not implemented in privx-sdk-go v1.29.0
			Roles: roles,
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

func (r *HostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *HostResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var scopePayload []string
	resp.Diagnostics.Append(data.Scope.ElementsAs(ctx, &scopePayload, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tagsPayload := make([]string, len(data.Tags.Elements()))
	resp.Diagnostics.Append(data.Tags.ElementsAs(ctx, &tagsPayload, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	addressesPayload := make([]hoststore.Address, len(data.Addresses.Elements()))
	resp.Diagnostics.Append(data.Addresses.ElementsAs(ctx, &addressesPayload, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var servicesPayload []hoststore.Service
	for _, service := range data.Services {
		servicesPayload = append(servicesPayload,
			hoststore.Service{
				Scheme:  hoststore.Scheme(service.Scheme.ValueString()),
				Address: hoststore.Address(service.Address.ValueString()),
				Port:    int(service.Port.ValueInt64()),
				// UseForPasswordRotation: service.UseForPasswordRotation.ValueBool(), // FIXME: Not implemented in privx-sdk-go v1.29.0
			})
	}

	var principalsPayload []hoststore.Principal
	for _, principal := range data.Principals {
		var rolesPayload []rolestore.RoleRef
		for _, role := range principal.Roles {
			rolesPayload = append(rolesPayload,
				rolestore.RoleRef{
					ID: role.ID.ValueString(),
				})
		}
		/* FIXME: object application not implemented, principal only takes []string.
		var applicationsPayload []hoststore.Application
		for _, application := range principal.Applications {
			applicationsPayload = append(applicationsPayload,
			hoststore.Application {
				Name: application.Name.ValueString(),
				Application: application.Applictaion.ValueString(),
				Arguments: application.Arguments.ValueString(),
				WorkingDirectory: application.WorkingDirectory.ValueString(),
			})
		}
		*/
		principalsPayload = append(principalsPayload,
			hoststore.Principal{
				ID:             principal.ID.ValueString(),
				UseUserAccount: principal.UseUserAccount.ValueBool(),
				Passphrase:     principal.Passphrase.ValueString(),
				Roles:          rolesPayload,
			})
	}

	var publicKeysPayload []hoststore.SSHPublicKey
	for _, SSHKey := range data.PublicKeys {
		publicKeysPayload = append(publicKeysPayload,
			hoststore.SSHPublicKey{
				Key: SSHKey.Key.ValueString(),
			})
	}

	host := hoststore.Host{
		AccessGroupID:       data.AccessGroupID.ValueString(),
		ExternalID:          data.ExternalID.ValueString(),
		InstanceID:          data.InstanceID.ValueString(),
		Name:                data.Name.ValueString(),
		ContactAdress:       data.ContactAddress.ValueString(),
		CloudProvider:       data.CloudProvider.ValueString(),
		CloudProviderRegion: data.CloudProviderRegion.ValueString(),
		DistinguishedName:   data.DistinguishedName.ValueString(),
		Organization:        data.Organization.ValueString(),
		OrganizationUnit:    data.OrganizationUnit.ValueString(),
		Zone:                data.Zone.ValueString(),
		HostType:            data.HostType.ValueString(),
		HostClassification:  data.HostClassification.ValueString(),
		Comment:             data.Comment.ValueString(),
		Tofu:                data.Tofu.ValueBool(),
		StandAlone:          data.StandAlone.ValueBool(),
		Audit:               data.Audit.ValueBool(),
		Scope:               scopePayload,
		Tags:                tagsPayload,
		Addresses:           addressesPayload,
		Services:            servicesPayload,
		Principals:          principalsPayload,
		PublicKeys:          publicKeysPayload,
	}

	tflog.Debug(ctx, fmt.Sprintf("hoststore.Host model used: %+v", host))

	err := r.client.UpdateHost(
		data.ID.ValueString(),
		&host)

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update host, got error: %s", err))
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *HostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *HostResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteHost(data.ID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete host, got error: %s", err))
		return
	}
}

func (r *HostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
