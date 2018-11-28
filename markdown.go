package x

import "gopkg.in/russross/blackfriday.v2"

//md "github.com/shurcooL/go/github_flavored_markdown"

// Markdown format input to html
func Markdown(markdown []byte) (html []byte) {
	html = blackfriday.Run(markdown, blackfriday.WithExtensions(blackfriday.CommonExtensions))
	return
}
