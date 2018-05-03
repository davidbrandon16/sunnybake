package Model

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"time"
	"encoding/json"
	"strings"
)

var domain string = "https://api.cloudinary.com/v1_1/"

// Cloudinary main struct
type Cloudinary struct {
	publicKey string
	secretKey string
	name      string
	urls      Option
}

// Option is the optional parameters custom struct
type Option map[string]string

// ErrorResp is the failed api request main struct
type ErrorResp struct {
	Message string `json:"message"`
}

// Create is creating a new cloudinary instance
func Create(public string, secret string, name string) *Cloudinary {
	return &Cloudinary{
		publicKey: public,
		secretKey: secret,
		name:      name,
		urls: Option{
			"upload": domain + name + "/image/upload",
		},
	}
}

func (c *Cloudinary) checkOptionsAreValid(options Option, validator []string) {
	for key := range options {
		if !validOption(key, validator) {
			panic("Upload paramater not valid: " + key)
		}
	}
}

var validOption = func(optionName string, val []string) bool {
	for _, param := range val {
		if param == optionName {
			return true
		}
	}
	return false
}

func (c *Cloudinary) sortParamsByKey(options Option) Option {
	options["timestamp"] = strconv.Itoa(int(time.Now().Unix()))
	result := Option{}
	sortedKeys := []string{}

	for key := range options {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		result[key] = options[key]
	}

	return result
}

func (c *Cloudinary) createSignature(options Option) string {
	signature := ""
	i := 0

	for key, value := range options {
		if i != 0 {
			signature += "&"
		}
		signature += key + "=" + value
		i++
	}

	hash := sha1.New()
	hash.Write([]byte(signature + c.secretKey))

	return hex.EncodeToString(hash.Sum(nil))
}

func (c *Cloudinary) send(url string, postParams url.Values, options Option) []byte {
	postParams.Add("api_key", c.publicKey)
	postParams.Add("signature", c.createSignature(options))
	postParams.Add("api_key", c.publicKey)
	postParams.Add("signature", c.createSignature(options))

	req, err := http.NewRequest("POST", c.urls[url], strings.NewReader(postParams.Encode()))

	if err != nil {
		panic(err)
	}

	client := http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return body
// KODINGAN DIA


	/*file, err := os.Open(postParams.Get("file_path"))
	if err != nil {
		panic(err)
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	fi, err := file.Stat()
	if err != nil {
		panic(err)
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", fi.Name())
	if err != nil {
		panic(err)
	}
	part.Write(fileContents)
	for key := range postParams {
		if(key != "file_path") {
			_ = writer.WriteField(key, postParams.Get(key))
		}
	}
	err = writer.Close()
	if err != nil {
		panic(err)
	}

	req, err :=  http.NewRequest("POST", c.urls[url], body)

	if err != nil {
		panic(err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	//respBody, _ := ioutil.ReadAll(resp.Body)

	return respBody*/
}
var uploadOptions []string = []string{
	"public_id",
	"use_filename",
	"unique_filename",
	"folder",
	"overwrite",
	"resource_type",
	"type",
	"tags",
	"context",
	"transformation",
	"format",
	"allowed_formats",
	"eager",
	"async",
	"eager_async",
	"proxy",
	"headers",
	"callback",
	"notification_url",
	"eager_notification_url",
	"faces",
	"image_metadata",
	"exif",
	"colors",
	"phash",
	"face_coordinates",
	"custom_coordinates",
	"backup",
	"return_delete_token",
	"invalidate",
	"discard_original_filename",
	"moderation",
	"upload_preset",
	"raw_convert",
	"categorization",
	"auto_tagging",
	"background_removal",
	"detection",
	"timestamp",

}
type Upload struct {
	PublicId         string    `json:"public_id"`
	Version          int       `json:"version"`
	Signature        string    `json:"signature"`
	Width            int       `json:"width"`
	Height           int       `json:"height"`
	Format           string    `json:"format"`
	ResourceType     string    `json:"resource_type"`
	CreatedAt        string    `json:"created_at"`
	Tags             []string  `json:"tags,omitempty"`
	Bytes            int       `json:"bytes"`
	Type             string    `json:"type"`
	Etag             string    `json:"etag"`
	Url              string    `json:"url"`
	SecureUrl        string    `json:"secure_url"`
	OriginalFilename string    `json:"original_filename"`
	Error            ErrorResp `json:"error,omitempty"`
}
func (c * Cloudinary) Upload(file string, options Option) (*Upload, error) {
	c.checkOptionsAreValid(options, uploadOptions)
	options = c.sortParamsByKey(options)

	form := url.Values{}
	for paramName, value := range options {
		form.Add(paramName, value)
	}
	form.Add("file", file)

	body := c.send("upload", form, options)

	upload := &Upload{}
	err := json.Unmarshal(body, upload)

	if err != nil {
		return nil, err
	}

	return upload, nil
}