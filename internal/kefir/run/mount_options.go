package run

import (
	"fmt"
	"kefir/pkg/ui"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type mountOptionDescriptionMountedTo struct {
	ContainerName string `json:"container"`
	ReadOnly      bool   `json:"readOnly,omitempty"`
	MountPath     string `json:"mountPath"`
	SubPath       string `json:"subPath,omitempty"`
}

func (m *mountOptionDescriptionMountedTo) String() string {
	return fmt.Sprintf("%s %s %s",
		m.ContainerName,
		IsRoIcons[m.ReadOnly],
		m.MountPath,
	)
}

type mountOptionDescription struct {
	MountedTo           []mountOptionDescriptionMountedTo `json:"mountedTo"`
	corev1.VolumeSource `json:",inline"`
}

func (m mountOptionDescription) String() string {
	s, _ := yaml.Marshal(m)

	return string(s)
}

// MountOptionImpl - реализация интерфейса ui.MountOption
type MountOptionImpl struct {
	Name      string
	MountPath string
	ReadOnly  bool
	// описание волума. включает в себя список контейнеров в который этот волум примонтирован, а также опции волума
	description mountOptionDescription
}

func (m *MountOptionImpl) Title() string {
	if m.MountPath == "" {

		return fmt.Sprintf("❌ %s", m.Name)
	}

	return fmt.Sprintf("%s %s → %s", IsRoIcons[m.ReadOnly], m.Name, m.MountPath)
}

func (m *MountOptionImpl) Description() string {
	return m.description.String()
}

func (m *MountOptionImpl) GetMountPath() string {
	return m.MountPath
}

func (m *MountOptionImpl) SetMountPath(path string) {
	m.MountPath = path
}

func (m *MountOptionImpl) GetReadOnly() bool {
	return m.ReadOnly
}

func (m *MountOptionImpl) SetReadOnly(readOnly bool) {
	m.ReadOnly = readOnly
}

func (m *MountOptionImpl) GetName() string {
	return m.Name
}

func (m *MountOptionImpl) Clone() ui.MountOption {
	return &MountOptionImpl{
		Name:        m.Name,
		MountPath:   m.MountPath,
		ReadOnly:    m.ReadOnly,
		description: m.description,
	}
}

var _ ui.MountOption = &MountOptionImpl{}

// GetAllMountOptions - из списка волумов Pod-а и всех найденных контейнеров пода составляет
// список имплементаций ui.MountOption
// список всех контейнеров нужен для более детальной информации в description
func GetAllMountOptions(volumes []corev1.Volume, allContainers []ContainerSelection) []ui.MountOption {
	ln := len(volumes)
	mountOptions := make([]ui.MountOption, ln)

	if ln == 0 {
		return mountOptions
	}

	for i, sv := range volumes {
		description := mountOptionDescription{VolumeSource: sv.VolumeSource}

		for _, c := range allContainers {
			for _, vm := range c.Container.VolumeMounts {
				if sv.Name == vm.Name {
					description.MountedTo = append(
						description.MountedTo, mountOptionDescriptionMountedTo{
							ContainerName: c.String(),
							ReadOnly:      vm.ReadOnly,
							MountPath:     vm.MountPath,
							SubPath:       vm.SubPath,
						})
				}
			}
		}

		mo := MountOptionImpl{
			Name:        sv.Name,
			MountPath:   "/run/volumes/" + sv.Name,
			ReadOnly:    true,
			description: description,
		}
		mountOptions[i] = &mo
	}

	return mountOptions
}
