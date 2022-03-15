# Bhojpur Render - Scene Graph

The `rendeer.js` is a lightweight 3D scene graph library, meant to be used in 3D web apps and games.
It is meant to be flexible and easy to tweak. It used the library `litegl.js` as a low level layer
for WebGL. It comes with some common useful classes like:

- Scene and SceneNode
- Camera
- Renderer
- ParticleEmissor

Since it uses litegl you have all the basic ones (Mesh, Shader, and Texture).

## Usage

First include the library and dependencies

```html
<script src="js/gl-matrix-min.js"></script>
<script src="js/litegl.js"></script>
<script src="js/rendeer.js"></script>
```

### Create the Scene

```js
var scene = new RD.Scene();
```

### Create the Renderer

```js
var context = GL.create({width: window.innerWidth, height:window.innerHeight});
var renderer = new RD.Renderer(context);
```

### Attach to DOM

```js
document.body.appendChild(renderer.canvas);
```

### Get User Input

```js
gl.captureMouse();
renderer.context.onmousedown = function(e) { ... }
renderer.context.onmousemove = function(e) { ... }

gl.captureKeys();
renderer.context.onkey = function(e) { ... }
```

### Set Camera

```js
var camera = new RD.Camera();
camera.perspective( 45, gl.canvas.width / gl.canvas.height, 1, 1000 );
camera.lookAt( [100,100,100],[0,0,0],[0,1,0] );
```

### Create and Register Mesh

```js
var mesh = GL.Mesh.fromURL("data/mesh.obj");
renderer.meshes["mymesh"] = mesh;
```

### Load and Register Texture

```js
var texture = GL.Texture.fromURL("mytexture.png", { minFilter: gl.LINEAR_MIPMAP_LINEAR, magFilter: gl.LINEAR });
renderer.textures["mytexture.png"] = texture;
```

### Compile and Register Shader

```js
var shader = new GL.Shader(vs_code, fs_code);
renderer.shaders["phong"] = shader;
```

### Add a Node to the Scene

```js
var node = new RD.SceneNode();
node.color = [1,0,0,1];
node.mesh = "mymesh";
node.texture = "mytexture.png";
node.shader = "phong";
node.position = [0,0,0];
node.scale([10,10,10]);
scene.root.addChild(node);
```

### Create Main Loop

```js
requestAnimationFrame(animate);
function animate() {
	requestAnimationFrame( animate );

	last = now;
	now = getTime();
	var dt = (now - last) * 0.001;
	renderer.render(scene, camera);
	scene.update(dt);
}
```
