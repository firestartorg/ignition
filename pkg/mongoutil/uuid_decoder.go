package mongoutil

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"reflect"
)

type UUIDDecoder struct{}

func (U UUIDDecoder) DecodeValue(context bsoncodec.DecodeContext, reader bsonrw.ValueReader, value reflect.Value) (err error) {
	if !value.CanSet() || value.Type() != uuidType {
		return bsoncodec.ValueDecoderError{Name: "uuidDecodeValue", Types: []reflect.Type{uuidType}, Received: value}
	}

	var data []byte
	var subtype byte
	switch vrType := reader.Type(); vrType {
	case bson.TypeBinary:
		data, subtype, err = reader.ReadBinary()
		if subtype != uuidBinaryType {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}
	case bson.TypeNull:
		err = reader.ReadNull()
	case bson.TypeUndefined:
		err = reader.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", vrType)
	}
	if err != nil {
		return
	}

	var uuidVal uuid.UUID
	uuidVal, err = uuid.FromBytes(data)
	if err != nil {
		return
	}
	value.Set(reflect.ValueOf(uuidVal))

	return
}

var _ bsoncodec.ValueDecoder = UUIDDecoder{}
