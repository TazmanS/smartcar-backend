package cars

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/TazmanS/smartcar-backend/internal/app"
)

type CVClient struct {
	app        *app.App
	httpClient *http.Client
}

func NewCVClient(app *app.App) *CVClient {
	return &CVClient{
		app:        app,
		httpClient: &http.Client{},
	}
}

func (c *CVClient) Detect(ctx context.Context, frame []byte) ([]byte, error) {
	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile(
		"file",
		"frame.jpg",
	)
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}

	if _, err := part.Write(frame); err != nil {
		return nil, fmt.Errorf("write frame: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	url := c.app.Config.PythonServerUrl + "/api/cv/detect"

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		&body,
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request to CV: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"CV returned status %d: %s",
			resp.StatusCode,
			string(data),
		)
	}

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read CV response: %w", err)
	}

	return result, nil
}
