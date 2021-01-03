package completeduploads

import "testing"

func TestService_IsAlreadyUploaded(t *testing.T) {
	// In memory repository is used for testing. DON'T USE IT IN PRODUCTION
	repo := NewInMemRepository()
	s := NewService(repo)

	t.Run("not uploaded file", func(t *testing.T) {
		filePath := "testdata/image.png"
		expected := false
		got, err := s.IsAlreadyUploaded(filePath)
		if err != nil {
			t.Errorf("error not expected at this stage: %v", err)
		}
		if got != expected {
			t.Errorf("got %t, want %t", got, expected)
		}
	})

	t.Run("already uploaded file", func(t *testing.T) {
		filePath := "testdata/image.png"
		expected := true
		item, err := NewCompletedUploadedFileItem(filePath)
		if err != nil {
			t.Errorf("error not expected at this stage: %v", err)
		}
		// add this item to the repository
		err = s.repo.Put(item)
		if err != nil {
			t.Errorf("error not expected at this stage: %v", err)
		}

		got, err := s.IsAlreadyUploaded(filePath)
		if err != nil {
			t.Errorf("error not expected at this stage: %v", err)
		}
		if got != expected {
			t.Errorf("got %t, want %t", got, expected)
		}
	})
}

func TestService_CacheAsAlreadyUploaded(t *testing.T) {
	repo := NewInMemRepository()
	s := NewService(repo)

	t.Run("existing file", func(t *testing.T) {
		filePath := "testdata/image.png"
		err := s.CacheAsAlreadyUploaded(filePath)
		if err != nil {
			t.Errorf("error not expected at this stage: %v", err)
		}
	})

	t.Run("not existing file", func(t *testing.T) {
		filePath := "testdata/image1.png"
		err := s.CacheAsAlreadyUploaded(filePath)
		if err == nil {
			t.Errorf("error was expected at this stage: %v", err)
		}
	})
}

func TestService_RemoveAsAlreadyUploaded(t *testing.T) {
	repo := NewInMemRepository()
	s := NewService(repo)

	t.Run("already uploaded file", func(t *testing.T) {
		filePath := "testdata/image.png"
		item, err := NewCompletedUploadedFileItem(filePath)
		if err != nil {
			t.Errorf("error not expected at this stage: %v", err)
		}
		// add this item to the repository
		err = s.repo.Put(item)
		if err != nil {
			t.Errorf("error not expected at this stage: %v", err)
		}

		err = s.RemoveAsAlreadyUploaded(filePath)
		if err != nil {
			t.Errorf("error not expected at this stage: %v", err)
		}

	})

	t.Run("not uploaded file", func(t *testing.T) {
		filePath := "testdata/image.png"

		err := s.RemoveAsAlreadyUploaded(filePath)
		if err == nil {
			t.Errorf("error was expected at this stage: %v", err)
		}

	})
}
