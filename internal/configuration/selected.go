package configuration

import (
	"fmt"
)

func GetSelectedId(bucket string) (string, error) {
	c, err := Read()
	if err != nil {
		return "", err
	}
	if bucket == "rule" {
		return c.Selected.Rule, nil
	} else if bucket == "proxy" {
		return c.Selected.Proxy, nil
	}
	return "", fmt.Errorf("error gettting selected id,type:%v err: invalid field", bucket)
}

func SetSelectedId(bucket, id string) error {
	c, err := Read()
	if err != nil {
		return err
	}
	if bucket == "rule" {
		c.Selected.Rule = id
	} else if bucket == "proxy" {
		c.Selected.Proxy = id
	} else {
		return fmt.Errorf("error seting selected id,type:%v err: invalid field", bucket)
	}
	err = Write(c)
	if err != nil {
		return err
	}
	return nil
}
