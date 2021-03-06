package config

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestEnvProviderSetUp(t *testing.T) {
	tests := []struct {
		name         string
		regs         []string
		wantRegistry map[string]Value
		wantRegs     []string
	}{
		{
			"Empty set",
			[]string{},
			map[string]Value{},
			[]string{},
		},
		{
			"Simple KV pair not preffixed with CONFIG_",
			[]string{"FOO=BAR"},
			map[string]Value{},
			[]string{},
		},
		{
			"Simple KV pair preffixed with CONFIG_",
			[]string{"CONFIG_FOO=BAR"},
			map[string]Value{"foo": "BAR"},
			[]string{"foo"},
		},
		{
			"A KV with nested key",
			[]string{"CONFIG_FOO_BAR_BAZ=moo"},
			map[string]Value{"foo.bar.baz": "moo"},
			[]string{"foo.bar.baz"},
		},
		{
			"A value with underscore",
			[]string{"CONFIG_FOO=bar_baz_moo"},
			map[string]Value{"foo": "bar_baz_moo"},
			[]string{"foo"},
		},
		{
			"A key with double underscore",
			[]string{"CONFIG_FOO__BAR=moo__baz"},
			map[string]Value{"foo_bar": "moo__baz"},
			[]string{"foo_bar"},
		},
	}

	t.Parallel()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			// Redefined function
			oldEnvVars := envVars
			envVars = func() []string { return testCase.regs }

			repo := NewRepository()
			prov, err := NewEnvProvider(repo, 0)
			if err != nil {
				t.Fatalf("Failed to initialize a new env provider: %s", err)
			}
			if err := prov.SetUp(repo); err != nil {
				t.Fatalf("Failed to set up env provider: %s", err)
			}

			gotRegs := flattenRepo(repo)
			for _, k := range testCase.wantRegs {
				provs, ok := gotRegs[k]
				if !ok {
					t.Fatalf("Failed to find a registration for key %q: got regs: %#v", k, gotRegs)
				}
				if !reflect.DeepEqual(provs, []Provider{prov}) {
					t.Fatalf("Unexpected provider list for key %q: %#v, want: %#v", k, provs, []Provider{prov})
				}
				delete(gotRegs, k)
			}
			if len(gotRegs) > 0 {
				extraKeys := make([]string, 0, len(gotRegs))
				for k := range gotRegs {
					extraKeys = append(extraKeys, k)
				}
				sort.Strings(extraKeys)
				t.Fatalf("Unexpected registration keys: %s", strings.Join(extraKeys, ", "))
			}

			if !reflect.DeepEqual(prov.registry, testCase.wantRegistry) {
				t.Fatalf("Unexpected state for CliProvider.registry: want: %#v, got: %#v", testCase.wantRegistry, prov.registry)
			}

			envVars = oldEnvVars
		})
	}
}
