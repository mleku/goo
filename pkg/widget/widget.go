package widget

import (
	"github.com/mleku/goo/pkg/interfaces"
	"lol.mleku.dev/chk"
)

// Re-export types from interfaces package for convenience
type (
	Point       = interfaces.Point
	Size        = interfaces.Size
	Constraints = interfaces.Constraints
	Box         = interfaces.Box
	Context     = interfaces.Context
	Widget      = interfaces.Widget
)

// NewConstraints creates constraints with min/max values and position
func NewConstraints(minWidth, minHeight, maxWidth, maxHeight, top, left float32) Constraints {
	return Constraints{
		MinWidth:  minWidth,
		MinHeight: minHeight,
		MaxWidth:  maxWidth,
		MaxHeight: maxHeight,
		Top:       top,
		Left:      left,
	}
}

// NewConstraintsNoPos creates constraints with min/max values and no specific position
func NewConstraintsNoPos(minWidth, minHeight, maxWidth, maxHeight float32) Constraints {
	return Constraints{
		MinWidth:  minWidth,
		MinHeight: minHeight,
		MaxWidth:  maxWidth,
		MaxHeight: maxHeight,
		Top:       0,
		Left:      0,
	}
}

// NewRigidConstraints creates constraints for a fixed size (rigid widget)
func NewRigidConstraints(width, height float32) Constraints {
	return Constraints{
		MinWidth:  width,
		MinHeight: height,
		MaxWidth:  width,
		MaxHeight: height,
		Top:       0,
		Left:      0,
	}
}

// NewRigidConstraintsAt creates constraints for a fixed size at specific position
func NewRigidConstraintsAt(width, height, top, left float32) Constraints {
	return Constraints{
		MinWidth:  width,
		MinHeight: height,
		MaxWidth:  width,
		MaxHeight: height,
		Top:       top,
		Left:      left,
	}
}

// NewFlexConstraints creates constraints for a flexible widget
func NewFlexConstraints(minWidth, minHeight, maxWidth, maxHeight float32) Constraints {
	return Constraints{
		MinWidth:  minWidth,
		MinHeight: minHeight,
		MaxWidth:  maxWidth,
		MaxHeight: maxHeight,
		Top:       0,
		Left:      0,
	}
}

// NewFlexConstraintsAt creates constraints for a flexible widget at specific position
func NewFlexConstraintsAt(minWidth, minHeight, maxWidth, maxHeight, top, left float32) Constraints {
	return Constraints{
		MinWidth:  minWidth,
		MinHeight: minHeight,
		MaxWidth:  maxWidth,
		MaxHeight: maxHeight,
		Top:       top,
		Left:      left,
	}
}

// NewBox creates a new box with the given position, size, and constraints
func NewBox(x, y, width, height float32, constraints Constraints) *Box {
	return &Box{
		Position:    Point{X: x, Y: y},
		Size:        Size{Width: width, Height: height},
		Constraints: constraints,
	}
}

// Direction specifies layout direction for containers
type Direction int

const (
	DirectionRow Direction = iota
	DirectionColumn
)

// FlexType specifies whether a widget is rigid or flexible
type FlexType int

const (
	FlexTypeRigid FlexType = iota
	FlexTypeFlex
)

// FlexChild represents a child widget in a flex container
type FlexChild struct {
	Widget Widget
	Type   FlexType
	Weight float32 // Only used for FlexTypeFlex
}

// NewRigidChild creates a rigid child widget
func NewRigidChild(widget Widget) FlexChild {
	return FlexChild{
		Widget: widget,
		Type:   FlexTypeRigid,
		Weight: 0,
	}
}

// NewFlexChild creates a flexible child widget with the given weight
func NewFlexChild(widget Widget, weight float32) FlexChild {
	return FlexChild{
		Widget: widget,
		Type:   FlexTypeFlex,
		Weight: weight,
	}
}

// Container is a widget that lays out children in rows or columns
type Container struct {
	Direction   Direction
	Children    []FlexChild
	constraints Constraints
}

// Row creates a new row container with default flexible constraints.
// Chain methods like Flex() or Rigid() to add children.
func Row(constraints ...Constraints) *Container {
	var c Constraints
	if len(constraints) > 0 {
		c = constraints[0]
	} else {
		c = NewFlexConstraints(0, 0, 1e9, 1e9)
	}
	return &Container{
		Direction:   DirectionRow,
		Children:    make([]FlexChild, 0),
		constraints: c,
	}
}

// Column creates a new column container with default flexible constraints.
// Chain methods like Flex() or Rigid() to add children.
func Column(constraints ...Constraints) *Container {
	var c Constraints
	if len(constraints) > 0 {
		c = constraints[0]
	} else {
		c = NewFlexConstraints(0, 0, 1e9, 1e9)
	}
	return &Container{
		Direction:   DirectionColumn,
		Children:    make([]FlexChild, 0),
		constraints: c,
	}
}

// NewContainer creates a new container with the specified direction.
// If no constraints are provided, uses default flexible constraints (0, 0, 1e9, 1e9).
func NewContainer(direction Direction, constraints ...Constraints) *Container {
	var c Constraints
	if len(constraints) > 0 {
		c = constraints[0]
	} else {
		c = NewFlexConstraints(0, 0, 1e9, 1e9)
	}
	return &Container{
		Direction:   direction,
		Children:    make([]FlexChild, 0),
		constraints: c,
	}
}

// AddChild adds a child widget to the container and returns the container for chaining
func (c *Container) AddChild(child FlexChild) *Container {
	c.Children = append(c.Children, child)
	return c
}

// Flex adds a flexible child with the specified weight to the container
func (c *Container) Flex(child Widget, weight float32) *Container {
	c.Children = append(c.Children, FlexChild{
		Widget: child,
		Type:   FlexTypeFlex,
		Weight: weight,
	})
	return c
}

// Rigid adds a rigid child to the container
func (c *Container) Rigid(child Widget) *Container {
	c.Children = append(c.Children, FlexChild{
		Widget: child,
		Type:   FlexTypeRigid,
		Weight: 0,
	})
	return c
}

// GetConstraints returns the container's constraints
func (c *Container) GetConstraints() Constraints {
	return c.constraints
}

// Render implements the Widget interface for Container
func (c *Container) Render(ctx *Context, box *Box) (usedSize Size, err error) {
	if len(c.Children) == 0 {
		return Size{}, nil
	}

	// Calculate layout based on direction
	switch c.Direction {
	case DirectionRow:
		return c.renderRow(ctx, box)
	case DirectionColumn:
		return c.renderColumn(ctx, box)
	default:
		return Size{}, errInvalidDirection
	}
}

// renderRow lays out children horizontally
func (c *Container) renderRow(ctx *Context, box *Box) (usedSize Size, err error) {
	availableWidth := box.Size.Width
	availableHeight := box.Size.Height

	// First pass: calculate rigid sizes and total flex weight
	var rigidWidth float32
	var totalFlexWeight float32
	var maxHeight float32

	for _, child := range c.Children {
		childConstraints := child.Widget.GetConstraints()

		if child.Type == FlexTypeRigid {
			rigidWidth += childConstraints.MinWidth
			if childConstraints.MinHeight > maxHeight {
				maxHeight = childConstraints.MinHeight
			}
		} else {
			totalFlexWeight += child.Weight
			if childConstraints.MinHeight > maxHeight {
				maxHeight = childConstraints.MinHeight
			}
		}
	}

	// Calculate remaining width for flex children
	flexWidth := availableWidth - rigidWidth
	if flexWidth < 0 {
		flexWidth = 0
	}

	// Second pass: render children
	var currentX float32
	var actualUsedWidth float32
	var actualMaxHeight float32

	for _, child := range c.Children {
		childConstraints := child.Widget.GetConstraints()
		var childWidth float32

		if child.Type == FlexTypeRigid {
			childWidth = childConstraints.MinWidth
		} else {
			if totalFlexWeight > 0 {
				childWidth = (flexWidth * child.Weight) / totalFlexWeight
				// Clamp to constraints
				if childWidth < childConstraints.MinWidth {
					childWidth = childConstraints.MinWidth
				}
				if childWidth > childConstraints.MaxWidth {
					childWidth = childConstraints.MaxWidth
				}
			} else {
				childWidth = childConstraints.MinWidth
			}
		}

		// Create child box
		childBox := &Box{
			Position: Point{
				X: box.Position.X + currentX,
				Y: box.Position.Y,
			},
			Size: Size{
				Width:  childWidth,
				Height: availableHeight,
			},
			Constraints: childConstraints,
		}

		// Create child context
		childCtx := &Context{
			WindowWidth:   ctx.WindowWidth,
			WindowHeight:  ctx.WindowHeight,
			ParentBox:     childBox,
			AvailableSize: childBox.Size,
		}

		// Render child
		childUsedSize, err := child.Widget.Render(childCtx, childBox)
		if chk.E(err) {
			return Size{}, err
		}

		currentX += childUsedSize.Width
		actualUsedWidth += childUsedSize.Width

		if childUsedSize.Height > actualMaxHeight {
			actualMaxHeight = childUsedSize.Height
		}
	}

	return Size{Width: actualUsedWidth, Height: actualMaxHeight}, nil
}

// renderColumn lays out children vertically
func (c *Container) renderColumn(ctx *Context, box *Box) (usedSize Size, err error) {
	availableWidth := box.Size.Width
	availableHeight := box.Size.Height

	// First pass: calculate rigid sizes and total flex weight
	var rigidHeight float32
	var totalFlexWeight float32
	var maxWidth float32

	for _, child := range c.Children {
		childConstraints := child.Widget.GetConstraints()

		if child.Type == FlexTypeRigid {
			rigidHeight += childConstraints.MinHeight
			if childConstraints.MinWidth > maxWidth {
				maxWidth = childConstraints.MinWidth
			}
		} else {
			totalFlexWeight += child.Weight
			if childConstraints.MinWidth > maxWidth {
				maxWidth = childConstraints.MinWidth
			}
		}
	}

	// Calculate remaining height for flex children
	flexHeight := availableHeight - rigidHeight
	if flexHeight < 0 {
		flexHeight = 0
	}

	// Second pass: render children
	var currentY float32
	var actualUsedHeight float32
	var actualMaxWidth float32

	for _, child := range c.Children {
		childConstraints := child.Widget.GetConstraints()
		var childHeight float32

		if child.Type == FlexTypeRigid {
			childHeight = childConstraints.MinHeight
		} else {
			if totalFlexWeight > 0 {
				childHeight = (flexHeight * child.Weight) / totalFlexWeight
				// Clamp to constraints
				if childHeight < childConstraints.MinHeight {
					childHeight = childConstraints.MinHeight
				}
				if childHeight > childConstraints.MaxHeight {
					childHeight = childConstraints.MaxHeight
				}
			} else {
				childHeight = childConstraints.MinHeight
			}
		}

		// Create child box
		childBox := &Box{
			Position: Point{
				X: box.Position.X,
				Y: box.Position.Y + currentY,
			},
			Size: Size{
				Width:  availableWidth,
				Height: childHeight,
			},
			Constraints: childConstraints,
		}

		// Create child context
		childCtx := &Context{
			WindowWidth:   ctx.WindowWidth,
			WindowHeight:  ctx.WindowHeight,
			ParentBox:     childBox,
			AvailableSize: childBox.Size,
		}

		// Render child
		childUsedSize, err := child.Widget.Render(childCtx, childBox)
		if chk.E(err) {
			return Size{}, err
		}

		currentY += childUsedSize.Height
		actualUsedHeight += childUsedSize.Height

		if childUsedSize.Width > actualMaxWidth {
			actualMaxWidth = childUsedSize.Width
		}
	}

	return Size{Width: actualMaxWidth, Height: actualUsedHeight}, nil
}

// RootWidget manages the root layout that spans the entire canvas
type RootWidget struct {
	child      Widget
	clearColor [4]float32
}

// Root creates a new root widget with the given child
func Root(child Widget) *RootWidget {
	return &RootWidget{
		child:      child,
		clearColor: [4]float32{0.0, 0.0, 0.0, 1.0}, // Default black
	}
}

// SetClearColor sets the background clear color for the root widget and returns the root for chaining
func (r *RootWidget) SetClearColor(red, green, blue, alpha float32) *RootWidget {
	r.clearColor = [4]float32{red, green, blue, alpha}
	return r
}

// GetConstraints returns unconstrained size (fills canvas)
func (r *RootWidget) GetConstraints() Constraints {
	return Constraints{
		MinWidth:  0,
		MinHeight: 0,
		MaxWidth:  1e9, // Very large number
		MaxHeight: 1e9,
	}
}

// Render implements the Widget interface for RootWidget
func (r *RootWidget) Render(ctx *Context, box *Box) (usedSize Size, err error) {
	if r.child == nil {
		return box.Size, nil
	}

	// Get child constraints to determine positioning
	childConstraints := r.child.GetConstraints()

	// Create a box that spans the entire canvas, but position child based on its constraints
	canvasWidth := float32(ctx.WindowWidth)
	canvasHeight := float32(ctx.WindowHeight)

	// Use constraint coordinates if specified, otherwise fill canvas
	childBox := &Box{
		Position: Point{
			X: childConstraints.Left,
			Y: childConstraints.Top,
		},
		Size: Size{
			Width:  canvasWidth - childConstraints.Left,
			Height: canvasHeight - childConstraints.Top,
		},
		Constraints: childConstraints,
	}

	// If child has specific size constraints, respect them
	if childConstraints.MaxWidth < childBox.Size.Width {
		childBox.Size.Width = childConstraints.MaxWidth
	}
	if childConstraints.MaxHeight < childBox.Size.Height {
		childBox.Size.Height = childConstraints.MaxHeight
	}
	if childConstraints.MinWidth > childBox.Size.Width {
		childBox.Size.Width = childConstraints.MinWidth
	}
	if childConstraints.MinHeight > childBox.Size.Height {
		childBox.Size.Height = childConstraints.MinHeight
	}

	// Create context for child
	childCtx := &Context{
		WindowWidth:   ctx.WindowWidth,
		WindowHeight:  ctx.WindowHeight,
		ParentBox:     childBox,
		AvailableSize: childBox.Size,
	}

	// Render child
	return r.child.Render(childCtx, childBox)
}

// OverlayWidget allows multiple widgets to be rendered on top of each other
type OverlayWidget struct {
	children    []Widget
	constraints Constraints
}

// Overlay creates a new overlay widget that renders children in sequence.
// If no constraints are provided, uses default flexible constraints (0, 0, 1e9, 1e9).
func Overlay(constraints ...Constraints) *OverlayWidget {
	var c Constraints
	if len(constraints) > 0 {
		c = constraints[0]
	} else {
		c = NewFlexConstraints(0, 0, 1e9, 1e9)
	}
	return &OverlayWidget{
		children:    make([]Widget, 0),
		constraints: c,
	}
}

// Child adds a child widget to be rendered on top of previous children and returns the overlay for chaining
func (o *OverlayWidget) Child(child Widget) *OverlayWidget {
	o.children = append(o.children, child)
	return o
}

// GetConstraints returns the overlay's constraints
func (o *OverlayWidget) GetConstraints() Constraints {
	return o.constraints
}

// Render implements the Widget interface for OverlayWidget
func (o *OverlayWidget) Render(ctx *Context, box *Box) (usedSize Size, err error) {
	var maxUsedSize Size

	// Render all children in sequence (later children paint over earlier ones)
	for _, child := range o.children {
		// Get child constraints to determine positioning and sizing
		childConstraints := child.GetConstraints()

		// Create child box based on its constraints
		childBox := &Box{
			Position: Point{
				X: box.Position.X + childConstraints.Left,
				Y: box.Position.Y + childConstraints.Top,
			},
			Size: Size{
				Width:  box.Size.Width - childConstraints.Left,
				Height: box.Size.Height - childConstraints.Top,
			},
			Constraints: childConstraints,
		}

		// For rigid widgets (min == max), use the exact constraint size
		// For flexible widgets, clamp to available space within constraints
		if childConstraints.MinWidth == childConstraints.MaxWidth {
			// Rigid width
			childBox.Size.Width = childConstraints.MinWidth
		} else {
			// Flexible width - clamp to constraints
			if childConstraints.MaxWidth < childBox.Size.Width {
				childBox.Size.Width = childConstraints.MaxWidth
			}
			if childConstraints.MinWidth > childBox.Size.Width {
				childBox.Size.Width = childConstraints.MinWidth
			}
		}

		if childConstraints.MinHeight == childConstraints.MaxHeight {
			// Rigid height
			childBox.Size.Height = childConstraints.MinHeight
		} else {
			// Flexible height - clamp to constraints
			if childConstraints.MaxHeight < childBox.Size.Height {
				childBox.Size.Height = childConstraints.MaxHeight
			}
			if childConstraints.MinHeight > childBox.Size.Height {
				childBox.Size.Height = childConstraints.MinHeight
			}
		}

		// Create child context
		childCtx := &Context{
			WindowWidth:   ctx.WindowWidth,
			WindowHeight:  ctx.WindowHeight,
			ParentBox:     childBox,
			AvailableSize: childBox.Size,
		}

		childUsedSize, err := child.Render(childCtx, childBox)
		if chk.E(err) {
			return Size{}, err
		}

		// Track the maximum used size
		if childUsedSize.Width > maxUsedSize.Width {
			maxUsedSize.Width = childUsedSize.Width
		}
		if childUsedSize.Height > maxUsedSize.Height {
			maxUsedSize.Height = childUsedSize.Height
		}
	}

	return maxUsedSize, nil
}

// Gravity specifies how a widget should be positioned within its container
type Gravity int

const (
	GravityCenter Gravity = iota
	GravityNorth
	GravitySouth
	GravityEast
	GravityWest
	GravityNorthEast
	GravityNorthWest
	GravitySouthEast
	GravitySouthWest
)

// DirectionWidget positions a single child widget using gravity-based positioning
type DirectionWidget struct {
	child       Widget
	gravity     Gravity
	constraints Constraints
}

// NewDirectionWidget creates a new direction widget with the specified gravity.
// If no constraints are provided, uses default flexible constraints (0, 0, 1e9, 1e9).
func NewDirectionWidget(child Widget, gravity Gravity, constraints ...Constraints) *DirectionWidget {
	var c Constraints
	if len(constraints) > 0 {
		c = constraints[0]
	} else {
		c = NewFlexConstraints(0, 0, 1e9, 1e9)
	}
	return &DirectionWidget{
		child:       child,
		gravity:     gravity,
		constraints: c,
	}
}

// Center creates a direction widget that centers its child.
// Equivalent to NewDirectionWidget(child, GravityCenter).
func Center(child Widget, constraints ...Constraints) *DirectionWidget {
	return NewDirectionWidget(child, GravityCenter, constraints...)
}

// GetConstraints returns the direction widget's constraints
func (d *DirectionWidget) GetConstraints() Constraints {
	return d.constraints
}

// Render implements the Widget interface for DirectionWidget
func (d *DirectionWidget) Render(ctx *Context, box *Box) (usedSize Size, err error) {
	if d.child == nil {
		return box.Size, nil
	}

	// Get child constraints
	childConstraints := d.child.GetConstraints()

	// Calculate child size (respecting rigid constraints)
	var childWidth, childHeight float32
	if childConstraints.MinWidth == childConstraints.MaxWidth {
		childWidth = childConstraints.MinWidth
	} else {
		childWidth = box.Size.Width
		if childWidth > childConstraints.MaxWidth {
			childWidth = childConstraints.MaxWidth
		}
		if childWidth < childConstraints.MinWidth {
			childWidth = childConstraints.MinWidth
		}
	}

	if childConstraints.MinHeight == childConstraints.MaxHeight {
		childHeight = childConstraints.MinHeight
	} else {
		childHeight = box.Size.Height
		if childHeight > childConstraints.MaxHeight {
			childHeight = childConstraints.MaxHeight
		}
		if childHeight < childConstraints.MinHeight {
			childHeight = childConstraints.MinHeight
		}
	}

	// Calculate position based on gravity
	var childX, childY float32
	switch d.gravity {
	case GravityCenter:
		childX = box.Position.X + (box.Size.Width-childWidth)/2
		childY = box.Position.Y + (box.Size.Height-childHeight)/2
	case GravityNorth:
		childX = box.Position.X + (box.Size.Width-childWidth)/2
		childY = box.Position.Y
	case GravitySouth:
		childX = box.Position.X + (box.Size.Width-childWidth)/2
		childY = box.Position.Y + box.Size.Height - childHeight
	case GravityEast:
		childX = box.Position.X + box.Size.Width - childWidth
		childY = box.Position.Y + (box.Size.Height-childHeight)/2
	case GravityWest:
		childX = box.Position.X
		childY = box.Position.Y + (box.Size.Height-childHeight)/2
	case GravityNorthEast:
		childX = box.Position.X + box.Size.Width - childWidth
		childY = box.Position.Y
	case GravityNorthWest:
		childX = box.Position.X
		childY = box.Position.Y
	case GravitySouthEast:
		childX = box.Position.X + box.Size.Width - childWidth
		childY = box.Position.Y + box.Size.Height - childHeight
	case GravitySouthWest:
		childX = box.Position.X
		childY = box.Position.Y + box.Size.Height - childHeight
	}

	// Create child box
	childBox := &Box{
		Position: Point{
			X: childX,
			Y: childY,
		},
		Size: Size{
			Width:  childWidth,
			Height: childHeight,
		},
		Constraints: childConstraints,
	}

	// Create child context
	childCtx := &Context{
		WindowWidth:   ctx.WindowWidth,
		WindowHeight:  ctx.WindowHeight,
		ParentBox:     childBox,
		AvailableSize: childBox.Size,
	}

	// Render child
	return d.child.Render(childCtx, childBox)
}
