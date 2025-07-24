package storage

import (
	supabasestorageuploader "github.com/adityarizkyramadhan/supabase-storage-uploader"
)

// New creates a new Supabase storage client
func New(url, key, bucket string) *supabasestorageuploader.Client {
	return supabasestorageuploader.New(url, key, bucket)
}
