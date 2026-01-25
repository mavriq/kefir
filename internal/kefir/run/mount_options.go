package run

import (
	"fmt"
	"kefir/pkg/ui"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type mountOptionDescriptionMountedTo struct {
	ContainerName string `json:"containerName"`
	ReadOnly      bool   `json:"readOnly,omitempty"`
	MountPath     string `json:"mountPath"`
	SubPath       string `json:"subPath,omitempty"`
}

func (m *mountOptionDescriptionMountedTo) String() string {
	b := &strings.Builder{}
	b.WriteString(m.ContainerName)

	if m.ReadOnly {
		fmt.Fprint(b, "(ro) ")
	} else {
		fmt.Fprint(b, "(rw) ")
	}

	b.WriteString(m.MountPath)

	if m.SubPath != "" {
		fmt.Fprintf(b, " [%s]", m.MountPath)
	}

	return b.String()
}

type mountOptionDescription struct {
	MountedTo           []mountOptionDescriptionMountedTo `json:"mountedTo"`
	corev1.VolumeSource `json:",inline"`
}

func (m mountOptionDescription) String() string {
	b := &strings.Builder{}
	if len(m.MountedTo) > 0 {
		b.WriteString("mountedTo:\n")
		for _, mt := range m.MountedTo {
			fmt.Fprintf(b, "  - %w\n", &mt)
		}
	}
	bb, _ := yaml.Marshal(m.VolumeSource)

	b.Write(bb)

	return b.String()
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
	readOnlyText := ""
	if m.ReadOnly {
		readOnlyText = "🔒"
	}

	mountPath := "[DO NOT MOUNT]"
	if m.MountPath == "" {
		mountPath = m.MountPath
	}

	return fmt.Sprintf("%s → %s%s", m.Name, mountPath, readOnlyText)
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

	// volumes := make(ui.MountOption, 0, len(volumes))
	// mountOptions := make([]*MountOptionImpl, 0, len(volumes))
	mountOptions := make([]ui.MountOption, 0, len(volumes))

	for i, sv := range volumes {
		description := mountOptionDescription{VolumeSource: sv.VolumeSource}
		// description.WriteString("used on containers:\n")

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

		mountOptions[i] = &MountOptionImpl{
			Name:        sv.Name,
			MountPath:   "/run/volumes/" + sv.Name,
			ReadOnly:    true,
			description: description,
		}
	}

	return mountOptions
}
