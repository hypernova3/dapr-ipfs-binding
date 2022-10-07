package main

import (
	"context"
	"encoding/json"
	"log"

	dapr "github.com/dapr/go-sdk/client"
)

const componentName = "myipfs"

func main() {
	client, err := dapr.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Dapr client: %v", err)
	}

	// Add a file
	res, err := client.InvokeBinding(context.TODO(), &dapr.BindingInvocation{
		Name:      componentName,
		Operation: "add",
		Data:      []byte("🚀"),
	})
	if err != nil {
		log.Fatalf("Failed to add file: %v", err)
	}
	resMap := map[string]any{}
	err = json.Unmarshal(res.Data, &resMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}
	log.Println("Added object (CIDv0):", resMap["path"])
	path := resMap["path"].(string)

	// Add a file with CIDv1
	res, err = client.InvokeBinding(context.TODO(), &dapr.BindingInvocation{
		Name:      componentName,
		Operation: "add",
		Data:      []byte("🚀"),
		Metadata: map[string]string{
			"cidVersion": "1",
		},
	})
	if err != nil {
		log.Fatalf("Failed to add file: %v", err)
	}
	err = json.Unmarshal(res.Data, &resMap)
	if err != nil {
		log.Fatalf("Failed to unmarshal response: %v", err)
	}
	log.Println("Added object (CIDv1):", resMap["path"])

	// Get file
	res, err = client.InvokeBinding(context.TODO(), &dapr.BindingInvocation{
		Name:      componentName,
		Operation: "get",
		Metadata: map[string]string{
			"path": path,
		},
	})
	if err != nil {
		log.Fatalf("Failed to get file: %v", err)
	}
	log.Printf("Read object at path '%s': '%s'", path, string(res.Data))

	// List files in folder
	res, err = client.InvokeBinding(context.TODO(), &dapr.BindingInvocation{
		Name:      componentName,
		Operation: "ls",
		Metadata: map[string]string{
			"path": "QmYwAPJzv5CZsnA625s3Xf2nemtYgPpHdWEz79ojWnPbdG",
		},
	})
	if err != nil {
		log.Fatalf("Failed to list files: %v", err)
	}
	log.Println("Files in folder:", string(res.Data))
}
