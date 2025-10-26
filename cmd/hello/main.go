package main

import (
	"github.com/go-gl/gl/all-core/gl"
	"github.com/mleku/goo/pkg/interfaces"
	"github.com/mleku/goo/pkg/widget"
	"github.com/mleku/goo/pkg/window"
	"lol.mleku.dev/chk"
)

// WidgetApp implements the window application
type WidgetApp struct {
	rootWidget *widget.RootWidget
}

// Init initializes the widget tree using the chained API with inline creation
func (app *WidgetApp) Init() (err error) {
	app.rootWidget = widget.Root(
		widget.Overlay().
			Child(
				widget.Column().
					Flex(
						widget.Row().
							Flex(widget.Fill(1.0, 0.0, 0.0, 1.0), 1.0).
							Flex(widget.Fill(1.0, 1.0, 0.0, 1.0), 1.0),
						1.0,
					).
					Flex(
						widget.Row().
							Flex(widget.Fill(0.0, 1.0, 0.0, 1.0), 1.0).
							Flex(widget.Fill(0.0, 0.0, 1.0, 1.0), 1.0),
						1.0,
					),
			).
			Child(
				widget.Center(
					widget.Fill(1.0, 1.0, 1.0, 0.75, widget.NewRigidConstraints(64, 64)),
				),
			),
	)

	return
}

// Render renders the widget tree
func (app *WidgetApp) Render(width, height int, mouseX, mouseY float64, cursorInWindow bool) (err error) {
	// Set the clear color to black
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Enable blending globally
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Enable(gl.SCISSOR_TEST)

	// Set up the projection matrix for 2D rendering
	// Use orthographic projection matching screen coordinates
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(width), 0, float64(height), -1, 1)

	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()

	// Create widget context with window dimensions
	widgetCtx := &interfaces.Context{
		WindowWidth:    width,  // Window logical size
		WindowHeight:   height, // Window logical size
		PaintedRegions: make([]interfaces.Rect, 0),
	}

	// Create a dummy box for the root widget
	rootBox := &interfaces.Box{}

	// Render the widget tree
	_, err = app.rootWidget.Render(widgetCtx, rootBox)
	if err != nil {
		return
	}

	// Draw crosshair at mouse cursor position only if cursor is in window
	if cursorInWindow {
		drawCrosshair(float32(mouseX), float32(height)-float32(mouseY), width, height)
	}

	return
}

// drawCrosshair draws a 1-pixel wide black crosshair at the specified position
func drawCrosshair(x, y float32, width, height int) {
	// Disable scissor test for crosshair to draw over everything
	gl.Disable(gl.SCISSOR_TEST)

	// Set line width to 1 pixel
	gl.LineWidth(1.0)

	// Set color to black
	gl.Color4f(0.0, 0.0, 0.0, 1.0)

	// Draw vertical line (full height)
	gl.Begin(gl.LINES)
	gl.Vertex2f(x, 0)
	gl.Vertex2f(x, float32(height))
	gl.End()

	// Draw horizontal line (full width)
	gl.Begin(gl.LINES)
	gl.Vertex2f(0, y)
	gl.Vertex2f(float32(width), y)
	gl.End()

	// Re-enable scissor test
	gl.Enable(gl.SCISSOR_TEST)
}

func main() {
	w, err := window.New(640, 480, "Fromage Widget Demo with GLFW")
	if chk.E(err) {
		return
	}

	app := &WidgetApp{}
	if err := app.Init(); chk.E(err) {
		return
	}

	if err := w.Run(app.Render); chk.E(err) {
		return
	}
}
