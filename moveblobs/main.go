package main

import (
	"context"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"log"
	"net/url"
)

// http://127.0.0.1:10000/devstoreaccount1/events/file-example_PDF_1MB.pdf

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
	panic(err)
}

func credentials() *azblob.SharedKeyCredential {
	accountName := "devstoreaccount1"
	accountKey := "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="

	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	handleError(err)

	return cred
}

func main() {

	sourceUrl, _ := url.Parse("http://127.0.0.1:10000/devstoreaccount1/events/file-example_PDF_1MB.pdf")
	sourceBlobURL := azblob.NewBlobURL(*sourceUrl, azblob.NewPipeline(credentials(), azblob.PipelineOptions{}))

	destinationUrl, _ := url.Parse("http://127.0.0.1:10000/devstoreaccount1/files/file-example_PDF_1MB.pdf")
	//deztinationBlobUrl := azblob.NewBlobURL(*destinationUrl, azblob.NewPipeline(credentials(), azblob.PipelineOptions{}))

	_, err := sourceBlobURL.StartCopyFromURL(context.Background(), *destinationUrl, azblob.Metadata{}, azblob.ModifiedAccessConditions{}, azblob.BlobAccessConditions{}, azblob.AccessTierNone, nil)

	if err != nil {
		log.Fatal(err)
	}
}
