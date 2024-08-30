package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ provider.Provider = &buildkitProvider{}
var _ provider.ProviderWithFunctions = &buildkitProvider{}

type buildkitProvider struct {
	version string
}

type buildkitProviderModel struct {
	BuildkitHost types.String        `tfsdk:"buildkit_host"`
	RegistryAuth []registryAuthModel `tfsdk:"registry_auth"`
}

type registryAuthModel struct {
	Address  types.String `tfsdk:"address"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

func (p *buildkitProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "buildkit"
	resp.Version = p.version
}

func (p *buildkitProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"buildkit_host": schema.StringAttribute{
				Optional:    true,
				Description: "The address of the BuildKit daemon. Defaults to 'unix:///var/run/buildkit/buildkitd.sock'.",
			},
			"registry_auth": schema.ListNestedAttribute{
				Optional:    true,
				Description: "Authentication configuration for Docker registries.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Required:    true,
							Description: "The address of the Docker registry.",
						},
						"username": schema.StringAttribute{
							Required:    true,
							Description: "The username for the Docker registry.",
						},
						"password": schema.StringAttribute{
							Required:    true,
							Sensitive:   true,
							Description: "The password for the Docker registry.",
						},
					},
				},
			},
		},
	}
}

func (p *buildkitProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config buildkitProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.BuildkitHost.IsNull() {
		config.BuildkitHost = types.StringValue("unix:///var/run/buildkit/buildkitd.sock")
	}

	// Validate registry auth configurations
	for _, auth := range config.RegistryAuth {
		if auth.Address.IsNull() || auth.Username.IsNull() || auth.Password.IsNull() {
			resp.Diagnostics.AddError(
				"Invalid Registry Auth Configuration",
				"All fields (address, username, password) must be provided for each registry_auth block.",
			)
			return
		}
	}

	// Make the config available to resources
	resp.ResourceData = config
}

func (p *buildkitProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// NewBuildkitImageResource,
	}
}

func (p *buildkitProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Add any data sources here
	}
}

func (p *buildkitProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		// NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &buildkitProvider{
			version: version,
		}
	}
}
