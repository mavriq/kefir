package configure

// import (
// 	_ "encoding/json"

// 	_ "gopkg.in/yaml.v3"
// )

type (
	Config []ImageItem

	ImageItem struct {
		Image  string   `json:"image"      yaml:"image"      mapstructure:"image"` // имя образа
		Tag    string   `json:",omitempty" yaml:",omitempty" mapstructure:",omitempty"`
		Hash   string   `json:",omitempty" yaml:",omitempty" mapstructure:",omitempty"` // sha256sum of image
		Labels []string `json:",omitempty" yaml:",omitempty" mapstructure:",omitempty"` // список тегов в качестве подсказки
	}
)

// func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
// 	var raw = make([]interface{})

// 	if err = unmarshal(&raw); err != nil {
// 		return err
// 	}

// }
