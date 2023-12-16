package environment

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type EnvironmentResourceModel struct {
	Id types.String `tfsdk:"id"`
}

func GetEnvironmentResourceSchema() *schema.Schema {
	return &schema.Schema{
		MarkdownDescription: "Environment resource",
	}
}
