---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "privx_carrier_config Data Source - terraform-provider-privx"
subcategory: ""
description: |-
  CarrierConfig DataSource
---

# privx_carrier_config (Data Source)

CarrierConfig DataSource

## Example Usage

```terraform
data "privx_carrier_config" "my_carrier_config" {
  trusted_client_id = "trusted_client_id_of_the_desired_carrier"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `trusted_client_id` (String) CarrierConfig ID

### Read-Only

- `carrier_config` (String) Carrier config
