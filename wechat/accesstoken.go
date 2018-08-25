package wechat

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"
)

// 微信授权access-token
type AccessToken struct {
	AppId     string
	AppSecret string
	TmpName   string
	LckName   string
}

func (at *AccessToken) Fresh() (string, error) {
	if at.TmpName == "" {
		at.TmpName = at.AppId + "-accesstoken.tmp"
	}
	if at.LckName == "" {
		at.LckName = at.TmpName + ".lck"
	}
	for {
		if at.locked() {
			time.Sleep(time.Second)
			continue
		}
		break
	}
	fi, err := os.Stat(at.TmpName)
	if err != nil && !os.IsExist(err) {
		return at.fetchAndStore()
	}
	expires := fi.ModTime().Add(2 * time.Hour).Unix()
	if expires <= time.Now().Unix() {
		return at.fetchAndStore()
	}
	tmp, err := os.Open(at.TmpName)
	if err != nil {
		return "", err
	}
	defer tmp.Close()
	data, err := ioutil.ReadAll(tmp)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (at *AccessToken) fetchAndStore() (string, error) {
	if err := at.lock(); err != nil {
		return "", err
	}
	defer at.unlock()
	token, err := at.fetch()
	if err != nil {
		return "", err
	}
	if err := at.store(token); err != nil {
		return "", err
	}
	return token, nil
}

func (at *AccessToken) store(token string) error {
	path := path.Dir(at.TmpName)
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	if !fi.IsDir() {
		return errors.New("path is not a directory")
	}
	tmp, err := os.OpenFile(at.TmpName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	defer tmp.Close()
	if _, err := tmp.Write([]byte(token)); err != nil {
		return err
	}
	return nil
}

func (at *AccessToken) fetch() (string, error) {
	rtn, err := get(fmt.Sprintf(
		"%stoken?grant_type=client_credential&appid=%s&secret=%s",
		UrlPrefix,
		at.AppId,
		at.AppSecret,
	))
	if err != nil {
		return "", err
	}
	return rtn.AccessToken, nil
}

func (at *AccessToken) unlock() error {
	return os.Remove(at.LckName)
}

func (at *AccessToken) lock() error {
	path := path.Dir(at.LckName)
	fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}
	if !fi.IsDir() {
		return errors.New("path is not a directory")
	}
	lck, err := os.Create(at.LckName)
	if err != nil {
		return err
	}
	lck.Close()
	return nil
}

func (at *AccessToken) locked() bool {
	_, err := os.Stat(at.LckName)
	return !os.IsNotExist(err)
}
