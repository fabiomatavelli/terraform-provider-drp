package drpv4

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"gitlab.com/rackn/provision/v4/api"
)

// expandStringList converts a interface{} to a []string
func expandStringList(v interface{}) []string {
	if v == nil {
		return nil
	}
	result := make([]string, len(v.([]interface{})))
	for i, s := range v.([]interface{}) {
		result[i] = s.(string)
	}
	return result
}

// expandMapInterface converts a map[string]interface{} to a map[string]string
func expandMapInterface(v interface{}) map[string]string {
	if v == nil {
		return nil
	}
	result := make(map[string]string)
	for k, s := range v.(map[string]interface{}) {
		result[k] = s.(string)
	}
	return result
}

// randomString generates a random string of length n
func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func resourceGenericConfigure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) *api.Client {
	if req.ProviderData == nil {
		tflog.Error(ctx, "Missing provider data")
		return nil
	}

	client, ok := req.ProviderData.(*Config)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *Config, got %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return nil
	}

	return client.session
}
