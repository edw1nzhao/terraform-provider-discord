package user

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
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

// userDataSource is the data source implementation.
type userDataSource struct {
	client *discord.Client
}

// userDataSourceModel maps the data source schema data.
type userDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Username    types.String `tfsdk:"username"`
	GlobalName  types.String `tfsdk:"global_name"`
	Avatar      types.String `tfsdk:"avatar"`
	Bot         types.Bool   `tfsdk:"bot"`
	Banner      types.String `tfsdk:"banner"`
	AccentColor types.Int64  `tfsdk:"accent_color"`
}

// NewUserDataSource returns a new user data source.
func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

// Metadata returns the data source type name.
func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Configure adds the provider configured client to the data source.
func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	d.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Schema defines the schema for the data source.
func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to get information about a Discord user.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the user.",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "The user's username.",
				Computed:    true,
			},
			"global_name": schema.StringAttribute{
				Description: "The user's display name.",
				Computed:    true,
			},
			"avatar": schema.StringAttribute{
				Description: "The user's avatar hash.",
				Computed:    true,
			},
			"bot": schema.BoolAttribute{
				Description: "Whether the user is a bot.",
				Computed:    true,
			},
			"banner": schema.StringAttribute{
				Description: "The user's banner hash.",
				Computed:    true,
			},
			"accent_color": schema.Int64Attribute{
				Description: "The user's banner color as an integer.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config userDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	u, err := d.client.GetUser(ctx, discord.Snowflake(config.ID.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Discord User",
			"Could not read user ID "+config.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	config.Username = types.StringValue(u.Username)
	config.Bot = types.BoolValue(u.Bot)

	if u.GlobalName != nil {
		config.GlobalName = types.StringValue(*u.GlobalName)
	} else {
		config.GlobalName = types.StringNull()
	}

	if u.Avatar != nil {
		config.Avatar = types.StringValue(*u.Avatar)
	} else {
		config.Avatar = types.StringNull()
	}

	if u.Banner != nil {
		config.Banner = types.StringValue(*u.Banner)
	} else {
		config.Banner = types.StringNull()
	}

	if u.AccentColor != nil {
		config.AccentColor = types.Int64Value(int64(*u.AccentColor))
	} else {
		config.AccentColor = types.Int64Value(0)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
