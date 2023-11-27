package environment

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

func GetEnvironmentResourceSchema() *schema.Schema {
	return &schema.Schema{
		MarkdownDescription: "Environment resource",
	}
}
