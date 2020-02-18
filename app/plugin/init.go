package plugin

func InitPlugins(plugins map[string]map[string]interface{}) {
	defaultPluginManager.Register(PluginDef{
		Type: TypeBook,
		Name: "kakao",
		NewFunc: func() Plugin {
			return &Kakao{}
		},
	})

	// init plugins
	for name, props := range plugins {
		if p, ok := defaultPluginManager.GetPluginByName(name); ok {
			p.SetProperties(props)
		}
	}
}
