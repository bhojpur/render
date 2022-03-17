# Bhojpur Render - Wireframe Drawing

The 3D `wireframe drawing` Library for Javascript.

## Getting Started

### Using

The coordinate space has a locked origin of `0,0,0` with the min/max boundaries of `-1,-1,-1` to
`+1,+1,+1`. The `Z` coordinate extends from `-1` (nearest) to `+1` (farthest).

There are four types of shapes; `line`, `cube`, `circle`, and `dot`. 
These can be transformed with the `Scale`, `Rotate`, and `Translate` functions.
Multiple shapes can be transformed by nesting in a `Begin/End` block.

#### A simple cube

```go
var p = new Pinhole();
p.drawCube(-0.3, -0.3, -0.3, 0.3, 0.3, 0.3);
p.render(canvasElement, {bgColor:'white'});
```

#### Rotate the cube

```go
var p = pinhole.New();
p.drawCube(-0.3, -0.3, -0.3, 0.3, 0.3, 0.3)
p.rotate(Math.PI/3, Math.PI/6, 0);
p.render(canvasElement, {bgColor:'white'});
```

#### Add, rotate, and transform a circle

```go
var p = pinhole.New();
p.drawCube(-0.3, -0.3, -0.3, 0.3, 0.3, 0.3)
p.rotate(Math.PI/3, Math.PI/6, 0);

p.begin()
p.drawCircle(0, 0, 0, 0.2)
p.rotate(0, math.Pi/2, 0)
p.translate(-0.6, -0.4, 0)
p.colorize("red");
p.end();

p.render(canvasElement, {bgColor:'white'});
```
