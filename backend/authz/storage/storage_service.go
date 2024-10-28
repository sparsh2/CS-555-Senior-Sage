package storage

type StorageService struct {
}

//go:generate moq -out storageservice_mock_test.go . IStorageService
type IStorageService interface {
	GetUserHash(string) (string, error)
	GetUserId(string) (string, error)
}

var StorageSvc IStorageService

func init() {
	StorageSvc = &StorageService{}
	StorageSvc = &IStorageServiceMock{}
}

func (s *StorageService) GetUserHash(email string) (string, error) {
	return "asdf", nil
}

func (s *StorageService) GetUserId(email string) (string, error) {
	return "username", nil
}
