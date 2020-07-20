package meilisearch

// APIWithIndexID is used to await an async update id response.
// Each apis that use an index internally implement this interface except APIUpdates.
type APIWithIndexID interface {
	IndexID() string
	Client() *Client
}

// APIIndexes index is an entity, like a table in SQL, with a specific schema definition. It gathers a collection of
// documents with the structure defined by the schema.
// An index is defined by an unique identifier uid that is generated by MeiliSearch (if none is given) on index
// creation. It also has a name to help you track your different indexes.
//
// Documentation: https://docs.meilisearch.com/references/indexes.html
type APIIndexes interface {

	// Get the index relative information.
	Get(uid string) (*Index, error)

	// List all indexes.
	List() ([]Index, error)

	// Create an index.
	// If no UID is specified in the request a randomly generated UID will be returned.
	// It's associated to the new index. This UID will be essential to make all request over the created index.
	// You can define your primary key during the index creation
	Create(request CreateIndexRequest) (*CreateIndexResponse, error)

	// Update an index name.
	UpdateName(uid string, name string) (*Index, error)

	// Update an index primary key.
	UpdatePrimaryKey(uid string, primaryKey string) (*Index, error)

	// Delete an index.
	Delete(uid string) (bool, error)
}

// APIDocuments are objects composed of fields containing any data.
//
// Documentation: https://docs.meilisearch.com/references/documents.html
type APIDocuments interface {

	// Get one document using its unique identifier.
	// documentPtr should be a pointer.
	Get(identifier string, documentPtr interface{}) error

	// Delete one document based on its unique identifier.
	Delete(identifier string) (*AsyncUpdateID, error)

	// Delete a selection of documents based on array of identifiers.
	Deletes(identifier []string) (*AsyncUpdateID, error)

	// List the documents in an unordered way.
	List(request ListDocumentsRequest, documentsPtr interface{}) error

	// AddOrReplace a list of documents, replace them if they already exist based on their unique identifiers.
	AddOrReplace(documentsPtr interface{}) (*AsyncUpdateID, error)

	// AddOrReplaceWithPrimaryKey do the same as AddOrReplace but will specify during the update to primaryKey to use for indexing
	AddOrReplaceWithPrimaryKey(documentsPtr interface{}, primaryKey string) (resp *AsyncUpdateID, err error)

	// AddOrUpdate a list of documents, update them if they already exist based on their unique identifiers.
	AddOrUpdate(documentsPtr interface{}) (*AsyncUpdateID, error)

	// AddOrUpdateWithPrimaryKey do the same as AddOrUpdate but will specify during the update to primaryKey to use for indexing
	AddOrUpdateWithPrimaryKey(documentsPtr interface{}, primaryKey string) (resp *AsyncUpdateID, err error)

	// DeleteAllDocuments in the specified index.
	DeleteAllDocuments() (*AsyncUpdateID, error)

	APIWithIndexID
}

// APISearch search through documents list in an index.
//
// Documentation: https://docs.meilisearch.com/references/search.html
type APISearch interface {

	// Search for documents matching a specific query in the given index.
	Search(params SearchRequest) (*SearchResponse, error)

	APIWithIndexID
}

// APIUpdates MeiliSearch is an asynchronous API. It means that the API does not behave as you would typically expect
// when handling the request's responses.
//
// Some actions are put in a queue and will be executed in turn (asynchronously). In this case, the server response
// contains the identifier to track the execution of the action.
//
// This API permit to get the state of an update id.
//
// Documentation: https://docs.meilisearch.com/references/updates.html
type APIUpdates interface {

	// Get the status of an update in a given index.
	Get(id int64) (*Update, error)

	// Get the status of all updates in a given index.
	List() ([]Update, error)

	APIWithIndexID
}

// APIKeys To communicate with MeiliSearch's RESTfull API most of the routes require an API key.
//
// Documentation: https://docs.meilisearch.com/references/keys.html
type APIKeys interface {

	// Get all keys.
	Get() (*Keys, error)
}

// APISettings allow to configure the MeiliSearch indexing & search behaviour.
//
// Documentation: https://docs.meilisearch.com/references/settings.html
type APISettings interface {
	GetAll() (*Settings, error)

	UpdateAll(request Settings) (*AsyncUpdateID, error)

	ResetAll() (*AsyncUpdateID, error)

	GetRankingRules() (*[]string, error)

	UpdateRankingRules([]string) (*AsyncUpdateID, error)

	ResetRankingRules() (*AsyncUpdateID, error)

	GetDistinctAttribute() (*string, error)

	UpdateDistinctAttribute(string) (*AsyncUpdateID, error)

	ResetDistinctAttribute() (*AsyncUpdateID, error)

	GetSearchableAttributes() (*[]string, error)

	UpdateSearchableAttributes([]string) (*AsyncUpdateID, error)

	ResetSearchableAttributes() (*AsyncUpdateID, error)

	GetDisplayedAttributes() (*[]string, error)

	UpdateDisplayedAttributes([]string) (*AsyncUpdateID, error)

	ResetDisplayedAttributes() (*AsyncUpdateID, error)

	GetStopWords() (*[]string, error)

	UpdateStopWords([]string) (*AsyncUpdateID, error)

	ResetStopWords() (*AsyncUpdateID, error)

	GetSynonyms() (*map[string][]string, error)

	UpdateSynonyms(map[string][]string) (*AsyncUpdateID, error)

	ResetSynonyms() (*AsyncUpdateID, error)

	GetAcceptNewFields() (*bool, error)

	UpdateAcceptNewFields(bool) (*AsyncUpdateID, error)

	GetAttributesForFaceting() (*[]string, error)

	UpdateAttributesForFaceting([]string) (*AsyncUpdateID, error)

	ResetAttributesForFaceting() (*AsyncUpdateID, error)
}

// APIStats retrieve statistic over all indexes or a specific index id.
//
// Documentation: https://docs.meilisearch.com/references/stats.html
type APIStats interface {

	// Get stats of an index.
	Get(indexUID string) (*StatsIndex, error)

	GetAll() (*Stats, error)
}

// APIHealth handle health of a MeiliSearch server.
//
// Documentation: https://docs.meilisearch.com/references/health.html
type APIHealth interface {

	// Get health of MeiliSearch server.
	Get() error

	// Update health of MeiliSearch server.
	Update(health bool) error
}

// APIVersion retrieve the version of MeiliSearch.
//
// Documentation: https://docs.meilisearch.com/references/version.html
type APIVersion interface {

	// Get version of MeiliSearch.
	Get() (*Version, error)
}
