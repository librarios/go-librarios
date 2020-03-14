package plugin

type PluginManager interface {
	Register(def PluginDef)
	GetPluginByName(name string) *Plugin
	GetPluginsByType(pluginType string) []Plugin
}

var (
	defaultPluginManager = DefaultPluginManager{
		plugins: make(map[string]Plugin),
	}
)

type DefaultPluginManager struct {
	plugins map[string]Plugin
}

// Register plugin
func (m DefaultPluginManager) Register(def PluginDef) {
	if def.NewFunc != nil {
		p := def.NewFunc()
		name := p.Name()
		m.plugins[name] = p
	}
}

func (m DefaultPluginManager) GetPluginByName(name string) (Plugin, bool) {
	p, ok := m.plugins[name]
	return p, ok
}

// GetPluginsByType returns plugins with matching pluginType
func (m DefaultPluginManager) GetPluginsByType(pluginType string) []Plugin {
	result := make([]Plugin, 0)
	for _, p := range m.plugins {
		if pluginType == p.Type() {
			result = append(result, p)
		}
	}
	return result
}

func GetPluginsByType(pluginType string) []Plugin {
	return defaultPluginManager.GetPluginsByType(pluginType)
}
