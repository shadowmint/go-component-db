package db

import (
	"ntoolkit/component"

	"github.com/jinzhu/gorm"
	"ntoolkit/errors"
)

type ObjectTemplate struct {
	gorm.Model
	Key  string
	Name string
	Data []byte
}

// ObjectStorageDb stores objects in a database as json chunks
type ObjectStorageDb struct {
	db *gorm.DB
}

// NewObjectStorageMemory returns a new instance that caches templates in a simple local hash.
func NewObjectStorageDb(db *gorm.DB) *ObjectStorageDb {
	db.AutoMigrate(&ObjectTemplate{})
	rtn := &ObjectStorageDb{db: db}
	return rtn
}

func (s *ObjectStorageDb) Set(obj *component.Object, factory *component.ObjectFactory) error {
	template, err := factory.Serialize(obj)
	if err != nil {
		return err
	}

	raw, err := component.ObjectTemplateAsJson(template)
	if err != nil {
		return err
	}

	// Update existing records, save new ones
	record := ObjectTemplate{}
	if s.db.Where("name = ?", obj.Name()).First(&record).RecordNotFound() {
		record.Name = obj.Name()
		record.Key = obj.ID()
		record.Data = raw
		s.db.Create(&record)
	} else {
		record.Data = raw
		record.Key = obj.ID()
		s.db.Save(record)
	}
	return nil
}

func (s *ObjectStorageDb) Clear(id string) error {
	return s.db.Unscoped().Where("name = ?", id).Delete(&ObjectTemplate{}).Error
}

func (s *ObjectStorageDb) Get(name string, factory *component.ObjectFactory) (*component.Object, error) {
	record := ObjectTemplate{}
	if s.db.Where("name = ?", name).First(&record).RecordNotFound() {
		return nil, errors.Fail(component.ErrNoMatch{}, nil, "No match for record")
	}

	template, err := component.ObjectTemplateFromJson(string(record.Data))
	if err != nil {
		return nil, err
	}

	object, err := factory.Deserialize(template)
	if err != nil {
		return nil, err
	}

	return object, nil
}

func (s *ObjectStorageDb) Has(name string) bool {
	count := 0
	if err := s.db.Where("name = ?", name).Find(&ObjectTemplate{}).Count(&count).Error; err != nil {
		return false
	}
	return count > 0
}

func (s *ObjectStorageDb) Getter() component.ObjectStorageGetter {
	return s
}

func (s *ObjectStorageDb) Setter() component.ObjectStorageSetter {
	return s
}
