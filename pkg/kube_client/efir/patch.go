package efir

import "encoding/json"

type (
	patch struct {
		Op    string     `json:"op"`
		Path  string     `json:"path"`
		Value patchValue `json:"value"`
	}
	patchValue struct {
		Name                string                  `json:"name"`
		Image               string                  `json:"image"`
		TargetContainerName string                  `json:"targetContainerName,omitempty"`
		Stdin               bool                    `json:"stdin"`
		Tty                 bool                    `json:"tty"`
		VolumeMounts        []patchValueVolumeMount `json:"volumeMounts"`
	}
	patchValueVolumeMount struct {
		Name      string `json:"name"`
		MountPath string `json:"mountPath"`
	}
)

func (e *Efir) GenPatch() (string, error) {
	p := patch{
		Op:   "add",
		Path: "/spec/ephemeralContainers/-",
		Value: patchValue{
			Name:                e.ephyContainerName,
			Image:               e.imageName,
			TargetContainerName: e.targetContainerName,
			Stdim:               true,
			Tty:                 true,
			VolumeMounts:        make([]patchValueVolumeMount, len(e.selectedMountOptions)),
		},
	}
	for i, mo := range e.selectedMountOptions {
		p.Value.VolumeMounts[i].patchValueVolumeMount.Name = mo.GetName()
		p.Value.VolumeMounts[i].patchValueVolumeMount.MountPath = mo.GetMountPath()
	}

	bb, err := json.Marshal([]patch{p})
	return string(bb), err
}
