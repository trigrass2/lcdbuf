
package main

import (
	"gitcafe.com/nuomi-studio/lcdbuf.git"
	"flag"
	"log"
)

func main() {
	file := flag.String("f", "", "from image file")
	tofile := flag.String("t", "", "to code file")
	flag.Parse()

	buf, err := lcdbuf.FromImageFile(*file)
	if err != nil {
		log.Println(err)
		return
	}

	buf.DumpCCode(*tofile)
}

