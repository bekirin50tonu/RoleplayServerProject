package services

import (
	"mime/multipart"
	"services/internal/storage/manager"
	"services/internal/storage/models"
)

type StorageService struct {
	storageManager *manager.StorageManager
}

func NewStorageService(manager *manager.StorageManager) (*StorageService, error) {

	return &StorageService{
		storageManager: manager,
	}, nil
}

func (s *StorageService) SaveImageToLocalStorage(disk, filename string, data []byte) (*models.Storage, error) {
	manager, err := s.storageManager.Get(disk)
	if err != nil {
		return nil, err
	}

	image, err := manager.Put(filename, data)
	if err != nil {
		return nil, err
	}

	return image, nil
}

func (s *StorageService) SaveImageToLocalStorageRest(diskName string, data *multipart.FileHeader) (*models.Storage, error) {
	filename := data.Filename

	disk, err := s.storageManager.Get(diskName)
	if err != nil {
		return nil, err
	}
	model, err := disk.Put(filename, data)

	return model, nil
}

func (s *StorageService) GetImage(diskName, id string) (*manager.StorageItem, error) {
	disk, err := s.storageManager.Get(diskName)
	if err != nil {
		return nil, err
	}
	file, err := disk.Get(id)
	if err != nil {
		return nil, err
	}

	return file, nil
}
