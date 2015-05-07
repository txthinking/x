package ant

import(
    "testing"
)

func TestIsEmail(t *testing.T){
    es := []string{
        "error@mail",
        "correct@mail.com",
        "_hi@mail.com",
        "aa_hi@qq.com",
        "a$_hi@qq.com",
        "!@@@1.com",
        "!~@@1.cm",
        "a@1.cm",
    }
    for _, v := range es{
        ok, err := IsEmail(v)
        if err != nil{
            t.Fatal(v, err)
        }
        t.Log(v, ok)
    }
}

func TestIsBankCard(t *testing.T){
    a := []int64{
        4512893900582108,
        6228480010323650910,
        6228480010323650919, // error
    }
    for _, v := range a{
        ok, err := IsBankCard(v)
        if err != nil{
            t.Fatal(v, err)
        }
        t.Log(v, ok)
    }
}

func TestIsChineseID(t *testing.T){
    a := []string{
        "61052819890402574X",
        "411081198804220861",
        "411081198804220851",
    }
    for _, v := range a{
        ok, err := IsChineseID(v)
        if err != nil{
            t.Fatal(v, err)
        }
        t.Log(v, ok)
    }
}

func TestIsChineseWords(t *testing.T){
    a := []string{
        "猪八戒",
        "xia往往",
    }
    for _, v := range a{
        ok, err := IsChineseWords(v)
        if err != nil{
            t.Fatal(v, err)
        }
        t.Log(v, ok)
    }
}
