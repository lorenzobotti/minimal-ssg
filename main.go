//TODO: move markdown parser to its own file
//TODO: numbered list support
//TODO: multiline code
//TODO: moving images

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
)

type Post struct {
	Title           string
	Content         string
	File            string
	PublicationDate string
}

type PostsSlice []Post

var htmlPostTemplate *template.Template
var htmlIndexTemplate *template.Template

func main() {
	//postTemplateFile, _ := ioutil.ReadFile("template.html")
	postTemplPath := flag.String("post-template", "templates/post_template.html", "")
	indexTemplPath := flag.String("index-template", "templates/index_template.html", "")
	outputFolder := flag.String("output-folder", "output/", "")
	postsFolder := flag.String("input-folder", "posts/", "")
	flag.Parse()
	var err error
	htmlPostTemplate, err = template.ParseFiles(path.Clean(*postTemplPath))
	if err != nil {
		fmt.Println("template file", *postTemplPath, "seems to not exist")
		panic(err)
	}
	htmlIndexTemplate, err = template.ParseFiles(path.Clean(*indexTemplPath))
	if err != nil {
		fmt.Println("template file", *indexTemplPath, "seems to not exist")
		panic(err)
	}

	fileNameRegex := regexp.MustCompile("(.*).md")

	postFiles, err := ioutil.ReadDir(path.Clean(*postsFolder))
	if err != nil {
		fmt.Println("error: input folder", postFiles, "does not exist!")
		panic(err)
	}
	posts := make([]Post, len(postFiles))

	//create the output folder
	os.Mkdir(*outputFolder, 0644)
	//we have to chmod the output folder to override the umask
	//otherwise we won't have permission to write to it
	//(unix only)
	err = os.Chmod(*outputFolder, 0755)
	if err != nil {
		fmt.Println("something's up with permissions in this folder.")
		fmt.Println("maybe check the umask?")
		panic(err)
	}
	//generate the post pages
	for i, post := range postFiles {
		fileName := post.Name()
		fileContent, err := ioutil.ReadFile(path.Join(*postsFolder, fileName))
		//skip this file if it can't be read
		if err != nil {
			fmt.Println(err)
			continue
		}
		posts[i] = markdownCompile(fileContent)

		//filename.md --> filename.html
		outputFile := fileNameRegex.ReplaceAllString(fileName, "$1.html")
		outputPath := path.Join(*outputFolder, outputFile)
		posts[i].File = outputFile
		fmt.Println(outputPath)

		//fmt.Println(outputPath)
		err = ioutil.WriteFile(outputPath, renderPost(posts[i]), 0644)
		if err != nil {
			fmt.Println(err)
		}
	}

	//sortPosts(posts)
	sort.Sort(sort.Reverse(PostsSlice(posts)))

	//generate the homepage
	var indexText bytes.Buffer
	htmlIndexTemplate.Execute(&indexText, posts)
	//fmt.Println(outputFolder + "index.html")
	ioutil.WriteFile(path.Join(*outputFolder, "index.html"), indexText.Bytes(), 0644)
}

func renderPost(post Post) []byte {
	var postText bytes.Buffer
	htmlPostTemplate.Execute(&postText, post)
	return postText.Bytes()
}

//implementing the sort.Interface interface for sorting PostsSlice by date

func (p PostsSlice) Len() int {
	return len(p)
}

func (p PostsSlice) Less(i, j int) bool {
	a := strings.Split(p[i].PublicationDate, "/")
	b := strings.Split(p[j].PublicationDate, "/")

	//if the dates are formatted incorrecly quit right away
	//if you try to access slice[i] with i < len(i), it panics
	if len(a) < 3 || len(b) < 3 {
		return true
	}

	//compare years, then months, then days
	for i := 2; i >= 0; i-- {
		aValue, _ := strconv.Atoi(a[i])
		bValue, _ := strconv.Atoi(b[i])
		if aValue < bValue {
			return true
		} else if aValue > bValue {
			return false
		}
	}

	//if we got here, the dates are identical
	return true
}

func (p PostsSlice) Swap(i, j int) {
	buffer := p[j]
	p[j] = p[i]
	p[i] = buffer
}
