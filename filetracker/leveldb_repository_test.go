package filetracker_test

import (
	filetracker2 "github.com/gphotosuploader/gphotos-uploader-cli/filetracker"
	"testing"

	"github.com/syndtr/goleveldb/leveldb/opt"
)

func TestLevelDBRepository_Get(t *testing.T) {
	testCases := []struct {
		name  string
		input string
		found bool
	}{
		{"Should success", ShouldSuccess, true},
		{"Should fail", ShouldMakeRepoFail, false},
	}

	repo := filetracker2.LevelDBRepository{
		DB: &mockedDB{},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, found := repo.Get(tc.input)
			if tc.found != found {
				t.Errorf("want: %t, got: %t", tc.found, found)
			}
		})
	}
}

func TestLevelDBRepository_Put(t *testing.T) {
	repo := filetracker2.LevelDBRepository{
		DB: mockedDB{},
	}
	if err := repo.Put("foo", filetracker2.TrackedFile{}); err != nil {
		t.Errorf("error was not expected, err: %s", err)
	}
}

func TestLevelDBRepository_Delete(t *testing.T) {
	repo := filetracker2.LevelDBRepository{
		DB: mockedDB{},
	}
	if err := repo.Delete("foo"); err != nil {
		t.Errorf("error was not expected, err: %s", err)
	}
}

func TestLevelDBRepository_Close(t *testing.T) {
	repo := filetracker2.LevelDBRepository{
		DB: mockedDB{},
	}
	if err := repo.Close(); err != nil {
		t.Errorf("error was not expected, err: %s", err)
	}
}

type mockedDB struct{}

func (m mockedDB) Get(key []byte, ro *opt.ReadOptions) ([]byte, error) {
	var a []byte
	if string(key) == ShouldMakeRepoFail {
		return a, ErrTestError
	}
	return a, nil
}

func (m mockedDB) Put(key []byte, item []byte, wo *opt.WriteOptions) error {
	return nil
}

func (m mockedDB) Delete(key []byte, wo *opt.WriteOptions) error {
	return nil
}

func (m mockedDB) Close() error {
	return nil
}
