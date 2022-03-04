# Bhojpur Render - WebAssmebly Canvas

The 2D Graphics `canvas` is a **WebAssembly** library for efficiently drawing on a HTML5 `canvas`
element within the web browser from Go without requiring calls back to Javascript to utilise
canvas drawing functions. The library provides the following features:

- Abstracts away the initial DOM interactions to setup the `canvas`.
- Creates the shadow image frame, and graphical Context to draw on it.
- Initializes basic font cache for text using truetype font.
- Sets up and handles `requestAnimationFrame` callback from the browser.

## Concepts

It takes an alternate approach to the current common methods for using `canvas`, allowing all
drawing primitives to be done totally with `Go` code, without calling `Javascript`.

### Standard syscall way

In a standard `wasm` application for canvas, the Go code must create a function that responds to `requestAnimationFrame` callbacks and renders the frame within that call. It interacts with
the canvas drawing primitives via the `syscall/js` functions and context switches. i.e.

```go
laserCtx.Call("beginPath")
laserCtx.Call("arc", gs.laserX, gs.laserY, gs.laserSize, 0, math.Pi*2, false)
laserCtx.Call("fill")
laserCtx.Call("closePath")
```

Downsides of this approach, are messy Javascript calls which can't easily be checked at compile time,
forcing a full redraw every frame, even if nothing changed on that canvas, or changes being much slower
than the requested frame rate.

### Go native

The `canvas` library allows all drawing to be done natively using `Go` by creating an entirely separate
image buffer which is drawn to using a 2D drawing library. This shadow `Image Buffer` can be updated at
whatever rate the developer deems appropriate, which may very well be slower than the web browser's
animation rate.

This shadow `Image Buffer` is then copied over to the web browser's canvas buffer during each `requestAnimationFrame` callback, at whatever rate the web browser requests. The handling of the callback
and copy is done automatically within the library.

Secondly, it allows the option of drawing to the `Image Buffer`, outside of the `requestAnimationFrame`
callback, if required. It is still best to do drawing within the `requestAnimationFrame` callback.

The `canvas` library provides several options to control all this, and take care of the web browser/dom interactions

- User specifies the Go render/draw callback method when calling the START function. This callback passes
the graphical context to the render routine.
- Render routine can choose to return whether any drawing took place. If it returns false, then the `requestAnimationFrame` callback does nothing, just returns immediately, saving CPU cycles. It allows the
drawing to be adaptive to the rate of data changes.
- The 'start' function accepts a maxFPS parameter. The canvas library will automatically throttle the
`requestAnimationFrame` callback to only do redraws or image buffer copies to this max rate. NOTE: it MAY
be slower depending on the Render time, and the requirements of the web browser doing other work. When a
tab is hidden, the browser regularly reduces and may even stop call to the animation callback. No
timing should be done in the render/draw routings.
- You may pass 'nil' for the render function. In this case all drawing happens totally under the 
control, outside of the library. This may be more useful in future when WASM supports proper threading.
Right now however, testing shows it is slower as all work is in one thread, and you lose the scheduling
benefits of the `requestAnimationFrame` call.

Drawing therefore, is pure **Go**. i.e.

```go
func Render(gc *img.GraphicContext) bool {
    // {some movement code removed for clarity, see example for full function}
    // draws red 🔴 laser
    gc.SetFillColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
    gc.SetStrokeColor(color.RGBA{0xff, 0x00, 0x00, 0xff})

    gc.BeginPath()
    gc.ArcTo(gs.laserX, gs.laserY, gs.laserSize, gs.laserSize, 0, math.Pi*2)
    gc.FillStroke()
    gc.Close()
return true  // Yes, we drew something, copy it over to the web browser
```

If you do want to render outside the animation loop, a simple way to cause the code to draw the frame
on schedule, independent from the web-browsers callbacks, is to use `time.Tick`.

If however your image is only updated from user input or some network activity, then it would be
straightforward to fire the redraw only when required from these inputs. This can be controlled within
the Render function, by just returning FALSE at the start. Nothing is draw, nor copied (saving CPU time)
and the previous frames data remains.

Compile with `GOOS=js GOARCH=wasm go build -o main.wasm`

Includes a [Bhojpur Web](https://github.com/bhojpur/web) configuration file to support WASM, so will serve
by just running `websvr` in the `internal/wasm` directory and opening browser to http://localhost:8080
