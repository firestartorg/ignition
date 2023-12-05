package mongoutil

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
)

var BsonRegistry = newRegistry()

func newRegistry() (registry *bsoncodec.Registry) {
	registry = bson.NewRegistry()
	registry.RegisterTypeEncoder(uuidType, UUIDEncoder{})
	registry.RegisterTypeDecoder(uuidType, UUIDDecoder{})

	return
}
