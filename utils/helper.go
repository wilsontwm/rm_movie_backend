package utils

import (
	//"fmt"
	"io/ioutil"
	"net/http"
)

// Get the param key in the URL
func GetParam(r *http.Request, key string) (result string) {
	keys, ok := r.URL.Query()[key]
	if ok && len(keys[0]) > 0 {
		result = keys[0]
	}

	return
}

// Helper function to download file from a url
func DownloadFile(url string, folder string, filename string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	fileBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(folder+"/"+filename, fileBytes, 0644)

	return err
}
