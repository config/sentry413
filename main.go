package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"os"
	"path/filepath"
	"time"

	"github.com/getsentry/sentry-go"
)

const (
	URL_TEMPLATE = "https://%s/api/%s/events/%s/attachments/?sentry_key=%s&sentry_version=7&sentry_client=custom-javascript"
)

type SentryCapture struct{}

func (s SentryCapture) Error() string {
	return "sentry capture"
}

func help() {
	fmt.Println("")
	fmt.Println("sentry413 [host] [project_id] [public_key] [environment] [/path/to/file.ext]")
	fmt.Println("using example DSN from https://develop.sentry.dev/sdk/overview/")
	fmt.Println("example: 	sentry413 o0.ingest.sentry.io 0 examplePublicKey staging /tmp/a.zip")
}

func main() {

	if len(os.Args) < 6 {
		fmt.Println("missing arguments")
		help()
		os.Exit(1)
	}

	host := os.Args[1]
	project := os.Args[2]
	publicKey := os.Args[3]
	environment := os.Args[4]
	pathToFile := os.Args[5]

	if host == "" || project == "" || publicKey == "" || pathToFile == "" {
		fmt.Println("missing required parameters")
		help()
		os.Exit(1)
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              fmt.Sprintf("https://%s@%s/%s", publicKey, host, project),
		Environment:      environment,
		Release:          "1.0" + "@" + "dev",
		Debug:            true,
		TracesSampleRate: 1,
	}); err != nil {
		fmt.Printf("error initializing sentry: %q\n", err)
		os.Exit(1)
	}

	id := ""
	sentryID := sentry.CaptureException(SentryCapture{})
	if sentryID != nil {
		id = string(*sentryID)
	}

	if id == "" {
		fmt.Println("sentry ID missing")
		os.Exit(1)
	}

	fmt.Printf("=> id:            %v\n", id)

	urlToSentry := fmt.Sprintf(URL_TEMPLATE, host, project, id, publicKey)

	err := UploadFileToSentry(pathToFile, urlToSentry)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("=> submitted")
}

func UploadFileToSentry(pathToFile, urlToSentry string) error {
	if pathToFile == "" {
		panic(fmt.Errorf("missing path to file"))
	}

	fmt.Printf("=> path           %q\n", pathToFile)

	if urlToSentry == "" {
		panic(fmt.Errorf("missing url"))
	}

	fmt.Printf("=> url:           %q\n", urlToSentry)

	r, err := os.Open(pathToFile)
	if err != nil {
		fmt.Printf("error opening file (%q): %q\n", pathToFile, err)
		return err
	}
	defer r.Close()

	client := &http.Client{
		Timeout:   180 * time.Second,
		Transport: http.DefaultTransport,
	}

	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	fw, err := w.CreateFormFile("file", filepath.Base(pathToFile))
	if err != nil {
		fmt.Printf("error creating multipart form: %q\n", err)
		return err
	}

	copied, err := io.Copy(fw, r)
	if err != nil {
		fmt.Printf("error copying data: %q\n", err)
		return err
	}

	if err := w.Close(); err != nil {
		fmt.Printf("error closing multiwriter: %q, continuing\n", err)
	}

	req, err := http.NewRequest(http.MethodPost, urlToSentry, &b)
	if err != nil {
		fmt.Printf("error creating new request: %q\n", err)
		return err
	}

	req.Header.Add("Content-Type", w.FormDataContentType())

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error POSTING request: %q\n", err)
		return err
	}
	defer resp.Body.Close()

	fmt.Printf("=> data size:     %v\n", copied)
	fmt.Printf("=> status code:   %v\n", resp.StatusCode)

	body, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Printf("error reading response: %q\n", err)
		return err
	}

	fmt.Println(string(body))

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("problem submitting file to sentry, status code: %v", resp.StatusCode)
	}

	return nil
}
