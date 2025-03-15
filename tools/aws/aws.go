package aws

import "context"

type IS3 interface {
	Upload(ctx context.Context, bucketName string, fileName string, fileBytes []byte, contentType string) error
}
