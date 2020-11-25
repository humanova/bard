// (c) 2020 Emir Erbasan (humanova)
// MIT License, see LICENSE for more details

package main

import (
	"bard/internal/post"
	"fmt"
	"github.com/kennygrant/sanitize"
	"html/template"
	"net/http"
)

type CmsPageData struct {
	Posts   []post.Data
	APIPath string
}

func getJSONRequestData(w http.ResponseWriter, r *http.Request, data interface{}) {
	err := DecodeJSONBody(w, r, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func (config Config) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var p post.NewData
	getJSONRequestData(w, r, &p)

	filename := fmt.Sprintf("%s/%s.md", config.PostDirectory, sanitize.Path(p.Title))

	err := post.CreatePost(filename, p.Title, p.Text)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't create file : %s", filename),
			http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully added new post  : %s", p.Title)
	fmt.Printf("added new post : %s\n", p.Title)
	return
}

func (config Config) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	var p post.Data
	getJSONRequestData(w, r, &p)

	filename := fmt.Sprintf("%s/%s.md", config.PostDirectory, sanitize.Path(p.Filename))

	err := post.UpdatePost(filename, p.Title, p.Text)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't update file : %s", p.Filename),
			http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully updated the post  : %s", p.Title)
	fmt.Printf("updated post : %s\n", p.Title)
	return
}

func (config Config) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	var p post.DeleteData
	getJSONRequestData(w, r, &p)

	filename := fmt.Sprintf("%s/%s.md", config.PostDirectory, sanitize.Path(p.Filename))

	err := post.RemovePost(filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("couldn't delete file : %s", p.Filename),
			http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Successfully deleted the post  : %s", p.Filename)
	fmt.Printf("deleted post : %s\n", p.Filename)
	return
}

func (config Config) cmsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	updatePostFilename := query.Get("update_post")

	var data CmsPageData
	var tpl *template.Template

	if updatePostFilename == "" {
		posts, err := post.GetPosts(config.PostDirectory)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		data = CmsPageData{posts, config.ListenPrefixPath}
		tpl = template.Must(template.ParseFiles("web/template/cms.gohtml"))

	} else {
		_post, err := post.GetPost(fmt.Sprintf("%s.md", updatePostFilename), config.PostDirectory)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}
		data = CmsPageData{[]post.Data{_post}, config.ListenPrefixPath}
		tpl = template.Must(template.ParseFiles("web/template/cms_update.gohtml"))
	}

	tpl.Execute(w, data)
	return
}
