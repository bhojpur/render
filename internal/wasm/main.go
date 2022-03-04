package main

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"net/http"

	ctxsvr "github.com/bhojpur/web/pkg/context"
	utils "github.com/bhojpur/web/pkg/core/utils"
	websvr "github.com/bhojpur/web/pkg/engine"
	"github.com/bhojpur/web/pkg/filter/cors"
	"github.com/bhojpur/web/pkg/synthesis"
	test "github.com/bhojpur/web/test"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "भोजपुर जिला घर बा, तब कौना बात के डर बा !!")
}

func namasteHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "नमस्ते, %s!", r.FormValue("नाम"))
}

func main() {
	websvr.InsertFilter("*", websvr.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// CORS Post method issue
	websvr.InsertFilter("*", websvr.BeforeRouter, func(ctx *ctxsvr.Context) {
		if ctx.Input.Method() == "OPTIONS" {
			ctx.WriteString("ok")
		}
	})

	//websvr.DelStaticPath("/static")
	websvr.SetStaticPath("/static", "../../pkg/webui")
	websvr.InsertFilter("/", websvr.BeforeRouter, StaticContentHandler) // must have this for default page
	websvr.InsertFilter("/*", websvr.BeforeRouter, StaticContentHandler)
	websvr.Run() // custom configuration read fron ../conf/app.conf file
	websvr.AddFuncMap("*", http.HandlerFunc(indexHandler))
	websvr.AddFuncMap("/अभिवादन/:नाम", http.HandlerFunc(namasteHandler))
	websvr.AddFuncMap("/data",
		http.FileServer(
			&synthesis.AssetFS{
				Asset:     test.Asset,
				AssetDir:  test.AssetDir,
				AssetInfo: test.AssetInfo,
				Prefix:    "data",
				Fallback:  "index.html",
			}))
}

func StaticContentHandler(ctx *ctxsvr.Context) {
	urlPath := ctx.Request.URL.Path
	path := "."
	if urlPath == "/" {
		path += "/index.html"
	} else {
		path += urlPath
	}

	if utils.FileExists(path) {
		http.ServeFile(ctx.ResponseWriter, ctx.Request, path)
	} else {
		http.ServeFile(ctx.ResponseWriter, ctx.Request, "./index.html")
	}
}
