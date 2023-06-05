package quickstart

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

// Azure Storage Quickstart Sample - Demonstrate how to upload, list, download, and delete blobs.
//
// Documentation References:
// - What is a Storage Account - https://docs.microsoft.com/azure/storage/common/storage-create-storage-account
// - Blob Service Concepts - https://docs.microsoft.com/rest/api/storageservices/Blob-Service-Concepts
// - Blob Service Go SDK API - https://godoc.org/github.com/Azure/azure-storage-blob-go
// - Blob Service REST API - https://docs.microsoft.com/rest/api/storageservices/Blob-Service-REST-API
// - Scalability and performance targets - https://docs.microsoft.com/azure/storage/common/storage-scalability-targets
// - Azure Storage Performance and Scalability checklist https://docs.microsoft.com/azure/storage/common/storage-performance-checklist
// - Storage Emulator - https://docs.microsoft.com/azure/storage/common/storage-use-emulator

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getClient(accountName string, accountKey string, useCliAuth bool) (*azblob.Client, error) {
	url := fmt.Sprintf("http://127.0.0.1:10000/%s", accountName)
	if useCliAuth {
		cred, err := azidentity.NewAzureCLICredential(nil)
		if err != nil {
			return nil, err
		}
		return azblob.NewClient(url, cred, nil)
	}
	cred, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}
	return azblob.NewClientWithSharedKeyCredential(url, cred, nil)
}

func main() {
	fmt.Printf("Azure Blob storage quick start sample\n")

	// TODO: replace <storage-account-name> with your actual storage account name
	//url := "https://<storage-account-name>.blob.core.windows.net/"
	//url := "http://127.0.0.1:10000/devstoreaccount1"

	accountName := "devstoreaccount1"
	accountKey := "Eby8vdM02xNOcqFlqUwJPLlmEtlCDXJ1OUzFT50uSRZ6IFsuFq2UVErCz4I6tq/K1SZFPTOtr/KBHBeksoGMGw=="

	ctx := context.Background()
	//
	//credential, err := azidentity.NewDefaultAzureCredential()
	//handleError(err)
	//
	client, err := getClient(accountName, accountKey, false)
	//handleError(err)

	// Create the container
	containerName := "quickstart-sample-container"
	fmt.Printf("Creating a container named %s\n", containerName)
	_, err = client.CreateContainer(ctx, containerName, nil)
	handleError(err)

	data := []byte("\nHello, world! This is a blob.\n")
	blobName := "sample-blob"

	// Upload to data to blob storage
	fmt.Printf("Uploading a blob named %s\n", blobName)
	_, err = client.UploadBuffer(ctx, containerName, blobName, data, &azblob.UploadBufferOptions{})
	handleError(err)

	// List the blobs in the container
	fmt.Println("Listing the blobs in the container:")

	pager := client.NewListBlobsFlatPager(containerName, &azblob.ListBlobsFlatOptions{
		Include: azblob.ListBlobsInclude{Snapshots: true, Versions: true},
	})

	for pager.More() {
		resp, err := pager.NextPage(context.TODO())
		handleError(err)

		for _, blob := range resp.Segment.BlobItems {
			fmt.Println(*blob.Name)
		}
	}

	// Download the blob
	get, err := client.DownloadStream(ctx, containerName, blobName, nil)
	handleError(err)

	downloadedData := bytes.Buffer{}
	retryReader := get.NewRetryReader(ctx, &azblob.RetryReaderOptions{})
	_, err = downloadedData.ReadFrom(retryReader)
	handleError(err)

	err = retryReader.Close()
	handleError(err)

	// Print the content of the blob we created
	fmt.Println("Blob contents:")
	fmt.Println(downloadedData.String())

	fmt.Printf("Press enter key to delete resources and exit the application.\n")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	fmt.Printf("Cleaning up.\n")

	// Delete the blob
	fmt.Printf("Deleting the blob " + blobName + "\n")

	_, err = client.DeleteBlob(ctx, containerName, blobName, nil)
	handleError(err)

	// Delete the container
	fmt.Printf("Deleting the container " + containerName + "\n")
	_, err = client.DeleteContainer(ctx, containerName, nil)
	handleError(err)
}
