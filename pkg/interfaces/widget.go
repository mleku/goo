package interfaces

// Point represents a 2D coordinate
type Point struct {
	X, Y float32
}

// Size represents width and height dimensions
type Size struct {
	Width, Height float32
}

// Constraints define minimum and maximum size limits and position within root widget
type Constraints struct {
	MinWidth, MinHeight float32
	MaxWidth, MaxHeight float32
	// Top/Left coordinates relative to root widget (0,0 = top-left of canvas)
	Top, Left float32
}

// Box represents the layout box for a widget with position and size
type Box struct {
	// Position relative to parent's top-left corner
	Position Point
	// Actual size of the box
	Size Size
	// Size constraints
	Constraints Constraints
}

// Rect represents a rectangular region
type Rect struct {
	X, Y          float32
	Width, Height float32
}

// Context provides the rendering context for widgets
type Context struct {
	// Window size
	WindowWidth, WindowHeight int
	// Parent box - widget's position is relative to this
	ParentBox *Box
	// Available space within parent
	AvailableSize Size
	// Painted regions to avoid double painting
	PaintedRegions []Rect
}

// Widget defines the interface that all widgets must implement
type Widget interface {
	// Render draws the widget within the given box and returns the actual size used
	Render(ctx *Context, box *Box) (usedSize Size, err error)
	// GetConstraints returns the size constraints for this widget
	GetConstraints() Constraints
}
