var IntroModule = {
	name: "intro",

	preferences: {
		show_intro_dialog: true
	},

	init: function()
	{
		if( this.preferences.show_intro_dialog !== false)
			this.showIntroDialog();	
	},

	showIntroDialog: function()
	{
		var dialog = new LiteGUI.Dialog("intro_dialog",{ width: 400, height: 400, closable: true });
		dialog.content.style.fontSize = "1em";
		dialog.content.style.backgroundColor = "white";
		dialog.content.innerHTML = ""+
			"<p class='center'><img width='127' height='38' target='_blank' src='https://static.bhojpur.net/image/logo.png'/></p>" + 
			"<p class='header center'>Welcome to Bhojpur Studio</p>" +
			"<p class='msg center'>An online 3D Modeling software based on <a target='_blank' href='https://www.bhojpur.net'>Bhojpur.NET Platform</a>.</p>" +
			"<p class='msg center'>Please feel free to visit <a target='_blank' href='https://www.bhojpur-consulting.com'>Bhojpur Consulting</a> site or <a target='_blank' href='https://desk.bhojpur-consulting.com'>Global Support Centre</a> or <a href='https://github.com/bhojpur/render'>GitHub</a> page, if you have any queries and/or suggestions.</p>";

		dialog.on_close = function()
		{
			IntroModule.preferences.show_intro_dialog = false;
		}
	
		dialog.addButton("Close");
		dialog.show();
		dialog.center();
		dialog.fadeIn();

		var links = dialog.content.querySelectorAll("a");
		for(var i = 0; i < links.length; i++)
			links[i].addEventListener("click",prevent_this, true);
		dialog.root.addEventListener("click",close_this);

		function prevent_this(e){
			e.stopImmediatePropagation();
			e.stopPropagation();
			return false;
		}

		function close_this(){
			dialog.close();
		}
	}
}

CORE.registerModule( IntroModule );