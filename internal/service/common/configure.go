package common

import (
	"fmt"

	"github.com/edw1nzhao/terraform-provider-discord/internal/conns"
	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ClientFromProviderData extracts the Discord client from provider data.
// Returns nil client if provider data is nil (during early provider configuration).
// Adds an error diagnostic if the type assertion fails.
func ClientFromProviderData(providerData any, diagnostics *diag.Diagnostics) *discord.Client {
	if providerData == nil {
		return nil
	}
	data, ok := providerData.(*conns.ProviderData)
	if !ok {
		diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *conns.ProviderData, got: %T", providerData),
		)
		return nil
	}
	return data.Client
}

// ProviderDataFromConfig extracts the full ProviderData from provider data.
// Used by resources that need both client and application ID.
func ProviderDataFromConfig(providerData any, diagnostics *diag.Diagnostics) *conns.ProviderData {
	if providerData == nil {
		return nil
	}
	data, ok := providerData.(*conns.ProviderData)
	if !ok {
		diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *conns.ProviderData, got: %T", providerData),
		)
		return nil
	}
	return data
}
