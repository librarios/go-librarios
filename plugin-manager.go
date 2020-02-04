package main

type PluginManager interface {
	Register(def PluginDef)
	GetPluginByName(name string) *Plugin
	GetPluginsByType(pluginType string) []Plugin
}

var (
	pluginManager = DefaultPluginManager{
		plugins: make(map[string]Plugin),
	}
)

type DefaultPluginManager struct {
	plugins map[string]Plugin
}

// Register plugin
func (m DefaultPluginManager) Register(def PluginDef) {
	if def.NewFunc != nil {
		plugin := def.NewFunc()
		name := plugin.Name()
		m.plugins[name] = plugin
	}
}

func (m DefaultPluginManager) GetPluginByName(name string) (Plugin, bool) {
	plugin, ok := m.plugins[name]
	return plugin, ok
}

// GetPluginsByType returns plugins with matching pluginType
func (m DefaultPluginManager) GetPluginsByType(pluginType string) []Plugin {
	result := make([]Plugin, 0)
	for _, plugin := range m.plugins {
		if pluginType == plugin.Type() {
			result = append(result, plugin)
		}
	}
	return result
}

