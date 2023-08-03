//
// Copyright (c) 2020 SSH Communications Security Inc.
//
// All rights reserved.
//

package hoststore

import "github.com/SSHcom/privx-sdk-go/api/rolestore"

// Source of host objects
type Source string

// Source constants
const (
	UI   = Source("UI")
	SCAN = Source("SCAN")
)

// Address is fully qualified domain names, IPv4 or IPv6 addresses of the host
type Address string

// Scheme of protocols allowed by the host
type Scheme string

// Scheme constants, all supported protocols
const (
	SSH   = Scheme("SSH")
	RDP   = Scheme("RDP")
	VNC   = Scheme("VNC")
	HTTP  = Scheme("HTTP")
	HTTPS = Scheme("HTTPS")
)

// Params struct for pagination queries
type Params struct {
	Offset  int    `json:"offset,omitempty"`
	Limit   int    `json:"limit,omitempty"`
	Sortdir string `json:"sortdir,omitempty"`
	Sortkey string `json:"sortkey,omitempty"`
	Filter  string `json:"filter,omitempty"`
	Query   string `json:"query,omitempty"`
}

// HostSearchObject host search object definition
type HostSearchObject struct {
	ID                    string   `json:"id,omitempty"`
	Keywords              string   `json:"keywords,omitempty"`
	ExternalID            string   `json:"external_id,omitempty"`
	InstanceID            string   `json:"instance_id,omitempty"`
	SourceID              string   `json:"source_id,omitempty"`
	Disabled              string   `json:"disabled,omitempty"`
	Deployable            bool     `json:"deployable,omitempty"`
	IgnoreDisabledSources bool     `json:"ignore_disabled_sources,omitempty"`
	Port                  []int    `json:"port,omitempty"`
	CommonName            []string `json:"common_name,omitempty"`
	Organization          []string `json:"organization,omitempty"`
	OrganizationalUnit    []string `json:"organizational_unit,omitempty"`
	Address               []string `json:"address,omitempty"`
	Service               []string `json:"service,omitempty"`
	Zone                  []string `json:"zone,omitempty"`
	HostType              []string `json:"host_type,omitempty"`
	HostClassification    []string `json:"host_classification,omitempty"`
	Role                  []string `json:"role,omitempty"`
	Scope                 []string `json:"scope,omitempty"`
	Tags                  []string `json:"tags,omitempty"`
	AccessGroupIDs        []string `json:"access_group_ids,omitempty"`
	CloudProviders        []string `json:"cloud_providers,omitempty"`
	CloudProviderRegions  []string `json:"cloud_provider_regions,omitempty"`
	Statuses              []string `json:"statuses,omitempty"`
	DistinguishedName     []string `json:"distinguished_name,omitempty"`
}

// HostDisabledRequest host disabled request definition
type HostDisabledRequest struct {
	Disabled bool `json:"disabled"`
}

// Service specify the service available on target host
type Service struct {
	Scheme  Scheme  `json:"service"`
	Address Address `json:"address"`
	Port    int     `json:"port"`
	Source  Source  `json:"source"`
}

// Principal of the target host
type Principal struct {
	ID             string              `json:"principal"`
	Roles          []rolestore.RoleRef `json:"roles"`
	Source         Source              `json:"source"`
	UseUserAccount bool                `json:"use_user_account"`
	Passphrase     string              `json:"passphrase"`
	Applications   []string            `json:"applications"`
}

// SSHPublicKey host public keys
type SSHPublicKey struct {
	Key         string `json:"key,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
}

// Status of the secret object
type Status struct {
	K string `json:"k,omitempty"`
	V string `json:"v,omitempty"`
}

// Host defines PrivX target
type Host struct {
	ID                  string         `json:"id,omitempty"`
	AccessGroupID       string         `json:"access_group_id,omitempty"`
	ExternalID          string         `json:"external_id,omitempty"`
	InstanceID          string         `json:"instance_id,omitempty"`
	SourceID            string         `json:"source_id,omitempty"`
	Name                string         `json:"common_name,omitempty"`
	ContactAdress       string         `json:"contact_address,omitempty"`
	CloudProvider       string         `json:"cloud_provider,omitempty"`
	CloudProviderRegion string         `json:"cloud_provider_region,omitempty"`
	Created             string         `json:"created,omitempty"`
	Updated             string         `json:"updated,omitempty"`
	UpdatedBy           string         `json:"updated_by,omitempty"`
	DistinguishedName   string         `json:"distinguished_name,omitempty"`
	Organization        string         `json:"organization,omitempty"`
	OrganizationUnit    string         `json:"organizational_unit,omitempty"`
	Zone                string         `json:"zone,omitempty"`
	HostType            string         `json:"host_type,omitempty"`
	HostClassification  string         `json:"host_classification,omitempty"`
	Comment             string         `json:"comment,omitempty"`
	Disabled            string         `json:"disabled,omitempty"`
	Deployable          bool           `json:"deployable,omitempty"`
	Tofu                bool           `json:"tofu,omitempty"`
	StandAlone          bool           `json:"stand_alone_host,omitempty"`
	Audit               bool           `json:"audit_enabled,omitempty"`
	Scope               []string       `json:"scope,omitempty"`
	Tags                []string       `json:"tags,omitempty"`
	Addresses           []Address      `json:"addresses,omitempty"`
	Services            []Service      `json:"services,omitempty"`
	Principals          []Principal    `json:"principals,omitempty"`
	PublicKeys          []SSHPublicKey `json:"ssh_host_public_keys,omitempty"`
	Status              []Status       `json:"status,omitempty"`
}

// SSHService default options
type SSHService struct {
	Shell        bool `json:"shell"`
	FileTransfer bool `json:"file_transfer"`
	Exec         bool `json:"exec"`
	Tunnels      bool `json:"tunnels"`
	Xeleven      bool `json:"x11"`
	Other        bool `json:"other"`
}

// RDPService default options
type RDPService struct {
	FileTransfer bool `json:"file_transfer"`
	Audio        bool `json:"audio"`
	Clipboard    bool `json:"clipboard"`
}

// WebService default options
type WebService struct {
	FileTransfer bool `json:"file_transfer"`
	Audio        bool `json:"audio"`
	Clipboard    bool `json:"clipboard"`
}

// DefaultServiceOptions default service options
type DefaultServiceOptions struct {
	SSH SSHService `json:"ssh"`
	RDP RDPService `json:"rdp"`
	Web WebService `json:"web"`
}

// Service creates a corresponding service definition
//   hosts.SSH.Service(...)
func (scheme Scheme) Service(addr Address, port int) Service {
	return Service{
		Scheme:  scheme,
		Address: addr,
		Port:    port,
		Source:  UI,
	}
}

// NewPrincipal creates a corresponding definition from roles
func NewPrincipal(id string, role ...rolestore.RoleRef) Principal {
	return Principal{
		ID:     id,
		Roles:  role,
		Source: UI,
	}
}
