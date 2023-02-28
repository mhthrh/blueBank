package Config

import (
	"bytes"
	"fmt"
	"github.com/mhthrh/BlueBank/Utils/CryptoUtil"
	"github.com/mhthrh/BlueBank/Utils/FileUtil"
	"github.com/spf13/viper"
)

type Config struct {
	Name string
	Type string
	Path string
}

func New(name, typ, path string) *Config {
	return &Config{
		Name: name,
		Type: typ,
		Path: path,
	}
}

func (c *Config) Initialize() error {
	viper.SetConfigName(c.Name)
	viper.SetConfigType(c.Type)

	err := viper.ReadConfig(bytes.NewBuffer(func() []byte {
		k := CryptoUtil.NewKey()
		k.Text, _ = File.New(c.Path, c.Name).Read()
		dec, _ := k.Decrypt()
		return []byte(dec)
	}()))
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	return nil
}
