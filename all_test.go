package ini4go

import (
    "testing"
    "fmt"
)

func TestRunes(t *testing.T) {

    sess, _ := NewSection(&Section{name:"haha"})


    fmt.Println(*sess)
}

