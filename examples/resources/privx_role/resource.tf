resource "privx_role" "foo" {
  name            = "test-dev-provider"
  comment         = ""
  access_group_id = "565381ce-0911-4ba8-8606-8eecd8074556"
  permissions     = []
  permit_agent    = false
  source_rules = jsonencode({
    type = "GROUP" // GROUP | RULE
    //source        = ""
    //search_string = ""
    match = "ANY" // ANY | ALL
    rules = []
  })
}
