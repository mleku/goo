package widget

import (
	"github.com/go-gl/gl/all-core/gl"
)

// Filler is a widget that fills its box with a solid color
type Filler struct {
	color [4]float32
}

// Fill creates a new Fill widget that fills its container with the specified color.
// The fill always fills to the edge of its box when calculated.
func Fill(red, green, blue, alpha float32) *Filler {
	return &Filler{
		color: [4]float32{red, green, blue, alpha},
	}
}

// SetColor updates the fill color
func (f *Filler) SetColor(red, green, blue, alpha float32) {
	f.color = [4]float32{red, green, blue, alpha}
}

// GetConstraints returns the size constraints for this Fill widget
func (f *Filler) GetConstraints() Constraints {
	// Fill widgets always have flexible constraints to fill their container
	return NewFlexConstraints(0, 0, 1e9, 1e9)
}

// Render implements the Widget interface for Fill
func (f *Filler) Render(ctx *Context, box *Box) (usedSize Size, err error) {
	// Convert to OpenGL coordinates (bottom-left origin) from top-left origin
	// Box Y is from top, convert to bottom-left
	bottomY := float32(ctx.WindowHeight) - box.Position.Y - box.Size.Height

	// Create vertices for the quad (exact pixel boundaries)
	x1, y1 := box.Position.X, float32(ctx.WindowHeight)-box.Position.Y
	x2, y2 := box.Position.X+box.Size.Width, float32(ctx.WindowHeight)-box.Position.Y
	x3, y3 := box.Position.X+box.Size.Width, float32(ctx.WindowHeight)-box.Position.Y-box.Size.Height
	x4, y4 := box.Position.X, float32(ctx.WindowHeight)-box.Position.Y-box.Size.Height

	// Enable scissor test with proper boundaries to enforce clipping
	// This prevents artifacts at the edges by ensuring rendering doesn't go beyond boundaries
	gl.Scissor(
		int32(box.Position.X),
		int32(bottomY),
		int32(box.Size.Width),
		int32(box.Size.Height),
	)

	// Set the color
	gl.Color4f(f.color[0], f.color[1], f.color[2], f.color[3])

	// Draw using immediate mode
	// The scissor test will clip these vertices to the exact boundaries
	gl.Begin(gl.QUADS)
	gl.Vertex2f(x1, y1)
	gl.Vertex2f(x2, y2)
	gl.Vertex2f(x3, y3)
	gl.Vertex2f(x4, y4)
	gl.End()

	return box.Size, nil
}
