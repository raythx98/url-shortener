package supabase

import (
	"os"

	"github.com/supabase-community/storage-go"
)

// TODO: Overly simplified implementation. Change to AWS S3 later.
var (
	// storageClient is a global variable that holds the storage client.
	storageClient *storage_go.Client
)

func New(supabaseSecret string) {
	//storageClient = storage_go.NewClient("https://caqzitwuslrszkfwbmve.supabase.co/storage/v1",
	//	supabaseSecret, nil)
	storageClient = storage_go.NewClient("https://caqzitwuslrszkfwbmve.supabase.co/storage/v1",
		supabaseSecret, nil)
}

func UpdateFile(fileName string, file *os.File) error {
	upsert := true
	contentType := "image/png"
	options := storage_go.FileOptions{
		ContentType: &contentType,
		Upsert:      &upsert,
	}
	_, err := storageClient.UpdateFile("qrs", fileName, file, options)
	return err
}
