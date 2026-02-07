package soundboard

import (
	"context"
	"fmt"
	"strings"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &soundboardSoundResource{}
	_ resource.ResourceWithConfigure   = &soundboardSoundResource{}
	_ resource.ResourceWithImportState = &soundboardSoundResource{}
)

// soundboardSoundResource is the resource implementation.
type soundboardSoundResource struct {
	client *discord.Client
}

// soundboardSoundModel maps the resource schema data.
type soundboardSoundModel struct {
	ID        types.String  `tfsdk:"id"`
	GuildID   types.String  `tfsdk:"guild_id"`
	Name      types.String  `tfsdk:"name"`
	Volume    types.Float64 `tfsdk:"volume"`
	EmojiID   types.String  `tfsdk:"emoji_id"`
	EmojiName types.String  `tfsdk:"emoji_name"`
	Available types.Bool    `tfsdk:"available"`
}

// NewSoundboardSoundResource is a helper function to simplify the provider implementation.
func NewSoundboardSoundResource() resource.Resource {
	return &soundboardSoundResource{}
}

// Metadata returns the resource type name.
func (r *soundboardSoundResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_soundboard_sound"
}

// nameValidator validates that the name is between 2 and 32 characters.
type nameValidator struct{}

func (v nameValidator) Description(_ context.Context) string {
	return "name must be between 2 and 32 characters"
}

func (v nameValidator) MarkdownDescription(_ context.Context) string {
	return "name must be between 2 and 32 characters"
}

func (v nameValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	val := req.ConfigValue.ValueString()
	if len(val) < 2 || len(val) > 32 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Name Length",
			fmt.Sprintf("Name must be between 2 and 32 characters, got %d.", len(val)),
		)
	}
}

// volumeValidator validates that the volume is between 0.0 and 1.0.
type volumeValidator struct{}

func (v volumeValidator) Description(_ context.Context) string {
	return "volume must be between 0.0 and 1.0"
}

func (v volumeValidator) MarkdownDescription(_ context.Context) string {
	return "volume must be between 0.0 and 1.0"
}

func (v volumeValidator) ValidateFloat64(_ context.Context, req validator.Float64Request, resp *validator.Float64Response) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	val := req.ConfigValue.ValueFloat64()
	if val < 0.0 || val > 1.0 {
		resp.Diagnostics.AddAttributeError(
			req.Path,
			"Invalid Volume",
			fmt.Sprintf("Volume must be between 0.0 and 1.0, got %f.", val),
		)
	}
}

// Schema defines the schema for the resource.
func (r *soundboardSoundResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord soundboard sound in a guild.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the soundboard sound.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"guild_id": schema.StringAttribute{
				Description: "The ID of the guild this soundboard sound belongs to.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The name of the soundboard sound (2-32 characters).",
				Required:    true,
				Validators: []validator.String{
					nameValidator{},
				},
			},
			"volume": schema.Float64Attribute{
				Description: "The volume of the soundboard sound (0.0 to 1.0). Defaults to 1.0.",
				Optional:    true,
				Computed:    true,
				Default:     float64default.StaticFloat64(1.0),
				Validators: []validator.Float64{
					volumeValidator{},
				},
			},
			"emoji_id": schema.StringAttribute{
				Description: "The ID of the custom emoji for this sound.",
				Optional:    true,
			},
			"emoji_name": schema.StringAttribute{
				Description: "The unicode emoji character for this sound.",
				Optional:    true,
			},
			"available": schema.BoolAttribute{
				Description: "Whether the sound is available for use.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *soundboardSoundResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// Create creates the resource and sets the initial Terraform state.
func (r *soundboardSoundResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan soundboardSoundModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.CreateSoundboardSoundParams{
		Name:  plan.Name.ValueString(),
		Sound: "data:audio/ogg;base64,", // Placeholder: sound upload not supported; use import.
	}

	volume := plan.Volume.ValueFloat64()
	params.Volume = &volume

	if !plan.EmojiID.IsNull() && !plan.EmojiID.IsUnknown() {
		eid := discord.Snowflake(plan.EmojiID.ValueString())
		params.EmojiID = &eid
	}
	if !plan.EmojiName.IsNull() && !plan.EmojiName.IsUnknown() {
		en := plan.EmojiName.ValueString()
		params.EmojiName = &en
	}

	sound, err := r.client.CreateGuildSoundboardSound(ctx, discord.Snowflake(plan.GuildID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Soundboard Sound",
			"Could not create soundboard sound: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(sound.SoundID.String())
	plan.Available = types.BoolValue(sound.Available)
	plan.Volume = types.Float64Value(sound.Volume)
	if sound.EmojiID != nil {
		plan.EmojiID = types.StringValue(sound.EmojiID.String())
	} else {
		plan.EmojiID = types.StringNull()
	}
	if sound.EmojiName != nil {
		plan.EmojiName = types.StringValue(*sound.EmojiName)
	} else {
		plan.EmojiName = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *soundboardSoundResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state soundboardSoundModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	sounds, err := r.client.ListGuildSoundboardSounds(ctx, discord.Snowflake(state.GuildID.ValueString()))
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Soundboard Sound",
			"Could not read soundboard sound: "+err.Error(),
		)
		return
	}

	var found *discord.SoundboardSound
	for _, s := range sounds {
		if s.SoundID.String() == state.ID.ValueString() {
			found = s
			break
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(found.Name)
	state.Volume = types.Float64Value(found.Volume)
	state.Available = types.BoolValue(found.Available)
	if found.EmojiID != nil {
		state.EmojiID = types.StringValue(found.EmojiID.String())
	} else {
		state.EmojiID = types.StringNull()
	}
	if found.EmojiName != nil {
		state.EmojiName = types.StringValue(*found.EmojiName)
	} else {
		state.EmojiName = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *soundboardSoundResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan soundboardSoundModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state soundboardSoundModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	volume := plan.Volume.ValueFloat64()
	params := &discord.ModifySoundboardSoundParams{
		Name:   &name,
		Volume: &volume,
	}

	if !plan.EmojiID.IsNull() && !plan.EmojiID.IsUnknown() {
		eid := discord.Snowflake(plan.EmojiID.ValueString())
		params.EmojiID = &eid
	}
	if !plan.EmojiName.IsNull() && !plan.EmojiName.IsUnknown() {
		en := plan.EmojiName.ValueString()
		params.EmojiName = &en
	}

	sound, err := r.client.ModifyGuildSoundboardSound(
		ctx,
		discord.Snowflake(plan.GuildID.ValueString()),
		discord.Snowflake(state.ID.ValueString()),
		params,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Soundboard Sound",
			"Could not update soundboard sound: "+err.Error(),
		)
		return
	}

	plan.ID = state.ID
	plan.Available = types.BoolValue(sound.Available)
	plan.Volume = types.Float64Value(sound.Volume)
	if sound.EmojiID != nil {
		plan.EmojiID = types.StringValue(sound.EmojiID.String())
	} else {
		plan.EmojiID = types.StringNull()
	}
	if sound.EmojiName != nil {
		plan.EmojiName = types.StringValue(*sound.EmojiName)
	} else {
		plan.EmojiName = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *soundboardSoundResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state soundboardSoundModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGuildSoundboardSound(
		ctx,
		discord.Snowflake(state.GuildID.ValueString()),
		discord.Snowflake(state.ID.ValueString()),
	)
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Soundboard Sound",
			"Could not delete soundboard sound: "+err.Error(),
		)
	}
}

// ImportState implements the import by guild_id/sound_id.
func (r *soundboardSoundResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Expected import ID format: guild_id/sound_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("guild_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
