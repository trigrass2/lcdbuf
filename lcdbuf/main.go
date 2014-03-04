
package main

import (
	"gitcafe.com/nuomi-studio/lcdbuf.git"
	"flag"
	"log"
)

func main() {
	file := flag.String("f", "", "from image file")
	cfile := flag.String("to-c", "", "to C code file")
	gofile := flag.String("to-go", "", "to Go code file")
	flag.Parse()

	buf, err := lcdbuf.FromImageFile(*file)
	if err != nil {
		log.Println(err)
		return
	}

	if *cfile != "" {
		buf.DumpCCode(*cfile)
	}
	if *gofile != "" {
		buf.DumpGoCode(*gofile)
	}
}

