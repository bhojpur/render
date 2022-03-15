# Bhojpur Render - Canvas2D to WebGL

The `Canvas2DtoWebGL.js` ports almost all the methods from the regular Canvas2D context
(CanvasRenderingContext2D) of HTML5 to WebGL calls. It allows to mix 3D in your 2D Canvas or
the opposite (to create GUIs, debug info, etc), and *in some cases* it could even improve the
performance of your regular Canvas2D (when working with big images). It uses `litegl.js` as
the base WebGL library. To improve performance it doesn't generate garbage (reuses the same
containers). It can work with non power of two textures (no mipmaps obviously).

## Fully Supported Functions

- translate, rotate, scale, transform, setTransform, save, restore
- clearRect
- strokeStyle, fillStyle, globalAlpha
- drawImage you can use images or canvas (video not tested)
- beginPath, lineTo, moveTo, closePath, stroke, rect, strokeRect, fillRect, arc
- fill (limited to convex shapes)
- createPattern with images
- bezierCurveTo and quadraticCurveTo
- fillText (it creates a texture atlas with all the characters)
- lineWidth (only one mode supported)
- imageSmoothingEnabled
- getImageData and putImageData (not fully tested)
- shadows (not blurred)
- createLinearGradient
- clip

## Not Supported Yet

- globalCompositeOperation
- concave polygon shapes
 
## Won't Be Supported

- Blurred shadows
  
It is easy to tweak, all the parameters are publicly available inside the context (matrix, colors, etc).
Textures handlers are cached inside the Image object itself, this means that reusing the same images
between different Canvas2D context would have a performance penalty.

## Extra Features

It not only provide the Canvas2D functions, it also comes with some extra ones that you can call directly
to speed up some processes (like setting colors, modifying the matrix) or even creating some FX that would
be hard in a regular Canvas (applying a shader to an image, colorizing an image).

## Simple Usage

Include the library and dependencies

```html
<script src="js/gl-matrix-min.js"></script>
<script src="js/litegl.min.js"></script>
<script src="js/Canvas2DtoWebGL.js"></script>
```

Once you have your Canvas created (and before a CanvasRenderingContext2D has been obtained), call this
function:

```js
var ctx = enableWebGLCanvas( mycanvas );
```

During your rendering you must call this two functions, it helps set the flags accordingly.

```js
ctx.start2D();

// your rendering code
// ...

ctx.finish2D();
```
