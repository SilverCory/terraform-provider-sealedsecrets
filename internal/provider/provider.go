package provider

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-scaffolding/encryption"
)

func init() {
	schema.DescriptionKind = schema.StringPlain
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"private_key": {
					Type:        schema.TypeString,
					Required:    true,
					Sensitive:   true,
					Description: "The private key used to encrypt the data",
					DefaultFunc: schema.EnvDefaultFunc("SECRET_PRIVATE_KEY", nil),
				},
			},
			ResourcesMap: map[string]*schema.Resource{
				"sealedsecrets_secret": resourceSealedSecretsSecret(),
			},
		}

		p.ConfigureContextFunc = configure(version, p)

		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(_ context.Context, r *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var privateKey string
		switch v := r.Get("private_key").(type) {
		case string:
			privateKey = v
		}

		if privateKey == "" {
			return nil, diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Missing private key",
				Detail:        "The private key is not supplied",
				AttributePath: cty.GetAttrPath("private_key"),
			}}
		}

		enc, err := encryption.New("", privateKey)
		if err != nil {
			return nil, diag.Diagnostics{diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       "Invalid private key",
				Detail:        err.Error(),
				AttributePath: cty.GetAttrPath("private_key"),
			}}
		}

		return enc, nil
	}
}
