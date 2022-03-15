# Bhojpur Render - Spatial Engine

A framework for 3D `geospatial visualization` in the web browser using Javascript.

## Getting started

The first step is to add the latest distribution to your website:

```html
<script src="path/to/spatial.min.js"></script>
<link rel="stylesheet" href="path/to/spatial.css">
```

From there you will have access to the `VIZI` namespace which you can use to interact with and set
up Spatial. You'll also want to add a HTML element that you want to contain your geospatial
visualisation:

```html
<div id="spatial"></div>
```

It's worth adding some CSS to the page to size the `spatial` element correctly, in this case
filling the entire page:

```css
* { margin: 0; padding: 0; }
html, body { height: 100%; overflow: hidden;}
#spatial { height: 100%; }
```

The next step is to set up an instance of the spatial `World` component and position it in
some city.

```javascript
// Manhattan
var coords = [40.739940, -73.988801];
var world = VIZI.world('spatial').setView(coords);
```

The first argument is the ID of the HTML element that you want to use as a container for the
spatial visualisation. Then, add some controls:

```javascript
VIZI.Controls.orbit().addTo(world);
```

And, a 2D basemap using tiles from CartoDB:

```javascript
VIZI.imageTileLayer('http://{s}.basemaps.cartocdn.com/light_nolabels/{z}/{x}/{y}.png', {
  attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors, &copy; <a href="http://cartodb.com/attributions">CartoDB</a>'
}).addTo(world);
```

At this point, you can take a look at your handywork and should be able to see a 2D map focussed
on the Manhattan area. You can move around using the mouse. If you want to be a bit more adventurous
then you can add 3D buildings using Mapzen vector tiles:

```javascript
VIZI.topoJSONTileLayer('https://vector.mapzen.com/osm/buildings/{z}/{x}/{y}.topojson?api_key=vector-tiles-NT5Emiw', {
  interactive: false,
  style: function(feature) {
    var height;

    if (feature.properties.height) {
      height = feature.properties.height;
    } else {
      height = 10 + Math.random() * 10;
    }

    return {
      height: height
    };
  },
  filter: function(feature) {
    // Don't show points
    return feature.geometry.type !== 'Point';
  },
  attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors, <a href="http://whosonfirst.mapzen.com#License">Who\'s on First</a>.'
}).addTo(world);
```

Refresh the page and you will see 3D buildings appear on top of the 2D basemap.
