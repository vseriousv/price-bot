package utils

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
)

func GetHttp(url string, token string) ([]byte, error) {
	reqStr := fmt.Sprintf(url)

	log.Print("[HTTP REQUEST TO] ", reqStr)
	client := &http.Client{}

	req, err := http.NewRequest("GET", reqStr, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	if err != nil {
		fmt.Println("[ERROR:GetHttp] ", err)
		return nil, err
	}

	get, err := client.Do(req)

	response, err2 := ioutil.ReadAll(get.Body)
	if err2 != nil {
		fmt.Println("[ERROR:GetHttp:ReadAll] ", err2)
		return nil, err2
	}

	if get.StatusCode >= 400 && get.StatusCode < 600 {
		fmt.Println("[ERROR:GetHttp:StatusCode] ", get.StatusCode)
		return nil, errors.New(fmt.Sprintf("%d", get.StatusCode))
	}

	return response, nil
}

func PostHttp(url string, body io.Reader) ([]byte, error) {
	reqStr := fmt.Sprintf(url)

	log.Print("[HTTP REQUEST TO] ", reqStr)

	post, err := http.Post(url, "application/json", body)
	if err != nil {
		fmt.Println("[ERROR:GetHttp] ", err)
		return nil, err
	}

	response, err2 := ioutil.ReadAll(post.Body)
	if err2 != nil {
		fmt.Println("[ERROR:GetHttp:ReadAll] ", err2)
		return nil, err2
	}

	if post.StatusCode >= 400 && post.StatusCode < 600 {
		fmt.Println("[ERROR:GetHttp:StatusCode] ", post.StatusCode)
		return nil, errors.New(fmt.Sprintf("%d", post.StatusCode))
	}

	return response, nil
}

type Number interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func Round[T Number](x T, f int) float64 {
	return math.Round(float64(x)*float64(f)) / float64(f)
}

func ToS[T Number](num T) string {
	return fmt.Sprintf("%.8f", num)
}

func SToInt64(str string) int64 {
	num, _ := strconv.Atoi(str)
	return int64(num)
}
