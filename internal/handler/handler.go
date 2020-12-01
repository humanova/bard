// (c) 2020 Emir Erbasan (humanova)
// MIT License, see LICENSE for more details

package handler

import (
	"bard/config"
	"bard/internal/post"
	"fmt"
	"github.com/kennygrant/sanitize"
	"html/template"
	"net/http"
)

type CmsPageData struct {
	Posts   []post.Data
	Page    string
	APIPath string
}

func getJSONRequestData(w http.ResponseWriter, r *http.Request, data interface{}) {
	err := DecodeJSONBody(w, r, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func CreatePostHandler(conf config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p post.NewData
		getJSONRequestData(w, r, &p)

		filename :=  sanitize.Path(fmt.Sprintf("%s.md",p.Title))
		path := fmt.Sprintf("%s/%s", conf.PostDirectory, filename)

		err := post.CreatePost(path, p.Title, p.Text)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't create file : %s", filename),
				http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Successfully added new post  : %s", p.Title)
		fmt.Printf("added new post : %s\n", p.Title)
		return
	}
}

func UpdatePostHandler(conf config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p post.Data
		getJSONRequestData(w, r, &p)

		filename := sanitize.Path(fmt.Sprintf("%s.md",p.Filename))
		path := fmt.Sprintf("%s/%s", conf.PostDirectory, filename)

		err := post.UpdatePost(path, p.Title, p.Text)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't update file : %s", p.Filename),
				http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Successfully updated the post  : %s", p.Title)
		fmt.Printf("updated post : %s\n", p.Title)
		return
	}
}

func DeletePostHandler(conf config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var p post.DeleteData
		getJSONRequestData(w, r, &p)

		filename := sanitize.Path(fmt.Sprintf("%s.md",p.Filename))
		path := fmt.Sprintf("%s/%s", conf.PostDirectory, filename)

		err := post.RemovePost(path)
		if err != nil {
			http.Error(w, fmt.Sprintf("couldn't delete file : %s", p.Filename),
				http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Successfully deleted the post  : %s", p.Filename)
		fmt.Printf("deleted post : %s\n", p.Filename)
		return
	}
}

// TODO : handle multiple templates (templ withing templ etc.)
func CmsHandler(conf config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		updatePostFilename := query.Get("update_post")

		var data CmsPageData
		tpl, err := template.ParseGlob("web/template/*.gohtml")
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		if updatePostFilename == "" {
			posts, err := post.GetPosts(conf.PostDirectory)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}
			data = CmsPageData{posts, "index", conf.ListenPrefixPath}

		} else {
			_post, err := post.GetPost(fmt.Sprintf("%s.md", updatePostFilename), conf.PostDirectory)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError)
				return
			}
			data = CmsPageData{[]post.Data{_post}, "update_post", conf.ListenPrefixPath}
		}

		tpl.ExecuteTemplate(w, "cms.gohtml", data)
		return
	}
}
