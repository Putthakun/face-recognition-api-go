package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type FaceClient interface {
	GetEmbedding(file *multipart.FileHeader) ([]float32, error)
	TriggerReload() error
}

type faceClient struct {
	baseURL string
	http    *http.Client
}

func NewFaceClient(baseURL string) FaceClient {
	return &faceClient{baseURL: baseURL, http: &http.Client{}}
}

type embeddingResponse struct {
	FaceDetected bool      `json:"face_detected"`
	Embedding    []float32 `json:"embedding"`
}

func (c *faceClient) GetEmbedding(fh *multipart.FileHeader) ([]float32, error) {
	f, err := fh.Open()
	if err != nil {
		return nil, fmt.Errorf("open uploaded file: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", fh.Filename)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(fw, f); err != nil {
		return nil, err
	}
	w.Close()

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/embeddings", &buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("face server request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("face server returned %d", resp.StatusCode)
	}

	var result embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if !result.FaceDetected {
		return nil, nil // signal: no face
	}
	return result.Embedding, nil
}

func (c *faceClient) TriggerReload() error {
	resp, err := c.http.Post(c.baseURL+"/api/reload", "application/json", nil)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
