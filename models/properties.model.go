package models

import "time"

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
	Name    string
	Version string
	Agent   struct {
		Tempo string
	}
}

type ExternalProperties struct {
	Services []*Service
}

type Properties struct {
	Internal   InternalProperties
	External   ExternalProperties
	Secrets    []*Secret
	CustomTags map[string]any
}
