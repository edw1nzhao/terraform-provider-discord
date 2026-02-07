package onboarding

import (
	"context"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure interface compliance.
var (
	_ resource.Resource                = &guildOnboardingResource{}
	_ resource.ResourceWithConfigure   = &guildOnboardingResource{}
	_ resource.ResourceWithImportState = &guildOnboardingResource{}
)

// NewGuildOnboardingResource returns a new resource for discord_guild_onboarding.
func NewGuildOnboardingResource() resource.Resource {
	return &guildOnboardingResource{}
}

type guildOnboardingResource struct {
	client *discord.Client
}

// guildOnboardingModel maps the Terraform schema to Go types.
type guildOnboardingModel struct {
	GuildID           types.String `tfsdk:"guild_id"`
	Enabled           types.Bool   `tfsdk:"enabled"`
	Mode              types.Int64  `tfsdk:"mode"`
	DefaultChannelIDs types.Set    `tfsdk:"default_channel_ids"`
	Prompts           types.List   `tfsdk:"prompts"`
}

// promptModel maps a single prompt entry.
type promptModel struct {
	ID           types.String `tfsdk:"id"`
	Type         types.Int64  `tfsdk:"type"`
	Title        types.String `tfsdk:"title"`
	SingleSelect types.Bool   `tfsdk:"single_select"`
	Required     types.Bool   `tfsdk:"required"`
	InOnboarding types.Bool   `tfsdk:"in_onboarding"`
	Options      types.List   `tfsdk:"options"`
}

// optionModel maps a single option inside a prompt.
type optionModel struct {
	ID            types.String `tfsdk:"id"`
	Title         types.String `tfsdk:"title"`
	Description   types.String `tfsdk:"description"`
	ChannelIDs    types.Set    `tfsdk:"channel_ids"`
	RoleIDs       types.Set    `tfsdk:"role_ids"`
	EmojiID       types.String `tfsdk:"emoji_id"`
	EmojiName     types.String `tfsdk:"emoji_name"`
	EmojiAnimated types.Bool   `tfsdk:"emoji_animated"`
}

// Attribute type maps for nested objects.

func optionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":             types.StringType,
		"title":          types.StringType,
		"description":    types.StringType,
		"channel_ids":    types.SetType{ElemType: types.StringType},
		"role_ids":       types.SetType{ElemType: types.StringType},
		"emoji_id":       types.StringType,
		"emoji_name":     types.StringType,
		"emoji_animated": types.BoolType,
	}
}

func promptAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":            types.StringType,
		"type":          types.Int64Type,
		"title":         types.StringType,
		"single_select": types.BoolType,
		"required":      types.BoolType,
		"in_onboarding": types.BoolType,
		"options":       types.ListType{ElemType: types.ObjectType{AttrTypes: optionAttrTypes()}},
	}
}

// Metadata sets the type name.
func (r *guildOnboardingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_guild_onboarding"
}

// Schema defines the schema.
func (r *guildOnboardingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord guild onboarding configuration.",
		Attributes: map[string]schema.Attribute{
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild. Acts as the resource ID.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether onboarding is enabled.",
				Required:    true,
			},
			"mode": schema.Int64Attribute{
				Description: "The onboarding mode (0 = ONBOARDING_DEFAULT, 1 = ONBOARDING_ADVANCED).",
				Optional:    true,
			},
			"default_channel_ids": schema.SetAttribute{
				Description: "Channel IDs that members get opted into automatically.",
				Required:    true,
				ElementType: types.StringType,
			},
			"prompts": schema.ListNestedAttribute{
				Description: "The onboarding prompts.",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The prompt ID. Computed by the API.",
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"type": schema.Int64Attribute{
							Description: "The prompt type (0 = MULTIPLE_CHOICE, 1 = DROPDOWN).",
							Required:    true,
						},
						"title": schema.StringAttribute{
							Description: "The prompt title.",
							Required:    true,
						},
						"single_select": schema.BoolAttribute{
							Description: "Whether users are limited to selecting one option.",
							Required:    true,
						},
						"required": schema.BoolAttribute{
							Description: "Whether the prompt is required.",
							Required:    true,
						},
						"in_onboarding": schema.BoolAttribute{
							Description: "Whether the prompt is present in the onboarding flow.",
							Required:    true,
						},
						"options": schema.ListNestedAttribute{
							Description: "The options for this prompt.",
							Required:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "The option ID. Computed by the API.",
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
									"title": schema.StringAttribute{
										Description: "The option title.",
										Required:    true,
									},
									"description": schema.StringAttribute{
										Description: "The option description.",
										Optional:    true,
									},
									"channel_ids": schema.SetAttribute{
										Description: "Channel IDs associated with this option.",
										Optional:    true,
										ElementType: types.StringType,
									},
									"role_ids": schema.SetAttribute{
										Description: "Role IDs associated with this option.",
										Optional:    true,
										ElementType: types.StringType,
									},
									"emoji_id": schema.StringAttribute{
										Description: "The emoji ID, if using a custom emoji.",
										Optional:    true,
									},
									"emoji_name": schema.StringAttribute{
										Description: "The emoji name if custom, or the unicode character.",
										Optional:    true,
									},
									"emoji_animated": schema.BoolAttribute{
										Description: "Whether the emoji is animated.",
										Optional:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure stores the provider data.
func (r *guildOnboardingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create applies the onboarding config via PUT (onboarding always exists as part of the guild).
func (r *guildOnboardingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan guildOnboardingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params, diags := expandOnboardingParams(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ob, err := r.client.ModifyGuildOnboarding(ctx, discord.Snowflake(plan.GuildID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error creating guild onboarding", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenOnboarding(ctx, ob, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read refreshes the Terraform state from the API.
func (r *guildOnboardingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state guildOnboardingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ob, err := r.client.GetGuildOnboarding(ctx, discord.Snowflake(state.GuildID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading guild onboarding", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenOnboarding(ctx, ob, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update modifies the onboarding config.
func (r *guildOnboardingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan guildOnboardingModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params, diags := expandOnboardingParams(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ob, err := r.client.ModifyGuildOnboarding(ctx, discord.Snowflake(plan.GuildID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError("Error updating guild onboarding", err.Error())
		return
	}

	resp.Diagnostics.Append(flattenOnboarding(ctx, ob, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete disables onboarding.
func (r *guildOnboardingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state guildOnboardingModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	enabled := false
	params := &discord.ModifyOnboardingParams{
		Enabled: &enabled,
	}

	_, err := r.client.ModifyGuildOnboarding(ctx, discord.Snowflake(state.GuildID.ValueString()), params)
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting guild onboarding", err.Error())
	}
}

// ImportState supports importing by guild_id.
func (r *guildOnboardingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), req.ID)...)
}

// --- Helper functions ---

func expandOnboardingParams(ctx context.Context, model *guildOnboardingModel) (*discord.ModifyOnboardingParams, diag.Diagnostics) {
	var diags diag.Diagnostics

	enabled := model.Enabled.ValueBool()
	params := &discord.ModifyOnboardingParams{
		Enabled: &enabled,
	}

	if !model.Mode.IsNull() && !model.Mode.IsUnknown() {
		v := int(model.Mode.ValueInt64())
		params.Mode = &v
	}

	// Default channel IDs.
	params.DefaultChannelIDs = expandSnowflakeSet(ctx, model.DefaultChannelIDs, &diags)
	if diags.HasError() {
		return nil, diags
	}

	// Prompts.
	var promptModels []promptModel
	diags.Append(model.Prompts.ElementsAs(ctx, &promptModels, false)...)
	if diags.HasError() {
		return nil, diags
	}

	prompts := make([]*discord.OnboardingPrompt, 0, len(promptModels))
	for _, pm := range promptModels {
		prompt := &discord.OnboardingPrompt{
			Type:         int(pm.Type.ValueInt64()),
			Title:        pm.Title.ValueString(),
			SingleSelect: pm.SingleSelect.ValueBool(),
			Required:     pm.Required.ValueBool(),
			InOnboarding: pm.InOnboarding.ValueBool(),
		}

		// If the prompt already has an ID from a prior apply, pass it through so Discord
		// can update in place rather than create a new one.
		if !pm.ID.IsNull() && !pm.ID.IsUnknown() {
			prompt.ID = discord.Snowflake(pm.ID.ValueString())
		}

		// Options.
		var optModels []optionModel
		diags.Append(pm.Options.ElementsAs(ctx, &optModels, false)...)
		if diags.HasError() {
			return nil, diags
		}

		opts := make([]*discord.OnboardingPromptOption, 0, len(optModels))
		for _, om := range optModels {
			opt := &discord.OnboardingPromptOption{
				Title: om.Title.ValueString(),
			}

			if !om.ID.IsNull() && !om.ID.IsUnknown() {
				opt.ID = discord.Snowflake(om.ID.ValueString())
			}

			if !om.Description.IsNull() && !om.Description.IsUnknown() {
				v := om.Description.ValueString()
				opt.Description = &v
			}

			opt.ChannelIDs = expandSnowflakeSet(ctx, om.ChannelIDs, &diags)
			opt.RoleIDs = expandSnowflakeSet(ctx, om.RoleIDs, &diags)

			if !om.EmojiID.IsNull() && !om.EmojiID.IsUnknown() {
				s := discord.Snowflake(om.EmojiID.ValueString())
				opt.EmojiID = &s
			}
			if !om.EmojiName.IsNull() && !om.EmojiName.IsUnknown() {
				v := om.EmojiName.ValueString()
				opt.EmojiName = &v
			}
			if !om.EmojiAnimated.IsNull() && !om.EmojiAnimated.IsUnknown() {
				v := om.EmojiAnimated.ValueBool()
				opt.EmojiAnimated = &v
			}

			opts = append(opts, opt)
		}
		prompt.Options = opts

		prompts = append(prompts, prompt)
	}
	params.Prompts = prompts

	return params, diags
}

func expandSnowflakeSet(ctx context.Context, set types.Set, diags *diag.Diagnostics) []discord.Snowflake {
	if set.IsNull() || set.IsUnknown() {
		return nil
	}
	var vals []string
	diags.Append(set.ElementsAs(ctx, &vals, false)...)
	result := make([]discord.Snowflake, len(vals))
	for i, v := range vals {
		result[i] = discord.Snowflake(v)
	}
	return result
}

func flattenOnboarding(ctx context.Context, ob *discord.GuildOnboarding, model *guildOnboardingModel) diag.Diagnostics {
	var diags diag.Diagnostics

	model.GuildID = types.StringValue(ob.GuildID.String())
	model.Enabled = types.BoolValue(ob.Enabled)
	model.Mode = types.Int64Value(int64(ob.Mode))

	// Default channel IDs.
	model.DefaultChannelIDs = flattenSnowflakeSet(ob.DefaultChannelIDs, &diags)

	// Prompts.
	promptVals := make([]attr.Value, 0, len(ob.Prompts))
	for _, p := range ob.Prompts {
		optVals := make([]attr.Value, 0, len(p.Options))
		for _, o := range p.Options {
			var desc types.String
			if o.Description != nil {
				desc = types.StringValue(*o.Description)
			} else {
				desc = types.StringNull()
			}

			channelIDs := flattenSnowflakeSet(o.ChannelIDs, &diags)
			roleIDs := flattenSnowflakeSet(o.RoleIDs, &diags)

			var emojiID types.String
			if o.EmojiID != nil {
				emojiID = types.StringValue(o.EmojiID.String())
			} else {
				emojiID = types.StringNull()
			}

			var emojiName types.String
			if o.EmojiName != nil {
				emojiName = types.StringValue(*o.EmojiName)
			} else {
				emojiName = types.StringNull()
			}

			var emojiAnimated types.Bool
			if o.EmojiAnimated != nil {
				emojiAnimated = types.BoolValue(*o.EmojiAnimated)
			} else {
				emojiAnimated = types.BoolNull()
			}

			optObj, d := types.ObjectValue(optionAttrTypes(), map[string]attr.Value{
				"id":             types.StringValue(o.ID.String()),
				"title":          types.StringValue(o.Title),
				"description":    desc,
				"channel_ids":    channelIDs,
				"role_ids":       roleIDs,
				"emoji_id":       emojiID,
				"emoji_name":     emojiName,
				"emoji_animated": emojiAnimated,
			})
			diags.Append(d...)
			optVals = append(optVals, optObj)
		}

		optsList, d := types.ListValue(types.ObjectType{AttrTypes: optionAttrTypes()}, optVals)
		diags.Append(d...)

		promptObj, d := types.ObjectValue(promptAttrTypes(), map[string]attr.Value{
			"id":            types.StringValue(p.ID.String()),
			"type":          types.Int64Value(int64(p.Type)),
			"title":         types.StringValue(p.Title),
			"single_select": types.BoolValue(p.SingleSelect),
			"required":      types.BoolValue(p.Required),
			"in_onboarding": types.BoolValue(p.InOnboarding),
			"options":       optsList,
		})
		diags.Append(d...)
		promptVals = append(promptVals, promptObj)
	}

	promptsList, d := types.ListValue(types.ObjectType{AttrTypes: promptAttrTypes()}, promptVals)
	diags.Append(d...)
	model.Prompts = promptsList

	return diags
}

func flattenSnowflakeSet(snowflakes []discord.Snowflake, diags *diag.Diagnostics) types.Set {
	if len(snowflakes) == 0 {
		return types.SetNull(types.StringType)
	}
	vals := make([]attr.Value, len(snowflakes))
	for i, s := range snowflakes {
		vals[i] = types.StringValue(s.String())
	}
	set, d := types.SetValue(types.StringType, vals)
	diags.Append(d...)
	return set
}
