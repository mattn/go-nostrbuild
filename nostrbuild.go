package nostrbuild

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

type Dimensions struct {
	Height int64 `json:"height"`
	Width  int64 `json:"width"`
}

type Metadata struct {
	Date_create              string `json:"date:create"`
	Date_modify              string `json:"date:modify"`
	Png_IHDR_bit_depth_orig  string `json:"png:IHDR.bit-depth-orig"`
	Png_IHDR_bitDepth        string `json:"png:IHDR.bit_depth"`
	Png_IHDR_color_type_orig string `json:"png:IHDR.color-type-orig"`
	Png_IHDR_colorType       string `json:"png:IHDR.color_type"`
	Png_IHDR_interlaceMethod string `json:"png:IHDR.interlace_method"`
	Png_IHDR_width_height    string `json:"png:IHDR.width,height"`
	Png_PLTE_numberColors    string `json:"png:PLTE.number_colors"`
	Png_pHYs                 string `json:"png:pHYs"`
	Png_sRGB                 string `json:"png:sRGB"`
}

type Responsive struct {
	One080p  string `json:"1080p"`
	Two40p   string `json:"240p"`
	Three60p string `json:"360p"`
	Four80p  string `json:"480p"`
	Seven20p string `json:"720p"`
}

type Data struct {
	Blurhash         string     `json:"blurhash"`
	Dimensions       Dimensions `json:"dimensions"`
	DimensionsString string     `json:"dimensionsString"`
	InputName        string     `json:"input_name"`
	Metadata         Metadata   `json:"metadata"`
	Mime             string     `json:"mime"`
	Name             string     `json:"name"`
	OriginalSha256   string     `json:"original_sha256"`
	Responsive       Responsive `json:"responsive"`
	Sha256           string     `json:"sha256"`
	Size             int64      `json:"size"`
	Thumbnail        string     `json:"thumbnail"`
	Type             string     `json:"type"`
	URL              string     `json:"url"`
}

type Response struct {
	Data    []Data `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func Upload(buf *bytes.Buffer) (*Response, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	part, err := w.CreateFormFile("fileToUpload", "fileToUpload")
	if err != nil {
		return nil, err
	}
	part.Write(buf.Bytes())
	err = w.Close()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, "https://nostr.build/api/v2/upload/files", &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	if resp.StatusCode != 200 {
		if b, err := io.ReadAll(resp.Body); err == nil {
			return nil, errors.New(string(b))
		}
	}

	var p Response
	err = json.NewDecoder(resp.Body).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
