package storage

type StorageService struct {
}

var StorageSvc *StorageService

func init() {
	StorageSvc = &StorageService{}
}

func (s *StorageService) GetUserHash(email string) (string, error) {
	return  "asdf", nil
}

func (s *StorageService) GetUserId(email string) ( string, error) {
	return  "username", nil
}