package cache

import (
	"bufio"
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"github.com/djherbis/times"
	"io/ioutil"
	"os"
	"sanakirjaorg-cli/log"
	"time"
)

// Age Cache age in minutes, currently 72 hours
const Age = 4320

func verifyCacheAge(cachePath string) (bool, error) {
	t, err := times.Stat(cachePath)
	if err != nil {
		return false, err
	}
	diff := time.Now().Sub(t.ChangeTime())
	return diff.Minutes() > Age, nil
}

type CacheExpiredError struct{}

func (m *CacheExpiredError) Error() string {
	return "Cache expired"
}

func GetCache(resource string, target interface{}) error {
	dirname, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	var shaSignature = sha1.New()
	shaSignature.Write([]byte(resource))
	contentHash := hex.EncodeToString(shaSignature.Sum(nil))
	var cacheDir = dirname + string(os.PathSeparator) + "sanakirja-cmd" + string(os.PathSeparator)
	log.Log(log.DEBUG, "Cache: "+cacheDir+contentHash)
	err2 := os.MkdirAll(cacheDir, os.ModePerm)
	if err2 != nil {
		return err2
	}
	valid, err := verifyCacheAge(cacheDir + contentHash)
	if err != nil {
		return err
	}
	if !valid {
		cacheContents, err3 := ioutil.ReadFile(cacheDir + contentHash)
		if err3 != nil {
			return err3
		}
		return json.NewDecoder(bytes.NewReader(cacheContents)).Decode(target)
	} else {
		return &CacheExpiredError{}
	}
}

func SaveCache(resource string, target interface{}) error {
	dirname, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	var shaSignature = sha1.New()
	shaSignature.Write([]byte(resource))
	contentHash := hex.EncodeToString(shaSignature.Sum(nil))
	var cacheDir = dirname + string(os.PathSeparator) + "sanakirja-cmd" + string(os.PathSeparator)
	log.Log(log.DEBUG, "Cache: "+cacheDir+contentHash)
	err2 := os.MkdirAll(cacheDir, os.ModePerm)
	if err2 != nil {
		return err2
	}
	fo, err := os.Create(cacheDir + contentHash)
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(fo)
	err = json.NewEncoder(writer).Encode(target)
	if err != nil {
		return err
	}
	err = writer.Flush()
	return err
}
