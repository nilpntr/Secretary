package cli

import (
	"github.com/spf13/pflag"
	"os"
	"strconv"
	"strings"
)

type EnvSettings struct {
	ExcludedNamespaces []string
	SyncDelay          int
}

func New() *EnvSettings {
	env := &EnvSettings{
		ExcludedNamespaces: parseStringSlice("EXCLUDED_NAMESPACES"),
		SyncDelay:          parseIntFallback("SYNC_DELAY", 15),
	}

	return env
}

func parseStringSlice(envName string) []string {
	val := os.Getenv(envName)
	if val == "" {
		return []string{}
	}
	return strings.Split(val, ",")
}

func parseIntFallback(envName string, fallback int) int {
	val, err := strconv.Atoi(os.Getenv(envName))
	if err != nil {
		return fallback
	}
	return val
}

func (e *EnvSettings) AddFlags(f *pflag.FlagSet) {
	f.StringSliceVarP(&e.ExcludedNamespaces, "excluded-namespaces", "", e.ExcludedNamespaces, "excluded namespaces")
	f.IntVarP(&e.SyncDelay, "sync-delay", "", e.SyncDelay, "sync delay in seconds")
}
