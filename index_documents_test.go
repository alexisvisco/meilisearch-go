package meilisearch

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex_AddDocuments(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
		documentsPtr interface{}
	}
	tests := []struct {
		name          string
		args          args
		wantResp      *AsyncUpdateID
		expectedError Error
	}{
		{
			name: "TestIndexBasicAddDocuments",
			args: args{
				UID:    "1",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexAddDocumentsWithCustomClient",
			args: args{
				UID:    "2",
				client: customClient,
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexMultipleAddDocuments",
			args: args{
				UID:    "2",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"ID": "1", "Name": "Alice In Wonderland"},
					{"ID": "123", "Name": "Pride and Prejudice"},
					{"ID": "456", "Name": "Le Petit Prince"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexBasicAddDocumentsWithIntID",
			args: args{
				UID:    "3",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"BookID": float64(123), "Title": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexAddDocumentsWithIntIDWithCustomClient",
			args: args{
				UID:    "4",
				client: customClient,
				documentsPtr: []map[string]interface{}{
					{"BookID": float64(123), "Title": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexMultipleAddDocumentsWithIntID",
			args: args{
				UID:    "5",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"BookID": float64(1), "Title": "Alice In Wonderland"},
					{"BookID": float64(123), "Title": "Pride and Prejudice"},
					{"BookID": float64(456), "Title": "Le Petit Prince", "Tag": "Conte"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.AddDocuments(tt.args.documentsPtr)
			require.GreaterOrEqual(t, gotResp.UpdateID, tt.wantResp.UpdateID)
			require.NoError(t, err)
			testWaitForPendingUpdate(t, i, gotResp)
			var documents []map[string]interface{}
			err = i.GetDocuments(&DocumentsRequest{
				Limit: 3,
			}, &documents)
			require.NoError(t, err)
			require.Equal(t, tt.args.documentsPtr, documents)
		})
	}
}

func TestIndex_AddDocumentsWithPrimaryKey(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
		documentsPtr interface{}
		primaryKey   string
	}
	tests := []struct {
		name     string
		args     args
		wantResp *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicAddDocumentsWithPrimaryKey",
			args: args{
				UID:    "1",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"key": "123", "Name": "Pride and Prejudice"},
				},
				primaryKey: "key",
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexAddDocumentsWithPrimaryKeyWithCustomClient",
			args: args{
				UID:    "2",
				client: customClient,
				documentsPtr: []map[string]interface{}{
					{"key": "123", "Name": "Pride and Prejudice"},
				},
				primaryKey: "key",
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexMultipleAddDocumentsWithPrimaryKey",
			args: args{
				UID:    "3",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"key": "1", "Name": "Alice In Wonderland"},
					{"key": "123", "Name": "Pride and Prejudice"},
					{"key": "456", "Name": "Le Petit Prince"},
				},
				primaryKey: "key",
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexAddDocumentsWithPrimaryKeyWithIntID",
			args: args{
				UID:    "4",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"key": float64(123), "Name": "Pride and Prejudice"},
				},
				primaryKey: "key",
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexMultipleAddDocumentsWithPrimaryKeyWithIntID",
			args: args{
				UID:    "5",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"key": float64(1), "Name": "Alice In Wonderland"},
					{"key": float64(123), "Name": "Pride and Prejudice"},
					{"key": float64(456), "Name": "Le Petit Prince"},
				},
				primaryKey: "key",
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.AddDocuments(tt.args.documentsPtr, tt.args.primaryKey)
			require.GreaterOrEqual(t, gotResp.UpdateID, tt.wantResp.UpdateID)
			require.NoError(t, err)

			testWaitForPendingUpdate(t, i, gotResp)

			var documents []map[string]interface{}
			err = i.GetDocuments(&DocumentsRequest{Limit: 3}, &documents)
			require.NoError(t, err)
			require.Equal(t, tt.args.documentsPtr, documents)
		})
	}
}

func TestIndex_AddDocumentsInBatches(t *testing.T) {
	type argsNoKey struct {
		UID          string
		client       *Client
		documentsPtr interface{}
		batchSize    int
	}

	type argsWithKey struct {
		UID          string
		client       *Client
		documentsPtr interface{}
		batchSize    int
		primaryKey   string
	}

	testsNoKey := []struct {
		name          string
		args          argsNoKey
		wantResp      []AsyncUpdateID
		expectedError Error
	}{
		{
			name: "TestIndexBasicAddDocumentsInBatches",
			args: argsNoKey{
				UID:    "0",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"ID": "122", "Name": "Pride and Prejudice"},
					{"ID": "123", "Name": "Pride and Prejudica"},
					{"ID": "124", "Name": "Pride and Prejudicb"},
					{"ID": "125", "Name": "Pride and Prejudicc"},
				},
				batchSize: 2,
			},
			wantResp: []AsyncUpdateID{
				{UpdateID: 0},
				{UpdateID: 1},
			},
		},
	}

	testsWithKey := []struct {
		name          string
		args          argsWithKey
		wantResp      []AsyncUpdateID
		expectedError Error
	}{
		{
			name: "TestIndexBasicAddDocumentsInBatches",
			args: argsWithKey{
				UID:    "0",
				client: defaultClient,
				documentsPtr: []map[string]interface{}{
					{"ID": "122", "Name": "Pride and Prejudice"},
					{"ID": "123", "Name": "Pride and Prejudica"},
					{"ID": "124", "Name": "Pride and Prejudicb"},
					{"ID": "125", "Name": "Pride and Prejudicc"},
				},
				batchSize:  2,
				primaryKey: "ID",
			},
			wantResp: []AsyncUpdateID{
				{UpdateID: 0},
				{UpdateID: 1},
			},
		},
	}

	for _, tt := range testsNoKey {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.AddDocumentsInBatches(tt.args.documentsPtr, tt.args.batchSize)
			require.Equal(t, gotResp, tt.wantResp)
			require.NoError(t, err)

			testWaitForPendingBatchUpdate(t, i, gotResp)

			var documents []map[string]interface{}
			err = i.GetDocuments(&DocumentsRequest{
				Limit: 4,
			}, &documents)

			require.NoError(t, err)
			require.Equal(t, tt.args.documentsPtr, documents)
		})
	}

	for _, tt := range testsWithKey {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotResp, err := i.AddDocumentsInBatches(tt.args.documentsPtr, tt.args.batchSize, tt.args.primaryKey)
			require.Equal(t, gotResp, tt.wantResp)
			require.NoError(t, err)

			testWaitForPendingBatchUpdate(t, i, gotResp)

			var documents []map[string]interface{}
			err = i.GetDocuments(&DocumentsRequest{
				Limit: 4,
			}, &documents)

			require.NoError(t, err)
			require.Equal(t, tt.args.documentsPtr, documents)
		})
	}
}

func TestIndex_DeleteAllDocuments(t *testing.T) {
	type args struct {
		UID    string
		client *Client
	}
	tests := []struct {
		name     string
		args     args
		wantResp *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicDeleteAllDocuments",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexDeleteAllDocumentsWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: customClient,
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			SetUpBasicIndex()
			gotResp, err := i.DeleteAllDocuments()
			require.NoError(t, err)
			require.Equal(t, gotResp, tt.wantResp)

			testWaitForPendingUpdate(t, i, gotResp)

			var documents interface{}
			err = i.GetDocuments(&DocumentsRequest{Limit: 5}, &documents)
			require.NoError(t, err)
			require.Empty(t, documents)
		})
	}
}

func TestIndex_DeleteOneDocument(t *testing.T) {
	type args struct {
		UID          string
		PrimaryKey   string
		client       *Client
		identifier   string
		documentsPtr interface{}
	}
	tests := []struct {
		name     string
		args     args
		wantResp *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicDeleteOneDocument",
			args: args{
				UID:        "1",
				client:     defaultClient,
				identifier: "123",
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexDeleteOneDocumentWithCustomClient",
			args: args{
				UID:        "2",
				client:     customClient,
				identifier: "123",
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexDeleteOneDocumentinMultiple",
			args: args{
				UID:        "3",
				client:     defaultClient,
				identifier: "456",
				documentsPtr: []map[string]interface{}{
					{"ID": "123", "Name": "Pride and Prejudice"},
					{"ID": "456", "Name": "Le Petit Prince"},
					{"ID": "1", "Name": "Alice In Wonderland"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexBasicDeleteOneDocumentWithIntID",
			args: args{
				UID:        "4",
				client:     defaultClient,
				identifier: "123",
				documentsPtr: []map[string]interface{}{
					{"BookID": 123, "Title": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexDeleteOneDocumentWithIntIDWithCustomClient",
			args: args{
				UID:        "5",
				client:     customClient,
				identifier: "123",
				documentsPtr: []map[string]interface{}{
					{"BookID": 123, "Title": "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
		{
			name: "TestIndexDeleteOneDocumentWithIntIDinMultiple",
			args: args{
				UID:        "6",
				client:     defaultClient,
				identifier: "456",
				documentsPtr: []map[string]interface{}{
					{"BookID": 123, "Title": "Pride and Prejudice"},
					{"BookID": 456, "Title": "Le Petit Prince"},
					{"BookID": 1, "Title": "Alice In Wonderland"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotAddResp, err := i.AddDocuments(tt.args.documentsPtr)
			require.GreaterOrEqual(t, gotAddResp.UpdateID, tt.wantResp.UpdateID)
			require.NoError(t, err)

			testWaitForPendingUpdate(t, i, gotAddResp)

			gotResp, err := i.DeleteDocument(tt.args.identifier)
			require.NoError(t, err)
			require.GreaterOrEqual(t, gotResp.UpdateID, tt.wantResp.UpdateID)

			testWaitForPendingUpdate(t, i, gotResp)

			var document []map[string]interface{}
			err = i.GetDocument(tt.args.identifier, &document)
			require.Error(t, err)
			require.Empty(t, document)
		})
	}
}

func TestIndex_DeleteDocuments(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
		identifier   []string
		documentsPtr []docTest
	}
	tests := []struct {
		name     string
		args     args
		wantResp *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicDeleteDocument",
			args: args{
				UID:        "1",
				client:     defaultClient,
				identifier: []string{"123"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexDeleteDocumentWithCustomClient",
			args: args{
				UID:        "2",
				client:     customClient,
				identifier: []string{"123"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexBasicDeleteDocument",
			args: args{
				UID:        "3",
				client:     defaultClient,
				identifier: []string{"123"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexBasicDeleteDocument",
			args: args{
				UID:        "4",
				client:     defaultClient,
				identifier: []string{"123", "456", "1"},
				documentsPtr: []docTest{
					{ID: "123", Name: "Pride and Prejudice"},
					{ID: "456", Name: "Le Petit Prince"},
					{ID: "1", Name: "Alice In Wonderland"},
				},
			},
			wantResp: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))

			gotAddResp, err := i.AddDocuments(tt.args.documentsPtr)
			require.NoError(t, err)

			testWaitForPendingUpdate(t, i, gotAddResp)

			gotResp, err := i.DeleteDocuments(tt.args.identifier)
			require.NoError(t, err)
			require.Equal(t, gotResp, tt.wantResp)

			testWaitForPendingUpdate(t, i, gotResp)

			var document docTest
			for _, identifier := range tt.args.identifier {
				err = i.GetDocument(identifier, &document)
				require.Error(t, err)
				require.Empty(t, document)
			}
		})
	}
}

func TestIndex_GetDocument(t *testing.T) {
	type args struct {
		UID         string
		client      *Client
		identifier  string
		documentPtr *docTestBooks
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "TestIndexBasicGetDocument",
			args: args{
				UID:         "indexUID",
				client:      defaultClient,
				identifier:  "123",
				documentPtr: &docTestBooks{},
			},
			wantErr: false,
		},
		{
			name: "TestIndexGetDocumentWithCustomClient",
			args: args{
				UID:         "indexUID",
				client:      customClient,
				identifier:  "123",
				documentPtr: &docTestBooks{},
			},
			wantErr: false,
		},
		{
			name: "TestIndexGetDocumentWithNoExistingDocument",
			args: args{
				UID:         "indexUID",
				client:      defaultClient,
				identifier:  "125",
				documentPtr: &docTestBooks{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))
			SetUpBasicIndex()

			require.Empty(t, tt.args.documentPtr)
			err := i.GetDocument(tt.args.identifier, tt.args.documentPtr)
			if tt.wantErr {
				require.Error(t, err)
				require.Empty(t, tt.args.documentPtr)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, tt.args.documentPtr)
				require.Equal(t, strconv.Itoa(tt.args.documentPtr.BookID), tt.args.identifier)
			}
		})
	}
}

func TestIndex_UpdateDocuments(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
		documentsPtr []docTestBooks
	}
	tests := []struct {
		name string
		args args
		want *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicUpdateDocument",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
				},
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
				},
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentOnMultipleDocuments",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
					{BookID: 1344, Title: "Harry Potter and the Half-Blood Prince"},
					{BookID: 4, Title: "The Hobbit"},
					{BookID: 42, Title: "The Great Gatsby"},
				},
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentWithNoExistingDocument",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 237, Title: "One Hundred Years of Solitude"},
				},
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentWithNoExistingMultipleDocuments",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 246, Title: "One Hundred Years of Solitude"},
					{BookID: 834, Title: "To Kill a Mockingbird"},
					{BookID: 44, Title: "Don Quixote"},
					{BookID: 594, Title: "The Great Gatsby"},
				},
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))
			SetUpBasicIndex()

			got, err := i.UpdateDocuments(tt.args.documentsPtr)
			require.NoError(t, err)
			require.Equal(t, got, tt.want)

			testWaitForPendingUpdate(t, i, got)

			var document docTestBooks
			for _, identifier := range tt.args.documentsPtr {
				err = i.GetDocument(strconv.Itoa(identifier.BookID), &document)
				require.NoError(t, err)
				require.Equal(t, identifier.BookID, document.BookID)
				require.Equal(t, identifier.Title, document.Title)
			}
		})
	}
}

func TestIndex_UpdateDocumentsWithPrimaryKey(t *testing.T) {
	type args struct {
		UID          string
		client       *Client
		documentsPtr []docTestBooks
		primaryKey   string
	}
	tests := []struct {
		name string
		args args
		want *AsyncUpdateID
	}{
		{
			name: "TestIndexBasicUpdateDocumentsWithPrimaryKey",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
				},
				primaryKey: "book_id",
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentsWithPrimaryKeyWithCustomClient",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
				},
				primaryKey: "book_id",
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentsWithPrimaryKeyOnMultipleDocuments",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
					{BookID: 1344, Title: "Harry Potter and the Half-Blood Prince"},
					{BookID: 4, Title: "The Hobbit"},
					{BookID: 42, Title: "The Great Gatsby"},
				},
				primaryKey: "book_id",
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentsWithPrimaryKeyWithNoExistingDocument",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 237, Title: "One Hundred Years of Solitude"},
				},
				primaryKey: "book_id",
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
		{
			name: "TestIndexUpdateDocumentsWithPrimaryKeyWithNoExistingMultipleDocuments",
			args: args{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 246, Title: "One Hundred Years of Solitude"},
					{BookID: 834, Title: "To Kill a Mockingbird"},
					{BookID: 44, Title: "Don Quixote"},
					{BookID: 594, Title: "The Great Gatsby"},
				},
				primaryKey: "book_id",
			},
			want: &AsyncUpdateID{
				UpdateID: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))
			SetUpBasicIndex()

			got, err := i.UpdateDocuments(tt.args.documentsPtr, tt.args.primaryKey)
			require.NoError(t, err)
			require.Equal(t, got, tt.want)

			testWaitForPendingUpdate(t, i, got)

			var document docTestBooks
			for _, identifier := range tt.args.documentsPtr {
				err = i.GetDocument(strconv.Itoa(identifier.BookID), &document)
				require.NoError(t, err)
				require.Equal(t, identifier.BookID, document.BookID)
				require.Equal(t, identifier.Title, document.Title)
			}
		})
	}
}

func TestIndex_UpdateDocumentsInBatches(t *testing.T) {
	type argsNoKey struct {
		UID          string
		client       *Client
		documentsPtr []docTestBooks
		batchSize    int
	}

	type argsWithKey struct {
		UID          string
		client       *Client
		documentsPtr []docTestBooks
		batchSize    int
		primaryKey   string
	}

	testsNoKey := []struct {
		name string
		args argsNoKey
		want []AsyncUpdateID
	}{
		{
			name: "TestIndexBatchUpdateDocuments",
			args: argsNoKey{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
					{BookID: 124, Title: "One Hundred Years of Solitude 2"},
				},
				batchSize: 1,
			},
			want: []AsyncUpdateID{
				{UpdateID: 1},
				{UpdateID: 2},
			},
		},
	}

	testsWithKey := []struct {
		name string
		args argsWithKey
		want []AsyncUpdateID
	}{
		{
			name: "TestIndexBatchUpdateDocuments",
			args: argsWithKey{
				UID:    "indexUID",
				client: defaultClient,
				documentsPtr: []docTestBooks{
					{BookID: 123, Title: "One Hundred Years of Solitude"},
					{BookID: 124, Title: "One Hundred Years of Solitude 2"},
				},
				batchSize:  1,
				primaryKey: "book_id",
			},
			want: []AsyncUpdateID{
				{UpdateID: 1},
				{UpdateID: 2},
			},
		},
	}

	for _, tt := range testsNoKey {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))
			SetUpBasicIndex()

			got, err := i.UpdateDocumentsInBatches(tt.args.documentsPtr, tt.args.batchSize)
			require.NoError(t, err)
			require.Equal(t, got, tt.want)

			testWaitForPendingBatchUpdate(t, i, got)

			var document docTestBooks
			for _, identifier := range tt.args.documentsPtr {
				err = i.GetDocument(strconv.Itoa(identifier.BookID), &document)
				require.NoError(t, err)
				require.Equal(t, identifier.BookID, document.BookID)
				require.Equal(t, identifier.Title, document.Title)
			}
		})
	}

	for _, tt := range testsWithKey {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.args.client
			i := c.Index(tt.args.UID)
			t.Cleanup(cleanup(c))
			SetUpBasicIndex()

			got, err := i.UpdateDocumentsInBatches(tt.args.documentsPtr, tt.args.batchSize, tt.args.primaryKey)
			require.NoError(t, err)
			require.Equal(t, got, tt.want)

			testWaitForPendingBatchUpdate(t, i, got)

			var document docTestBooks
			for _, identifier := range tt.args.documentsPtr {
				err = i.GetDocument(strconv.Itoa(identifier.BookID), &document)
				require.NoError(t, err)
				require.Equal(t, identifier.BookID, document.BookID)
				require.Equal(t, identifier.Title, document.Title)
			}
		})
	}
}
