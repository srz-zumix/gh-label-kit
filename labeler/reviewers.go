package labeler

import (
	"maps"
	"slices"
)

func CollectCodeownersSet(labels []string, cfg LabelerConfig) map[string]struct{} {
	ownerSet := make(map[string]struct{})
	for _, label := range labels {
		if lc, ok := cfg[label]; ok {
			for _, owner := range lc.Codeowners {
				ownerSet[owner] = struct{}{}
			}
		}
	}
	return ownerSet
}

func CollectCodeowners(labels []string, cfg LabelerConfig) []string {
	ownerSet := CollectCodeownersSet(labels, cfg)
	return slices.Collect(maps.Keys(ownerSet))
}
