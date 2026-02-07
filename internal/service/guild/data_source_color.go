package guild

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource = &colorDataSource{}
)

// colorDataSource is the data source implementation.
type colorDataSource struct{}

// colorDataSourceModel maps the data source schema data.
type colorDataSourceModel struct {
	Hex types.String `tfsdk:"hex"`
	RGB types.Object `tfsdk:"rgb"`
	Int types.Int64  `tfsdk:"int"`
}

// NewColorDataSource returns a new color data source.
func NewColorDataSource() datasource.DataSource {
	return &colorDataSource{}
}

// Metadata returns the data source type name.
func (d *colorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_color"
}

// rgbAttrTypes returns the attr.Type map for the rgb object.
func rgbAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"r": types.Int64Type,
		"g": types.Int64Type,
		"b": types.Int64Type,
	}
}

// Schema defines the schema for the data source.
func (d *colorDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A utility data source that converts a hex color string or RGB values to a Discord color integer. " +
			"This is not backed by an API call. Provide either hex or rgb.",
		Attributes: map[string]schema.Attribute{
			"hex": schema.StringAttribute{
				Description: "The hex color string (e.g. \"#FF0000\").",
				Optional:    true,
			},
			"rgb": schema.SingleNestedAttribute{
				Description: "The RGB color values.",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"r": schema.Int64Attribute{
						Description: "Red component (0-255).",
						Required:    true,
					},
					"g": schema.Int64Attribute{
						Description: "Green component (0-255).",
						Required:    true,
					},
					"b": schema.Int64Attribute{
						Description: "Blue component (0-255).",
						Required:    true,
					},
				},
			},
			"int": schema.Int64Attribute{
				Description: "The computed color as an integer.",
				Computed:    true,
			},
		},
	}
}

// Read computes the color integer from hex or rgb inputs.
func (d *colorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config colorDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	hasHex := !config.Hex.IsNull() && !config.Hex.IsUnknown()
	hasRGB := !config.RGB.IsNull() && !config.RGB.IsUnknown()

	if !hasHex && !hasRGB {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"At least one of 'hex' or 'rgb' must be provided.",
		)
		return
	}

	var colorInt int64

	if hasHex {
		hex := config.Hex.ValueString()
		hex = strings.TrimPrefix(hex, "#")
		if len(hex) != 6 {
			resp.Diagnostics.AddError(
				"Invalid Hex Color",
				fmt.Sprintf("Hex color must be 6 characters (with optional # prefix), got: %q", config.Hex.ValueString()),
			)
			return
		}
		var r, g, b int64
		_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid Hex Color",
				fmt.Sprintf("Could not parse hex color %q: %s", config.Hex.ValueString(), err.Error()),
			)
			return
		}
		colorInt = (r << 16) | (g << 8) | b
	} else {
		// Parse RGB from the object.
		attrs := config.RGB.Attributes()
		rVal, ok := attrs["r"].(types.Int64)
		if !ok {
			resp.Diagnostics.AddError("Invalid RGB", "Could not read 'r' value.")
			return
		}
		gVal, ok := attrs["g"].(types.Int64)
		if !ok {
			resp.Diagnostics.AddError("Invalid RGB", "Could not read 'g' value.")
			return
		}
		bVal, ok := attrs["b"].(types.Int64)
		if !ok {
			resp.Diagnostics.AddError("Invalid RGB", "Could not read 'b' value.")
			return
		}

		r := rVal.ValueInt64()
		g := gVal.ValueInt64()
		b := bVal.ValueInt64()

		if r < 0 || r > 255 || g < 0 || g > 255 || b < 0 || b > 255 {
			resp.Diagnostics.AddError(
				"Invalid RGB Values",
				"RGB values must be between 0 and 255.",
			)
			return
		}

		colorInt = (r << 16) | (g << 8) | b
	}

	config.Int = types.Int64Value(colorInt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
