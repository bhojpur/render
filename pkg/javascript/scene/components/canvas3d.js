///@INFO: UNCOMMON
/**
* Allows to render 2d canvas primitives, but they are rendered into a plane that can be positioned in 3D space.
* It also supports to store the texture so it can be used in another material.
* 
* The CANVAS2D mode renders busing a native Canvas2D, which has all the features but it could be slower because it has to upload the full canvas every frame.
* The WEBGL mode renders the canvas using WebGL calls, it is faster but the quality is worse and some features are not available (but you can render other textures as images)
* To fill the canvas you must have a Script in the same node, that contains a method called OnRenderCanvas
* @class Canvas3D
* @namespace ONE.Components
* @constructor
* @param {String} object to configure from
*/
function Canvas3D(o)
{
	this.enabled = true;

	this.mode = 1;
	this.width = 512;
	this.height = 512;
	this.texture_name = ":canvas3D";
	this.visible = true;
	this.input_active = true; //used for ONE.GUI methods
	this.use_node_material = false;
	this.generate_mipmaps = false;
	this.max_interactive_distance = 100; //distance beyong which the mouse is no longer projected
	this.high_precision = false; //use a texture format of more than one byte per channel
	this.opacity = 1.0;

	this._clear_buffer = true; //not public, just here in case somebody wants it
	this._skip_backside = true;
	this._texture = null;
	this._fbo = null;
	this._RI = null;
	this._standard_material = null;

	this._mouse = vec3.create();

	this._is_mouse_inside = false;

	this._local_mouse = {
		mousex: 0,
		mousey: 0,
		buttons: 0
	};

	this._local_mouse_click = {
		mousex: 0,
		mousey: 0
	}

	if(o)
		this.configure(o);
}

Canvas3D.icon = "mini-icon-brush.png";

Canvas3D.MODE_CANVAS2D = 1; //renders to a canvas2D, then uploads the texture to the GPU after every frame
Canvas3D.MODE_WEBGL = 2;	//renders to a WebGLTexture translating the Canvas2D calls to WebGL (using Canvas2DtoWebGL)
Canvas3D.MODE_IMMEDIATE = 3; //renders directly to current viewport (using Canvas2DtoWebGL)

Canvas3D["@mode"] = { type:"enum", values: { "Canvas2D":Canvas3D.MODE_CANVAS2D, "WebGL":Canvas3D.MODE_WEBGL, "Immediate": Canvas3D.MODE_IMMEDIATE } };
Canvas3D["@width"] = { type:"number", step:1, precision:0 };
Canvas3D["@height"] = { type:"number", step:1, precision:0 };
Canvas3D["@texture_name"] = { type:"string" };

Object.defineProperty( Canvas3D.prototype, "texture", {
	set: function(){
		throw("Canvas3D texture cannot be set manually");
	},
	get: function(){
		return this._texture;
	},
	enumerable: false
});

Canvas3D.prototype.onAddedToScene = function(scene)
{
	LEvent.bind(scene, ONE.EVENT.READY_TO_RENDER, this.onRender,this);
	LEvent.bind(scene, ONE.EVENT.AFTER_RENDER_INSTANCES, this.onRender,this);
}

Canvas3D.prototype.onRemovedFromScene = function(scene)
{
	LEvent.unbind(scene, ONE.EVENT.READY_TO_RENDER, this.onRender,this);
	LEvent.unbind(scene, ONE.EVENT.AFTER_RENDER_INSTANCES, this.onRender,this);
}

Canvas3D.prototype.onAddedToNode = function( node )
{
	if(!this.texture_name)
		this.texture_name = ":canvas3D";

	LEvent.bind( node, ONE.EVENT.COLLECT_RENDER_INSTANCES, this.onCollectInstances, this );
}

Canvas3D.prototype.onRemovedFromNode = function( node )
{
	LEvent.unbind( node, ONE.EVENT.COLLECT_RENDER_INSTANCES, this.onCollectInstances, this );
}

//called before rendering scene
Canvas3D.prototype.onRender = function(e)
{
	var camera = ONE.Renderer._current_camera;
	if(!this.enabled || !camera || !camera.checkLayersVisibility( this._root.layers ) )
		return;

	if(	(e == ONE.EVENT.READY_TO_RENDER && ( this.mode == Canvas3D.MODE_CANVAS2D || this.mode == Canvas3D.MODE_WEBGL)) || 
		(e == ONE.EVENT.AFTER_RENDER_INSTANCES && this.mode == Canvas3D.MODE_IMMEDIATE)
	)
	{
		this.drawCanvas();
	}
}

Canvas3D.prototype.drawCanvas = function()
{
	var w = this.width|0;
	var h = this.height|0;

	//create resources
	if( this.mode == Canvas3D.MODE_CANVAS2D )
	{
		if(!this._canvas)
			this._canvas = document.createElement("canvas");
		if(this._canvas.width != w)
			this._canvas.width = w;
		if(this._canvas.height != h)
			this._canvas.height = h;
	}

	var type = this.high_precision ? gl.HIGH_PRECISION_FORMAT : gl.UNSIGNED_BYTE;

	if(this.mode != Canvas3D.MODE_IMMEDIATE)
	{
		if(!this._texture || this._texture.width != w || this._texture.height != h || this._texture.type != type)
			this._texture = new GL.Texture(w,h,{ type: type, format: GL.RGBA, filter: GL.LINEAR, wrap: GL.CLAMP_TO_EDGE });
	}

	//project mouse into the canvas plane
	if(this.visible)
		this.projectMouse();

	//render the canvas
	if( this.mode == Canvas3D.MODE_CANVAS2D )
	{
		var ctx = this._canvas.getContext("2d");
		if(this._clear_buffer)
			ctx.clearRect(0,0,this._canvas.width,this._canvas.height); //clear
		ONE.GUI._ctx = ctx;
		this._root.processActionInComponents("onRenderCanvas",[ctx,this._canvas,this._mouse,this]);
		ONE.GUI._ctx = gl;
		this._texture.uploadImage( this._canvas );
	}
	else if ( this.mode == Canvas3D.MODE_WEBGL )
	{
		var ctx = gl;
		if(!this._fbo)
			this._fbo = new GL.FBO();
		this._fbo.setTextures([this._texture]);
		this._fbo.bind();
		gl.start2D();
		if(this._clear_buffer)
		{
			gl.clearColor(0,0,0,0);
			gl.clear(GL.COLOR_BUFFER_BIT);
		}
		ONE.GUI._ctx = gl;
		this._root.processActionInComponents("onRenderCanvas",[ctx,this._texture,this._mouse,this]);
		gl.finish2D();
		this._fbo.unbind();
	}
	else if ( this.mode == Canvas3D.MODE_IMMEDIATE )
	{
		var ctx = gl;
		gl.start2D();
		ONE.GUI._ctx = gl;

		//pass MVP matrix
		var mvp = this._mvp;
		if(!mvp)
			mvp = this._mvp = mat4.create();
		mat4.identity(mvp);
		if(this._root.transform)
			mat4.multiply(mvp,mvp, this._root.transform.getGlobalMatrixRef() );
		mat4.scale(mvp,mvp,[1/this.width,-1/this.height,1]);
		mat4.translate(mvp,mvp,[this.width*-0.5,this.height*-0.5,0]);
		var camera = ONE.Renderer._current_camera;
		mat4.multiply( mvp, camera._viewprojection_matrix, mvp );
		//mat4.multiply( mvp, mvp, camera._viewprojection_matrix );
		//mat4.identity(mvp);
		gl.WebGLCanvas.set3DMatrix( mvp );
		if(!this._canvas_info)
			this._canvas_info = { width: 0, height: 0 };
		this._canvas_info.width = this.width;
		this._canvas_info.height = this.height;
		gl.disable( gl.CULL_FACE );
		gl.enable( gl.DEPTH_TEST );
		gl.depthFunc( gl.LEQUAL );
		gl.globalAlpha = this.opacity;

		this._root.processActionInComponents("onRenderCanvas",[ctx,this._canvas_info,this._mouse,this]);

		gl.globalAlpha = 1;
		gl.finish2D();
		gl.depthFunc( gl.LESS );
		gl.WebGLCanvas.set3DMatrix(null);
		return;
	}

	//process and share the texture
	if(this._texture)
	{
		if(this.generate_mipmaps && isPowerOfTwo(w) && isPowerOfTwo(h) )
		{
			this._texture.setParameter( GL.TEXTURE_MIN_FILTER, GL.LINEAR_MIPMAP_LINEAR );
			gl.generateMipmap( gl.TEXTURE_2D );
		}
		else
			this._texture.setParameter( GL.TEXTURE_MIN_FILTER, GL.LINEAR );
		ONE.RM.registerResource( this.texture_name || ":canvas3D", this._texture );
	}

	//restore stuff
	if( this._prev_mouse )
	{
		 ONE.Input.Mouse = this._prev_mouse;
		 this._prev_mouse = null;
	}
	if( ONE.Input.current_click && this._prev_click_mouse )
	{
		ONE.Input.current_click = this._prev_click_mouse;
		this._prev_click_mouse = null;
	}
}

Canvas3D.prototype.onCollectInstances = function(e,instances)
{
	if(!this.enabled || !this.visible || !this._texture || this.mode == Canvas3D.MODE_IMMEDIATE)
		return;

	if(!this._RI)
		this._RI = new ONE.RenderInstance();
	var RI = this._RI;
	var material = null;
	if(this.use_node_material)
		material = this._root.getMaterial();
	if(!material)
		material = this._standard_material;
	if(!material)
		material = this._standard_material = new ONE.MaterialClasses.StandardMaterial({ flags: { ignore_lights: true, cast_shadows: false }, blend_mode: ONE.Blend.ALPHA });

	if(!this.use_node_material)
	{
		material.opacity = this.opacity;
		material.blend_mode = material.opacity < 1 ? ONE.Blend.ALPHA : ONE.Blend.NORMAL;
	}

	material.setTexture("color", this.texture_name || ":canvas3D" );
	var sampler = material.textures["color"];

	RI.fromNode( this._root );
	RI.setMaterial( material );

	if(!this._mesh)
		this._mesh = GL.Mesh.plane();
	RI.setMesh(this._mesh);
	instances.push( RI );

	return instances;
}


Canvas3D.prototype.clear = function( redraw )
{
	if( this.mode == Canvas3D.MODE_CANVAS2D )
	{
		var ctx = this._canvas.getContext("2d");
		ctx.clearRect(0,0,this._canvas.width,this._canvas.height); //clear
	}
	else if( this.mode == Canvas3D.MODE_WEBGL )
	{
		if(this._texture)
			this._texture.fill([0,0,0,0]);
	}
	if(redraw)
		this.onRender();
}

Canvas3D.prototype.projectMouse = function()
{
	var camera = ONE.Renderer._main_camera;
	if(!camera)
		return;

	//Canvas Plane
	if(!this.root.transform)
	{
		this._mouse[0] = ONE.Input.Mouse.canvasx;
		this._mouse[1] = ONE.Input.Mouse.canvasy;
		this._mouse[2] = 0;
		this._is_mouse_inside = true;
		return;
	}

	var cam_dist = 0;
	cam_dist = vec3.distance( camera.getEye(), this.root.transform.getGlobalPosition() );
	var too_far = cam_dist > this.max_interactive_distance;

	this._is_mouse_inside = false;

	var x = ONE.Input.Mouse.canvasx;
	var y = ONE.Input.Mouse.canvasy;
	var w = this.width|0;
	var h = this.height|0;

	if( !this.input_active || too_far )
	{
		x = -1;
		y = -1;
		this._mouse[0] = x;
		this._mouse[1] = y;
		this._mouse[2] = -1;
	}
	else
	{
		this._mouse[0] = -1;
		this._mouse[1] = -1;
		this._mouse[2] = cam_dist;

		var ray = camera.getRay( x, y );
		if(!ray) //??
			return;

		var temp = vec3.create();
		var plane_normal = this.root.transform.localVectorToGlobal( ONE.FRONT, temp );

		if( !this._skip_backside || vec3.dot( ray.direction, plane_normal ) > 0.0 )
		{
			var local_origin = this.root.transform.globalToLocal( ray.origin, temp );
			var local_direction = this.root.transform.globalVectorToLocal( ray.direction );

			if( geo.testRayPlane( local_origin, local_direction, ONE.ZEROS, ONE.FRONT, this._mouse ) )
			{
				this._mouse[0] = (this._mouse[0] + 0.5) * w;
				this._mouse[1] = (this._mouse[1] + 0.5) * h;
				//flip Y
				this._mouse[1] = h - this._mouse[1];
			}
		}
	}

	//mark the mouse as inside
	if( this._mouse[0] >= 0 && this._mouse[0] < w &&
		this._mouse[1] >= 0 && this._mouse[1] < h )
		this._is_mouse_inside = true;

	//hacks to work with the ONE.GUI...
	//*
	this._local_mouse.mousex = this._local_mouse.x = this._mouse[0];
	this._local_mouse.mousey = this._local_mouse.y = this._mouse[1];
	this._prev_mouse = ONE.Input.Mouse;
	ONE.Input.Mouse = this._local_mouse;

	if( ONE.Input.current_click )
	{
		this._local_mouse_click.mousex = this._local_mouse.x = this._mouse[0];
		this._local_mouse_click.mousey = this._local_mouse.y = this._mouse[1];
		this._prev_click_mouse = ONE.Input.current_click;
		ONE.Input.current_click = this._local_mouse_click;
	}
	//*/
}

/*
Canvas3D.prototype.getResources = function(res)
{
	if( this.material && this.material.constructor === String )
		res[this.material] = ONE.Material;
	return res;
}

Canvas3D.prototype.onResourceRenamed = function (old_name, new_name, resource)
{
	if(this.material == old_name)
		this.material = new_name;
}
*/

ONE.registerComponent( Canvas3D );