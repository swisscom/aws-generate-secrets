package generatesecrets

type GeneratorConfig struct {
	Length  int     `yaml:"length"`
	Chars   string  `yaml:"chars"`
	Charset Charset `yaml:"charset"`
}

type Charset string

const (
	CharsetStd       Charset = "standard"
	CharsetSpecial   Charset = "special"
	CharsetLowercase Charset = "lowercase"
)

type Secret struct {
	Name      string                     `yaml:"name"`
	Generator map[string]GeneratorConfig `yaml:"generator"`
}

type Config struct {
	Secrets []Secret `yaml:"secrets"`
	Profile string   `yaml:"profile"`
}
