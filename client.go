package meilisearch

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/valyala/fasthttp"
)

// ClientConfig configure the Client
type ClientConfig struct {

	// Host is the host of your Meilisearch database
	// Example: 'http://localhost:7700'
	Host string

	// APIKey is optional
	APIKey string

	// Timeout is optional
	Timeout time.Duration
}

type WaitParams struct {
	Context  context.Context
	Interval time.Duration
}

// ClientInterface is interface for all Meilisearch client
type ClientInterface interface {
	Index(uid string) *Index
	GetIndex(indexID string) (resp *Index, err error)
	GetRawIndex(uid string) (resp map[string]interface{}, err error)
	GetIndexes() (resp []*Index, err error)
	GetRawIndexes() (resp []map[string]interface{}, err error)
	CreateIndex(config *IndexConfig) (resp *Task, err error)
	DeleteIndex(uid string) (resp *Task, err error)
	CreateKey(request *Key) (resp *Key, err error)
	GetKey(identifier string) (resp *Key, err error)
	GetKeys() (resp *ResultKey, err error)
	UpdateKey(identifier string, request *Key) (resp *Key, err error)
	DeleteKey(identifier string) (resp bool, err error)
	GetStats() (resp *Stats, err error)
	CreateDump() (resp *Dump, err error)
	GetDumpStatus(dumpUID string) (resp *Dump, err error)
	Version() (*Version, error)
	GetVersion() (resp *Version, err error)
	Health() (*Health, error)
	IsHealthy() bool
	GetTask(taskID int64) (resp *Task, err error)
	GetTasks() (resp *ResultTask, err error)
	WaitForTask(task *Task, options ...WaitParams) (*Task, error)
	GenerateTenantToken(searchRules map[string]interface{}, options *TenantTokenOptions) (resp string, err error)
}

var _ ClientInterface = &Client{}

// NewFastHTTPCustomClient creates Meilisearch with custom fasthttp.Client
func NewFastHTTPCustomClient(config ClientConfig, client *fasthttp.Client) *Client {
	c := &Client{
		config:     config,
		httpClient: client,
	}
	return c
}

// NewClient creates Meilisearch with default fasthttp.Client
func NewClient(config ClientConfig) *Client {
	client := &fasthttp.Client{
		Name: "meilisearch-client",
		// Reuse the most recently-used idle connection.
		ConnPoolStrategy: fasthttp.LIFO,
	}
	c := &Client{
		config:     config,
		httpClient: client,
	}
	return c
}

func (c *Client) Version() (resp *Version, err error) {
	resp = &Version{}
	req := internalRequest{
		endpoint:            "/version",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Version",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetVersion() (resp *Version, err error) {
	return c.Version()
}

func (c *Client) GetStats() (resp *Stats, err error) {
	resp = &Stats{}
	req := internalRequest{
		endpoint:            "/stats",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetStats",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) CreateKey(request *Key) (resp *Key, err error) {
	parsedRequest := convertKeyToParsedKey(*request)
	resp = &Key{}
	req := internalRequest{
		endpoint:            "/keys",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         &parsedRequest,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusCreated},
		functionName:        "CreateKey",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetKey(identifier string) (resp *Key, err error) {
	resp = &Key{}
	req := internalRequest{
		endpoint:            "/keys/" + identifier,
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetKey",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetKeys() (resp *ResultKey, err error) {
	resp = &ResultKey{}
	req := internalRequest{
		endpoint:            "/keys",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetKeys",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) UpdateKey(identifier string, request *Key) (resp *Key, err error) {
	parsedRequest := convertKeyToParsedKey(*request)
	resp = &Key{}
	req := internalRequest{
		endpoint:            "/keys/" + identifier,
		method:              http.MethodPatch,
		contentType:         contentTypeJSON,
		withRequest:         &parsedRequest,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "UpdateKey",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) DeleteKey(identifier string) (resp bool, err error) {
	req := internalRequest{
		endpoint:            "/keys/" + identifier,
		method:              http.MethodDelete,
		withRequest:         nil,
		withResponse:        nil,
		acceptedStatusCodes: []int{http.StatusNoContent},
		functionName:        "DeleteKey",
	}
	if err := c.executeRequest(req); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Client) Health() (resp *Health, err error) {
	resp = &Health{}
	req := internalRequest{
		endpoint:            "/health",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "Health",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) IsHealthy() bool {
	if _, err := c.Health(); err != nil {
		return false
	}
	return true
}

func (c *Client) CreateDump() (resp *Dump, err error) {
	resp = &Dump{}
	req := internalRequest{
		endpoint:            "/dumps",
		method:              http.MethodPost,
		contentType:         contentTypeJSON,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusAccepted},
		functionName:        "CreateDump",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetDumpStatus(dumpUID string) (resp *Dump, err error) {
	resp = &Dump{}
	req := internalRequest{
		endpoint:            "/dumps/" + dumpUID + "/status",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetDumpStatus",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetTask(taskID int64) (resp *Task, err error) {
	resp = &Task{}
	req := internalRequest{
		endpoint:            "/tasks/" + strconv.FormatInt(taskID, 10),
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetTask",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) GetTasks() (resp *ResultTask, err error) {
	resp = &ResultTask{}
	req := internalRequest{
		endpoint:            "/tasks",
		method:              http.MethodGet,
		withRequest:         nil,
		withResponse:        &resp,
		acceptedStatusCodes: []int{http.StatusOK},
		functionName:        "GetTasks",
	}
	if err := c.executeRequest(req); err != nil {
		return nil, err
	}
	return resp, nil
}

// WaitForTask waits for a task to be processed
//
// The function will check by regular interval provided in parameter interval
// the TaskStatus.
// If no ctx and interval are provided WaitForTask will check each 50ms the
// status of a task.
func (c *Client) WaitForTask(task *Task, options ...WaitParams) (*Task, error) {
	if options == nil {
		ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
		defer cancelFunc()
		options = []WaitParams{
			{
				Context:  ctx,
				Interval: time.Millisecond * 50,
			},
		}
	}
	for {
		if err := options[0].Context.Err(); err != nil {
			return nil, err
		}
		getTask, err := c.GetTask(task.UID)
		if err != nil {
			return nil, err
		}
		if getTask.Status != TaskStatusEnqueued && getTask.Status != TaskStatusProcessing {
			return getTask, nil
		}
		time.Sleep(options[0].Interval)
	}
}

// Generate a JWT token for the use of multitenancy
//
// SearchRules parameters is mandatory and should contains the rules to be enforced at search time for all or specific
// accessible indexes for the signing API Key.
// ExpiresAt options is a time.Time when the key will expire. Note that if an ExpiresAt value is included it should be in UTC time.
// ApiKey options is the API key parent of the token. If you leave it empty the client API Key will be used.
func (c *Client) GenerateTenantToken(SearchRules map[string]interface{}, Options *TenantTokenOptions) (resp string, err error) {
	// Validate the arguments
	if SearchRules == nil {
		return "", fmt.Errorf("GenerateTenantToken: The search rules added in the token generation must be of type array or object")
	}
	if (Options == nil || Options.APIKey == "") && c.config.APIKey == "" {
		return "", fmt.Errorf("GenerateTenantToken: The API key used for the token generation must exist and be a valid Meilisearch key")
	}
	if Options != nil && !Options.ExpiresAt.IsZero() && Options.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("GenerateTenantToken: When the expiresAt field in the token generation has a value, it must be a date set in the future")
	}

	var secret string
	if Options == nil || Options.APIKey == "" {
		secret = c.config.APIKey
	} else {
		secret = Options.APIKey
	}

	// For HMAC signing method, the key should be any []byte
	hmacSampleSecret := []byte(secret)

	// Create the claims
	claims := TenantTokenClaims{}
	if Options != nil && !Options.ExpiresAt.IsZero() {
		claims.StandardClaims = jwt.StandardClaims{
			ExpiresAt: Options.ExpiresAt.Unix(),
		}
	}
	claims.APIKeyPrefix = secret[:8]
	claims.SearchRules = SearchRules

	// Create a new token object, specifying signing method and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)

	return tokenString, err
}

// This function allows the user to create a Key with an ExpiresAt in time.Time
// and transform the Key structure into a KeyParsed structure to send the time format
// managed by Meilisearch
func convertKeyToParsedKey(key Key) (resp KeyParsed) {
	resp = KeyParsed{Description: key.Description, Actions: key.Actions, Indexes: key.Indexes}

	// Convert time.Time to *string to feat the exact ISO-8601
	// format of Meilisearch
	if !key.ExpiresAt.IsZero() {
		const Format = "2006-01-02T15:04:05"
		timeParsedToString := key.ExpiresAt.Format(Format)
		resp.ExpiresAt = &timeParsedToString
	}
	return resp
}
