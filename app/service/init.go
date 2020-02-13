package service

func InitPlugins(plugins map[string]map[string]interface{}) {
	pluginManager.Register(kakaoDef)

	// init plugins
	for name, props := range plugins {
		if plugin, ok := pluginManager.GetPluginByName(name); ok {
			plugin.SetProperties(props)
		}
	}

}