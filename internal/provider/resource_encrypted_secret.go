package provider

import (
	"context"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-scaffolding/encryption"
)

func resourceEncryptedSecret() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "A secret, shh!",

		CreateContext: resourceEncryptedSecretRead,
		ReadContext:   resourceEncryptedSecretRead,
		UpdateContext: resourceEncryptedSecretRead,
		DeleteContext: resourceEncryptedSecretDelete,

		Schema: map[string]*schema.Schema{
			"encrypted_secret": {
				Description: "The base64 encrypted secret",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}

func resourceEncryptedSecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceEncryptedSecretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
