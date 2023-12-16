package organization_identity

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OrganizationIdentityResourceModel struct {
	Id types.String `tfsdk:"id"`
}

func GetOrganizationIdentityResourceSchema() *schema.Schema {
	return &schema.Schema{
		MarkdownDescription: "OrganizationIdentity resource",
	}
}
