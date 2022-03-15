# Bhojpur Render - Scene Graph

The `litescene` is a scene graph library for WebGL with a component based hierarchical node system.
It comes with a realistic rendering pipeline and some interesting components to make it easier to
build and share scenes.

- Component based node system
- Realistic rendering pipeline, it supports shadows, reflections, textures for all properties, etc
- Material system that automatically computes the best shader, making it easy to control properties
- Resources Manager to load and store any kind of resource (textures, meshes, etc)
- Serializing methods to convert any Scene to JSON
- Parser for most common file formats
- Easy to embed

It uses its own low-level library called `litegl.js`.

## Simple Usage

Include the library and dependencies

```html
<script src="dist/gl-matrix-min.js"></script>
<script src="dist/litegl.min.js"></script>
<script src="js/litescene.js"></script>
```

Create the Context

```js
var player = new LS.Player({
	width:800, height:600,
	resources: "resources/",
	shaders: "data/shaders.xml"
});
```

Attach to Canvas to the DOM:

```js
document.getElementById("mycontainer").appendChild( player.canvas )
```

or, you can pass the canvas in the player settings as { canvas: my_canvas_element }

Load the scene and play it:

```js
player.loadScene("scene.json");
```
