package main

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
)

const (
	BUCKET_NAME = "weibin-around-bucket"
)

// More sample code can be found at https://github.com/GoogleCloudPlatform/golang-samples/blob/master/storage/objects/main.go
func saveToGCS(r io.Reader, objectName string) (string, error) {

	// Get an emtpy context
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}

	// TODO: may configure time out in the context, e.g.
	// ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	// defer cancel() // https://golang.org/pkg/context/#CancelFunc
	// Then the time limit is configured into wc through `wc := object.NewWriter(ctx)` below
	object := client.Bucket(BUCKET_NAME).Object(objectName)

	// wc can be thought of as the remote file that we are writing to
	wc := object.NewWriter(ctx)
	if _, err := io.Copy(wc, r); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	// The file is publicly readable.
	// TODO: Improve this with better access control, ideally only allow access from this service.
	if err := object.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	attrs, err := object.Attrs(ctx)
	if err != nil {
		return "", err
	}

	// Return the URL where the media is stored.
	fmt.Printf("Image is saved to GCS: %s\n", attrs.MediaLink)
	return attrs.MediaLink, nil
}
