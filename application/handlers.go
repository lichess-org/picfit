package application

import (
	"encoding/json"
	"github.com/mholt/binding"
	"github.com/thoas/muxer"
	"github.com/thoas/picfit/image"
	"net/http"
)

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "404 not found", http.StatusNotFound)
	})
}

type Handler func(muxer.Response, *Request, *Application)

var ImageHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, true, true)

	if err != nil {
		panic(err)
	}

	res.SetHeaders(file.Headers, true)
	res.ResponseWriter.Write(file.Content())
}

var UploadHandler = func(res muxer.Response, req *http.Request, app *Application) {
	if app.SourceStorage == nil {
		res.Abort(500, "Your application doesn't have a source storage")
		return
	}

	var err error

	multipartForm := new(MultipartForm)
	errs := binding.Bind(req, multipartForm)
	if errs.Handle(res) {
		return
	}

	err = multipartForm.Upload(app.SourceStorage)

	if err != nil {
		panic(err)
	}

	file := &image.ImageFile{
		Filepath: multipartForm.Data.Filename,
		Storage:  app.SourceStorage,
	}

	content, err := json.Marshal(map[string]string{
		"filename": file.Filename(),
		"path":     file.Path(),
		"url":      file.URL(),
	})

	if err != nil {
		panic(err)
	}

	res.ContentType("application/json")
	res.ResponseWriter.Write(content)
}

var GetHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, false, false)

	if err != nil {
		panic(err)
	}

	content, err := json.Marshal(map[string]string{
		"filename": file.Filename(),
		"path":     file.Path(),
		"url":      file.URL(),
	})

	if err != nil {
		panic(err)
	}

	res.ContentType("application/json")
	res.ResponseWriter.Write(content)
}

var RedirectHandler Handler = func(res muxer.Response, req *Request, app *Application) {
	file, err := app.ImageFileFromRequest(req, false, false)

	if err != nil {
		panic(err)
	}

	res.PermanentRedirect(file.URL())
}
