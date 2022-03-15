# Bhojpur Render - Pixel Library

The `litepixel.js` is 2D GPU accelerated Games library. It uses WebGL to accelerate graphic rendering.
It is built on top of `litegl.js`.

## Simple Usage

### Include Library Dependencies

```html
<script src="js/gl-matrix-min.js"></script>
<script src="js/litegl.js"></script>
<script src="js/litepixel.js"></script>
```

### Create the Stage

```js
var stage = new LitePixel.Stage();
```

### Create the Renderer

```js
var renderer = new LitePixel.Renderer(window.innerWidth, window.innerHeight);
```

### Attach to DOM

```js
document.body.appendChild(renderer.canvas);
```

### Hook Events

```js
```

### Get user Input

```js
gl.captureMouse();
renderer.context.onmousedown = function(e) { ... }
renderer.context.onmousemove = function(e) { ... }

gl.captureKeys();
renderer.context.onkey = function(e) { ... }
```

### Add Sprite

```js
player = Sprite.fromImage("astronaut.png");
player.position.set([240,300]);
player.scale.set([2, 2]);
stage.addChild(player);
```

### Create Main Loop

```js
requestAnimationFrame(animate);
function animate() {
	requestAnimationFrame( animate );

	last = now;
	now = getTime();
	var dt = (now - last) * 0.001;
	renderer.render(stage);
	stage.update(dt);
}
```
