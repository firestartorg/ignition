package mongo

// Config is the configuration for the mongo app.
type Config struct {
	// ConnectionString is the connection string.
	ConnectionString string
	// Database is the Database name.
	Database string
	// Collection is the collection name to collection mapping.
	Collection map[string]string
}
