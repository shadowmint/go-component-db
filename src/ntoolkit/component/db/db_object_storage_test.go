package db_test

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"ntoolkit/assert"
	"ntoolkit/component"
	"ntoolkit/component/db"
)

func textFixture(orm *gorm.DB, T *assert.T) *component.ObjectStorage {
	dbstorage := db.NewObjectStorageDb(orm)
	root := component.NewObject("Root")

	factory := component.NewObjectFactory()
	factory.Register(&StatefulComponent{})
	factory.Register(&StatelessComponent{})

	active := component.NewObjectStorageRuntime(root)
	templates := component.NewObjectStorageMemory()
	templates.CanSet = false

	stack := component.NewObjectStorageStack()
	stack.Add(dbstorage)
	stack.Add(templates)

	storage := component.NewObjectStorage(factory, active, stack)

	object1 := component.NewObject("Hello")
	object1.AddComponent(&StatefulComponent{})
	object1.AddComponent(&StatelessComponent{})

	object2 := component.NewObject("World")
	object2.AddComponent(&StatefulComponent{})
	object2.AddComponent(&StatelessComponent{})

	T.Assert(templates.Set(object1, factory) == nil)
	T.Assert(storage.Add(object2) == nil)

	return storage
}

func TestGetComponents(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		orm, err := gorm.Open("sqlite3", "test.db")
		if err != nil {
			panic("failed to connect database")
		}

		defer orm.Close()

		storage := textFixture(orm, T)

		obj1, err := storage.Get("Hello")

		T.Assert(err == nil)
		T.Assert(obj1 != nil)

		var ref *StatefulComponent
		err = obj1.Find(&ref)
		T.Assert(err == nil)
		ref.state.Name = "Dog"
		ref.state.Value = "Cat"

		ref = nil
		err = obj1.Find(&ref)
		T.Assert(err == nil)
		T.Assert(ref.state.Name == "Dog")
		T.Assert(ref.state.Value == "Cat")

		obj2, err := storage.Get("World")
		T.Assert(err == nil)
		T.Assert(obj2 != nil)

		T.Assert(storage.SetActive("Hello", false) == nil)
		T.Assert(storage.SetActive("World", false) == nil)

		obj1, err = storage.Get("Hello")
		T.Assert(err == nil)
		T.Assert(obj1 != nil)

		ref = nil
		err = obj1.Find(&ref)
		T.Assert(err == nil)
		T.Assert(ref != nil)
		T.Assert(ref.state.Name == "Dog")
		T.Assert(ref.state.Value == "Cat")

		obj2, err = storage.Get("World")
		T.Assert(err == nil)
		T.Assert(obj2 != nil)

		storage.Drop("World")
		obj2, err = storage.Get("World")
		T.Assert(err != nil)
		T.Assert(obj2 == nil)
	})
}
