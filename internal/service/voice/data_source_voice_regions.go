package voice

import (
	"context"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &voiceRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &voiceRegionsDataSource{}
)

// voiceRegionsDataSource is the data source implementation.
type voiceRegionsDataSource struct {
	client *discord.Client
}

// voiceRegionsDataSourceModel maps the data source schema data.
type voiceRegionsDataSourceModel struct {
	Regions types.List `tfsdk:"regions"`
}

// NewVoiceRegionsDataSource returns a new voice regions data source.
func NewVoiceRegionsDataSource() datasource.DataSource {
	return &voiceRegionsDataSource{}
}

// regionObjectType returns the attr.Type map for a voice region object.
func regionObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.StringType,
		"name":       types.StringType,
		"optimal":    types.BoolType,
		"deprecated": types.BoolType,
	}
}

// Metadata returns the data source type name.
func (d *voiceRegionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_voice_regions"
}

// Configure adds the provider configured client to the data source.
func (d *voiceRegionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Schema defines the schema for the data source.
func (d *voiceRegionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get a list of available voice regions.",
		Attributes: map[string]schema.Attribute{
			"regions": schema.ListNestedAttribute{
				Description: "A list of available voice regions.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique ID for the region.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the region.",
							Computed:    true,
						},
						"optimal": schema.BoolAttribute{
							Description: "True for the closest server to the current user's client.",
							Computed:    true,
						},
						"deprecated": schema.BoolAttribute{
							Description: "Whether this is a deprecated voice region.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *voiceRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	regions, err := d.client.ListVoiceRegions(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Voice Regions",
			"Could not read voice regions: "+err.Error(),
		)
		return
	}

	regionObjects := make([]attr.Value, 0, len(regions))
	for _, r := range regions {
		obj, diags := types.ObjectValue(
			regionObjectType(),
			map[string]attr.Value{
				"id":         types.StringValue(r.ID),
				"name":       types.StringValue(r.Name),
				"optimal":    types.BoolValue(r.Optimal),
				"deprecated": types.BoolValue(r.Deprecated),
			},
		)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		regionObjects = append(regionObjects, obj)
	}

	regionsList, diags := types.ListValue(types.ObjectType{AttrTypes: regionObjectType()}, regionObjects)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := voiceRegionsDataSourceModel{
		Regions: regionsList,
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
