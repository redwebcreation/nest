package common

import (
	"github.com/redwebcreation/nest/util"
	"gopkg.in/yaml.v3"
	"testing"
)

func TestServiceMap_IncludeService(t *testing.T) {
	f := util.TmpFile()
	_, err := f.WriteString("image: nginx\nlistening_on: 5016\n")
	if err != nil {
		t.Fatal(err)
	}

	sm := ServiceMap{
		"example": {
			Include: f.Name(),
		},
	}

	marshalled, err := yaml.Marshal(&sm)
	if err != nil {
		t.Fatal(err)
	}

	err = yaml.Unmarshal(marshalled, &sm)
	if err != nil {
		t.Errorf("unmarshalling failed: %s", err)
	}

	if sm["example"].Image != "nginx" {
		t.Errorf("expected image to be nginx, got %s", sm["example"].Image)
	}

	if sm["example"].ListeningOn != "5016" {
		t.Errorf("expected listening_on to be 5016, got %s", sm["example"].ListeningOn)
	}
}
