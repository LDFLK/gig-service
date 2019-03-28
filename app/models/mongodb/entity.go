package mongodb

import (
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Entity struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Title     string        `json:"title" bson:"title"`
	Content   string        `json:"content" bson:"content"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" bson:"updated_at"`
}

func newEntityCollection() *Collection {
	return NewCollectionSession("entities")
}

/**
AddEntity insert a new Entity into database and returns
last inserted entity on success.
 */
func AddEntity(m Entity) (entity Entity, err error) {
	c := newEntityCollection()
	defer c.Close()
	m.ID = bson.NewObjectId()
	m.CreatedAt = time.Now()
	return m, c.Session.Insert(m)
}

/**
UpdateEntity update a Entity into database and returns
last nil on success.
 */
func (m Entity) UpdateEntity() error {
	c := newEntityCollection()
	defer c.Close()

	err := c.Session.Update(bson.M{
		"_id": m.ID,
	}, bson.M{
		"$set": bson.M{
			"title": m.Title, "content": m.Content, "updatedAt": time.Now()},
	})
	return err
}

/**
DeleteEntity Delete Entity from database and returns
last nil on success.
 */
func (m Entity) DeleteEntity() error {
	c := newEntityCollection()
	defer c.Close()

	err := c.Session.Remove(bson.M{"_id": m.ID})
	return err
}

/**
GetEntities Get all Entities from database and returns
list of Entity on success
 */
func GetEntities(search string) ([]Entity, error) {
	var (
		entities []Entity
		err      error
	)

	c := newEntityCollection()
	defer c.Close()

	err = c.Session.Find(bson.M{"title": bson.RegEx{".*" + search + "*.", ""}}).All(&entities)
	return entities, err
}

/**
GetEntity Get a Entity from database and returns
a Entity on success
 */
func GetEntity(id bson.ObjectId) (Entity, error) {
	var (
		entity Entity
		err    error
	)

	c := newEntityCollection()
	defer c.Close()

	err = c.Session.Find(bson.M{"_id": id}).One(&entity)
	return entity, err
}