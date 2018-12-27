package mapprovider

// MapProvider provides a map for memory storarage of configuration settings
type MapProvider struct {
	Map map[string]string
}

// Provide implements the Provider interface
func (m MapProvider) Provide() (map[string]string, error) {
	return m.Map, nil
}
