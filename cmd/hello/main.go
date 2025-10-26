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

// Init initializes the widget tree
func (app *WidgetApp) Init() (err error) {
	// Create four fill widgets with different colors
	redFill := widget.NewFlexFill(
		1.0, 0.0, 0.0, 1.0, // Red color (RGBA)
		0, 0, // Min width/height (flexible)
		1e9, 1e9, // Max width/height (very large)
	)

	yellowFill := widget.NewFlexFill(
		1.0, 1.0, 0.0, 1.0, // Yellow color (RGBA)
		0, 0, // Min width/height (flexible)
		1e9, 1e9, // Max width/height (very large)
	)

	greenFill := widget.NewFlexFill(
		0.0, 1.0, 0.0, 1.0, // Green color (RGBA)
		0, 0, // Min width/height (flexible)
		1e9, 1e9, // Max width/height (very large)
	)

	blueFill := widget.NewFlexFill(
		0.0, 0.0, 1.0, 1.0, // Blue color (RGBA)
		0, 0, // Min width/height (flexible)
		1e9, 1e9, // Max width/height (very large)
	)

	// Create first row container (red and yellow)
	topRow := widget.NewContainer(
		widget.DirectionRow,
		widget.NewFlexConstraints(0, 0, 1e9, 1e9), // Flexible constraints
	)
	topRow.AddChild(widget.NewFlexChild(redFill, 1.0))    // Equal weight
	topRow.AddChild(widget.NewFlexChild(yellowFill, 1.0)) // Equal weight

	// Create second row container (green and blue)
	bottomRow := widget.NewContainer(
		widget.DirectionRow,
		widget.NewFlexConstraints(0, 0, 1e9, 1e9), // Flexible constraints
	)
	bottomRow.AddChild(widget.NewFlexChild(greenFill, 1.0)) // Equal weight
	bottomRow.AddChild(widget.NewFlexChild(blueFill, 1.0))  // Equal weight

	// Create main column container
	mainColumn := widget.NewContainer(
		widget.DirectionColumn,
		widget.NewFlexConstraints(0, 0, 1e9, 1e9), // Flexible constraints
	)
	mainColumn.AddChild(widget.NewFlexChild(topRow, 1.0))    // Equal weight
	mainColumn.AddChild(widget.NewFlexChild(bottomRow, 1.0)) // Equal weight

	// Create a white box with fixed 64x64 size (no position needed)
	// Using 0.5 alpha to test alpha blending
	whiteBox := widget.NewRigidFill(
		1.0, 1.0, 1.0, 0.75, // White with 0.75 alpha
		64, 64, // Fixed 64x64 size
	)

	// Wrap the white box in a DirectionWidget with center gravity
	centeredWhiteBox := widget.NewDirectionWidget(
		whiteBox,
		widget.GravityCenter,
		widget.NewFlexConstraints(0, 0, 1e9, 1e9), // Flexible constraints to fill available space
	)

	// Create overlay widget to demonstrate overpainting
	overlay := widget.NewOverlayWidget(
		widget.NewFlexConstraints(0, 0, 1e9, 1e9), // Flexible constraints
	)

	// Add the flex layout first (background)
	overlay.AddChild(mainColumn)
	// Add the centered white box second (foreground - will paint over the flex layout)
	overlay.AddChild(centeredWhiteBox)

	// Create root widget with the overlay as child
	app.rootWidget = widget.NewRootWidget(overlay)

	return
}

// Render renders the widget tree
func (app *WidgetApp) Render(width, height int, mouseX, mouseY float64) (err error) {
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

	// Draw crosshair at mouse cursor position
	drawCrosshair(float32(mouseX), float32(height)-float32(mouseY), width, height)

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
