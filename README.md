# Bhojpur Render - Model Processing Engine

The Bhojpur Render is a high performance 2D/3D graphic rendering engine applied within the [Bhojpur.NET Platform](https://github.com/bhojpur/platform) for delivery of collaborative and distributed applications or services. It features `client-side` and `server-side` framework for `desktop`, `web`, and `mobile` application development. 

## Key Features

This rendering framework provides the following capabilities

- client-side rendering library for desktop applications (e.g., macOS, Linux, Windows)
- server-side rendering library for web-based delivery

## Pre-requisites

The `Go` >= 1.17 is primary requirement.

- [wasm_exec.js](https://github.com/tinygo-org/tinygo/blob/release/targets/wasm_exec.js) for helper functions. Ensure you are using the same version of `wasm_exec.js` as the version of `tinygo` you are using to compile.

### [WebAssembly](https://github.com/WebAssembly)

You must install [TinyGo](https://github.com/tinygo-org/tinygo) to be able ro build `wasm` applications
using `Go` programming language.

```bash
$ brew tap tinygo-org/tools
$ brew install tinygo
$ tinygo version
```

Generate the `main.wasm` Javascript for your application using following command.

```bash
$ tinygo build -o main.wasm -target wasm ./main.go
```

Now, create an `index.html` page that loads the `main.wasm` file. The general steps required
to run the WebAssembly file in the web browser includes loading it into JavaScript with
`WebAssembly.instantiateStreaming`, or `WebAssembly.instantiate` in some browsers:

```js
const go = new Go(); // Defined in wasm_exec.js
const WASM_URL = 'wasm.wasm';

var wasm;

if ('instantiateStreaming' in WebAssembly) {
	WebAssembly.instantiateStreaming(fetch(WASM_URL), go.importObject).then(function (obj) {
		wasm = obj.instance;
		go.run(wasm);
	})
} else {
	fetch(WASM_URL).then(resp =>
		resp.arrayBuffer()
	).then(bytes =>
		WebAssembly.instantiate(bytes, go.importObject).then(function (obj) {
			wasm = obj.instance;
			go.run(wasm);
		})
	)
}
```

If you have used explicit exports, you can call them by invoking them under the `wasm.exports`
namespace. See the [`export`](./export/wasm.js) directory for an example of this.

In addition to the `javascript`, it is important that the `wasm` file is served with the
[`Content-Type`](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Content-Type)
header set to `application/wasm`. Without it, most web browsers won't run it.

Finally, host your `wasm` application using `rendersvr` CLI tool

```bash
$ rendersvr wasm
```
