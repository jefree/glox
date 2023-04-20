package main

import (
	"fmt"
	"os"
)

func main() {

	file, err := os.OpenFile("../delete_me.go", os.O_WRONLY|os.O_CREATE, 0660)

	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	fmt.Fprint(file, "package main\n")

}
