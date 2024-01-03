# Config

Config is a package that provides a way to load layered configurations from multiple sources.

## Usage

The config package provides a way to load configuration from multiple sources. The configuration is loaded in layers,
with each layer overriding the previous layer. The layers are loaded in the following order:

1. Environment variables starting with `IGNITION_`
2. Configuration file (environment specific if available `GO_ENV`)
3. Configuration file
4. Default configuration

*Each layer can be overridden by the next layer. For example, if a configuration key is set in the environment variables,*
*it will override the same key in the configuration file.*

The configuration file is loaded from the following locations:

1. `./appsettings.yaml`
2. `./appsettings.yml`
3. `./appsettings.json`

from the current working directory of the application (can be overridden using the `IGNITION_CONFIG_PATH` environment variable).

### Example

```go
package main

import (
  "gitlab.com/firestart/ignition/pkg/config"
)

type Config struct {
  Name string
  TestSsl NestedConfig // Nested configuration is supported
}

type NestedConfig struct {
  Value string
}

func main() {
  configuration, _ := config.LoadConfig() // Load dynamic configuration

  // Staticly override configuration
  configuration.Set("Name:Test", "Test") // Sub configuration can be set using a colon

  var cfg Config
  _ = configuration.Unpack(&cfg) // Unpack the configuration into a struct
}
```

### Naming conventions

The configuration package uses multiple naming conventions for configuration keys. The following naming conventions are used:

- `snake_case` - Used for yaml configuration files
- `PascalCase` - Used for json configuration files and struct fields
- `PascalCase` - Used for keys in the configuration with an additional colon (`:`) as a separator for nested configurations (e.g. `Name:Test`)
- `SCREAMING_SNAKE_CASE` - Used for environment variables with a prefix of `IGNITION_` and a double underscore as a separator (`__`) for nested configurations (e.g. `IGNITION_NAME__TEST`)

To illustrate these naming conventions, the following configuration would be valid for the [example above](#example):

```yaml
name: Test
test_ssl:
  value: test
```

```json
{
  "Name": "Test",
  "TestSsl": {
    "Value": "test"
  }
}
```

```go
configuration.Set("TestSsl:Value", "test")
```

```bash
IGNITION_NAME=Test
IGNITION_TEST_SSL__VALUE=test
```

The environment variable naming convention follow the same rules as the [ASP.Net Core configuration](https://learn.microsoft.com/en-us/aspnet/core/fundamentals/configuration/?view=aspnetcore-8.0#naming-of-environment-variables).