package mongo

import (
	"context"
	"errors"
	"github.com/rs/zerolog/log"
	"gitlab.com/firestart/ignition/pkg/injector"
	"gitlab.com/firestart/ignition/pkg/mongoutil"
	"gitlab.com/firestart/ignition/x/application"
	"gitlab.com/firestart/ignition/x/application/extensions/monitor"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Name is the name of the injectable
const Name = "ignition/db-mongo"

type client struct {
	// client is the mongo client to use
	client *mongo.Client
	// database is the default database client
	database *mongo.Database

	// opt
	opt settings
}

// WithMongoClient adds a mongo client to the application.
func WithMongoClient(opts ...Option) application.Option {
	return func(app application.App) {
		// Create the settings
		s := newSettings(opts...)
		// Add provide the client
		injector.ProvideNamed(app.Injector, Name, provideFactory(s))

		app.AddShutdownHook(func(ctx context.Context, app application.App) error {
			cl, err := injector.GetNamed[client](app.Injector, Name)
			if err != nil {
				return err
			}

			err = cl.client.Disconnect(ctx)
			if err != nil {
				return err
			}

			return nil
		})

		app.AddHook(monitor.HookReadiness, func(ctx context.Context, app application.App) error {
			cl, err := injector.GetNamed[client](app.Injector, Name)
			if err != nil {
				return err
			}

			err = cl.client.Ping(ctx, nil)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("mongo health check failed")
				return err
			}

			return nil
		})
	}
}

func provideFactory(s settings) func(inj *injector.Injector) (client, error) {
	return func(inj *injector.Injector) (client, error) {
		// If the Config option is set, load the Config
		if s.Config {
			conf, err := injector.ExtractConfig[Config](inj, "Mongo")
			if err != nil {
				panic(err)
			}

			s.ClientOptions = options.Client().ApplyURI(conf.ConnectionString).SetRegistry(mongoutil.BsonRegistry)
			s.Database = conf.Database
			s.Collections = conf.Collection
		}

		cl, err := mongo.Connect(context.Background(), s.ClientOptions)
		if err != nil {
			return client{}, err
		}

		var db *mongo.Database
		if s.Database != "" {
			db = cl.Database(s.Database)
		}

		return client{cl, db, s}, nil
	}
}

// WithDefaultMongoClient adds a mongo client to the application.
func WithDefaultMongoClient() application.Option {
	return WithMongoClient(FromSettings())
}

// NamedCollection gets a Collections collection client
func NamedCollection(inj *injector.Injector, name string, opts ...*options.CollectionOptions) (*mongo.Collection, error) {
	cl, err := injector.GetNamed[client](inj, Name)
	if err != nil {
		return nil, err
	}

	collection, ok := cl.opt.Collections[name]
	if !ok {
		return nil, errors.New("collection name not found")
	}

	return cl.database.Collection(collection, opts...), nil
}

// Collection gets a collection client.
// It's recommended to use NamedCollection instead.
func Collection(inj *injector.Injector, collection string, opts ...*options.CollectionOptions) (*mongo.Collection, error) {
	cl, err := injector.GetNamed[client](inj, Name)
	if err != nil {
		return nil, err
	}

	return cl.database.Collection(collection, opts...), nil
}

// Database gets a Database client.
// It's not recommended to use this function.
func Database(inj *injector.Injector, name string, opts ...*options.DatabaseOptions) (*mongo.Database, error) {
	cl, err := injector.GetNamed[client](inj, Name)
	if err != nil {
		return nil, err
	}

	return cl.client.Database(name, opts...), nil
}
