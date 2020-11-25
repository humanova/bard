// (c) 2020 Emir Erbasan (humanova)
// MIT License, see LICENSE for more details

package post

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

type Data struct {
	Title    string
	Text     string
	Filename string
}

type NewData struct {
	Title string
	Text  string
}

type DeleteData struct {
	Filename string
}

func CreatePost(filename string, title string, text string) error {
	pText := fmt.Sprintf("---\ntitle: %s\n---\n%s", title, text)

	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(pText)
	if err != nil {
		return err
	}
	return nil
}

func UpdatePost(filename string, title string, text string) error {
	pText := fmt.Sprintf("---\ntitle: %s\n---\n%s", title, text)

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(pText)
	if err != nil {
		return err
	}
	return nil
}

func RemovePost(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

func GetPosts(postDirectory string) ([]Data, error) {
	var posts []Data

	files, err := ioutil.ReadDir(postDirectory)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		p, err := GetPost(file.Name(), postDirectory)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, nil
}

func GetPost(filename string, postDirectory string) (Data, error) {
	f, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", postDirectory, filename))
	if err != nil {
		p := Data{Filename: filename, Title: "[error]", Text: "[error]"}
		return p, err
	}
	content := string(f)

	titleRegex, _ := regexp.Compile("title: (.*)")
	title := strings.Split(titleRegex.FindString(content), ": ")[1]
	text := strings.Split(content, "\n---\n")[1]

	p := Data{Filename: filename, Title: title, Text: text}
	return p, nil
}
