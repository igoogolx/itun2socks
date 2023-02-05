package configuration

import (
	"github.com/igoogolx/itun2socks/configuration/configuration-types"
)

func GetSetting() (configuration_types.SettingCfg, error) {
	c, err := Read()
	if err != nil {
		return configuration_types.SettingCfg{}, err
	}
	return c.Setting, nil
}

func SetSetting(data configuration_types.SettingCfg) error {
	c, err := Read()
	if err != nil {
		return err
	}
	c.Setting = data
	err = Write(c)
	if err != nil {
		return err
	}
	return nil
}
