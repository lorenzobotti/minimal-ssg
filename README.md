# Currently name-less static site generator

This is a minimal [static site generator](https://en.wikipedia.org/wiki/Web_template_system#Static_site_generators) and regex-based [Markdown](https://www.markdownguide.org/getting-started/) parser built as an exercise in golang. It has no dependencies other than the Go standard library.

## How to compile

```
git clone this-repo
cd this-repo
go build *.go -o ssg
```

## Flags / how to use
Flags:

```
--input-folder
   location of the md files for the posts
   default: posts/

--output-folder
   where to save to output html files
   default: output/

--post-template
   location of the template file for the post page
   default: templates/post_template.html

--index-template
   location of the template file for the homepage
   default: templates/index_template.html
```

Example usage:
`./ssg --input-folder="md-files/" --output-folder="mysite/"`

## How to write a blog post
The first large heading (`# title`) in each post is used as the title that shows up in the homepage.
In every markdown file, the program reads the first line to get the pubblication date. the line must be formatted like this:
```
[date: 2020-08-23]
```
(It's yyyy-mm-dd)
If this line does not appear, no big deal, the markdown parser does its job as if nothing had happened.
For now, the parser does not implement the complete markdown syntax, missing things such as `====` and `-----` for headings, numbering images and ignoring markdown syntax in inline code snippets.

## Example post file
```
[date: 2020-08-09]
# Title of the post
this is a paragraph, **bold text**, *italic text*
* list element
* other list element
   * what is this? owo

[link to a cool song](https://www.youtube.com/watch?v=LykGbQy4wQA)

## subheading
> quote from famous guy
```

## Changing the templates
The templates files are written in the [Go template syntax](https://golang.org/pkg/text/template/#pkg-overview)
One of the plans for this project was to create my own templating engine but I ended up not having the time.