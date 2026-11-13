package handler

func responseHeadersByName(merged map[string]namedAssetHeader) map[string]string {
	result := make(map[string]string, len(merged))
	for _, header := range merged {
		result[header.name] = header.value
	}
	return result
}
