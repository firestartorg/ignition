package mongoutil

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"reflect"
)

type UUIDEncoder struct{}

func (U UUIDEncoder) EncodeValue(context bsoncodec.EncodeContext, writer bsonrw.ValueWriter, value reflect.Value) error {
	if !value.IsValid() || value.Type() != uuidType {
		return bsoncodec.ValueEncoderError{Name: "uuidEncodeValue", Types: []reflect.Type{uuidType}, Received: value}
	}
	b := value.Interface().(uuid.UUID)
	return writer.WriteBinaryWithSubtype(b[:], uuidBinaryType)
}

var _ bsoncodec.ValueEncoder = UUIDEncoder{}
