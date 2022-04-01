///@INFO: UNCOMMON
/* pointCloud.js */

function PointCloud(o)
{
	this.enabled = true;
	this.max_points = 1024;
	this.mesh = null; //use a mesh
	this._points = [];

	this.size = 1;

	//material
	this.texture = null;
	this.global_opacity = 1;
	this.color = vec3.fromValues(1,1,1);
	this.additive_blending = false;

	this.use_node_material = false; 
	this.in_world_coordinates = false;
	this.sort_in_z = false; //slower
	this.use_perspective = false;

	this.serialize_points = true; //disable if the points shouldnt be stored (if they are generated by another component on runtime)

	if(o)
		this.configure(o);

	this._render_instance = new ONE.RenderInstance(null, this);
	this._min = vec3.create();
	this._max = vec3.create();

	//debug
	/*
	for(var i = 0; i < 100; i++)
	{
		var pos = vec3.create();
		vec3.random( pos );
		vec3.scale( pos, pos, 50 * Math.random() );
		this.addPoint( pos, [Math.random(),1,1,1], 1 + Math.random() * 2);
	}
	*/

	if(global.gl)
		this.createMesh();

}
PointCloud.icon = "mini-icon-points.png";
PointCloud["@texture"] = { widget: "texture" };
PointCloud["@color"] = { widget: "color" };
PointCloud["@num_points"] = { widget: "info" };
PointCloud["@size"] = { type: "number", step: 0.001, precision: 3 };

PointCloud.default_color = vec4.fromValues(1,1,1,1);

Object.defineProperty( PointCloud.prototype, "num_points", {
	set: function(v) {},
	get: function() { return this._points.length; },
	enumerable: true
});

Object.defineProperty( PointCloud.prototype, "mustUpdate", {
	set: function(v) { this._mustUpdate = v; },
	get: function() { return this._mustUpdate; },
	enumerable: false
});


PointCloud.prototype.addPoint = function( position, color, size, frame_id )
{
	var data = new Float32Array(3+4+2+1); //+1 extra por distance
	data.set(position,0);
	if(color)
		data.set(color,3);
	else
		data.set( PointCloud.default_color, 3 );
	if(size !== undefined)
		data[7] = size;
	else
		data[7] = 1;
	if(frame_id != undefined )
		data[8] = frame_id;
	else
		data[8] = 0;

	this._points.push( data );
	this._must_update = true;

	return this._points.length - 1;
}

PointCloud.prototype.clear = function()
{
	this._points.length = 0;
}

PointCloud.prototype.reset = PointCloud.prototype.clear;

PointCloud.prototype.setPoint = function(id, position, color, size, frame_id )
{
	var data = this._points[id];
	if(!data) return;

	if(position)
		data.set(position,0);
	if(color)
		data.set(color,3);
	if(size !== undefined )
		data[7] = size;
	if(frame_id !== undefined )
		data[8] = frame_id;
	else
		data[8] = id;

	this._must_update = true;
}

//returns a point structure as [x,y,z, r,g,b,a, size, index]
PointCloud.prototype.getPoint = function(id)
{
	return this._points[id];
}

/*
PointCloud.prototype.setPointsFromMesh = function( mesh )
{
	//TODO
}
*/

PointCloud.prototype.removePoint = function(id)
{
	this._points.splice(id,1);
	this._must_update = true;
}


PointCloud.prototype.onAddedToNode = function(node)
{
	LEvent.bind(node, "collectRenderInstances", this.onCollectInstances, this);
}

PointCloud.prototype.onRemovedFromNode = function(node)
{
	LEvent.unbind(node, "collectRenderInstances", this.onCollectInstances, this);
}

PointCloud.prototype.getResources = function(res)
{
	if(this.mesh) res[ this.emissor_mesh ] = Mesh;
	if(this.texture) res[ this.texture ] = Texture;
}

PointCloud.prototype.onResourceRenamed = function (old_name, new_name, resource)
{
	if(this.mesh == old_name)
		this.mesh = new_name;
	if(this.texture == old_name)
		this.texture = new_name;
}

PointCloud.prototype.createMesh = function ()
{
	var max_points = this.max_points|0;
	if( this._mesh_max_points == max_points)
		return;

	this._vertices = new Float32Array(max_points * 3); 
	this._colors = new Float32Array(max_points * 4);
	this._extra2 = new Float32Array(max_points * 2); //size and texture frame

	var default_size = 1;
	for(var i = 0; i < max_points; i++)
	{
		this._colors.set( PointCloud.default_color, i*4);
		this._extra2[i*2] = default_size;
		//this._extra2[i*2+1] = 0;
	}

	if(this._mesh)
		this._mesh.deleteBuffers();

	this._mesh = new GL.Mesh();
	this._mesh.addBuffers({ vertices:this._vertices, colors: this._colors, extra2: this._extra2 }, null, gl.STREAM_DRAW);
	this._mesh_max_points = max_points;
}

PointCloud.prototype.updateMesh = function ( camera )
{
	if(!this._points.length)
		return;

	var max_points = this.max_points|0;
	if( this._mesh_max_points != max_points) 
		this.createMesh();

	var points = this._points;
	if(this.sort_in_z && camera)
	{
		var center = camera.getEye(); 
		var front = camera.getFront();

		points = this._points.concat(); //copy array
		var plane = geo.createPlane(center, front); //compute camera plane
		var den = Math.sqrt(plane[0]*plane[0] + plane[1]*plane[1] + plane[2]*plane[2]); //delta
		for(var i = 0; i < points.length; ++i)
			points[i][9] = Math.abs(vec3.dot(points[i].subarray(0,3),plane) + plane[3])/den;

		points.sort(function(a,b) { return a[9] < b[9] ? 1 : (a[9] > b[9] ? -1 : 0); });
	}

	var min = this._min;
	var max = this._max;

	//update mesh
	var i = 0, f = 0;
	var vertices = this._vertices;
	var colors = this._colors;
	var extra2 = this._extra2;

	min[0] = max[0] = points[0][0];
	min[1] = max[1] = points[0][1];
	min[2] = max[2] = points[0][2];

	for(var iPoint = 0; iPoint < points.length; ++iPoint)
	{
		if( iPoint*3 >= vertices.length)
			break; //too many points

		var p = points[iPoint];
		vertices[ iPoint*3 ] = p[0];
		vertices[ iPoint*3+1 ] = p[1];
		vertices[ iPoint*3+2 ] = p[2];

		if(p[0] < min[0]) min[0] = p[0];
		if(p[1] < min[1]) min[1] = p[1];
		if(p[2] < min[2]) min[2] = p[2];
		if(p[0] > max[0]) max[0] = p[0];
		if(p[1] > max[1]) max[1] = p[1];
		if(p[2] > max[2]) max[2] = p[2];

		var c = p.subarray(3,7);
		colors[iPoint * 4] = p[3]; 
		colors[iPoint * 4+1] = p[4];
		colors[iPoint * 4+2] = p[5];
		colors[iPoint * 4+3] = p[6];

		extra2[iPoint * 2] = p[7];
		extra2[iPoint * 2 + 1] = p[8];
	}

	//upload geometry
	this._mesh.vertexBuffers["vertices"].data = vertices;
	this._mesh.vertexBuffers["vertices"].upload();

	this._mesh.vertexBuffers["colors"].data = colors;
	this._mesh.vertexBuffers["colors"].upload();

	this._mesh.vertexBuffers["extra2"].data = extra2;
	this._mesh.vertexBuffers["extra2"].upload();
}

PointCloud._identity = mat4.create();

PointCloud.prototype.onCollectInstances = function(e, instances, options)
{
	if(!this._root) return;

	if(this._points.length == 0 || !this.enabled)
		return;

	var camera = ONE.Renderer._current_camera;

	if(this._must_update)
		this.updateMesh( camera );

	if(!this._material)
		this._material = new ONE.StandardMaterial();

	var material = this._material;

	material.color.set(this.color);
	material.opacity = this.global_opacity;
	material.setTexture( ONE.Material.COLOR, this.texture, { uvs: Material.COORDS_UV_POINTCOORD } );
	material.blend_mode = this.additive_blending ? ONE.Blend.ADD : ONE.Blend.ALPHA;
	material.flags.depth_write = !this.additive_blending;
	material.translucency = 1;
	material.flags.ignore_lights = false;

	if(!this._mesh)
		return null;

	var RI = this._render_instance;
	RI.fromNode( this._root );

	if(this.in_world_coordinates && this._root.transform )
		RI.matrix.set( this._root.transform._global_matrix );
	else
		mat4.copy( RI.matrix, PointCloud._identity );

	var material = (this._root.material && this.use_node_material) ? this._root.getMaterial() : this._material;
	mat4.multiplyVec3(RI.center, RI.matrix, vec3.create());

	RI.uniforms.u_point_size = this.size;

	if( this.use_perspective )
	{
		//enable extra2
		RI.addShaderBlock( ONE.Shaders.extra2_block );
		//enable point particles 
		RI.addShaderBlock( pointparticles_block );
	}
	else
	{
		RI.removeShaderBlock( ONE.Shaders.extra2_block );
		RI.removeShaderBlock( pointparticles_block );
	}

	RI.setMaterial( material );
	RI.setMesh( this._mesh, gl.POINTS );
	var primitives = this._points.length;
	if(primitives > this._vertices.length / 3)
		primitives = this._vertices.length / 3;

	BBox.setMinMax( this._render_instance.oobb, this._min, this._max );

	RI.setRange(0, primitives );
	instances.push(RI);
}

PointCloud.prototype.serialize = function()
{
	var o = ONE.cloneObject(this);
	o.object_class = "PointCloud";

	if(this.uid) //special case, not enumerable
		o.uid = this.uid;
	if(this.serialize_points)
	{
		var points = Array(this._points.length);
		for(var i = 0; i < this._points.length; i++)
		{
			points[i] = typedArrayToArray( this._points[i] );
		}
		o.points = points;
	}
	return o;
}

PointCloud.prototype.configure = function(o)
{
	if(!o)
		return;
	if(o.uid) 
		this.uid = o.uid;

	if(o.points && o.serialize_points)
	{
		this._points = Array( o.points.length );
		for(var i = 0; i < o.points.length; i++)
			this._points[i] = new Float32Array( o.points[i] );
		o.points = null;
	}
	ONE.cloneObject( o, this );
}


ONE.registerComponent( PointCloud );