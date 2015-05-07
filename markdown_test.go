package ant

import(
    "testing"
)

func TestMarkdown(t *testing.T){
    m := []byte(`
* aaaa
* bbbb
* ssss
`)
    h := Markdown(m)
    t.Log(h)
}
