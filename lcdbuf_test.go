
package lcdbuf

import (
	"log"
	"testing"

	"gitcafe.com/nuomi-studio/pcf.git"
)

func TestXxx(t *testing.T) {
	buf := New(128, 64)
	font, _ := pcf.Open("wenquanyi_13px.pcf")
	Draw(buf, PCFText(font, "测试哦"), 0, 0, nil, nil)
	Draw(buf, PCFText(font, "Hello"), 0, 16, nil, nil)
	buf.DumpAscii("out")
	log.Println("end")
}

