package completeduploads

import "testing"

func TestHash(t *testing.T) {
	filePath := "testdata/image.png"
	expected := uint32(1210845310)

	got, err := Hash(filePath)
	if err != nil {
		t.Errorf("error not expected at this stage: %v", err)
	}

	if got != expected {
		t.Errorf("got %d, want %d", got, expected)
	}
}

func TestCompletedUploadedFileItem_GetTrackedHash(t *testing.T) {
	filePath := "testdata/image.png"
	expected := "1210845310"

	item, err := NewCompletedUploadedFileItem(filePath)
	if err != nil {
		t.Errorf("error not expected at this stage: %v", err)
	}

	got := item.GetTrackedHash()
	if got != expected {
		t.Errorf("got %s, want %s", got, expected)
	}
}
