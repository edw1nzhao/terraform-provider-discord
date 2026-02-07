package message

import (
	"context"
	"strings"

	"github.com/edw1nzhao/terraform-provider-discord/internal/discord"
	"github.com/edw1nzhao/terraform-provider-discord/internal/service/common"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &messageResource{}
	_ resource.ResourceWithConfigure   = &messageResource{}
	_ resource.ResourceWithImportState = &messageResource{}
)

// messageResource is the resource implementation.
type messageResource struct {
	client *discord.Client
}

// embedModel maps the embed block schema data.
type embedModel struct {
	Title        types.String `tfsdk:"title"`
	Description  types.String `tfsdk:"description"`
	URL          types.String `tfsdk:"url"`
	Color        types.Int64  `tfsdk:"color"`
	FooterText   types.String `tfsdk:"footer_text"`
	ImageURL     types.String `tfsdk:"image_url"`
	ThumbnailURL types.String `tfsdk:"thumbnail_url"`
}

// messageModel maps the resource schema data.
type messageModel struct {
	ID        types.String `tfsdk:"id"`
	ChannelID types.String `tfsdk:"channel_id"`
	Content   types.String `tfsdk:"content"`
	TTS       types.Bool   `tfsdk:"tts"`
	Embed     []embedModel `tfsdk:"embed"`
}

// NewMessageResource is a helper function to simplify the provider implementation.
func NewMessageResource() resource.Resource {
	return &messageResource{}
}

// Metadata returns the resource type name.
func (r *messageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_message"
}

// Schema defines the schema for the resource.
func (r *messageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Discord message in a channel.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the message.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"channel_id": schema.StringAttribute{
				Description: "The ID of the channel to send the message in.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Description: "The content of the message.",
				Optional:    true,
			},
			"tts": schema.BoolAttribute{
				Description: "Whether this is a text-to-speech message.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolRequiresReplace{},
				},
			},
		},
		Blocks: map[string]schema.Block{
			"embed": schema.ListNestedBlock{
				Description: "Embedded rich content.",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"title": schema.StringAttribute{
							Description: "Title of the embed.",
							Optional:    true,
						},
						"description": schema.StringAttribute{
							Description: "Description of the embed.",
							Optional:    true,
						},
						"url": schema.StringAttribute{
							Description: "URL of the embed.",
							Optional:    true,
						},
						"color": schema.Int64Attribute{
							Description: "Color code of the embed.",
							Optional:    true,
						},
						"footer_text": schema.StringAttribute{
							Description: "Footer text of the embed.",
							Optional:    true,
						},
						"image_url": schema.StringAttribute{
							Description: "Image URL of the embed.",
							Optional:    true,
						},
						"thumbnail_url": schema.StringAttribute{
							Description: "Thumbnail URL of the embed.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

// boolRequiresReplace is a plan modifier that requires replacement when the value changes.
type boolRequiresReplace struct{}

func (m boolRequiresReplace) Description(_ context.Context) string {
	return "Requires replacement when the value changes."
}

func (m boolRequiresReplace) MarkdownDescription(_ context.Context) string {
	return "Requires replacement when the value changes."
}

func (m boolRequiresReplace) PlanModifyBool(_ context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	if req.StateValue.IsNull() {
		return
	}
	if !req.PlanValue.Equal(req.StateValue) {
		resp.RequiresReplace = true
	}
}

// Configure adds the provider configured client to the resource.
func (r *messageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.client = common.ClientFromProviderData(req.ProviderData, &resp.Diagnostics)
}

// buildEmbeds converts the embed models to Discord API embed objects.
func buildEmbeds(embeds []embedModel) []*discord.Embed {
	if len(embeds) == 0 {
		return nil
	}
	result := make([]*discord.Embed, 0, len(embeds))
	for _, e := range embeds {
		embed := &discord.Embed{}
		if !e.Title.IsNull() && !e.Title.IsUnknown() {
			t := e.Title.ValueString()
			embed.Title = &t
		}
		if !e.Description.IsNull() && !e.Description.IsUnknown() {
			d := e.Description.ValueString()
			embed.Description = &d
		}
		if !e.URL.IsNull() && !e.URL.IsUnknown() {
			u := e.URL.ValueString()
			embed.URL = &u
		}
		if !e.Color.IsNull() && !e.Color.IsUnknown() {
			c := int(e.Color.ValueInt64())
			embed.Color = &c
		}
		if !e.FooterText.IsNull() && !e.FooterText.IsUnknown() {
			embed.Footer = &discord.EmbedFooter{
				Text: e.FooterText.ValueString(),
			}
		}
		if !e.ImageURL.IsNull() && !e.ImageURL.IsUnknown() {
			u := e.ImageURL.ValueString()
			embed.Image = &discord.EmbedImage{
				URL: &u,
			}
		}
		if !e.ThumbnailURL.IsNull() && !e.ThumbnailURL.IsUnknown() {
			u := e.ThumbnailURL.ValueString()
			embed.Thumbnail = &discord.EmbedImage{
				URL: &u,
			}
		}
		result = append(result, embed)
	}
	return result
}

// flattenEmbeds converts Discord API embed objects to embed models.
func flattenEmbeds(embeds []*discord.Embed) []embedModel {
	if len(embeds) == 0 {
		return nil
	}
	result := make([]embedModel, 0, len(embeds))
	for _, e := range embeds {
		m := embedModel{}
		if e.Title != nil {
			m.Title = types.StringValue(*e.Title)
		} else {
			m.Title = types.StringNull()
		}
		if e.Description != nil {
			m.Description = types.StringValue(*e.Description)
		} else {
			m.Description = types.StringNull()
		}
		if e.URL != nil {
			m.URL = types.StringValue(*e.URL)
		} else {
			m.URL = types.StringNull()
		}
		if e.Color != nil {
			m.Color = types.Int64Value(int64(*e.Color))
		} else {
			m.Color = types.Int64Null()
		}
		if e.Footer != nil {
			m.FooterText = types.StringValue(e.Footer.Text)
		} else {
			m.FooterText = types.StringNull()
		}
		if e.Image != nil && e.Image.URL != nil {
			m.ImageURL = types.StringValue(*e.Image.URL)
		} else {
			m.ImageURL = types.StringNull()
		}
		if e.Thumbnail != nil && e.Thumbnail.URL != nil {
			m.ThumbnailURL = types.StringValue(*e.Thumbnail.URL)
		} else {
			m.ThumbnailURL = types.StringNull()
		}
		result = append(result, m)
	}
	return result
}

// Create creates the resource and sets the initial Terraform state.
func (r *messageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan messageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.CreateMessageParams{}

	if !plan.Content.IsNull() && !plan.Content.IsUnknown() {
		c := plan.Content.ValueString()
		params.Content = &c
	}

	if !plan.TTS.IsNull() && !plan.TTS.IsUnknown() {
		tts := plan.TTS.ValueBool()
		params.TTS = &tts
	}

	params.Embeds = buildEmbeds(plan.Embed)

	msg, err := r.client.CreateMessage(ctx, discord.Snowflake(plan.ChannelID.ValueString()), params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Message",
			"Could not create message: "+err.Error(),
		)
		return
	}

	plan.ID = types.StringValue(msg.ID.String())
	plan.TTS = types.BoolValue(msg.TTS)
	if msg.Content != "" {
		plan.Content = types.StringValue(msg.Content)
	}
	if len(msg.Embeds) > 0 {
		plan.Embed = flattenEmbeds(msg.Embeds)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *messageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state messageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	msg, err := r.client.GetChannelMessage(
		ctx,
		discord.Snowflake(state.ChannelID.ValueString()),
		discord.Snowflake(state.ID.ValueString()),
	)
	if err != nil {
		if discord.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error Reading Message",
			"Could not read message: "+err.Error(),
		)
		return
	}

	state.Content = types.StringValue(msg.Content)
	if msg.Content == "" {
		state.Content = types.StringNull()
	}
	state.TTS = types.BoolValue(msg.TTS)
	if len(msg.Embeds) > 0 {
		state.Embed = flattenEmbeds(msg.Embeds)
	} else {
		state.Embed = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *messageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan messageModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state messageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := &discord.EditMessageParams{}

	if !plan.Content.IsNull() && !plan.Content.IsUnknown() {
		c := plan.Content.ValueString()
		params.Content = &c
	} else {
		empty := ""
		params.Content = &empty
	}

	params.Embeds = buildEmbeds(plan.Embed)

	msg, err := r.client.EditMessage(
		ctx,
		discord.Snowflake(plan.ChannelID.ValueString()),
		discord.Snowflake(state.ID.ValueString()),
		params,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Message",
			"Could not update message: "+err.Error(),
		)
		return
	}

	plan.ID = state.ID
	plan.Content = types.StringValue(msg.Content)
	if msg.Content == "" {
		plan.Content = types.StringNull()
	}
	plan.TTS = types.BoolValue(msg.TTS)
	if len(msg.Embeds) > 0 {
		plan.Embed = flattenEmbeds(msg.Embeds)
	} else {
		plan.Embed = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *messageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state messageModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteMessage(
		ctx,
		discord.Snowflake(state.ChannelID.ValueString()),
		discord.Snowflake(state.ID.ValueString()),
	)
	if err != nil {
		if discord.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError(
			"Error Deleting Message",
			"Could not delete message: "+err.Error(),
		)
	}
}

// ImportState implements the import by channel_id/message_id.
func (r *messageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			"Expected import ID format: channel_id/message_id",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("channel_id"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
