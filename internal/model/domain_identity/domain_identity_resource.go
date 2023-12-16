package domain_identity

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DomainIdentityResourceModel struct {
	Id types.String `tfsdk:"id"`
}

func GetDomainIdentityResourceSchema() *schema.Schema {
	return &schema.Schema{
		MarkdownDescription: "DomainIdentity resource",
	}
}
