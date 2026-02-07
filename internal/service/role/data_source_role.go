package role

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
	_ datasource.DataSource              = &roleDataSource{}
	_ datasource.DataSourceWithConfigure = &roleDataSource{}
)

// roleDataSource is the data source implementation.
type roleDataSource struct {
	client *discord.Client
}

// roleDataSourceModel maps the data source schema data.
type roleDataSourceModel struct {
	GuildID     types.String `tfsdk:"guild_id"`
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Permissions types.String `tfsdk:"permissions"`
	Color       types.Int64  `tfsdk:"color"`
	Hoist       types.Bool   `tfsdk:"hoist"`
	Mentionable types.Bool   `tfsdk:"mentionable"`
	Position    types.Int64  `tfsdk:"position"`
	Managed     types.Bool   `tfsdk:"managed"`
}

// NewRoleDataSource returns a new role data source.
func NewRoleDataSource() datasource.DataSource {
	return &roleDataSource{}
}

// Metadata returns the data source type name.
func (d *roleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_role"
}

// Configure adds the provider configured client to the data source.
func (d *roleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Schema defines the schema for the data source.
func (d *roleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Discord role by ID or name within a guild.",
		Attributes: map[string]schema.Attribute{
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild to look up the role in.",
				Required:    true,
			},
			"id": schema.StringAttribute{
				Description: "The ID of the role. At least one of id or name must be provided.",
				Optional:    true,
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the role. At least one of id or name must be provided.",
				Optional:    true,
				Computed:    true,
			},
			"permissions": schema.StringAttribute{
				Description: "The permission bit set for this role.",
				Computed:    true,
			},
			"color": schema.Int64Attribute{
				Description: "The integer representation of the role color.",
				Computed:    true,
			},
			"hoist": schema.BoolAttribute{
				Description: "Whether this role is hoisted (displayed separately in the sidebar).",
				Computed:    true,
			},
			"mentionable": schema.BoolAttribute{
				Description: "Whether this role is mentionable.",
				Computed:    true,
			},
			"position": schema.Int64Attribute{
				Description: "The position of this role.",
				Computed:    true,
			},
			"managed": schema.BoolAttribute{
				Description: "Whether this role is managed by an integration.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *roleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config roleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that at least one of id or name is provided.
	if (config.ID.IsNull() || config.ID.ValueString() == "") &&
		(config.Name.IsNull() || config.Name.ValueString() == "") {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"At least one of 'id' or 'name' must be provided to look up a role.",
		)
		return
	}

	roles, err := d.client.GetGuildRoles(ctx, discord.Snowflake(config.GuildID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Discord Guild Roles",
			"Could not read guild roles for guild "+config.GuildID.ValueString()+": "+err.Error(),
		)
		return
	}

	var found *discord.Role

	// If id is provided, look up by ID.
	if !config.ID.IsNull() && config.ID.ValueString() != "" {
		for _, r := range roles {
			if r.ID.String() == config.ID.ValueString() {
				found = r
				break
			}
		}
	} else {
		// Otherwise look up by name.
		for _, r := range roles {
			if r.Name == config.Name.ValueString() {
				found = r
				break
			}
		}
	}

	if found == nil {
		resp.Diagnostics.AddError(
			"Role Not Found",
			"Could not find a role matching the given criteria in guild "+config.GuildID.ValueString()+".",
		)
		return
	}

	config.ID = types.StringValue(found.ID.String())
	config.Name = types.StringValue(found.Name)
	config.Permissions = types.StringValue(found.Permissions)
	config.Color = types.Int64Value(int64(found.Color))
	config.Hoist = types.BoolValue(found.Hoist)
	config.Mentionable = types.BoolValue(found.Mentionable)
	config.Position = types.Int64Value(int64(found.Position))
	config.Managed = types.BoolValue(found.Managed)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
