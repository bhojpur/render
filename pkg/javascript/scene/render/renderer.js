///@INFO: BASE

//************************************
/**
* The Renderer is in charge of generating one frame of the scene. Contains all the passes and intermediate functions to create the frame.
*
* @class Renderer
* @namespace LS
* @constructor
*/

//passes
var COLOR_PASS = ONE.COLOR_PASS = { name: "color", id: 1 };
var SHADOW_PASS = ONE.SHADOW_PASS = { name: "shadow", id: 2 };
var PICKING_PASS = ONE.PICKING_PASS = { name: "picking", id: 3 };

//render events
EVENT.BEFORE_RENDER = "beforeRender";
EVENT.READY_TO_RENDER = "readyToRender";
EVENT.RENDER_SHADOWS = "renderShadows";
EVENT.AFTER_VISIBILITY = "afterVisibility";
EVENT.RENDER_REFLECTIONS = "renderReflections";
EVENT.BEFORE_RENDER_MAIN_PASS = "beforeRenderMainPass";
EVENT.ENABLE_FRAME_CONTEXT = "enableFrameContext";
EVENT.SHOW_FRAME_CONTEXT = "showFrameContext";
EVENT.AFTER_RENDER = "afterRender";
EVENT.BEFORE_RENDER_FRAME = "beforeRenderFrame";
EVENT.BEFORE_RENDER_SCENE = "beforeRenderScene";
EVENT.COMPUTE_VISIBILITY = "computeVisibility";
EVENT.AFTER_RENDER_FRAME = "afterRenderFrame";
EVENT.AFTER_RENDER_SCENE = "afterRenderScene";
EVENT.RENDER_HELPERS = "renderHelpers";
EVENT.RENDER_PICKING = "renderPicking";
EVENT.BEFORE_SHOW_FRAME_CONTEXT = "beforeShowFrameContext";
EVENT.BEFORE_CAMERA_ENABLED = "beforeCameraEnabled";
EVENT.AFTER_CAMERA_ENABLED = "afterCameraEnabled";
EVENT.BEFORE_RENDER_INSTANCES = "beforeRenderInstances";
EVENT.RENDER_INSTANCES = "renderInstances";
EVENT.RENDER_SCREEN_SPACE = "renderScreenSpace";
EVENT.AFTER_RENDER_INSTANCES = "afterRenderInstances";
EVENT.RENDER_GUI = "renderGUI";
EVENT.FILL_SCENE_UNIFORMS = "fillSceneUniforms";
EVENT.AFTER_COLLECT_DATA = "afterCollectData";
EVENT.PREPARE_MATERIALS = "prepareMaterials";

var Renderer = {

	default_render_settings: new ONE.RenderSettings(), //overwritten by the global info or the editor one
	default_material: new ONE.StandardMaterial(), //used for objects without material

	global_aspect: 1, //used when rendering to a texture that doesnt have the same aspect as the screen
	default_point_size: 1, //point size in pixels (could be overwritte by render instances)

	render_profiler: false,

	_global_viewport: vec4.create(), //the viewport we have available to render the full frame (including subviewports), usually is the 0,0,gl.canvas.width,gl.canvas.height
	_full_viewport: vec4.create(), //contains info about the full viewport available to render (current texture size or canvas size)

	//temporal info during rendering
	_current_scene: null,
	_current_render_settings: null,
	_current_camera: null,
	_current_target: null, //texture where the image is being rendered
	_current_pass: COLOR_PASS, //object containing info about the pass
	_current_layers_filter: 0xFFFF,// do a & with this to know if it must be rendered
	_global_textures: {}, //used to speed up fetching global textures
	_global_shader_blocks: [], //used to add extra shaderblocks to all objects in the scene (it gets reseted every frame)
	_global_shader_blocks_flags: 0, 
	_reverse_faces: false,
	_in_player: true, //true if rendering in the player

	_queues: [], //render queues in order

	_main_camera: null,

	_visible_cameras: null,
	_active_lights: null, //array of lights that are active in the scene
	_visible_instances: null,
	_visible_materials: [],
	_near_lights: [],
	_active_samples: [],

	//stats
	_frame_time: 0,
	_frame_cpu_time: 0,
	_rendercalls: 0, //calls to instance.render
	_rendered_instances: 0, //instances processed
	_rendered_passes: 0,
	_frame: 0,
	_last_time: 0,

	//using timer queries
	gpu_times: {
		total: 0,
		shadows: 0,
		reflections: 0,
		main: 0,
		postpo: 0,
		gui: 0
	},

	//to measure performance
	timer_queries_enabled: true,
	_timer_queries: {},
	_waiting_queries: false,

	//settings
	_collect_frequency: 1, //used to reuse info (WIP)

	//reusable locals
	_view_matrix: mat4.create(),
	_projection_matrix: mat4.create(),
	_viewprojection_matrix: mat4.create(),
	_2Dviewprojection_matrix: mat4.create(),

	_temp_matrix: mat4.create(),
	_temp_cameye: vec3.create(),
	_identity_matrix: mat4.create(),
	_uniforms: {},
	_samplers: [],
	_instancing_data: [],

	//safety
	_is_rendering_frame: false,
	_ignore_reflection_probes: false,

	//debug
	allow_textures: true,
	_sphere_mesh: null,
	_debug_instance: null,

	//fixed texture slots for global textures
	SHADOWMAP_TEXTURE_SLOT: 7,
	ENVIRONMENT_TEXTURE_SLOT: 6,
	IRRADIANCE_TEXTURE_SLOT: 5,
	LIGHTPROJECTOR_TEXTURE_SLOT: 4,
	LIGHTEXTRA_TEXTURE_SLOT: 3,

	//used in special cases
	BONES_TEXTURE_SLOT: 3,
	MORPHS_TEXTURE_SLOT: 2,
	MORPHS_TEXTURE2_SLOT: 1,

	//called from...
	init: function()
	{
		//create some useful textures: this is used in case a texture is missing
		this._black_texture = new GL.Texture(1,1, { pixel_data: [0,0,0,255] });
		this._gray_texture = new GL.Texture(1,1, { pixel_data: [128,128,128,255] });
		this._white_texture = new GL.Texture(1,1, { pixel_data: [255,255,255,255] });
		this._normal_texture = new GL.Texture(1,1, { pixel_data: [128,128,255,255] });
		this._white_cubemap_texture = new GL.Texture(1,1, { texture_type: gl.TEXTURE_CUBE_MAP, pixel_data: (new Uint8Array(6*4)).fill(255) });
		this._missing_texture = this._gray_texture;
		var internal_textures = [ this._black_texture, this._gray_texture, this._white_texture, this._normal_texture, this._missing_texture ];
		internal_textures.forEach(function(t){ t._is_internal = true; });
		ONE.ResourcesManager.textures[":black"] = this._black_texture;
		ONE.ResourcesManager.textures[":gray"] = this._gray_texture;
		ONE.ResourcesManager.textures[":white"] = this._white_texture;
		ONE.ResourcesManager.textures[":flatnormal"] = this._normal_texture;

		//some global meshes could be helpful: used for irradiance probes
		this._sphere_mesh = GL.Mesh.sphere({ size:1, detail:32 });

		//draw helps rendering debug stuff
		if(ONE.Draw)
		{
			ONE.Draw.init();
			ONE.Draw.onRequestFrame = function() { ONE.GlobalScene.requestFrame(); }
		}

		//enable webglCanvas lib so it is easy to render in 2D
		if(global.enableWebGLCanvas && !gl.canvas.canvas2DtoWebGL_enabled)
			global.enableWebGLCanvas( gl.canvas );

		// we use fixed slots to avoid changing texture slots all the time
		// from more common to less (to avoid overlappings with material textures)
		// the last slot is reserved for litegl binding stuff
		var max_texture_units = this._max_texture_units = gl.getParameter( gl.MAX_TEXTURE_IMAGE_UNITS );
		this.SHADOWMAP_TEXTURE_SLOT = max_texture_units - 2;
		this.ENVIRONMENT_TEXTURE_SLOT = max_texture_units - 3;
		this.IRRADIANCE_TEXTURE_SLOT = max_texture_units - 4;

		this.BONES_TEXTURE_SLOT = max_texture_units - 5;
		this.MORPHS_TEXTURE_SLOT = max_texture_units - 6;
		this.MORPHS_TEXTURE2_SLOT = max_texture_units - 7;

		this.LIGHTPROJECTOR_TEXTURE_SLOT = max_texture_units - 8;
		this.LIGHTEXTRA_TEXTURE_SLOT = max_texture_units - 9;

		this._active_samples.length = max_texture_units;

		this.createRenderQueues();

		this._full_viewport.set([0,0,gl.drawingBufferWidth,gl.drawingBufferHeight]);

		this._uniforms.u_viewport = gl.viewport_data;
		this._uniforms.environment_texture = this.ENVIRONMENT_TEXTURE_SLOT;
		this._uniforms.u_clipping_plane = vec4.create();
	},

	reset: function()
	{
	},

	//used to clear the state
	resetState: function()
	{
		this._is_rendering_frame = false;
		this._reverse_faces = false;
	},

	//used to store which is the current full viewport available (could be different from the canvas in case is a FBO or the camera has a partial viewport)
	setFullViewport: function(x,y,w,h)
	{
		if(arguments.length == 0) //restore
		{
			this._full_viewport[0] = this._full_viewport[1] = 0;
			this._full_viewport[2] = gl.drawingBufferWidth;
			this._full_viewport[3] = gl.drawingBufferHeight;
		}
		else if(x.constructor === Number)
		{
			this._full_viewport[0] = x; this._full_viewport[1] = y; this._full_viewport[2] = w; this._full_viewport[3] = h;
		}
		else if(x.length)
			this._full_viewport.set(x);
	},

	/**
	* Renders the current scene to the screen
	* Many steps are involved, from gathering info from the scene tree, generating shadowmaps, setup FBOs, render every camera
	* If you want to change the rendering pipeline, do not overwrite this function, try to understand it first, otherwise you will miss lots of features
	*
	* @method render
	* @param {Scene} scene
	* @param {RenderSettings} render_settings
	* @param {Array} [cameras=null] if no cameras are specified the cameras are taken from the scene
	*/
	render: function( scene, render_settings, cameras )
	{
		scene = scene || ONE.GlobalScene;

		if( this._is_rendering_frame )
		{
			console.error("Last frame didn't finish and a new one was issued. Remember that you cannot call ONE.Renderer.render from an event dispatched during the render, this would cause a recursive loop. Call ONE.Renderer.reset() to clear from an error.");
			//this._is_rendering_frame = false; //for safety, we setting to false 
			return;
		}

		//init frame
		this._is_rendering_frame = true;
		render_settings = render_settings || this.default_render_settings;
		this._current_render_settings = render_settings;
		this._current_scene = scene;
		this._main_camera = cameras ? cameras[0] : null;
		scene._frame += 1; //done at the beginning just in case it crashes
		this._frame += 1;
		scene._must_redraw = false;

		var start_time = getTime();
		this._frame_time = start_time - this._last_time;
		this._last_time = start_time;
		this._rendercalls = 0;
		this._rendered_instances = 0;
		this._rendered_passes = 0;
		this._global_shader_blocks.length = 0;
		this._global_shader_blocks_flags = 0;
		for(var i in this._global_textures)
			this._global_textures[i] = null;
		if(!this._current_pass)
			this._current_pass = COLOR_PASS;
		this._reverse_faces = false;

		//extract info about previous frame
		this.resolveQueries();

		//to restore from a possible exception (not fully tested, remove if problem)
		if(!render_settings.ignore_reset)
			ONE.RenderFrameContext.reset();

		if(gl.canvas.canvas2DtoWebGL_enabled)
			gl.resetTransform(); //reset 

		ONE.GUI.ResetImmediateGUI(true);//just to let the GUI ready

		//force fullscreen viewport
		if( !render_settings.keep_viewport )
		{
			gl.viewport(0, 0, gl.drawingBufferWidth, gl.drawingBufferHeight );
			this.setFullViewport(0, 0, gl.drawingBufferWidth, gl.drawingBufferHeight); //assign this as the full available viewport
		}
		else
			this.setFullViewport( gl.viewport_data );
		this._global_viewport.set( gl.viewport_data );

		//Event: beforeRender used in actions that could affect which info is collected for the rendering
		this.startGPUQuery( "beforeRender" );
		LEvent.trigger( scene, EVENT.BEFORE_RENDER, render_settings );
		this.endGPUQuery();

		//get render instances, cameras, lights, materials and all rendering info ready (computeVisibility)
		this.processVisibleData( scene, render_settings, cameras );

		//Define the main camera, the camera that should be the most important (used for LOD info, or shadowmaps)
		cameras = cameras && cameras.length ? cameras : scene._cameras;//this._visible_cameras;
		if(cameras.length == 0)
			throw("no cameras");
		this._visible_cameras = cameras; //the cameras being rendered
		this._main_camera = cameras[0];

		//Event: readyToRender when we have all the info to render
		LEvent.trigger( scene, EVENT.READY_TO_RENDER, render_settings );

		//remove the lights that do not lay in front of any camera (this way we avoid creating shadowmaps)
		//TODO

		//Event: renderShadowmaps helps to generate shadowMaps that need some camera info (which could be not accessible during processVisibleData)
		this.startGPUQuery("shadows");
		LEvent.trigger(scene, EVENT.RENDER_SHADOWS, render_settings );
		this.endGPUQuery();

		//Event: afterVisibility allows to cull objects according to the main camera
		LEvent.trigger(scene, EVENT.AFTER_VISIBILITY, render_settings );

		//Event: renderReflections in case some realtime reflections are needed, this is the moment to render them inside textures
		this.startGPUQuery("reflections");
		LEvent.trigger(scene, EVENT.RENDER_REFLECTIONS, render_settings );
		this.endGPUQuery();

		//Event: beforeRenderMainPass in case a last step is missing
		LEvent.trigger(scene, EVENT.BEFORE_RENDER_MAIN_PASS, render_settings );

		//enable global FX context
		if(render_settings.render_fx)
			LEvent.trigger( scene, EVENT.ENABLE_FRAME_CONTEXT, render_settings );

		//render what every camera can see
		if(this.onCustomRenderFrameCameras)
			this.onCustomRenderFrameCameras( cameras, render_settings );
		else
			this.renderFrameCameras( cameras, render_settings );

		//keep original viewport
		if( render_settings.keep_viewport )
			gl.setViewport( this._global_viewport );

		//disable and show final FX context
		if(render_settings.render_fx)
		{
			this.startGPUQuery("postpo");
			LEvent.trigger( scene, EVENT.SHOW_FRAME_CONTEXT, render_settings );
			this.endGPUQuery();
		}

		//renderGUI
		this.startGPUQuery("gui");
		this.renderGUI( render_settings );
		this.endGPUQuery();

		//profiling must go here
		this._frame_cpu_time = getTime() - start_time;
		if( ONE.Draw ) //developers may decide not to include ONE.Draw
			this._rendercalls += ONE.Draw._rendercalls; ONE.Draw._rendercalls = 0; //stats are not centralized

		//Event: afterRender to give closure to some actions
		LEvent.trigger( scene, EVENT.AFTER_RENDER, render_settings ); 
		this._is_rendering_frame = false;

		//coroutines
		ONE.triggerCoroutines("render");

		if(this.render_profiler)
			this.renderProfiler();
	},

	/**
	* Calls renderFrame of every camera in the cameras list (triggering the appropiate events)
	*
	* @method renderFrameCameras
	* @param {Array} cameras
	* @param {RenderSettings} render_settings
	*/
	renderFrameCameras: function( cameras, render_settings )
	{
		var scene = this._current_scene;

		//for each camera
		for(var i = 0; i < cameras.length; ++i)
		{
			var current_camera = cameras[i];

			LEvent.trigger(scene, EVENT.BEFORE_RENDER_FRAME, render_settings );
			LEvent.trigger(current_camera, EVENT.BEFORE_RENDER_FRAME, render_settings );
			LEvent.trigger(current_camera, EVENT.ENABLE_FRAME_CONTEXT, render_settings );

			//main render
			this.startGPUQuery("main");
			if(this.onCustomRenderFrame)
				this.onCustomRenderFrame( current_camera, render_settings ); 
			else
				this.renderFrame( current_camera, render_settings ); 
			this.endGPUQuery();

			//show buffer on the screen
			this.startGPUQuery("postpo");
			LEvent.trigger(current_camera, EVENT.SHOW_FRAME_CONTEXT, render_settings );
			LEvent.trigger(current_camera, EVENT.AFTER_RENDER_FRAME, render_settings );
			LEvent.trigger(scene, EVENT.AFTER_RENDER_FRAME, render_settings );
			this.endGPUQuery();
		}
	},

	/**
	* renders the view from one camera to the current viewport (could be the screen or a texture)
	*
	* @method renderFrame
	* @param {Camera} camera 
	* @param {Object} render_settings [optional]
	* @param {Scene} scene [optional] this can be passed when we are rendering a different scene from ONE.GlobalScene (used in renderMaterialPreview)
	*/
	renderFrame: function ( camera, render_settings, scene )
	{
		render_settings = render_settings || this.default_render_settings;

		//get all the data
		if(scene) //in case we use another scene than the default one
		{
			scene._frame++;
			this.processVisibleData( scene, render_settings );
		}
		this._current_scene = scene = scene || this._current_scene; //ugly, I know

		//set as active camera and set viewport
		this.enableCamera( camera, render_settings, render_settings.skip_viewport, scene ); 

		//clear buffer
		this.clearBuffer( camera, render_settings );

		//send before events
		LEvent.trigger(scene, EVENT.BEFORE_RENDER_SCENE, camera );
		LEvent.trigger(this, EVENT.BEFORE_RENDER_SCENE, camera );

		//in case the user wants to filter instances
		LEvent.trigger(this, EVENT.COMPUTE_VISIBILITY, this._visible_instances );

		//here we render all the instances
		if(this.onCustomRenderInstances)
			this.onCustomRenderInstances( render_settings, this._visible_instances );
		else
			this.renderInstances( render_settings, this._visible_instances );

		//send after events
		LEvent.trigger( scene, EVENT.AFTER_RENDER_SCENE, camera );
		LEvent.trigger( this, EVENT.AFTER_RENDER_SCENE, camera );
		if(this.onRenderScene)
			this.onRenderScene( camera, render_settings, scene);

		//render helpers (guizmos)
		if(render_settings.render_helpers)
		{
			if(GL.FBO.current) //rendering to multibuffer gives warnings if the shader outputs to a single fragColor
				GL.FBO.current.toSingle(); //so we disable multidraw for debug rendering (which uses a single render shader)
			LEvent.trigger(this, EVENT.RENDER_HELPERS, camera );
			LEvent.trigger(scene, EVENT.RENDER_HELPERS, camera );
			if(GL.FBO.current)
				GL.FBO.current.toMulti();
		}
	},

	//shows a RenderFrameContext to the viewport (warning, some components may do it bypassing this function)
	showRenderFrameContext: function( render_frame_context, camera )
	{
		//if( !this._current_render_settings.onPlayer)
		//	return;
		LEvent.trigger(this, EVENT.BEFORE_SHOW_FRAME_CONTEXT, render_frame_context );
		render_frame_context.show();
	},

	/**
	* Sets camera as the current camera, sets the viewport according to camera info, updates matrices, and prepares ONE.Draw
	*
	* @method enableCamera
	* @param {Camera} camera
	* @param {RenderSettings} render_settings
	*/
	enableCamera: function(camera, render_settings, skip_viewport, scene )
	{
		scene = scene || this._current_scene || ONE.GlobalScene;

		LEvent.trigger( camera, EVENT.BEFORE_CAMERA_ENABLED, render_settings );
		LEvent.trigger( scene, EVENT.BEFORE_CAMERA_ENABLED, camera );

		//assign viewport manually (shouldnt use camera.getLocalViewport to unify?)
		var startx = this._full_viewport[0];
		var starty = this._full_viewport[1];
		var width = this._full_viewport[2];
		var height = this._full_viewport[3];
		if(width == 0 && height == 0)
		{
			console.warn("enableCamera: full viewport was 0, assigning to full viewport");
			width = gl.viewport_data[2];
			height = gl.viewport_data[3];
		}

		var final_x = Math.floor(width * camera._viewport[0] + startx);
		var final_y = Math.floor(height * camera._viewport[1] + starty);
		var final_width = Math.ceil(width * camera._viewport[2]);
		var final_height = Math.ceil(height * camera._viewport[3]);

		if(!skip_viewport)
		{
			//force fullscreen viewport?
			if(render_settings && render_settings.ignore_viewports )
			{
				camera.final_aspect = this.global_aspect * camera._aspect * (width / height);
				gl.viewport( this._full_viewport[0], this._full_viewport[1], this._full_viewport[2], this._full_viewport[3] );
			}
			else
			{
				camera.final_aspect = this.global_aspect * camera._aspect * (final_width / final_height); //what if we want to change the aspect?
				gl.viewport( final_x, final_y, final_width, final_height );
			}
		}
		camera._last_viewport_in_pixels.set( gl.viewport_data );

		//recompute the matrices (view,proj and viewproj)
		camera.updateMatrices();

		//store matrices locally
		mat4.copy( this._view_matrix, camera._view_matrix );
		mat4.copy( this._projection_matrix, camera._projection_matrix );
		mat4.copy( this._viewprojection_matrix, camera._viewprojection_matrix );

		//safety in case something went wrong in the camera
		for(var i = 0; i < 16; ++i)
			if( isNaN( this._viewprojection_matrix[i] ) )
				console.warn("warning: viewprojection matrix contains NaN when enableCamera is used");

		//2D Camera: TODO: MOVE THIS SOMEWHERE ELSE
		mat4.ortho( this._2Dviewprojection_matrix, -1, 1, -1, 1, 1, -1 );

		//set as the current camera
		this._current_camera = camera;
		ONE.Camera.current = camera;
		this._current_layers_filter = render_settings ? camera.layers & render_settings.layers : camera.layers;

		//Draw allows to render debug info easily
		if(ONE.Draw)
		{
			ONE.Draw.reset(); //clear 
			ONE.Draw.setCamera( camera );
		}

		LEvent.trigger( camera, EVENT.AFTER_CAMERA_ENABLED, render_settings );
		LEvent.trigger( scene, EVENT.AFTER_CAMERA_ENABLED, camera ); //used to change stuff according to the current camera (reflection textures)
	},

	/**
	* Returns the camera active
	*
	* @method getCurrentCamera
	* @return {Camera} camera
	*/
	getCurrentCamera: function()
	{
		return this._current_camera;
	},

	/**
	* clear color using camera info ( background color, viewport scissors, clear depth, etc )
	*
	* @method clearBuffer
	* @param {Camera} camera
	* @param {ONE.RenderSettings} render_settings
	*/
	clearBuffer: function( camera, render_settings )
	{
		if( render_settings.ignore_clear || (!camera.clear_color && !camera.clear_depth) )
			return;

		//scissors test for the gl.clear, otherwise the clear affects the full viewport
		gl.scissor( gl.viewport_data[0], gl.viewport_data[1], gl.viewport_data[2], gl.viewport_data[3] );
		gl.enable(gl.SCISSOR_TEST);

		//clear color buffer 
		gl.colorMask( true, true, true, true );
		gl.clearColor( camera.background_color[0], camera.background_color[1], camera.background_color[2], camera.background_color[3] );

		//clear depth buffer
		gl.depthMask( true );

		//to clear the stencil
		gl.enable( gl.STENCIL_TEST );
		gl.clearStencil( 0x0 );

		//do the clearing
		if(GL.FBO.current)
			GL.FBO.current.toSingle();
		gl.clear( ( camera.clear_color ? gl.COLOR_BUFFER_BIT : 0) | (camera.clear_depth ? gl.DEPTH_BUFFER_BIT : 0) | gl.STENCIL_BUFFER_BIT );
		if(GL.FBO.current)
			GL.FBO.current.toMulti();

		//in case of multibuffer we want to clear with black the secondary buffers with black
		if( GL.FBO.current )
			GL.FBO.current.clearSecondary( ONE.ZEROS4 );
		/*
		if( fbo && fbo.color_textures.length > 1 && gl.extensions.WEBGL_draw_buffers )
		{
			var ext = gl.extensions.WEBGL_draw_buffers;
			var new_order = [gl.NONE];
			for(var i = 1; i < fbo.order.length; ++i)
				new_order.push(fbo.order[i]);
			ext.drawBuffersWEBGL( new_order );
			gl.clearColor( 0,0,0,0 );
			gl.clear( gl.COLOR_BUFFER_BIT );
			GL.FBO.current.toMulti();
		}
		*/

		gl.disable( gl.SCISSOR_TEST );
		gl.disable( gl.STENCIL_TEST );
	},

	//creates the separate render queues for every block of instances
	createRenderQueues: function()
	{
		this._queues.length = 0;

		this._renderqueue_background = this.addRenderQueue( new ONE.RenderQueue( ONE.RenderQueue.BACKGROUND, ONE.RenderQueue.NO_SORT, { name: "BACKGROUND" } ));
		this._renderqueue_geometry = this.addRenderQueue( new ONE.RenderQueue( ONE.RenderQueue.GEOMETRY, ONE.RenderQueue.SORT_NEAR_TO_FAR, { name: "GEOMETRY" } ));
		this._renderqueue_transparent = this.addRenderQueue( new ONE.RenderQueue( ONE.RenderQueue.TRANSPARENT, ONE.RenderQueue.SORT_FAR_TO_NEAR, { name: "TRANSPARENT" } ));
		this._renderqueue_readback = this.addRenderQueue( new ONE.RenderQueue( ONE.RenderQueue.READBACK_COLOR, ONE.RenderQueue.SORT_FAR_TO_NEAR , { must_clone_buffers: true, name: "READBACK" }));
		this._renderqueue_overlay = this.addRenderQueue( new ONE.RenderQueue( ONE.RenderQueue.OVERLAY, ONE.RenderQueue.SORT_BY_PRIORITY, { name: "OVERLAY" }));
	},

	addRenderQueue: function( queue )
	{
		var index = Math.floor(queue.value * 0.1);
		if( this._queues[ index ] )
			console.warn("Overwritting render queue:", queue.name );
		this._queues[ index ] = queue;
		return queue;
	},

	//clears render queues and inserts objects according to their settings
	updateRenderQueues: function( camera, instances )
	{
		//compute distance to camera
		var camera_eye = camera.getEye( this._temp_cameye );
		for(var i = 0, l = instances.length; i < l; ++i)
		{
			var instance = instances[i];
			if(instance)
				instance._dist = vec3.dist( instance.center, camera_eye );
		}

		var queues = this._queues;

		//clear render queues
		for(var i = 0; i < queues.length; ++i)
			if(queues[i])
				queues[i].clear();

		//add to their queues
		for(var i = 0, l = instances.length; i < l; ++i)
		{
			var instance = instances[i];
			if( !instance || !instance.material || !instance._is_visible )
				continue;
			this.addInstanceToQueue( instance );
		}

		//sort queues
		for(var i = 0, l = queues.length; i < l; ++i)
		{
			var queue = queues[i];
			if(!queue || !queue.sort_mode || !queue.instances.length)
				continue;
			queue.sort();
		}
	},

	addInstanceToQueue: function(instance)
	{
		var queues = this._queues;
		var queue = null;
		var queue_index = -1;

		if( instance.material.queue == RenderQueue.AUTO || instance.material.queue == null ) 
		{
			if( instance.material._render_state.blend )
				queue = this._renderqueue_transparent;
			else
				queue = this._renderqueue_geometry;
		}
		else
		{
			//queue index use the tens digit
			queue_index = Math.floor( instance.material.queue * 0.1 );
			queue = queues[ queue_index ];
		}

		if( !queue ) //create new queue
		{
			queue = new ONE.RenderQueue( queue_index * 10 + 5, ONE.RenderQueue.NO_SORT );
			queues[ queue_index ] = queue;
		}

		if(queue)
			queue.add( instance );
		return queue;
	},

	/**
	* To set gl state to a known and constant state in every render pass
	*
	* @method resetGLState
	* @param {RenderSettings} render_settings
	*/
	resetGLState: function( render_settings )
	{
		render_settings = render_settings || this._current_render_settings;

		//maybe we should use this function instead
		//ONE.RenderState.reset(); 

		gl.enable( gl.CULL_FACE );
		gl.frontFace(gl.CCW);

		gl.colorMask(true,true,true,true);

		gl.enable( gl.DEPTH_TEST );
		gl.depthFunc( gl.LESS );
		gl.depthMask(true);

		gl.disable( gl.BLEND );
		gl.blendFunc( gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA );

		gl.disable( gl.STENCIL_TEST );
		gl.stencilMask( 0xFF );
		gl.stencilOp( gl.KEEP, gl.KEEP, gl.KEEP );
		gl.stencilFunc( gl.ALWAYS, 1, 0xFF );
	},

	/**
	* Calls the render method for every RenderInstance (it also takes into account events and frustrum culling)
	*
	* @method renderInstances
	* @param {RenderSettings} render_settings
	* @param {Array} instances array of RIs, if not specified the last visible_instances are rendered
	*/
	renderInstances: function( render_settings, instances, scene )
	{
		scene = scene || this._current_scene;
		if(!scene)
		{
			console.warn("ONE.Renderer.renderInstances: no scene found in ONE.Renderer._current_scene");
			return 0;
		}

		this._rendered_passes += 1;

		var pass = this._current_pass;
		var camera = this._current_camera;
		var camera_index_flag = camera._rendering_index != -1 ? (1<<(camera._rendering_index)) : 0;
		var apply_frustum_culling = render_settings.frustum_culling;
		var frustum_planes = camera.updateFrustumPlanes();
		var layers_filter = this._current_layers_filter = camera.layers & render_settings.layers;

		//globals
		this._uniforms.u_viewport = gl.viewport_data;

		LEvent.trigger( scene, EVENT.BEFORE_RENDER_INSTANCES, render_settings );
		//scene.triggerInNodes( EVENT.BEFORE_RENDER_INSTANCES, render_settings );

		//reset state of everything!
		this.resetGLState( render_settings );

		LEvent.trigger( scene, EVENT.RENDER_INSTANCES, render_settings );
		LEvent.trigger( this, EVENT.RENDER_INSTANCES, render_settings );

		//reset again!
		this.resetGLState( render_settings );

		/*
		var render_instance_func = pass.render_instance;
		if(!render_instance_func)
			return 0;
		*/

		var render_instances = instances || this._visible_instances;

		//global samplers
		this.bindSamplers( this._samplers );

		var instancing_data = this._instancing_data;

		//compute visibility pass: checks which RIs are visible from this camera according to its flags, layers and AABB
		for(var i = 0, l = render_instances.length; i < l; ++i)
		{
			//render instance
			var instance = render_instances[i];
			var node_flags = instance.node.flags;
			instance._is_visible = false;

			//hidden nodes
			if( pass == SHADOW_PASS && !(instance.material.flags.cast_shadows) )
				continue;
			if( pass == PICKING_PASS && node_flags.selectable === false )
				continue;
			if( (layers_filter & instance.layers) === 0 )
				continue;

			//done here because sometimes some nodes are moved in this action
			if(instance.onPreRender)
				if( instance.onPreRender( render_settings ) === false)
					continue;

			if(!instance.material) //in case something went wrong...
				continue;

			var material = camera._overwrite_material || instance.material;

			if(material.opacity <= 0) //TODO: remove this, do it somewhere else
				continue;

			//test visibility against camera frustum
			if( apply_frustum_culling && instance.use_bounding && !material.flags.ignore_frustum )
			{
				if(geo.frustumTestBox( frustum_planes, instance.aabb ) == CLIP_OUTSIDE )
					continue;
			}

			//save visibility info
			instance._is_visible = true;
			if(camera_index_flag) //shadowmap cameras dont have an index
				instance._camera_visibility |= camera_index_flag;
		}

		//separate in render queues, and sort them according to distance or priority
		this.updateRenderQueues( camera, render_instances, render_settings );

		var start = this._rendered_instances;
		var debug_instance = this._debug_instance;

		//process render queues
		for(var j = 0; j < this._queues.length; ++j)
		{
			var queue = this._queues[j];
			if(!queue || !queue.instances.length || !queue.enabled) //empty
				continue;

			//used to change RenderFrameContext stuff (cloning textures for refraction, etc)
			if(queue.start( pass, render_settings ) == false)
				continue;

			var render_instances = queue.instances;

			//for each render instance
			for(var i = 0, l = render_instances.length; i < l; ++i)
			{
				//render instance
				var instance = render_instances[i];

				//used to debug
				if(instance == debug_instance)
				{
					console.log(debug_instance);
					debugger; 
				}

				if( !instance._is_visible || !instance.mesh )
					continue;

				this._rendered_instances += 1;

				var material = camera._overwrite_material || instance.material;

				if( pass == PICKING_PASS && material.renderPickingInstance )
					material.renderPickingInstance( instance, render_settings, pass );
				else if( material.renderInstance )
					material.renderInstance( instance, render_settings, pass );
				else
					continue;

				//some instances do a post render action (DEPRECATED)
				if(instance.onPostRender)
					instance.onPostRender( render_settings );
			}

			queue.finish( pass, render_settings );
		}

		this.resetGLState( render_settings );

		LEvent.trigger( scene, EVENT.RENDER_SCREEN_SPACE, render_settings);

		//restore state
		this.resetGLState( render_settings );

		LEvent.trigger( scene, EVENT.AFTER_RENDER_INSTANCES, render_settings );
		LEvent.trigger( this, EVENT.AFTER_RENDER_INSTANCES, render_settings );

		//and finally again
		this.resetGLState( render_settings );

		return this._rendered_instances - start;
	},

	/*
	groupingInstances: function(instances)
	{
		//TODO: if material supports instancing WIP
		var instancing_supported = gl.webgl_version > 1 || gl.extensions["ANGLE_instanced_arrays"];
		if( instancing_supported && material._allows_instancing && !instance._shader_blocks.length )
		{
			var instancing_ri_info = null;
			if(!instancing_data[ material._index ] )
				instancing_data[ material._index ] = instancing_ri_info = [];
			instancing_ri_info.push( instance );
		}
	},
	*/

	renderGUI: function( render_settings )
	{
		//renders GUI items using mostly the Canvas2DtoWebGL library
		gl.viewport( this._full_viewport[0], this._full_viewport[1], this._full_viewport[2], this._full_viewport[3] ); //assign full viewport always?
		if(gl.start2D) //in case we have Canvas2DtoWebGL installed (it is optional)
			gl.start2D();
		if( render_settings.render_gui )
		{
			if( LEvent.hasBind( this._current_scene, EVENT.RENDER_GUI ) ) //to avoid forcing a redraw if no gui is set
			{
				if(ONE.GUI)
					ONE.GUI.ResetImmediateGUI(); //mostly to change the cursor (warning, true to avoid forcing redraw)
				LEvent.trigger( this._current_scene, EVENT.RENDER_GUI, gl );
			}
		}
		if( this.on_render_gui ) //used by the editor (here to ignore render_gui flag)
			this.on_render_gui( render_settings );
		if( gl.finish2D )
			gl.finish2D();
	},

	/**
	* returns a list of all the lights overlapping this instance (it uses sperical bounding so it could returns lights that are not really overlapping)
	* It is used by the multipass lighting to iterate 
	*
	* @method getNearLights
	* @param {RenderInstance} instance the render instance
	* @param {Array} result [optional] the output array
	* @return {Array} array containing a list of ONE.Light affecting this RenderInstance
	*/
	getNearLights: function( instance, result )
	{
		result = result || [];

		result.length = 0; //clear old lights

		//it uses the lights gathered by prepareVisibleData
		var lights = this._active_lights;
		if(!lights || !lights.length)
			return result;

		//Compute lights affecting this RI (by proximity, only takes into account spherical bounding)
		result.length = 0;
		var numLights = lights.length;
		for(var j = 0; j < numLights; j++)
		{
			var light = lights[j];
			//same layer?
			if( (light.illuminated_layers & instance.layers) == 0 || (light.illuminated_layers & this._current_camera.layers) == 0)
				continue;
			var light_intensity = light.computeLightIntensity();
			//light intensity too low?
			if(light_intensity < 0.0001)
				continue;
			var light_radius = light.computeLightRadius();
			var light_pos = light.position;
			//overlapping?
			if( light_radius == -1 || instance.overlapsSphere( light_pos, light_radius ) )
				result.push( light );
		}

		return result;
	},

	regenerateShadowmaps: function( scene, render_settings )
	{
		scene = scene || this._current_scene;
		render_settings = render_settings || this.default_render_settings;
		LEvent.trigger( scene, EVENT.RENDER_SHADOWS, render_settings );
		for(var i = 0; i < this._active_lights.length; ++i)
		{
			var light = this._active_lights[i];
			light.prepare( render_settings );
			light.onGenerateShadowmap();
		}
	},

	mergeSamplers: function( samplers, result )
	{
		result = result || [];
		result.length = this._max_texture_units;

		for(var i = 0; i < result.length; ++i)
		{
			for(var j = samplers.length - 1; j >= 0; --j)
			{
				if(	samplers[j][i] )
				{
					result[i] = samplers[j][i];
					break;
				}
			}
		}

		return result;
	},

	//to be sure we dont have anything binded
	clearSamplers: function()
	{
		for(var i = 0; i < this._max_texture_units; ++i)
		{
			gl.activeTexture(gl.TEXTURE0 + i);
			gl.bindTexture( gl.TEXTURE_2D, null );
			gl.bindTexture( gl.TEXTURE_CUBE_MAP, null );
			this._active_samples[i] = null;
		}
	},

	bindSamplers: function( samplers )
	{
		if(!samplers)
			return;

		var allow_textures = this.allow_textures; //used for debug

		for(var slot = 0; slot < samplers.length; ++slot)
		{
			var sampler = samplers[slot];
			if(!sampler) 
				continue;

			//REFACTOR THIS
			var tex = null;
			if(sampler.constructor === String || sampler.constructor === GL.Texture) //old way
			{
				tex = sampler;
				sampler = null;
			}
			else if(sampler.texture)
				tex = sampler.texture;
			else //dont know what this var type is?
			{
				//continue; //if we continue the sampler slot will remain empty which could lead to problems
			}

			if( tex && tex.constructor === String)
				tex = ONE.ResourcesManager.textures[ tex ];
			if(!allow_textures)
				tex = null;

			if(!tex)
			{
				if(sampler)
				{
					switch( sampler.missing )
					{
						case "black": tex = this._black_texture; break;
						case "white": tex = this._white_texture; break;
						case "gray": tex = this._gray_texture; break;
						case "normal": tex = this._normal_texture; break;
						case "cubemap": tex = this._white_cubemap_texture; break;
						default: 
							if(sampler.is_cubemap) //must be manually specified
								tex = this._white_cubemap_texture;
							else
								tex = this._missing_texture;
					}
				}
				else
					tex = this._missing_texture;
			}

			//avoid to read from the same texture we are rendering to (generates warnings)
			if(tex._in_current_fbo) 
				tex = this._missing_texture;

			tex.bind( slot );
			this._active_samples[slot] = tex;

			//texture properties
			if(sampler)// && sampler._must_update ) //disabled because samplers ALWAYS must set to the value, in case the same texture is used in several places in the scene
			{
				if(sampler.minFilter)
				{
					if( sampler.minFilter !== gl.LINEAR_MIPMAP_LINEAR || (GL.isPowerOfTwo( tex.width ) && GL.isPowerOfTwo( tex.height )) )
						gl.texParameteri(tex.texture_type, gl.TEXTURE_MIN_FILTER, sampler.minFilter);
				}
				if(sampler.magFilter)
					gl.texParameteri(tex.texture_type, gl.TEXTURE_MAG_FILTER, sampler.magFilter);
				if(sampler.wrap)
				{
					gl.texParameteri(tex.texture_type, gl.TEXTURE_WRAP_S, sampler.wrap);
					gl.texParameteri(tex.texture_type, gl.TEXTURE_WRAP_T, sampler.wrap);
				}
				if(sampler.anisotropic != null && gl.extensions.EXT_texture_filter_anisotropic )
					gl.texParameteri(tex.texture_type, gl.extensions.EXT_texture_filter_anisotropic.TEXTURE_MAX_ANISOTROPY_EXT, sampler.anisotropic );

				//sRGB textures must specified ON CREATION, so no
				//if(sampler.anisotropic != null && gl.extensions.EXT_sRGB )
				//sampler._must_update = false;
			}
		}
	},

	//Called at the beginning of processVisibleData 
	fillSceneUniforms: function( scene, render_settings )
	{
		//global uniforms
		var uniforms = scene._uniforms;
		uniforms.u_time = scene._time || getTime() * 0.001;
		uniforms.u_ambient_light = scene.info ? scene.info.ambient_color : vec3.create();

		this._samplers.length = 0;

		//clear globals
		this._global_textures.environment = null;

		//fetch global textures
		if(scene.info)
		for(var i in scene.info.textures)
		{
			var texture = ONE.getTexture( scene.info.textures[i] );
			if(!texture)
				continue;

			var slot = 0;
			if( i == "environment" )
				slot = ONE.Renderer.ENVIRONMENT_TEXTURE_SLOT;
			else
				continue; 

			var type = (texture.texture_type == gl.TEXTURE_2D ? "_texture" : "_cubemap");
			if(texture.texture_type == gl.TEXTURE_2D)
			{
				texture.bind(0);
				texture.setParameter( gl.TEXTURE_MIN_FILTER, gl.LINEAR ); //avoid artifact
			}
			this._samplers[ slot ] = texture;
			scene._uniforms[ i + "_texture" ] = slot; 
			scene._uniforms[ i + type ] = slot; //LEGACY

			if( i == "environment" )
				this._global_textures.environment = texture;
		}

		LEvent.trigger( scene, EVENT.FILL_SCENE_UNIFORMS, scene._uniforms );
	},	

	/**
	* Collects and process the rendering instances, cameras and lights that are visible
	* Its a prepass shared among all rendering passes
	* Called ONCE per frame from ONE.Renderer.render before iterating cameras
	* Warning: rendering order is computed here, so it is shared among all the cameras (TO DO, move somewhere else)
	*
	* @method processVisibleData
	* @param {Scene} scene
	* @param {RenderSettings} render_settings
	* @param {Array} cameras in case you dont want to use the scene cameras
	*/
	processVisibleData: function( scene, render_settings, cameras, instances, skip_collect_data )
	{
		//options = options || {};
		//options.scene = scene;
		var frame = scene._frame;
		instances = instances || scene._instances;

		this._current_scene = scene;
		//compute global scene info
		this.fillSceneUniforms( scene, render_settings );

		//update info about scene (collecting it all or reusing the one collected in the frame before)
		if(!skip_collect_data)
		{
			if( this._frame % this._collect_frequency == 0)
				scene.collectData( cameras );
			LEvent.trigger( scene, EVENT.AFTER_COLLECT_DATA, scene );
		}

		//set cameras: use the parameters ones or the ones found in the scene
		cameras = (cameras && cameras.length) ? cameras : scene._cameras;
		if( cameras.length == 0 )
		{
			console.error("no cameras found");
			return;
		}
				
		//find which materials are going to be seen
		var materials = this._visible_materials; 
		materials.length = 0;

		//prepare cameras: TODO: sort by priority
		for(var i = 0, l = cameras.length; i < l; ++i)
		{
			var camera = cameras[i];
			camera._rendering_index = i;
			camera.prepare();
			if(camera.overwrite_material)
			{
				var material = camera.overwrite_material.constructor === String ? ONE.ResourcesManager.resources[ camera.overwrite_material ] : camera.overwrite_material;
				if(material)
				{
					camera._overwrite_material = material;
					materials.push( material );
				}
			}
			else
				camera._overwrite_material = null;
		}

		//define the main camera (the camera used for some algorithms)
		if(!this._main_camera)
		{
			if( cameras.length )
				this._main_camera = cameras[0];
			else
				this._main_camera = new ONE.Camera(); // ??
		}

		//nearest reflection probe to camera
		var nearest_reflection_probe = scene.findNearestReflectionProbe( this._main_camera.getEye() );

		//process instances
		this.processRenderInstances( instances, materials, scene, render_settings );

		//store all the info
		this._visible_instances = scene._instances;
		this._active_lights = scene._lights;
		this._visible_cameras = cameras; 
		//this._visible_materials = materials;

		//prepare lights (collect data and generate shadowmaps)
		for(var i = 0, l = this._active_lights.length; i < l; ++i)
			this._active_lights[i].prepare( render_settings );

		LEvent.trigger( scene, EVENT.AFTER_COLLECT_DATA, scene );
	},

	//this processes the instances 
	processRenderInstances: function( instances, materials, scene, render_settings )
	{
		materials = materials || this._visible_materials;
		var frame = scene._frame;
		render_settings = render_settings || this._current_render_settings;

		//process render instances (add stuff if needed, gather materials)
		for(var i = 0, l = instances.length; i < l; ++i)
		{
			var instance = instances[i];
			if(!instance)
				continue;

			var node_flags = instance.node.flags;

			if(!instance.mesh)
			{
				console.warn("RenderInstance must always have mesh");
				continue;
			}

			//materials
			if(!instance.material)
				instance.material = this.default_material;

			if( instance.material._last_frame_update != frame )
			{
				instance.material._last_frame_update = frame;
				materials.push( instance.material );
			}

			//add extra info: distance to main camera (used for sorting)
			instance._dist = 0;

			//find nearest reflection probe
			if( scene._reflection_probes.length && !this._ignore_reflection_probes )
				instance._nearest_reflection_probe = scene.findNearestReflectionProbe( instance.center ); //nearest_reflection_probe;
			else
				instance._nearest_reflection_probe = null;

			//change conditionaly
			if(render_settings.force_wireframe && instance.primitive != gl.LINES ) 
			{
				instance.primitive = gl.LINES;
				if(instance.mesh)
				{
					if(!instance.mesh.indexBuffers["wireframe"])
						instance.mesh.computeWireframe();
					instance.index_buffer = instance.mesh.indexBuffers["wireframe"];
				}
			}

			//clear camera visibility mask (every flag represents a camera index)
			instance._camera_visibility = 0|0;
			instance.index = i;
		}

		//prepare materials 
		for(var i = 0; i < materials.length; ++i)
		{
			var material = materials[i];
			material._index = i;
			if( material.prepare )
				material.prepare( scene );
		}

		LEvent.trigger( scene, EVENT.PREPARE_MATERIALS );
	},

	/**
	* Renders a frame into a texture (could be a cubemap, in which case does the six passes)
	*
	* @method renderInstancesToRT
	* @param {Camera} cam
	* @param {Texture} texture
	* @param {RenderSettings} render_settings
	*/
	renderInstancesToRT: function( cam, texture, render_settings, instances )
	{
		render_settings = render_settings || this.default_render_settings;
		this._current_target = texture;
		var scene = ONE.Renderer._current_scene;
		texture._in_current_fbo = true;

		if(texture.texture_type == gl.TEXTURE_2D)
		{
			this.enableCamera(cam, render_settings);
			texture.drawTo( inner_draw_2d );
		}
		else if( texture.texture_type == gl.TEXTURE_CUBE_MAP)
			this.renderToCubemap( cam.getEye(), texture.width, texture, render_settings, cam.near, cam.far );
		this._current_target = null;
		texture._in_current_fbo = false;

		function inner_draw_2d()
		{
			ONE.Renderer.clearBuffer( cam, render_settings );
			/*
			gl.clearColor(cam.background_color[0], cam.background_color[1], cam.background_color[2], cam.background_color[3] );
			if(render_settings.ignore_clear != true)
				gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
			*/
			//render scene
			ONE.Renderer.renderInstances( render_settings, instances );
		}
	},

	/**
	* Renders the current scene to a cubemap centered in the given position
	*
	* @method renderToCubemap
	* @param {vec3} position center of the camera where to render the cubemap
	* @param {number} size texture size
	* @param {Texture} texture to reuse the same texture
	* @param {RenderSettings} render_settings
	* @param {number} near
	* @param {number} far
	* @return {Texture} the resulting texture
	*/
	renderToCubemap: function( position, size, texture, render_settings, near, far, background_color, instances )
	{
		size = size || 256;
		near = near || 1;
		far = far || 1000;

		if(render_settings && render_settings.constructor !== ONE.RenderSettings)
			throw("render_settings parameter must be ONE.RenderSettings.");

		var eye = position;
		if( !texture || texture.constructor != GL.Texture)
			texture = null;

		var scene = this._current_scene;
		if(!scene)
			scene = this._current_scene = ONE.GlobalScene;

		var camera = this._cubemap_camera;
		if(!camera)
			camera = this._cubemap_camera = new ONE.Camera();
		camera.configure({ fov: 90, aspect: 1.0, near: near, far: far });

		texture = texture || new GL.Texture(size,size,{texture_type: gl.TEXTURE_CUBE_MAP, minFilter: gl.NEAREST});
		this._current_target = texture;
		texture._in_current_fbo = true; //block binding this texture during rendering of the reflection

		texture.drawTo( function(texture, side) {

			var info = ONE.Camera.cubemap_camera_parameters[side];
			if(texture._is_shadowmap || !background_color )
				gl.clearColor(0,0,0,0);
			else
				gl.clearColor( background_color[0], background_color[1], background_color[2], background_color[3] );
			gl.clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT);
			camera.configure({ eye: eye, center: [ eye[0] + info.dir[0], eye[1] + info.dir[1], eye[2] + info.dir[2]], up: info.up });

			ONE.Renderer.enableCamera( camera, render_settings, true );
			ONE.Renderer.renderInstances( render_settings, instances, scene );
		});

		this._current_target = null;
		texture._in_current_fbo = false;
		return texture;
	},

	/**
	* Returns the last camera that falls into a given screen position
	*
	* @method getCameraAtPosition
	* @param {number} x in canvas coordinates (0,0 is bottom-left)
	* @param {number} y in canvas coordinates (0,0 is bottom-left)
	* @param {Scene} scene if not specified last rendered scene will be used
	* @return {Camera} the camera
	*/
	getCameraAtPosition: function(x,y, cameras)
	{
		cameras = cameras || this._visible_cameras;
		if(!cameras || !cameras.length)
			return null;

		for(var i = cameras.length - 1; i >= 0; --i)
		{
			var camera = cameras[i];
			if(!camera.enabled || camera.render_to_texture)
				continue;

			if( camera.isPoint2DInCameraViewport(x,y) )
				return camera;
		}
		return null;
	},

	setRenderPass: function( pass )
	{
		if(!pass)
			pass = COLOR_PASS;
		this._current_pass = pass;
	},

	addImmediateRenderInstance: function( instance )
	{
		if(!instance.material)
			return;

		//this is done in collect so...
		instance.updateAABB(); 

		//add material to the list of visible materials
		if( instance.material._last_frame_update != this._frame )
		{
			instance.material._last_frame_update = this._frame;
			this._visible_materials.push( instance.material );
			if( instance.material.prepare )
				instance.material.prepare( this._current_scene );
		}

		this.addInstanceToQueue( instance );

		this._visible_instances.push( instance );
	},
	
	/**
	* Enables a ShaderBlock ONLY DURING THIS FRAME
	* must be called during frame rendering (event like fillSceneUniforms)
	*
	* @method enableFrameShaderBlock
	* @param {String} shader_block_name
	*/
	enableFrameShaderBlock: function( shader_block_name, uniforms, samplers )
	{
		var shader_block = shader_block_name.constructor === ONE.ShaderBlock ? shader_block_name : ONE.Shaders.getShaderBlock( shader_block_name );

		if( !shader_block || this._global_shader_blocks_flags & shader_block.flag_mask )
			return; //already added

		this._global_shader_blocks.push( shader_block );
		this._global_shader_blocks_flags |= shader_block.flag_mask;

		//add uniforms to renderer uniforms?
		if(uniforms)
			for(var i in uniforms)
				this._uniforms[i] = uniforms[i];

		if(samplers)
			for(var i = 0; i < samplers.length; ++i)
				if( samplers[i] )
					this._samplers[i] = samplers[i];
	},

	/**
	* Disables a ShaderBlock ONLY DURING THIS FRAME
	* must be called during frame rendering (event like fillSceneUniforms)
	*
	* @method disableFrameShaderBlock
	* @param {String} shader_block_name
	*/
	disableFrameShaderBlock:  function( shader_block_name, uniforms, samplers )
	{
		var shader_block = shader_block_name.constructor === ONE.ShaderBlock ? shader_block_name : ONE.Shaders.getShaderBlock( shader_block_name );
		if( !shader_block || !(this._global_shader_blocks_flags & shader_block.flag_mask) )
			return; //not active

		var index = this._global_shader_blocks.indexOf( shader_block );
		if(index != -1)
			this._global_shader_blocks.splice( index, 1 );
		this._global_shader_blocks_flags &= ~( shader_block.flag_mask ); //disable bit
	},

	//time queries for profiling
	_current_query: null,

	startGPUQuery: function( name )
	{
		if(!gl.extensions["disjoint_timer_query"] || !this.timer_queries_enabled) //if not supported
			return;
		if(this._waiting_queries)
			return;
		var ext = gl.extensions["disjoint_timer_query"];
		var query = this._timer_queries[ name ];
		if(!query)
			query = this._timer_queries[ name ] = ext.createQueryEXT();
		ext.beginQueryEXT( ext.TIME_ELAPSED_EXT, query );
		this._current_query = query;
	},

	endGPUQuery: function()
	{
		if(!gl.extensions["disjoint_timer_query"] || !this.timer_queries_enabled) //if not supported
			return;
		if(this._waiting_queries)
			return;
		var ext = gl.extensions["disjoint_timer_query"];
		ext.endQueryEXT( ext.TIME_ELAPSED_EXT );
		this._current_query = null;
	},

	resolveQueries: function()
	{
		if(!gl.extensions["disjoint_timer_query"] || !this.timer_queries_enabled) //if not supported
			return;

		//var err = gl.getError();
		//if(err != gl.NO_ERROR)
		//	console.log("GL_ERROR: " + err );

		var ext = gl.extensions["disjoint_timer_query"];

		var last_query = this._timer_queries["gui"];
		if(!last_query)
			return;

		var available = ext.getQueryObjectEXT( last_query, ext.QUERY_RESULT_AVAILABLE_EXT );
		if(!available)
		{
			this._waiting_queries = true;
			return;
		}
	
		var disjoint = gl.getParameter( ext.GPU_DISJOINT_EXT );
		if(!disjoint)
		{
			var total = 0;
			for(var i in this._timer_queries)
			{
				var query = this._timer_queries[i];
				// See how much time the rendering of the object took in nanoseconds.
				var timeElapsed = ext.getQueryObjectEXT( query, ext.QUERY_RESULT_EXT ) * 10e-6; //to milliseconds;
				total += timeElapsed;
				this.gpu_times[ i ] = timeElapsed;
				//ext.deleteQueryEXT(query);
				//this._timer_queries[i] = null;
			}
			this.gpu_times.total = total;
		}

		this._waiting_queries = false;
	},

	profiler_text: [],

	renderProfiler: function()
	{
		if(!gl.canvas.canvas2DtoWebGL_enabled)
			return;

		var text = this.profiler_text;
		var ext = gl.extensions["disjoint_timer_query"];

		if(this._frame % 5 == 0)
		{
			text.length = 0;
			var fps = 1000 / this._frame_time;
			text.push( fps.toFixed(2) + " FPS" );
			text.push( "CPU: " + this._frame_cpu_time.toFixed(2) + " ms" );
			text.push( " - Passes: " + this._rendered_passes );
			text.push( " - RIs: " + this._rendered_instances );
			text.push( " - Draws: " + this._rendercalls );

			if( ext )
			{
				text.push( "GPU: " + this.gpu_times.total.toFixed(2) );
				text.push( " - PreRender: " + this.gpu_times.beforeRender.toFixed(2) );
				text.push( " - Shadows: " + this.gpu_times.shadows.toFixed(2) );
				text.push( " - Reflections: " + this.gpu_times.reflections.toFixed(2) );
				text.push( " - Scene: " + this.gpu_times.main.toFixed(2) );
				text.push( " - Postpo: " + this.gpu_times.postpo.toFixed(2) );
				text.push( " - GUI: " + this.gpu_times.gui.toFixed(2) );
			}
			else
				text.push( "GPU: ???");
		}

		var ctx = gl;
		ctx.save();
		ctx.translate( gl.canvas.width - 200, gl.canvas.height - 280 );
		ctx.globalAlpha = 0.7;
		ctx.font = "14px Tahoma";
		ctx.fillStyle = "black";
		ctx.fillRect(0,0,200,280);
		ctx.fillStyle = "white";
		ctx.fillText( "Profiler", 20, 20 );
		ctx.fillStyle = "#AFA";
		for(var i = 0; i < text.length; ++i)
			ctx.fillText( text[i], 20,50 + 20 * i );
		ctx.restore();
	},

	/**
	* Renders one texture into another texture, it allows to apply a shader
	*
	* @method blit
	* @param {GL.Texture} source
	* @param {GL.Texture} destination
	* @param {GL.Shader} shader [optional] shader to apply, it must use the GL.Shader.QUAD_VERTEX_SHADER as vertex shader
	* @param {Object} uniforms [optional] uniforms for the shader
	*/
	blit: function( source, destination, shader, uniforms )
	{
		if(!source || !destination)
			throw("data missing in blit");

		if(source != destination)
		{
			destination.drawTo( function(){
				source.toViewport( shader, uniforms );
			});
			return;
		}

		if(!shader)
			throw("blitting texture to the same texture doesnt makes sense unless a shader is specified");

		var temp = GL.Texture.getTemporary( source.width, source.height, source );
		source.copyTo( temp );
		temp.copyTo( source, shader, uniforms );
		GL.Texture.releaseTemporary( temp );
	}
};

//Add to global Scope
ONE.Renderer = Renderer;
