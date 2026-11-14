// Package generator: generate shell functions.
package generator

import (
	"fmt"
	"io"
	"strings"

	"github.com/aprx/around/internal/config"
)

func Generate(configurations []config.Configuration, w io.Writer) error {
	for _, conf := range configurations {
		allkeys := []string{}
		name := conf.Name
		_, _ = fmt.Fprintf(w, "## Source: %s\n", conf.Source)
		_, _ = fmt.Fprintf(w, "conf_%s() {\n", name)
		for _, v := range conf.Variables {
			allkeys = append(allkeys, v.Name)
			_, _ = fmt.Fprintf(w, "  export %s=%s\n", v.Name, v.Value)
		}
		if len(conf.SSHKeys) > 0 {
			name = fmt.Sprintf("[K:%d]%s", len(conf.SSHKeys), name)
			_, _ = fmt.Fprint(w, "  conf_start_ssh_agent\n")
			for _, key := range conf.SSHKeys {
				_, _ = fmt.Fprintf(w, "  %s | ssh-add -\n", key)
			}
		}
		_, _ = fmt.Fprintf(w, "  _around_set_env %s\n", name)
		_, _ = fmt.Fprintf(w, "  _around_set_clean_env %s\n", strings.Join(allkeys, ":"))
		_, _ = fmt.Fprint(w, "}\n\n")
	}
	return nil
}
