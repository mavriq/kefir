package efir

import (
	"context"
	"fmt"
	"kefir/pkg/kube_client"
	"kefir/pkg/ui"
	"slices"
	"strings"

	"github.com/containers/image/v5/docker/reference"
	corev1 "k8s.io/api/core/v1"
)

type Efir struct {
	imageName            string
	ns                   string
	podName              string
	ephyContainerName    string
	targetContainerName  string
	selectedMountOptions []ui.MountOption
}

func New(imageName, ns, podName, ephyContainerName string,
	targetContainer *corev1.Container,
	selectedMountOptions []ui.MountOption,
) (e *Efir, err error) {
	if ephyContainerName == "" {
		ref, err := reference.ParseNormalizedNamed(imageName)

		if err != nil {
			return nil, err
		}

		pathSlice := strings.Split(reference.Path(ref), "/")
		ephyContainerName = addSuffix(pathSlice[len(pathSlice)-1])
	}

	e = &Efir{
		imageName:         imageName,
		ns:                ns,
		podName:           podName,
		ephyContainerName: ephyContainerName,
		selectedMountOptions: slices.DeleteFunc(
			selectedMountOptions, func(m ui.MountOption) bool {
				return m.GetMountPath() != ""
			}),
	}

	if targetContainer != nil {
		e.targetContainerName = targetContainer.Name
	}

	return e, nil
}

func (e *Efir) Noop(ctx context.Context) (err error) {
	var patch string

	if patch, err = e.GenPatch(); err != nil {
		return err
	}

	cmd1 := []string{
		"kubectl", "patch", "pod",
		"-n", e.ns,
		"--subresource=ephemeralcontainers",
		"--type='json'",
		"-p", fmt.Sprintf("<(echo '%s')", patch),
	}

	fmt.Printf(strings.Join(cmd1, " "))

	return nil
}

func (p *Efir) Exec(ctx context.Context) (err error) {
	_, err = kube_client.ClientSet()

	return err
	// TODO: запилить
}
