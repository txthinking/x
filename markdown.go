package ant

import (
	//"bytes"
	MD "github.com/russross/blackfriday"
	//md "github.com/shurcooL/go/github_flavored_markdown"
)

func Markdown(markdown []byte) (html []byte) {
	var renderer MD.Renderer
	var hf, ef int

	hf |= MD.HTML_SKIP_HTML
	hf |= MD.HTML_SKIP_STYLE
	//hf |= MD.HTML_SANITIZE_OUTPUT
	hf |= MD.HTML_SAFELINK
	hf |= MD.HTML_NOFOLLOW_LINKS
	hf |= MD.HTML_HREF_TARGET_BLANK
	//hf |= MD.HTML_GITHUB_BLOCKCODE

	hf |= MD.HTML_USE_SMARTYPANTS
	hf |= MD.HTML_SMARTYPANTS_FRACTIONS
	//hf |= MD.HTML_FOOTNOTE_RETURN_LINKS
	//hf |= MD.HTML_TOC
	//hf |= MD.HTML_OMIT_CONTENTS

	ef |= MD.EXTENSION_TABLES
	ef |= MD.EXTENSION_FENCED_CODE
	ef |= MD.EXTENSION_AUTOLINK
	ef |= MD.EXTENSION_STRIKETHROUGH
	ef |= MD.EXTENSION_HARD_LINE_BREAK
	ef |= MD.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK
	ef |= MD.EXTENSION_HEADER_IDS

	//ef |= MD.EXTENSION_TITLEBLOCK
	//ef |= MD.EXTENSION_FOOTNOTES
	//ef |= MD.EXTENSION_NO_INTRA_EMPHASIS
	//ef |= MD.EXTENSION_LAX_HTML_BLOCKS
	renderer = MD.HtmlRenderer(hf, "", "")
	html = MD.Markdown(markdown, renderer, ef)

	//html = bytes.Replace(html, []byte("<code>"), []byte("<code class=\"prettyprint\">"), -1);
	return
}
