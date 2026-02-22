package docker_client

type (
	ImageRecord struct {
		Image  string   `json:"image"      yaml:"image"      mapstructure:"image"` // имя образа
		Tag    string   `json:",omitempty" yaml:",omitempty" mapstructure:",omitempty"`
		Hash   string   `json:",omitempty" yaml:",omitempty" mapstructure:",omitempty"` // sha256sum of image
		Labels []string `json:",omitempty" yaml:",omitempty" mapstructure:",omitempty"` // список тегов в качестве подсказки
	}
)
