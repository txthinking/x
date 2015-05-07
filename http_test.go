package ant

import(
    "testing"
)

func TestMultipartFormDataFromFile(t *testing.T){
    r, err := MultipartFormDataFromFile(
        map[string][]string{
            "a": []string{"哈啊好"},
        },
        map[string][]string{
            "f": []string{"/etc/hosts"},
        },
        "aasfsafsfsfsfsdf",
    )
    if err != nil{
        t.Fatal(r, err)
    }
    t.Log(r)
}

