package extract

import "testing"

func TestExtractJsonRecord(t *testing.T) {
	var test_val string = `{"foo":"bar"}`
	ret_val, err := ExtractJsonRecord(test_val)
	if err != nil {
		t.Error("Err not nill")
	}
	if ret_val["foo"] != "bar" {
		t.Error("foo not equal bar")
	}
}
