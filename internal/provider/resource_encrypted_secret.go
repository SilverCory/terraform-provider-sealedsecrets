package provider

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-scaffolding/encryption"
)

func resourceSealedSecretsSecret() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A secret, shh!",

		CreateContext: resourceSealedSecretsSecretRead,
		ReadContext:   resourceSealedSecretsSecretRead,
		UpdateContext: resourceSealedSecretsSecretRead,
		DeleteContext: resourceSealedSecretsSecretDelete,

		Schema: map[string]*schema.Schema{
			"encrypted_secret": {
				Description: "The base64 encoded GPG encrypted secret",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   false,
			},
			"value": {
				Description: "The secret in the flesh",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceSealedSecretsSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	enc := meta.(*encryption.Encryptor)

	var encrypredString string
	switch v := d.Get("encrypted_secret").(type) {
	case string:
		encrypredString = v
	}

	out, hash, err := enc.DecryptString(encrypredString)
	if err != nil {
		return diag.Diagnostics{{
			Severity:      diag.Error,
			Summary:       "Unable to decrypt string",
			Detail:        err.Error(),
			AttributePath: cty.GetAttrPath("encrypted_secret"),
		}}
	}

	d.SetId(hash)
	d.Set("value", out)

	return nil
}

func resourceSealedSecretsSecretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
