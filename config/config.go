package config

// Configer defines how to get and set value from configuration raw data.
type Configer interface {
	Set(key, val string) error   //support section::key type in given key when using ini type.
	String(key string) string    //support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	Strings(key string) []string //get string slice
	Int(key string) (int, error)
	Int64(key string) (int64, error)
	Bool(key string) (bool, error)
	Float(key string) (float64, error)
	DefaultString(key string, defaultVal string) string      // support section::key type in key string when using ini and json type; Int,Int64,Bool,Float,DIY are same.
	DefaultStrings(key string, defaultVal []string) []string //get string slice
	DefaultInt(key string, defaultVal int) int
	DefaultInt64(key string, defaultVal int64) int64
	DefaultBool(key string, defaultVal bool) bool
	DefaultFloat(key string, defaultVal float64) float64
	DIY(key string) (interface{}, error)
	GetSection(section string) (map[string]string, error)
	SaveConfigFile(filename string) error
}

type serviceConfig struct {
	innerConfig Configer
}

func newServiceConfig(appConfigProvider, appConfigPath string) (*serviceConfig, error) {
	ac, err := NewConfig(appConfigProvider, appConfigPath)
	if err != nil {
		return nil, err
	}
	return &serviceConfig{ac}, nil
}

func (b *serviceConfig) Set(key, val string) error {
	if err := b.innerConfig.Set(key, val); err != nil {
		return err
	}
	return b.innerConfig.Set(key, val)
}

func (b *serviceConfig) String(key string) string {
	if v := b.innerConfig.String(key); v != "" {
		return v
	}
	return b.innerConfig.String(key)
}

func (b *serviceConfig) Strings(key string) []string {
	if v := b.innerConfig.Strings(key); len(v) > 0 {
		return v
	}
	return b.innerConfig.Strings(key)
}

func (b *serviceConfig) Int(key string) (int, error) {
	if v, err := b.innerConfig.Int(key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int(key)
}

func (b *serviceConfig) Int64(key string) (int64, error) {
	if v, err := b.innerConfig.Int64(key); err == nil {
		return v, nil
	}
	return b.innerConfig.Int64(key)
}

func (b *serviceConfig) Bool(key string) (bool, error) {
	if v, err := b.innerConfig.Bool(key); err == nil {
		return v, nil
	}
	return b.innerConfig.Bool(key)
}

func (b *serviceConfig) Float(key string) (float64, error) {
	if v, err := b.innerConfig.Float(key); err == nil {
		return v, nil
	}
	return b.innerConfig.Float(key)
}

func (b *serviceConfig) DefaultString(key string, defaultVal string) string {
	if v := b.String(key); v != "" {
		return v
	}
	return defaultVal
}

func (b *serviceConfig) DefaultStrings(key string, defaultVal []string) []string {
	if v := b.Strings(key); len(v) != 0 {
		return v
	}
	return defaultVal
}

func (b *serviceConfig) DefaultInt(key string, defaultVal int) int {
	if v, err := b.Int(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *serviceConfig) DefaultInt64(key string, defaultVal int64) int64 {
	if v, err := b.Int64(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *serviceConfig) DefaultBool(key string, defaultVal bool) bool {
	if v, err := b.Bool(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *serviceConfig) DefaultFloat(key string, defaultVal float64) float64 {
	if v, err := b.Float(key); err == nil {
		return v
	}
	return defaultVal
}

func (b *serviceConfig) DIY(key string) (interface{}, error) {
	return b.innerConfig.DIY(key)
}

func (b *serviceConfig) GetSection(section string) (map[string]string, error) {
	return b.innerConfig.GetSection(section)
}

func (b *serviceConfig) SaveConfigFile(filename string) error {
	return b.innerConfig.SaveConfigFile(filename)
}

func init() {

}