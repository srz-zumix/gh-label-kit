package labeler

import (
	"maps"
	"slices"
)

func CollectCodeowners(labels []string, cfg LabelerConfig) []string {
	ownerSet := make(map[string]struct{})
	for _, label := range labels {
		if lc, ok := cfg[label]; ok {
			for _, owner := range lc.Codeowners {
				ownerSet[owner] = struct{}{}
			}
		}
	}
	return slices.Collect(maps.Keys(ownerSet))
}
