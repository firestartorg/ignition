package mongoutil

import (
	"github.com/google/uuid"
	"reflect"
)

var uuidType = reflect.TypeOf(uuid.UUID{})
var uuidBinaryType byte = 4
