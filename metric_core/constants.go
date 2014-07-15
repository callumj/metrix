package metric_core

import (
	"fmt"
)

var (
	SourcesKey = "sources"
)

func KeySourcesKey(source string) string {
	if len(source) == 0 {
		source = "default"
	}
	return fmt.Sprintf("keys:%s", source)
}

func VersionSourcesKey(source string) string {
	if len(source) == 0 {
		source = "default"
	}

	return fmt.Sprintf("versions:%v", source)
}
