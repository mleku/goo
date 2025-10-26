package window

import (
	"runtime"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"lol.mleku.dev/chk"
	"lol.mleku.dev/log"
)

// Window manages the OpenGL window and application lifecycle
type Window struct {
	width            int
	height           int
	title            string
	window           *glfw.Window
	running          bool
	canvasWidth      int
	canvasHeight     int
	frameCount       int
	skipResizeFrames bool
	resizeThreshold  int
	mouseX           float64
	mouseY           float64
	cursorInWindow   bool
}

func init() {
	runtime.LockOSThread()
}

// New creates a new window with the given configuration
func New(width, height int, title string) (w *Window, err error) {
	w = &Window{
		width:            width,
		height:           height,
		title:            title,
		canvasWidth:      width,
		canvasHeight:     height,
		resizeThreshold:  8,
		skipResizeFrames: true,
	}
	return
}

// Run starts the window and runs the application main loop
func (w *Window) Run(renderFunc func(windowWidth, windowHeight int, mouseX, mouseY float64, cursorInWindow bool) error) (err error) {
	if err = glfw.Init(); chk.E(err) {
		return
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	// Don't set OpenGLProfile - use compatibility profile for immediate mode
	glfw.WindowHint(glfw.Resizable, glfw.True)

	w.window, err = glfw.CreateWindow(w.width, w.height, w.title, nil, nil)
	if chk.E(err) {
		return
	}
	defer w.window.Destroy()

	w.window.MakeContextCurrent()

	if err = gl.Init(); chk.E(err) {
		return
	}

	// Set the viewport
	gl.Viewport(0, 0, int32(w.width), int32(w.height))

	// Enable scissor test for clipping
	gl.Enable(gl.SCISSOR_TEST)

	// Initialize canvas dimensions
	w.canvasWidth, w.canvasHeight = w.window.GetFramebufferSize()

	// Set mouse cursor position callback
	w.window.SetCursorPosCallback(func(window *glfw.Window, xpos, ypos float64) {
		w.mouseX = xpos
		w.mouseY = ypos
		log.D.Ln("Cursor position:", xpos, ypos)
	})

	// Set keyboard callback
	w.window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		log.D.Ln("Key event: key=", key, "scancode=", scancode, "action=", action, "mods=", mods)
	})

	// Set mouse button callback
	w.window.SetMouseButtonCallback(func(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		log.D.Ln("Mouse button: button=", button, "action=", action, "mods=", mods)
	})

	// Set scroll callback
	w.window.SetScrollCallback(func(window *glfw.Window, xoffset, yoffset float64) {
		log.D.Ln("Scroll: xoffset=", xoffset, "yoffset=", yoffset)
	})

	// Set character input callback
	w.window.SetCharCallback(func(window *glfw.Window, char rune) {
		log.D.Ln("Character input:", string(char))
	})

	// Set cursor enter/leave callback
	w.window.SetCursorEnterCallback(func(window *glfw.Window, entered bool) {
		w.cursorInWindow = entered
		if entered {
			log.D.Ln("Cursor entered window")
		} else {
			log.D.Ln("Cursor left window")
		}
	})

	w.running = true
	for !w.window.ShouldClose() && w.running {
		// Get window size (logical size in screen coordinates)
		windowWidth, windowHeight := w.window.GetSize()

		// Get framebuffer/canvas size (actual rendering surface)
		canvasWidth, canvasHeight := w.window.GetFramebufferSize()

		// Increment frame counter
		w.frameCount++

		// Update viewport if canvas size changed
		if canvasWidth != w.canvasWidth || canvasHeight != w.canvasHeight {
			gl.Viewport(0, 0, int32(canvasWidth), int32(canvasHeight))
			w.canvasWidth = canvasWidth
			w.canvasHeight = canvasHeight
		}

		// Render with window dimensions and mouse position
		if err = renderFunc(windowWidth, windowHeight, w.mouseX, w.mouseY, w.cursorInWindow); chk.E(err) {
			return
		}

		w.window.SwapBuffers()

		glfw.PollEvents()

	}

	return
}

// Stop stops the main loop
func (w *Window) Stop() {
	w.running = false
}

// GetWindow returns the underlying GLFW window
func (w *Window) GetWindow() *glfw.Window {
	return w.window
}
