// +build darwin

package darwin

import "testing"

func TestMemoryGenerator(t *testing.T) {
	g := &MemoryGenerator{}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	for _, name := range []string{
		"total",
		"free",
		"cached",
		"active",
		"inactive",
		"used",
	} {
		if v, ok := values["memory."+name]; !ok {
			t.Errorf("memory should has %s", name)
		} else {
			t.Logf("memory '%s' collected: %+v", name, v)
		}
	}
}
