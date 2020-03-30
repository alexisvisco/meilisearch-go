package meilisearch

import (
	"testing"
)

type docTest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func TestClientDocuments_Get(t *testing.T) {
	var indexUID = "TestClientDocuments_Get"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]interface{}{
			docTest{ID: "123", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateIDSimplified(indexUID, updateIDRes)

	var doc docTest
	if err = client.
		Documents(indexUID).
		Get("123", &doc); err != nil {
		t.Fatal(err)
	}

	expect := docTest{ID: "123", Name: "nestle"}
	if doc != expect {
		t.Errorf("%v != %v", doc, expect)
	}
}

func TestClientDocuments_Delete(t *testing.T) {
	var indexUID = "TestClientDocuments_Delete"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]interface{}{
			docTest{ID: "123", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateIDSimplified(indexUID, updateIDRes)

	updateIDRes, err = client.Documents(indexUID).Delete("123")

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateIDSimplified(indexUID, updateIDRes)

	var doc docTest
	err = client.Documents(indexUID).Get("123", &doc)

	if err.(*Error).ErrCode != ErrCodeResponseStatusCode {
		t.Fatal(err)
	}
}

func TestClientDocuments_Deletes(t *testing.T) {
	var indexUID = "deletes"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]interface{}{
			docTest{ID: "123", Name: "nestle"},
			docTest{ID: "456", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateIDSimplified(indexUID, updateIDRes)

	updateIDRes, err = client.Documents(indexUID).Deletes([]string{"123", "456"})

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateIDSimplified(indexUID, updateIDRes)

	var doc docTest
	err = client.Documents(indexUID).Get("123", &doc)

	if err.(*Error).ErrCode != ErrCodeResponseStatusCode {
		t.Fatal(err)
	}
}

func TestClientDocuments_List(t *testing.T) {
	var indexUID = "TestClientDocuments_List"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]interface{}{
			docTest{ID: "123", Name: "nestle"},
			docTest{ID: "456", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateIDSimplified(indexUID, updateIDRes)

	var list []docTest
	err = client.Documents(indexUID).List(ListDocumentsRequest{
		Offset: 0,
		Limit:  100,
	}, &list)

	if err != nil {
		t.Fatal(err)
	}

	// tests are running in parallel so there can be more than 1 docs
	if len(list) < 1 {
		t.Fatal("number of doc should be at least 1")
	}
}

func TestClientDocuments_AddOrReplace(t *testing.T) {
	var indexUID = "TestClientDocuments_AddOrReplace"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrReplace([]docTest{
			{ID: "123", Name: "nestle"},
			{ID: "456", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateIDSimplified(indexUID, updateIDRes)

	var list []docTest
	err = client.Documents(indexUID).List(ListDocumentsRequest{
		Offset: 0,
		Limit:  100,
	}, &list)

	if err != nil {
		t.Fatal(err)
	}

	// tests are running in parallel so there can be more than 1 docs
	if len(list) < 2 {
		t.Fatal("number of doc should be at least 1")
	}
}

func TestClientDocuments_AddOrUpdate(t *testing.T) {
	var indexUID = "TestClientDocuments_AddOrUpdate"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]docTest{
			{ID: "123", Name: "nestle"},
			{ID: "456", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateIDSimplified(indexUID, updateIDRes)

	var list []docTest
	err = client.Documents(indexUID).List(ListDocumentsRequest{
		Offset: 0,
		Limit:  100,
	}, &list)

	if err != nil {
		t.Fatal(err)
	}

	// tests are running in parallel so there can be more than 1 docs
	if len(list) < 2 {
		t.Fatal("number of doc should be at least 1")
	}
}

func TestClientDocuments_DeleteAllDocuments(t *testing.T) {
	var indexUID = "TestClientDocuments_DeleteAllDocuments"

	var client = NewClient(Config{
		Host: "http://localhost:7700",
	})

	_, err := client.Indexes().Create(CreateIndexRequest{
		UID: indexUID,
	})

	if err != nil {
		t.Fatal(err)
	}

	updateIDRes, err := client.
		Documents(indexUID).
		AddOrUpdate([]interface{}{
			docTest{ID: "123", Name: "nestle"},
			docTest{ID: "456", Name: "nestle"},
		})

	if err != nil {
		t.Fatal(err)
	}

	client.AwaitAsyncUpdateIDSimplified(indexUID, updateIDRes)

	_, err = client.Documents(indexUID).DeleteAllDocuments()

	if err != nil {
		t.Fatal(err)
	}
}
