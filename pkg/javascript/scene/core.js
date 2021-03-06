//Global Scope
//better array conversion to string for serializing
if( !Uint8Array.prototype.toJSON )
{
	var typed_arrays = [ Uint8Array, Int8Array, Uint16Array, Int16Array, Uint32Array, Int32Array, Float32Array, Float64Array ];
	typed_arrays.forEach( function(v) { 
		v.prototype.toJSON = function(){ return Array.prototype.slice.call(this); }
	});
}

if( typeof(GL) === "undefined" )
	console.error("LiteScene requires to include LiteGL first.");


/**
* LS is the global scope for the global functions and containers of LiteScene
*
///@INFO: BASE
* @class  LS
* @module LS
*/

var ONE = {

	Version: 0.6,

	//systems: defined in their own files
	ResourcesManager: null,
	Picking: null,
	Player: null,
	GUI: null,
	Network: null,
	Input: null,
	Renderer: null,
	Physics: null,
	ShadersManager: null,
	Formats: null,
	Tween: null,

	//containers
	Classes: {}, //maps classes name like "Prefab" or "Animation" to its namespace "ONE.Prefab". Used in Formats and ResourceManager when reading classnames from JSONs or WBin.
	ResourceClasses: {}, //classes that can contain a resource of the system
	ResourceClasses_by_extension: {}, //used to associate JSONs to resources, only used by GRAPHs
	Globals: {}, //global scope to share info among scripts

	/**
	* Contains all the registered components
	* 
	* @property Components
	* @type {Object}
	* @default {}
	*/
	Components: {},

	/**
	* Contains all the registered material classes
	* 
	* @property MaterialClasses
	* @type {Object}
	* @default {}
	*/
	MaterialClasses: {},

	//vars used for uuid genereration
	_last_uid: 1,
	_uid_prefix: "@", //WARNING: must be one character long
	debug: false, //enable to see verbose output
	allow_static: true, //used to disable static instances in the editor
	allow_scripts: true, //if true, then Script components and Graphs can contain code

	//for HTML GUI
	_gui_element: null,
	_gui_style: null,

	/**
	* Generates a UUID based in the user-agent, time, random and sequencial number. Used for Nodes and Components.
	* @method generateUId
	* @return {string} uuid
	*/
	generateUId: function ( prefix, suffix ) {
		prefix = prefix || "";
		suffix = suffix || "";
		var str = this._uid_prefix + prefix + (window.navigator.userAgent.hashCode() % 0x1000000).toString(16) + "-"; //user agent
		str += (GL.getTime()|0 % 0x1000000).toString(16) + "-"; //date
		str += Math.floor((1 + Math.random()) * 0x1000000).toString(16) + "-"; //rand
		str += (this._last_uid++).toString(16); //sequence
		str += suffix;
		return str; 
	},

	/**
	* validates name string to ensure there is no forbidden characters
	* valid characters are letters, numbers, spaces, dash, underscore and dot
	* @method validateName
	* @param {string} name
	* @return {boolean} 
	*/
	validateName: function(v)
	{
		var exp = /^[a-z\s0-9-_.]+$/i; //letters digits and dashes
		return v.match(exp);
	},

	valid_property_types: ["String","Number","Boolean","color","vec2","vec3","vec4","quat","mat3","mat4","Resource","Animation","Texture","Prefab","Mesh","ShaderCode","node","component"],
	
	//used when creating a property to a component, to see if the type is valid
	validatePropertyType: function(v)
	{
		if(	this.valid_property_types.indexOf(v) == -1 )
		{
			console.error( v + " is not a valid property value type." );
			return false;
		}
		return true;
	},

	_catch_exceptions: false, //used to try/catch all possible callbacks (used mostly during development inside an editor) It is linked to LScript too

	/**
	* Register a component (or several) so it is listed when searching for new components to attach
	*
	* @method registerComponent
	* @param {Component} component component class to register
	* @param {String} old_classname [optional] the name of the component that this class replaces (in case you are renaming it)
	*/
	registerComponent: function( component, old_classname ) { 

		//to know from what file does it come from
		//console.log( document.currentScript.src );

		//allows to register several at the same time
		var name = ONE.getClassName( component );

		if(old_classname && old_classname.constructor !== String)
			throw("old_classname must be null or a String");

		//save the previous class in case we are replacing it
		var old_class = this.Components[ old_classname || name ];

		//register
		this.Components[ name ] = component; 
		component.is_component = true;	
		component.resource_type = "Component";

		//Helper: checks for errors
		if( !!component.prototype.onAddedToNode != !!component.prototype.onRemovedFromNode ||
			!!component.prototype.onAddedToScene != !!component.prototype.onRemovedFromScene )
			console.warn("%c Component "+name+" could have a bug, check events: " + name , "font-size: 2em");
		if( component.prototype.getResources && !component.prototype.onResourceRenamed )
			console.warn("%c Component "+name+" could have a bug, it uses resources but doesnt implement onResourceRenamed, this could lead to problems when resources are renamed.", "font-size: 1.2em");
		//if( !component.prototype.serialize || !component.prototype.configure )
		//	console.warn("%c Component "+name+" could have a bug, it doesnt have a serialize or configure method. No state will be saved.", "font-size: 1.2em");

		//add stuff to the class
		if(!component.actions)
			component.actions = {};

		//add default methods
		ONE.extendClass( component, ONE.BaseComponent );
		BaseComponent.addExtraMethods( component );

		if( ONE.debug )
		{
			var c = new component();
			var r = c.serialize();
			if(!r.object_class)
				console.warn("%c Component "+name+" could have a bug, serialize() method has object_class missing.", "font-size: 1.2em");
		}

		//event
		LEvent.trigger(LS, "component_registered", component ); 

		if(ONE.GlobalScene) //because main components are create before the global scene is created
		{
			this.replaceComponentClass( ONE.GlobalScene, old_classname || name, name );
			if( old_classname != name )
				this.unregisterComponent( old_classname );
		}

	},

	/**
	* Unregisters a component from the system (although existing instances are kept in the scene)
	*
	* @method unregisterComponent
	* @param {String} name the name of the component to unregister
	*/
	unregisterComponent: function( name ) { 
		//not found
		if(!this.Components[name])
			return;
		//delete from the list of component (existing components will still exists)
		delete this.Components[name];
	},


	/**
	* Tells you if one class is a registered component class
	*
	* @method isClassComponent
	* @param {ComponentClass} comp component class to evaluate
	* @return {boolean} true if the component class is registered
	*/
	isClassComponent: function( comp_class )
	{
		var name = this.getClassName( comp_class );
		return !!this.Components[name];
	},

	/**
	* Replaces all components of one class in the scene with components of another class
	*
	* @method replaceComponentClass
	* @param {Scene} scene where to apply the replace
	* @param {String} old_class_name name of the class to be replaced
	* @param {String} new_class_name name of the class that will be used instead
	* @return {Number} the number of components replaced
	*/
	replaceComponentClass: function( scene, old_class_name, new_class_name )
	{
		var proposed_class = new_class_name.constructor === String ? ONE.Components[ new_class_name ] : new_class_name;
		if(!proposed_class)
			return 0;

		//this may be a problem if the old class has ben unregistered...
		var old_class = null;
		
		if(	old_class_name.constructor === String )
		{
			old_class = ONE.Components[ old_class_name ];
			if( old_class )
				old_class_name = ONE.getClassName( old_class );
		}

		var num = 0;

		for(var i = 0; i < scene._nodes.length; ++i)
		{
			var node = scene._nodes[i];

			//search in current components
			for(var j = 0; j < node._components.length; ++j)
			{
				var comp = node._components[j];
				var comp_name = comp.constructor === ONE.MissingComponent ? comp._comp_class : ONE.getObjectClassName( comp );

				//it it is the exact same class then skip it
				if( comp.constructor === proposed_class )
					continue;

				//if this component is neither the old comp nor the new one
				if( comp_name != old_class_name && comp_name != new_class_name ) 
					continue;

				var info = comp.serialize();
				node.removeComponent( comp );

				var new_comp = null;
				try
				{
					new_comp = new proposed_class();
				}
				catch (err)
				{
					console.error("Error in replacement component constructor");
					console.error(err);
					continue;
				}
				node.addComponent( new_comp, j < node._components.length ? j : undefined );
				new_comp.configure( info );
				num++;
			}
		}

		return num;
	},

	/**
	* Register a resource class so we know which classes could be use as resources
	*
	* @method registerResourceClass
	* @param {ComponentClass} c component class to register
	*/
	registerResourceClass: function( resourceClass )
	{
		var class_name = ONE.getClassName( resourceClass );
		this.ResourceClasses[ class_name ] = resourceClass;
		this.Classes[ class_name ] = resourceClass;
		resourceClass.is_resource = true;

		if( !resourceClass.FORMAT )
			resourceClass.FORMAT = {};

		resourceClass.FORMAT.resource = class_name;
		resourceClass.FORMAT.resourceClass = resourceClass;

		if( resourceClass.EXTENSION )
			resourceClass.FORMAT.extension = resourceClass.EXTENSION.toLowerCase();

		var extension = resourceClass.FORMAT.extension;
		if(!extension && resourceClass != ONE.Resource )
			console.warn("Resource without extension info? " + class_name );
		else
		{
			this.ResourceClasses_by_extension[ extension ] = resourceClass;
			resourceClass.EXTENSION = extension;
		}

		if( ONE.Formats && extension )
		{
			var format_info = ONE.Formats.supported[ extension ];
			if( !format_info )
				ONE.Formats.supported[ extension ] = format_info = resourceClass.FORMAT;
			else
			{
				if(!format_info.resourceClass)
					format_info.resourceClass = resourceClass;
				//else if(format_info.resourceClass != resourceClass) //animations and prefab use the same file extension
				//	console.warn("format has resourceClass that do not match this resource: ", ONE.getClassName(format_info.resourceClass), ONE.getClassName(resourceClass) );
			}
		}

		//some validation here? maybe...
	},

	//Coroutines that allow to work with async functions
	coroutines: {},

	addWaitingCoroutine: function( resolve, event )
	{
		event = event || "render";
		var coroutines = this.coroutines[ event ];
		if(!coroutines)
			coroutines = this.coroutines[ event ] = [];
		coroutines.push( resolve );
	},

	triggerCoroutines: function( event, data )
	{
		event = event || "render";
		var coroutines = this.coroutines[ event ];
		if(!coroutines)
			return;
		for(var i = 0; i < coroutines.length; ++i)
			ONE.safeCall( coroutines[i], data ); //call resolves
		coroutines.length = 0;
	},

	createCoroutine: function( event )
	{
		return new Promise(function(resolve){
			ONE.addWaitingCoroutine( resolve, event );
		});
	},

	/**
	* Returns a Promise that will be fulfilled once the time has passed
	* @method sleep
	* @param {Number} ms time in milliseconds
	* @return {Promise} 
	*/
	sleep: function(ms) {
	  return new Promise( function(resolve){ setTimeout(resolve, ms); });
	},

	/**
	* Returns a Promise that will be fulfilled when the next frame is rendered
	* @method nextFrame
	* @return {Promise} 
	*/
	nextFrame: function( skip_request )
	{
		if(!skip_request)
			ONE.GlobalScene.requestFrame();
		return new Promise(function(resolve){
			ONE.addWaitingCoroutine( resolve, "render" );
		});
	},

	/**
	* Is a wrapper for callbacks that throws an LS "exception" in case something goes wrong (needed to catch the error from the system and editor)
	* @method safeCall
	* @param {function} callback
	* @param {array} params
	* @param {object} instance
	*/
	safeCall: function(callback, params, instance)
	{
		if(!ONE.catch_exceptions)
		{
			if(instance)
				return callback.apply( instance, params );
			return callback( params );
		}

		try
		{
			return callback.apply( instance, params );
		}
		catch (err)
		{
			LEvent.trigger(LS,"exception",err);
			//test this
			//throw new Error( err.stack );
			console.error( err.stack );
		}
	},

	/**
	* Is a wrapper for setTimeout that throws an LS "code_error" in case something goes wrong (needed to catch the error from the system)
	* @method setTimeout
	* @param {function} callback
	* @param {number} time in ms
	* @param {number} timer_id
	*/
	setTimeout: function(callback, time)
	{
		if(!ONE.catch_exceptions)
			return setTimeout( callback,time );

		try
		{
			return setTimeout( callback,time );
		}
		catch (err)
		{
			LEvent.trigger(LS,"exception",err);
		}
	},

	/**
	* Is a wrapper for setInterval that throws an LS "code_error" in case something goes wrong (needed to catch the error from the system)
	* @method setInterval
	* @param {function} callback
	* @param {number} time in ms
	* @param {number} timer_id
	*/
	setInterval: function(callback, time)
	{
		if(!ONE.catch_exceptions)
			return setInterval( callback,time );

		try
		{
			return setInterval( callback,time );
		}
		catch (err)
		{
			LEvent.trigger(LS,"exception",err);
		}
	},

	/**
	* copy the properties (methods and properties) of origin class into target class
	* @method extendClass
	* @param {Class} target
	* @param {Class} origin
	*/
	extendClass: function( target, origin ) {
		for(var i in origin) //copy class properties
		{
			if(target.hasOwnProperty(i))
				continue;
			target[i] = origin[i];
		}

		if(origin.prototype) //copy prototype properties
			for(var i in origin.prototype) //only enumerables
			{
				if(!origin.prototype.hasOwnProperty(i)) 
					continue;

				if(target.prototype.hasOwnProperty(i)) //avoid overwritting existing ones
					continue;

				//copy getters 
				if(origin.prototype.__lookupGetter__(i))
					target.prototype.__defineGetter__(i, origin.prototype.__lookupGetter__(i));
				else 
					target.prototype[i] = origin.prototype[i];

				//and setters
				if(origin.prototype.__lookupSetter__(i))
					target.prototype.__defineSetter__(i, origin.prototype.__lookupSetter__(i));
			}
	},

	/**
	* Clones an object (no matter where the object came from)
	* - It skip attributes starting with "_" or "jQuery" or functions
	* - it tryes to see which is the best copy to perform
	* - to the rest it applies JSON.parse( JSON.stringify ( obj ) )
	* - use it carefully
	* @method cloneObject
	* @param {Object} object the object to clone
	* @param {Object} target=null optional, the destination object
	* @param {bool} recursive=false optional, if you want to encode objects recursively
	* @param {bool} only_existing=false optional, only assign to methods existing in the target object
	* @param {bool} encode_objets=false optional, if a special object is found, encode it as ["@ENC",node,object]
	* @return {Object} returns the cloned object (target if it is specified)
	*/
	cloneObject: function( object, target, recursive, only_existing, encode_objects )
	{
		if(object === undefined)
			return undefined;
		if(object === null)
			return null;

		//base type
		switch( object.constructor )
		{
			case String:
			case Number:
			case Boolean:
				return object;
		}

		//typed array
		if( object.constructor.BYTES_PER_ELEMENT )
		{
			if(!target)
				return new object.constructor( object );
			if(target.set)
				target.set(object);
			else if(target.construtor === Array)
			{
				for(var i = 0; i < object.length; ++i)
					target[i] = object[i];
			}
			else
				throw("cloneObject: target has no set method");
			return target;
		}

		var o = target;
		if(o === undefined || o === null)
		{
			if(object.constructor === Array)
				o = [];
			else
				o = {};
		}

		//copy every property of this object
		for(var i in object)
		{
			if(i[0] == "@" || i[0] == "_" || i.substr(0,6) == "jQuery") //skip vars with _ (they are private) or '@' (they are definitions)
				continue;

			if(only_existing && !target.hasOwnProperty(i) && !target.__proto__.hasOwnProperty(i) ) //target[i] === undefined)
				continue;

			//if(o.constructor === Array) //not necessary
			//	i = parseInt(i);

			var v = object[i];
			if(v == null)
				o[i] = null;			
			else if ( isFunction(v) ) //&& Object.getOwnPropertyDescriptor(object, i) && Object.getOwnPropertyDescriptor(object, i).get )
				continue;//o[i] = v;
			else if (v.constructor === File ) 
				o[i] = null;
			else if (v.constructor === Number || v.constructor === String || v.constructor === Boolean ) //elemental types
				o[i] = v;
			else if( v.buffer && v.byteLength && v.buffer.constructor === ArrayBuffer ) //typed arrays are ugly when serialized
			{
				if(o[i] && v && only_existing) 
				{
					if(o[i].length == v.length) //typed arrays force to fit in the same container
						o[i].set( v );
				}
				else
					o[i] = new v.constructor(v); //clone typed array
			}
			else if ( v.constructor === Array ) //clone regular array (container and content!)
			{
				//not safe to use concat or slice(0) because it doesnt clone content, only container
				if( o[i] && o[i].set && o[i].length >= v.length ) //reuse old container
				{
					o[i].set(v);
					continue;
				}

				if( encode_objects && target && v[0] == "@ENC" ) //encoded object (component, node...)
				{
					var decoded_obj = ONE.decodeObject(v);
					o[i] = decoded_obj;
					if(!decoded_obj) //object not found
					{
						if( ONE._pending_encoded_objects )
							ONE._pending_encoded_objects.push([o,i,v]);
						else
							console.warn( "Object UID referencing object not found in the scene:", v[2] );
					}
				}
				else
					o[i] = ONE.cloneObject( v ); 
			}
			else //Objects: 
			{
				if( v.constructor.is_resource )
				{
					console.error("Resources cannot be saved as a property of a component nor script, they must be saved individually as files in the file system. If assigning them to a component/script use private variables (name start with underscore) to avoid being serialized.");
					continue;
				}

				if( encode_objects && !target )
				{
					o[i] = ONE.encodeObject(v);
					continue;
				}

				if( v.constructor !== Object && !target && !v.toJSON )
				{
					console.warn("Cannot clone internal classes:", ONE.getObjectClassName( v )," When serializing an object I found a var with a class that doesnt support serialization. If this var shouldnt be serialized start the name with underscore.'");
					continue;
				}

				if( v.toJSON )
					o[i] = v.toJSON();
				else if( recursive )
					o[i] = ONE.cloneObject( v, null, true );
				else {
					if(v.constructor !== Object && ONE.Classes[ ONE.getObjectClassName(v) ])
						console.warn("Cannot clone internal classes:", ONE.getObjectClassName(v)," When serializing an object I found a var with a class that doesnt support serialization. If this var shouldnt be serialized start the name with underscore.'" );

					if(ONE.catch_exceptions)
					{
						try
						{
							//prevent circular recursions //slow but safe
							o[i] = JSON.parse( JSON.stringify(v) );
						}
						catch (err)
						{
							console.error(err);
						}
					}
					else //slow but safe
					{
						o[i] = JSON.parse( JSON.stringify(v) );
					}
				}
			}
		}
		return o;
	},

	encodeObject: function( obj )
	{
		if( !obj || obj.constructor === Number || obj.constructor === String || obj.constructor === Boolean || obj.constructor === Object ) //regular objects
			return obj;
		if( obj.constructor.is_component && obj._root) //in case the value of this property is an actual component in the scene
			return [ "@ENC", ONE.TYPES.COMPONENT, obj.getLocator(), ONE.getObjectClassName( obj ) ];
		if( obj.constructor == ONE.SceneNode && obj._in_tree) //in case the value of this property is an actual node in the scene
			return [ "@ENC", ONE.TYPES.SCENENODE, obj.uid ];
		if( obj.constructor == ONE.Scene)
			return [ "@ENC", ONE.TYPES.SCENE, obj.fullpath ]; //weird case
		if( obj.serialize || obj.toJSON )
		{
			//return obj.serialize ? obj.serialize() : obj.toJSON(); //why not this?
			return [ "@ENC", ONE.TYPES.OBJECT, obj.serialize ? obj.serialize() : obj.toJSON(), ONE.getObjectClassName( obj ) ];
		}
		console.warn("Cannot clone internal classes:", ONE.getObjectClassName( obj )," When serializing an object I found a property with a class that doesnt support serialization. If this property shouldn't be serialized start the name with underscore.'");
		return null;
	},

	decodeObject: function( data )
	{
		if(!data || data.constructor !== Array || data[0] != "@ENC" )
			return null;

		switch( data[1] )
		{
			case ONE.TYPES.COMPONENT: 
			case "node":  //legacy
			case ONE.TYPES.SCENENODE: 
				var obj = LSQ.get( data[2] );
				if( obj )
					return obj;
				return null;
				break;
				//return  break;
			case ONE.TYPES.SCENE: return null; break; //weird case
			case ONE.TYPES.OBJECT: 
			default:
				if( !data[2] || !data[2].object_class )
					return null;
				var ctor = ONE.Classes[ data[2].object_class ];
				if(!ctor)
					return null;
				var v = new ctor();
				v.configure(data[2]);
				return v;
			break;
		}
		return null;
	},

	//used during scene configure because when configuring objects they may have nodes encoded in the properties
	resolvePendingEncodedObjects: function()
	{
		if(!ONE._pending_encoded_objects)
		{
			console.warn("no pending enconded objects");
			return;
		}
		for(var i = 0; i < ONE._pending_encoded_objects.length; ++i)
		{
			var pending = ONE._pending_encoded_objects[i];
			var decoded_object = ONE.decodeObject(pending[2]);
			if(decoded_object)
				pending[0][ pending[1] ] = decoded_object;
			else
				console.warn("Decoded object not found when configuring from JSON");
		}
		ONE._pending_encoded_objects = null;
	},

	switchGlobalScene: function( scene )
	{
		if(scene === ONE.GlobalScene)
			return;
		if( scene === null || scene.constructor !== ONE.Scene )
			throw("Not an scene");
		var old_scene = ONE.GlobalScene;
		LEvent.trigger( LS, "global_scene_changed", scene );
		ONE.GlobalScene = scene;
	},

	/**
	* Clears all the uids inside this object and children (it also works with serialized object)
	* @method clearUIds
	* @param {Object} root could be a node or an object from a node serialization
	*/
	clearUIds: function( root, uids_removed )
	{
		uids_removed = uids_removed || {};

		if(root.uid)
		{
			uids_removed[ root.uid ] = root;
			delete root.uid;
		}

		//remove for embeded materials
		if(root.material && root.material.uid)
		{
			uids_removed[ root.material.uid ] = root.material;
			delete root.material.uid;
		}

		var components = root.components;
		if(!components && root.getComponents)
			components = root.getComponents();

		if(!components)
			return;

		if(components)
		{
			for(var i in components)
			{
				var comp = components[i];
				if(comp[1].uid)
				{
					uids_removed[ comp[1].uid ] = comp[1];
					delete comp[1].uid;
				}
				if(comp[1]._uid)
				{
					uids_removed[ comp[1]._uid ] = comp[1];
					delete comp[1]._uid;
				}
			}
		}

		var children = root.children;
		if(!children && root.getChildren)
			children = root.getChildren();

		if(!children)
			return;

		for(var i in children)
			ONE.clearUIds( children[i], uids_removed );

		return uids_removed;
	},


	/**
	* Returns an object class name (uses the constructor toString)
	* @method getObjectClassName
	* @param {Object} the object to see the class name
	* @return {String} returns the string with the name
	*/
	getObjectClassName: function(obj)
	{
		if (!obj)
			return;

		if(obj.constructor.fullname) //this is to overwrite the common name "Prefab" for a global name "ONE.Prefab"
			return obj.constructor.fullname;

		if(obj.constructor.name)
			return obj.constructor.name;

		var arr = obj.constructor.toString().match(
			/function\s*(\w+)/);

		if (arr && arr.length == 2) {
			return arr[1];
		}
	},

	/**
	* Returns an string with the class name
	* @method getClassName
	* @param {Object} class object
	* @return {String} returns the string with the name
	*/
	getClassName: function(obj)
	{
		if (!obj)
			return;

		//from function info, but not standard
		if(obj.name)
			return obj.name;

		//from sourcecode
		if(obj.toString) {
			var arr = obj.toString().match(
				/function\s*(\w+)/);
			if (arr && arr.length == 2) {
				return arr[1];
			}
		}
	},

	/**
	* Returns the public properties of one object and the type (not the values)
	* @method getObjectProperties
	* @param {Object} object
	* @return {Object} returns object with attribute name and its type
	*/
	//TODO: merge this with the locator stuff
	getObjectProperties: function( object )
	{
		if(object.getPropertiesInfo)
			return object.getPropertiesInfo();
		var class_object = object.constructor;
		if(class_object.properties)
			return class_object.properties;

		var o = {};
		for(var i in object)
		{
			//ignore some
			if(i[0] == "_" || i[0] == "@" || i.substr(0,6) == "jQuery") //skip vars with _ (they are private)
				continue;

			if(class_object != Object)
			{
				var hint = class_object["@"+i];
				if(hint && hint.type)
				{
					o[i] = hint.type;
					continue;
				}
			}

			var v = object[i];
			if(v == null)
				o[i] = null;
			else if ( isFunction(v) )//&& Object.getOwnPropertyDescriptor(object, i) && Object.getOwnPropertyDescriptor(object, i).get )
				continue; //o[i] = v;
			else if (  v.constructor === Boolean )
				o[i] = ONE.TYPES.BOOLEAN;
			else if (  v.constructor === Number )
				o[i] = ONE.TYPES.NUMBER;
			else if ( v.constructor === String )
				o[i] = ONE.TYPES.STRING;
			else if ( v.buffer && v.buffer.constructor === ArrayBuffer ) //typed array
			{
				if(v.length == 2)
					o[i] = ONE.TYPES.VEC2;
				else if(v.length == 3)
					o[i] = ONE.TYPES.VEC3;
				else if(v.length == 4)
					o[i] = ONE.TYPES.VEC4;
				else if(v.length == 9)
					o[i] = ONE.TYPES.MAT3;
				else if(v.length == 16)
					o[i] = ONE.TYPES.MAT4;
				else
					o[i] = 0;
			}
			else
				o[i] = 0;
		}
		return o;
	},

	//TODO: merge this with the locator stuff
	setObjectProperty: function( obj, name, value )
	{
		if(obj.setProperty)
			return obj.setProperty(name, value);
		obj[ name ] = value; //clone????
		if(obj.onPropertyChanged)
			obj.onPropertyChanged( name, value );
	},

	/**
	* Register a Material class so it is listed when searching for new materials to attach
	*
	* @method registerMaterialClass
	* @param {ComponentClass} comp component class to register
	*/
	registerMaterialClass: function( material_class )
	{ 
		var class_name = ONE.getClassName( material_class );

		//register
		this.MaterialClasses[ class_name ] = material_class;
		this.Classes[ class_name ] = material_class;

		//add extra material methods
		ONE.extendClass( material_class, Material );

		//event
		LEvent.trigger( LS, "materialclass_registered", material_class );
		material_class.resource_type = "Material";
		material_class.is_material = true;
	},

	/**
	* Returns an script context using the script name (not the node name), usefull to pass data between scripts.
	*
	* @method getScript
	* @param {String} name the name of the script according to the Script component.
	* @return {Object} the context of the script.
	*/
	getScript: function( name )
	{
		var script = ONE.Script.active_scripts[name];
		if(script)
			return script.context;
		return null;
	},

	getDebugRender: function()
	{
		if(!ONE.debug_render)
			ONE.debug_render = new ONE.DebugRender();
		return ONE.debug_render;
	},

	//we do it in a function to make it more standard and traceable
	//used by the editor to show the error
	dispatchCodeError: function( err, line, resource, extra )
	{
		var error_info = { error: err, line: line, resource: resource, extra: extra };
		console.error(error_info);
		LEvent.trigger( this, "code_error", error_info );
	},

	//used to tell the editor it must remove the error marker
	dispatchNoErrors: function( resource, extra )
	{
		var info = { resource: resource, extra: extra };
		LEvent.trigger( this, "code_no_errors", info );
	},

	convertToString: function( data )
	{
		if(!data)
			return "";
		if(data.constructor === String)
			return data;
		if(data.constructor === Object)
			return JSON.stringify( object.serialize ? object.serialize() : object );
		if(data.constructor === ArrayBuffer)
			data = new Uint8Array(data);
		return String.fromCharCode.apply(null,data);
	},

	
	/**
	* Checks if this locator belongs to a property inside a prefab, which could be tricky in some situations
	* @method checkLocatorBelongsToPrefab
	**/
	checkLocatorBelongsToPrefab: function( locator, root )
	{
		root = root || ONE.GlobalScene.root;
		var property_path = locator.split("/");

		if( !property_path.length )
			return null;

		var node = LSQ.get( property_path[0], root );
		if(!node)
			return null;

		return node.insidePrefab();
	},

	/**
	* Used to convert locators so instead of using UIDs for properties it uses Names
	* This is used when you cannot rely on the UIDs because they belong to prefabs that could change them
	* @method convertLocatorFromUIDsToName
	* @param {String} locator string with info about a property (p.e. "my_node/Transform/y")
	* @param {boolean} use_basename if you want to just use the node name, othewise it uses the fullname (name with path)
	* @param {ONE.SceneNode} root
	* @return {String} the result name without UIDs
	*/
	convertLocatorFromUIDsToName: function( locator, use_basename, root )
	{
		root = root || ONE.GlobalScene.root;
		var property_path = locator.split("/");

		if( !property_path.length )
			return null;

		if( property_path[0][0] !== ONE._uid_prefix && ( property_path.length == 1 || property_path[1][0] !== ONE._uid_prefix))
			return null; //is already using names

		var node = LSQ.get( property_path[0], root );
		if(!node)
		{
			console.warn("getIDasName: node not found in ONE.GlobalScene: " + property_path[0] );
			return false;
		}

		if(!node.name)
		{
			console.warn("getIDasName: node without name?");
			return false;
		}

		//what about the component?
		if( property_path.length > 1 && property_path[1][0] == ONE._uid_prefix )
		{
			var comp = ONE.GlobalScene.findComponentByUId( property_path[1] );
			if(comp)
			{
				var comp_name = comp.constructor.name;
				if(comp_name == "Script" || comp_name == "ScriptFromFile")
				{
					comp_name = comp.name;
					if( comp_name == "unnamed" )
					{
						console.error("converting component UIDs to component name, but property belongs to a Script without name. You must name the script to avoid errors.");
						comp_name = comp.constructor.name;
					}
				}
				property_path[1] = comp_name;	
			}
		}

		var result = property_path.concat();
		if(use_basename)
			result[0] = node.name;
		else
			result[0] = node.fullname;
		return result.join("/");
	},

	/**
	* clears the global scene and the resources manager
	*
	* @method reset
	*/
	reset: function()
	{
		ONE.GlobalScene.clear();
		ONE.ResourcesManager.reset();
		LEvent.trigger( LS, "reset" );
	},

	log: function()
	{
		console.log.apply( console, arguments );
	},

	stringToValue: function( v )
	{
		var value = v;
		try
		{
			value = JSON.parse(v);
		}
		catch (err)
		{
			console.error( "Not a valid value: " + v );
		}
		return value;
	},

	isValueOfType: function( value, type )
	{
		if(value === null || value === undefined)
		{
			switch (type)
			{
				case "float": 
				case "sampler2D": 
				case "samplerCube":
				case ONE.TYPES.NUMBER: 
				case ONE.TYPES.VEC2: 
				case ONE.TYPES.VEC3:
				case ONE.TYPES.VEC4:
				case ONE.TYPES.COLOR:
				case ONE.TYPES.COLOR4:
				case "mat3": 
				case "mat4":
					return false;
			}
			return true;
		}

		switch (type)
		{
			//used to validate shaders
			case "float": 
			case "sampler2D": 
			case "samplerCube":
			case ONE.TYPES.NUMBER: return isNumber(value);
			case ONE.TYPES.VEC2: return value.length === 2;
			case ONE.TYPES.VEC3: return value.length === 3;
			case ONE.TYPES.VEC4: return value.length === 4;
			case ONE.TYPES.COLOR: return value.length === 3;
			case ONE.TYPES.COLOR4: return value.length === 4;
			case "mat3": return value.length === 9;
			case "mat4": return value.length === 16;
		}
		return true;
	},

	//solution from http://stackoverflow.com/questions/979975/how-to-get-the-value-from-the-url-parameter
	queryString: function () {
	  // This function is anonymous, is executed immediately and 
	  // the return value is assigned to QueryString!
	  var query_string = {};
	  var query = window.location.search.substring(1);
	  var vars = query.split("&");
	  for (var i=0;i<vars.length;i++) {
		var pair = vars[i].split("=");
			// If first entry with this name
		if (typeof query_string[pair[0]] === "undefined") {
		  query_string[pair[0]] = decodeURIComponent(pair[1]);
			// If second entry with this name
		} else if (typeof query_string[pair[0]] === "string") {
		  var arr = [ query_string[pair[0]],decodeURIComponent(pair[1]) ];
		  query_string[pair[0]] = arr;
			// If third or later entry with this name
		} else {
		  query_string[pair[0]].push(decodeURIComponent(pair[1]));
		}
	  } 
		return query_string;
	}(),

	downloadFile: function( filename, data, dataType )
	{
		if(!data)
		{
			console.warn("No file provided to download");
			return;
		}

		if(!dataType)
		{
			if(data.constructor === String )
				dataType = 'text/plain';
			else
				dataType = 'application/octet-stream';
		}

		var file = null;
		if(data.constructor !== File && data.constructor !== Blob)
			file = new Blob( [ data ], {type : dataType});
		else
			file = data;

		var url = URL.createObjectURL( file );
		var element = document.createElement("a");
		element.setAttribute('href', url);
		element.setAttribute('download', filename );
		element.style.display = 'none';
		document.body.appendChild(element);
		element.click();
		document.body.removeChild(element);
		setTimeout( function(){ URL.revokeObjectURL( url ); }, 1000*60 ); //wait one minute to revoke url
	}
}

var LS = ONE; //LEGACY

//ensures no exception is catched by the system (useful for developers)
Object.defineProperty( LS, "catch_exceptions", { 
	set: function(v){ 
		this._catch_exceptions = v; 
		LScript.catch_exceptions = v; 
		LScript.catch_important_exceptions = v;
	},
	get: function() { return this._catch_exceptions; },
	enumerable: true
});

//ensures no exception is catched by the system (useful for developers)
Object.defineProperty( LS, "block_scripts", { 
	set: function(v){ 
		ONE._block_scripts = v; 
		LScript.block_execution = v; 
	},
	get: function() { 
		return !!ONE._block_scripts;
	},
	enumerable: true
});


//Add some classes
ONE.Classes.WBin = ONE.WBin = global.WBin = WBin;

/**
* LSQ allows to set or get values easily from the global scene, using short strings as identifiers
* similar to jQuery and the DOM:  LSQ("nod_name").material = ...
*
* @class  LSQ
*/
function LSQ(v)
{
	return LSQ.get(v);
}

/**
* Assigns a value to a property of one node in the scene, just by using a string identifier
* Example:  LSQ.set("mynode|a_child/MeshRenderer/enabled",false);
*
* @method set
* @param {String} locator the locator string identifying the property
* @param {*} value value to assign to property
*/
LSQ.set = function( locator, value, root, scene )
{
	scene = scene || ONE.GlobalScene;
	if(!root)
		scene.setPropertyValue( locator, value );
	else
	{
		if(root.constructor === ONE.SceneNode)
		{
			var path = locator.split("/");
			var node = root.findNodeByUId( path[0] );
			if(!node)
				return null;
			return node.setPropertyValueFromPath( path.slice(1), value );
		}
	}

	scene.requestFrame();
}

/**
* Retrieves the value of a property of one node in the scene, just by using a string identifier
* Example: var value = LSQ.get("mynode|a_child/MeshRenderer/enabled");
*
* @method get
* @param {String} locator the locator string identifying the property
* @return {*} value of the property
*/
LSQ.get = function( locator, root, scene )
{
	if(!locator) //sometimes we have a var with a locator that is null
		return null;
	scene = scene || ONE.GlobalScene;
	var info;
	if(!root)
		info = scene.getPropertyInfo( locator );
	else
	{
		if(root.constructor === ONE.SceneNode)
		{
			var path = locator.split("/");
			var node = root.findNodeByUId( path[0] );
			if(!node)
				return null;
			info = node.getPropertyInfoFromPath( path.slice(1) );
		}
	}
	if(info)
		return info.value;
	return null;
}

/**
* Shortens a locator that uses unique identifiers to a simpler one, but be careful, because it uses names instead of UIDs it could point to the wrong property
* Example: "@NODE--a40661-1e8a33-1f05e42-56/@COMP--a40661-1e8a34-1209e28-57/size" -> "node|child/Collider/size"
*
* @method shortify
* @param {String} locator the locator string to shortify
* @return {String} the locator using names instead of UIDs
*/
LSQ.shortify = function( locator, scene )
{
	if(!locator)
		return;

	var t = locator.split("/");
	var node = null;

	//already short
	if( t[0][0] != ONE._uid_prefix )
		return locator;

	scene = scene || ONE.GlobalScene;

	node = scene._nodes_by_uid[ t[0] ];
	if(!node) //node not found
		return locator;

	t[0] = node.getPathName();
	if(t[1])
	{
		if( t[1][0] == ONE._uid_prefix )
		{
			var compo = node.getComponentByUId(t[1]);
			if(compo)
				t[1] = ONE.getObjectClassName( compo );
		}
	}
	return t.join("/");
}

/**
* Assigns a value using the getLocatorInfo object instead of searching it again
* This is faster but if the locator points to a different object it wont work.
*
* @method setFromInfo
* @param {Object} info information of a location (obtain using scene.getLocatorInfo
* @param {*} value to assign
*/
LSQ.setFromInfo = function( info, value )
{
	if(!info || !info.target)
		return;
	var target = info.target;
	if( target.setPropertyValue  )
		if( target.setPropertyValue( info.name, value ) === true )
			return target;
	if( target[ info.name ] === undefined && info.value === undefined )
		return;
	target[ info.name ] = value;	
}

LSQ.getFromInfo = function( info )
{
	if(!info || !info.target)
		return;
	var target = info.target;
	var varname = info.name;
	var v = undefined;
	if( target.getPropertyValue )
		v = target.getPropertyValue( varname );
	if( v === undefined && target[ varname ] === undefined )
		return null;
	return v !== undefined ? v : target[ varname ];
}

//register resource classes
if(global.GL)
{
	GL.Mesh.EXTENSION = "wbin";
	GL.Texture.EXTENSION = "png";

	ONE.registerResourceClass( GL.Mesh );
	ONE.registerResourceClass( GL.Texture );

	ONE.Mesh = GL.Mesh;
	ONE.Texture = GL.Texture;
	ONE.Buffer = GL.Buffer;
	//ONE.Shader = GL.Shader; //this could be confussing since there is also ShaderBlocks etc in LiteScene
}


global.LSQ = LSQ;
