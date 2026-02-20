package properties

import (
	"os"
	"time"

	"github.com/gustavo000/goLibGustavo/models"
	"github.com/gustavo000/goLibGustavo/models/rest"
)

type InternalProperties struct {
	Server struct {
		Port           string
		BasePath       string
		Environment    string
		Release        string
		Namespace      string
		IsLocal        bool
		StartUpTime    time.Time
		LibraryVersion string
	}
	Name           string
	Version        string
	AuthMiddleware bool
}

type ExternalProperties struct {
	Services []*rest.Service
}

type Properties struct {
	Internal   InternalProperties
	External   ExternalProperties
	Secrets    []*models.Secret
	CustomTags map[string]any
}

var globalProps *Properties

type PropsOption func(*Properties)

func NewProperties(opts ...PropsOption) *Properties {
	globalProps = &Properties{}
	globalProps.Internal.Server.StartUpTime = time.Now()
	for _, fn := range opts {
		fn(globalProps)
	}
	return globalProps
}

func GetProperty() *Properties {
	return globalProps
}

func WithPort(v string) PropsOption {
	return func(properties *Properties) {
		properties.Internal.Server.Port = v
	}
}

func WithBasePath(v string) PropsOption {
	return func(properties *Properties) {
		properties.Internal.Server.BasePath = v
	}
}

func WithServiceName(v string) PropsOption {
	return func(properties *Properties) {
		properties.Internal.Name = v
	}
}

func WithVersion(v string) PropsOption {
	return func(properties *Properties) {
		properties.Internal.Version = v
	}
}

func WithServices(services ...*rest.Service) PropsOption {
	return func(properties *Properties) {
		properties.External.Services = services
	}
}

func WithEnvironment(defaultValue string) PropsOption {
	return func(properties *Properties) {
		properties.Internal.Server.Environment = defaultValue
		properties.Internal.Server.IsLocal = true
		if value, ok := os.LookupEnv("ENVIRONMENT"); ok {
			properties.Internal.Server.Environment = value
			properties.Internal.Server.IsLocal = !ok
		}
	}
}

func WithRelease(defaultValue string) PropsOption {
	return func(properties *Properties) {
		properties.Internal.Server.Release = defaultValue
		if value, ok := os.LookupEnv("RELEASE"); ok {
			properties.Internal.Server.Release = value
		}
	}
}

func WithNameSpace(defaultValue string) PropsOption {
	return func(properties *Properties) {
		properties.Internal.Server.Namespace = defaultValue
		if value, ok := os.LookupEnv("NAMESPACE"); ok {
			properties.Internal.Server.Namespace = value
		}
	}
}

func WithCustomTag(tag string, value any) PropsOption {
	return func(properties *Properties) {
		if properties.CustomTags == nil {
			properties.CustomTags = make(map[string]any)
		}
		properties.CustomTags[tag] = value
	}
}

func (p Properties) IsLocal() bool {
	return p.Internal.Server.IsLocal
}

func (p Properties) GetPort() string {
	return p.Internal.Server.Port
}

func (p Properties) GetName() string {
	return p.Internal.Name
}

func (p Properties) GetNamespace() string {
	return p.Internal.Server.Namespace
}

func (p Properties) GetRelease() string {
	return p.Internal.Server.Release
}

func (p Properties) GetEnv() string {
	return p.Internal.Server.Environment
}

func (p Properties) GetBasePath() string {
	return p.Internal.Server.BasePath
}

func (p Properties) GetVersion() string {
	return p.Internal.Version
}

func (p Properties) GetCustomTag(tag string) any {
	return p.CustomTags[tag]
}

func (p Properties) GetCustomTagString(tag string) string {
	defer panicHandler()
	return p.CustomTags[tag].(string)
}

func (p Properties) GetStartUpTime() time.Time {
	return p.Internal.Server.StartUpTime
}

func (p Properties) GetLibraryVersion() string {
	return p.Internal.Server.LibraryVersion
}

func panicHandler() {
	recover()
}
