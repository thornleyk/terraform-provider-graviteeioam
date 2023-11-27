package organization

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type OrganizationResourceModel struct {
	Id types.String `tfsdk:"id"`
}

func GetEnvironmentResourceSchema() *schema.Schema {
	return &schema.Schema{
		MarkdownDescription: "Organization resource",
	}
}
