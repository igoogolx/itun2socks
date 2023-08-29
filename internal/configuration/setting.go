package configuration

func GetSetting() (SettingCfg, error) {
	c, err := Read()
	if err != nil {
		return SettingCfg{}, err
	}
	return c.Setting, nil
}

func SetSetting(data SettingCfg) error {
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
