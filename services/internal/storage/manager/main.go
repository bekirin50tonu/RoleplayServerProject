package manager

import (
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"mime/multipart"
	"os"
	"path"
	"services/internal/storage/models"
	"services/internal/storage/repositories"
	"services/pkg/database"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StorageItem struct {
	FileName  string
	Disk      string
	Size      int64
	CreatedAt time.Time
	Extension string
	File      *os.File
}

type StorageDisk struct {
	name         string
	rootLocation string
	Location     string
	repository   *repositories.StorageRepository
}

type StorageManagerConfig struct {
	DbConfig     database.DatabaseConfig
	RootLocation string
	Disks        map[string]*StorageDisk
}

type StorageManager struct {
	config StorageManagerConfig
}

func NewStorageDiskSystem(config StorageManagerConfig) (*StorageManager, error) {
	// TODO: Initialize Folder Process.
	db, err := database.New(config.DbConfig)
	if err != nil {
		return nil, err
	}

	err = os.Mkdir(config.RootLocation, os.ModeDir)
	if err != nil {
		if i := os.IsExist(err); !i {
			return nil, err
		}
	}
	for k, v := range config.Disks {
		rep, err := repositories.NewStorageRepository(db.Database, k)
		if err != nil {
			return nil, err
		}
		config.Disks[k].repository = rep
		config.Disks[k].name = k
		config.Disks[k].rootLocation = config.RootLocation
		err = os.Mkdir(path.Join(config.RootLocation, v.Location), os.ModeAppend)
		if err != nil {
			if i := os.IsExist(err); !i {
				return nil, err
			}
		}
	}

	//

	return &StorageManager{
		config: config,
	}, nil
}

func (s *StorageManager) Get(name string) (*StorageDisk, error) {
	disk, ok := s.config.Disks[name]
	if !ok {
		return nil, errors.New("Storage Disk Didn't Match Given Name Parameter. Get Problem.")
	}
	return disk, nil

}

func (s *StorageManager) Delete(name string) error {
	_, ok := s.config.Disks[name]
	if !ok {
		return errors.New("Storage Disk Didn't Match Given Name Parameter.Delete Problem.")
	}
	delete(s.config.Disks, name)
	return nil
}

func (s *StorageDisk) Put(filename string, file any) (*models.Storage, error) {
	var res StorageItem
	res.Extension = path.Ext(filename)
	res.Disk = s.name

	alg := fnv.New32()
	_, err := alg.Write([]byte(filename + time.Now().GoString()))
	if err != nil {
		return nil, err
	}
	res.FileName = fmt.Sprintf("%v", alg.Sum32())

	filePath := path.Join(s.rootLocation, s.Location)
	fileDest := path.Join(filePath, res.FileName)

	fb, ok := file.([]byte)
	if ok {
		err = os.WriteFile(fileDest, fb, os.ModeAppend)
		if err != nil {
			return nil, err
		}

		f, err := os.Open(fileDest)
		if err != nil {
			return nil, err
		}
		res.File = f
	}

	fh, ok := file.(*multipart.FileHeader)
	if ok {
		srcFile, err := fh.Open()
		if err != nil {
			return nil, err
		}

		defer srcFile.Close()

		srcFile.Seek(0, io.SeekStart)

		dstFile, err := os.Create(fileDest)
		if err != nil {
			return nil, err
		}

		size, err := io.Copy(dstFile, srcFile)

		res.File = dstFile
		res.Size = size
	}
	f, err := os.OpenFile(fileDest, os.Getegid(), os.ModeAppend)
	if err != nil {
		return nil, err
	}

	finfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	res.CreatedAt = finfo.ModTime()
	res.Size = finfo.Size()

	storage := models.NewStorage(res.FileName, res.Disk, res.Extension, res.Size)
	img, err := s.repository.Create(storage)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func (s *StorageDisk) Get(id string) (*StorageItem, error) {

	query, err := s.repository.FindOneWithParameters(primitive.M{
		"filename": id,
	})
	if err != nil {
		return nil, err
	}
	filePath := path.Join(s.rootLocation, s.Location, query.Filename)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	var res StorageItem

	res.File = file
	res.FileName = query.Filename
	res.Disk = query.Disk
	res.CreatedAt = query.CreatedAt
	res.Extension = query.Extension

	return &res, nil
}
