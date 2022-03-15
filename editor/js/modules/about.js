var AboutModule = {
	name: "about",

	init: function()
	{
		LiteGUI.menubar.add("About", { callback: function() { 
			var dialog = new LiteGUI.Dialog({ title: "About Bhojpur Studio", closable: true, width: 400, height: 240} );
			dialog.content.style.fontSize = "1em";
			dialog.content.style.backgroundColor = "black";
			dialog.content.innerHTML = "<p>Bhojpur Studio version "+CORE.config.version+"</p><p>Copyright &copy; 2018 by <a href='https://www.bhojpur-consulting.com' target='_blank'>Bhojpur Consulting Private Limited</a>, India.</p><p>All rights reserved.</p>";
			dialog.show();
		}});
	}
}

CORE.registerModule( AboutModule );