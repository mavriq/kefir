package image_list

// import (
// 	_ "encoding/json"

// 	_ "gopkg.in/yaml.v3"
// )

type (
	Config []ImageItem

	ImageItem struct {
		Image string   `json:"image" yaml:"image" mapstructure:"image"`                            // имя образа
		Tags  []string `json:"tags,omitempty" yaml:"tags,omitempty" mapstructure:"tags,omitempty"` // список тегов в качестве подсказки
	}
)

// func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
// 	var raw = make([]interface{})

// 	if err = unmarshal(&raw); err != nil {
// 		return err
// 	}

// }
