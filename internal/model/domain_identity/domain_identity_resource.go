package domain_identity

import "github.com/hashicorp/terraform-plugin-framework/resource/schema"

func GetDomainIdentityResourceSchema() *schema.Schema {
	return &schema.Schema{
		MarkdownDescription: "DomainIdentity resource",
	}
}
