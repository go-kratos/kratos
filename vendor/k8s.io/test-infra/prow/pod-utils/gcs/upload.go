/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gcs

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"

	"cloud.google.com/go/storage"
	"github.com/sirupsen/logrus"

	"k8s.io/test-infra/prow/errorutil"
)

// UploadFunc knows how to upload into an object
type UploadFunc func(obj *storage.ObjectHandle) error

// Upload uploads all of the data in the
// uploadTargets map to GCS in parallel. The map is
// keyed on GCS path under the bucket
func Upload(bucket *storage.BucketHandle, uploadTargets map[string]UploadFunc) error {
	errCh := make(chan error, len(uploadTargets))
	group := &sync.WaitGroup{}
	group.Add(len(uploadTargets))
	for dest, upload := range uploadTargets {
		obj := bucket.Object(dest)
		logrus.WithField("dest", dest).Info("Queued for upload")
		go func(f UploadFunc, obj *storage.ObjectHandle, name string) {
			defer group.Done()
			if err := f(obj); err != nil {
				errCh <- err
			}
			logrus.WithField("dest", name).Info("Finished upload")
		}(upload, obj, dest)
	}
	group.Wait()
	close(errCh)
	if len(errCh) != 0 {
		var uploadErrors []error
		for err := range errCh {
			uploadErrors = append(uploadErrors, err)
		}
		return fmt.Errorf("encountered errors during upload: %v", uploadErrors)
	}

	return nil
}

// FileUpload returns an UploadFunc which copies all
// data from the file on disk to the GCS object
func FileUpload(file string) UploadFunc {
	return func(obj *storage.ObjectHandle) error {
		reader, err := os.Open(file)
		if err != nil {
			return err
		}

		uploadErr := DataUpload(reader)(obj)
		closeErr := reader.Close()

		return errorutil.NewAggregate(uploadErr, closeErr)
	}
}

// DataUpload returns an UploadFunc which copies all
// data from src reader into GCS
func DataUpload(src io.Reader) UploadFunc {
	return func(obj *storage.ObjectHandle) error {
		writer := obj.NewWriter(context.Background())
		_, copyErr := io.Copy(writer, src)
		closeErr := writer.Close()

		return errorutil.NewAggregate(copyErr, closeErr)
	}
}
