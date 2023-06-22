package filetracker_test

import (
	"errors"
	filetracker2 "github.com/gphotosuploader/gphotos-uploader-cli/filetracker"
	"testing"
)

const (
	ShouldSuccess      = "testdata/image.jpg"
	ShouldMakeHashFail = "testdata/non-existent"
	ShouldMakeRepoFail = "should-make-repo-fail"
)

var (
	// ErrTestError denotes an error raised by the test.
	ErrTestError = errors.New("error")
)

func TestFileTracker_MarkAsUploaded(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should success", ShouldSuccess, false},
		{"Should fail if repo fails", ShouldMakeRepoFail, true},
		{"Should fail if Hasher fails", ShouldMakeHashFail, true},
	}

	ft := filetracker2.New(&mockedRepository{})
	ft.Hasher = &mockedHasher{"test-file-hash"}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ft.MarkAsUploaded(tc.input)
			assertExpectedError(t, tc.isErrExpected, err)
		})
	}
}

func TestFileTracker_IsUploaded(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		want  bool
	}{
		{"Should return true if file is in the repo", ShouldSuccess, true},
		{"Should return false if file is not in the repo", ShouldMakeRepoFail, false},
		{"Should return false if Hasher fails", ShouldMakeHashFail, false},
	}

	ft := filetracker2.New(&mockedRepository{
		valueInRepo: filetracker2.NewTrackedFile("test-file-hash"),
	})
	ft.Hasher = &mockedHasher{"test-file-hash"}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := ft.IsUploaded(tc.input)
			if tc.want != got {
				t.Errorf("want: %t, got: %t", tc.want, got)
			}
		})
	}
}

func TestFileTracker_UnmarkAsUploaded(t *testing.T) {
	testCases := []struct {
		name          string
		input         string
		isErrExpected bool
	}{
		{"Should success", ShouldSuccess, false},
		{"Should success", ShouldMakeRepoFail, true},
	}

	ft := filetracker2.New(&mockedRepository{})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ft.UnmarkAsUploaded(tc.input)
			assertExpectedError(t, tc.isErrExpected, err)
		})
	}
}

func TestFileTracker_Close(t *testing.T) {
	testCases := []struct {
		name          string
		scenario      string
		isErrExpected bool
	}{
		{"Should success", ShouldSuccess, false},
	}

	ft := filetracker2.New(&mockedRepository{})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := ft.Close()
			assertExpectedError(t, tc.isErrExpected, err)
		})
	}
}

type mockedRepository struct {
	valueInRepo filetracker2.TrackedFile
}

func (m mockedRepository) Get(key string) (filetracker2.TrackedFile, bool) {
	switch key {
	case ShouldMakeRepoFail:
		return filetracker2.TrackedFile{}, false
	default:
		return m.valueInRepo, true
	}
}

func (m mockedRepository) Put(key string, item filetracker2.TrackedFile) error {
	if key == ShouldMakeRepoFail {
		return ErrTestError
	}
	return nil
}

func (m mockedRepository) Delete(key string) error {
	if key == ShouldMakeRepoFail {
		return ErrTestError
	}
	return nil
}

func (m mockedRepository) Close() error {
	return nil
}

type mockedHasher struct {
	hash string
}

func (mh mockedHasher) Hash(filename string) (string, error) {
	switch filename {
	case ShouldMakeHashFail:
		return "", ErrTestError
	default:
		return mh.hash, nil
	}
}

func assertExpectedError(t *testing.T, errExpected bool, err error) {
	if errExpected && err == nil {
		t.Fatalf("error was expected, but not produced")
	}
	if !errExpected && err != nil {
		t.Fatalf("error was not expected, err: %s", err)
	}
}
