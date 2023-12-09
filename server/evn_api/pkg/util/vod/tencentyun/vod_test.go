package tencentyun

import "testing"

func TestVod(t *testing.T) {
	res := vodUpload()
	if res == nil {
		t.Logf("Success\n%v\n", res)
	} else {
		t.Fatalf("Failed\n%v\n", res)
	}
}
