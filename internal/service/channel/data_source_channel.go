package channel

import (
	"context"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &channelDataSource{}
	_ datasource.DataSourceWithConfigure = &channelDataSource{}
)

// channelDataSource is the data source implementation.
type channelDataSource struct {
	client *discord.Client
}

// channelDataSourceModel maps the data source schema data.
type channelDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	GuildID  types.String `tfsdk:"guild_id"`
	Name     types.String `tfsdk:"name"`
	Type     types.Int64  `tfsdk:"type"`
	Topic    types.String `tfsdk:"topic"`
	Position types.Int64  `tfsdk:"position"`
	NSFW     types.Bool   `tfsdk:"nsfw"`
	ParentID types.String `tfsdk:"parent_id"`
}

// NewChannelDataSource returns a new channel data source.
func NewChannelDataSource() datasource.DataSource {
	return &channelDataSource{}
}

// Metadata returns the data source type name.
func (d *channelDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_channel"
}

// Configure adds the provider configured client to the data source.
func (d *channelDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Schema defines the schema for the data source.
func (d *channelDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about a Discord channel.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the channel.",
				Required:    true,
			},
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild this channel belongs to.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the channel.",
				Computed:    true,
			},
			"type": schema.Int64Attribute{
				Description: "The type of the channel.",
				Computed:    true,
			},
			"topic": schema.StringAttribute{
				Description: "The channel topic.",
				Computed:    true,
			},
			"position": schema.Int64Attribute{
				Description: "The sorting position of the channel.",
				Computed:    true,
			},
			"nsfw": schema.BoolAttribute{
				Description: "Whether the channel is NSFW.",
				Computed:    true,
			},
			"parent_id": schema.StringAttribute{
				Description: "The ID of the parent category for a channel.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *channelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config channelDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ch, err := d.client.GetChannel(ctx, discord.Snowflake(config.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Discord Channel",
			"Could not read channel ID "+config.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	config.Type = types.Int64Value(int64(ch.Type))
	config.NSFW = types.BoolValue(ch.NSFW)

	if ch.GuildID != nil {
		config.GuildID = types.StringValue(ch.GuildID.String())
	} else {
		config.GuildID = types.StringNull()
	}

	if ch.Name != nil {
		config.Name = types.StringValue(*ch.Name)
	} else {
		config.Name = types.StringNull()
	}

	if ch.Topic != nil {
		config.Topic = types.StringValue(*ch.Topic)
	} else {
		config.Topic = types.StringNull()
	}

	if ch.Position != nil {
		config.Position = types.Int64Value(int64(*ch.Position))
	} else {
		config.Position = types.Int64Value(0)
	}

	if ch.ParentID != nil {
		config.ParentID = types.StringValue(ch.ParentID.String())
	} else {
		config.ParentID = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
