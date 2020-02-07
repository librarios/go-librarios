package main

const (
	PluginTypeBook = "book"
)

type Plugin interface {
	Type() string
	Name() string
	SetProperties(map[string]interface{})
}

type NewPluginFunc func() Plugin

type PluginDef struct {
	Type    string
	Name    string
	NewFunc NewPluginFunc
}
