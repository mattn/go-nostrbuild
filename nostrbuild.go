package nostrbuild

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"

	"github.com/nbd-wtf/go-nostr"
)

type Dimensions struct {
	Height int64 `json:"height"`
	Width  int64 `json:"width"`
}

type Metadata struct {
	DateCreate             string `json:"date:create"`
	DateModify             string `json:"date:modify"`
	PngIHDRBitDepthOrig    string `json:"png:IHDR.bit-depth-orig"`
	PngIHDRBitDepth        string `json:"png:IHDR.bit_depth"`
	PngIHDRColorTypeOrig   string `json:"png:IHDR.color-type-orig"`
	PngIHDRColorType       string `json:"png:IHDR.color_type"`
	PngIHDRInterlaceMethod string `json:"png:IHDR.interlace_method"`
	PngIHDRWidthHeight     string `json:"png:IHDR.width,height"`
	PngPLTENumberColors    string `json:"png:PLTE.number_colors"`
	PngPHYS                string `json:"png:pHYs"`
	PngSRGB                string `json:"png:sRGB"`
}

type Responsive struct {
	S1080p string `json:"1080p"`
	S240p  string `json:"240p"`
	S360p  string `json:"360p"`
	S480p  string `json:"480p"`
	S720p  string `json:"720p"`
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
	Data    []Data `json:"data,omitempty"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func Upload(buf *bytes.Buffer, f func(ev *nostr.Event) error) (*Response, error) {
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

	postUrl := "https://nostr.build/api/v2/upload/files"

	req, err := http.NewRequest(http.MethodPost, postUrl, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	if f != nil {
		var ev nostr.Event
		ev.Tags = ev.Tags.AppendUnique(nostr.Tag{"u", postUrl})
		ev.Tags = ev.Tags.AppendUnique(nostr.Tag{"method", "POST"})
		ev.Kind = 27235
		ev.CreatedAt = nostr.Now()
		err = f(&ev)
		if err != nil {
			return nil, err
		}
		b, err := ev.MarshalJSON()
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Nostr "+base64.StdEncoding.EncodeToString(b))
	}

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

func Delete(deleteUrl string, f func(ev *nostr.Event) error) (*Response, error) {
	if u, err := url.Parse(deleteUrl); err == nil {
		u.Host = "nostr.build"
		u.Path = path.Join("/api/v2/nip96/upload", path.Base(deleteUrl))
		deleteUrl = u.String()
		println(deleteUrl)
	}

	req, err := http.NewRequest(http.MethodDelete, deleteUrl, nil)
	if err != nil {
		return nil, err
	}

	if f != nil {
		var ev nostr.Event
		ev.Tags = ev.Tags.AppendUnique(nostr.Tag{"u", deleteUrl})
		ev.Tags = ev.Tags.AppendUnique(nostr.Tag{"method", "DELETE"})
		ev.Kind = 27235
		ev.CreatedAt = nostr.Now()
		err = f(&ev)
		if err != nil {
			return nil, err
		}
		b, err := ev.MarshalJSON()
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Nostr "+base64.StdEncoding.EncodeToString(b))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
