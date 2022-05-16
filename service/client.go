package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
)

// create a list of stylebook endpoints
var stylebookEndpoints = []string{
	"stylebooks",
	"configpacks",
	"jobs",
}

// URLResourceToBodyResource map of urlResource to bodyResource
var URLResourceToBodyResource = map[string]string{
	// FIXME: API Problem: Some API resources do not exactly match with that of body resources, esp., stylebook APIs
	"stylebooks":  "stylebook",
	"configpacks": "configpack",
	"jobs":        "job",
}

// NitroRequestParams is a struct to hold the parameters for a Nitro request
type NitroRequestParams struct {
	ResourcePath       string
	Method             string
	Headers            map[string]string
	Resource           string
	ResourceData       interface{}
	SuccessStatusCodes []int
	ActionParams       string
}

//NitroParams encapsulates options to create a NitroClient
type NitroParams struct {
	Host         string
	HostLocation string
	ID           string
	Secret       string
	CustomerID   string
	// SslVerify     bool
	// Timeout       int
	// RootCAPath    string
	// ServerName    string
	Headers       map[string]string
	LogLevel      string
	JSONLogFormat bool
}

//NitroClient has methods to configure the NetScaler
//It abstracts the REST operations of the NITRO API
type NitroClient struct {
	host         string
	hostLocation string
	id           string
	secret       string
	client       *http.Client
	customerID   string
	// sessionidMux sync.RWMutex
	// sessionid    string
	// timeout      int
	headers             map[string]string
	logger              hclog.Logger
	accessToken         string
	sessionID           string
	ActivityTimeout     int
	StylebookJobTimeout int
}

//NewNitroClientFromParams returns a usable NitroClient. Does not check validity of supplied parameters
func NewNitroClientFromParams(params NitroParams) (*NitroClient, error) {
	u, err := url.Parse(params.Host)
	if err != nil {
		return nil, fmt.Errorf("Supplied URL %s is not a URL", params.Host)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, fmt.Errorf("Supplied Host %s does not have a HTTP/HTTPS scheme", params.Host)
	}
	c := new(NitroClient)
	c.host = params.Host
	c.id = params.ID
	c.secret = params.Secret
	c.headers = params.Headers
	c.hostLocation = params.HostLocation
	c.customerID = params.CustomerID
	c.client = &http.Client{}
	c.ActivityTimeout = 120     // seconds to wait for activity to complete (manged_device)
	c.StylebookJobTimeout = 120 // seconds to wait for activity to complete (configpacks)

	// Get New Token
	if err := c.setNewToken(); err != nil {
		return nil, err
	}

	// Get Session ID for v2 API
	if err := c.setSessionID(); err != nil {
		return nil, err
	}

	level := hclog.LevelFromString(params.LogLevel)
	if level == hclog.NoLevel {
		level = hclog.Off
	}
	logLevel, ok := os.LookupEnv("NS_LOG")
	if ok {
		lvl := hclog.LevelFromString(logLevel)
		if lvl != hclog.NoLevel {
			level = lvl
		} else {
			log.Printf("nitro-go: NS_LOG not set to a valid log level (%s), defaulting to %d", logLevel, level)
		}
	}
	//c.logger = hclog.Default()
	c.logger = hclog.New(&hclog.LoggerOptions{
		Name:            "citrixadm-client",
		Level:           level,
		Color:           hclog.AutoColor,
		JSONFormat:      params.JSONLogFormat,
		IncludeLocation: true,
	})
	return c, nil
}

// GetSessionID returns a new access token
func (c *NitroClient) setSessionID() error {
	// endpoint := fmt.Sprintf("%s/nitro/v2/config/login", c.host)

	resource := "login"
	resourceData := map[string]interface{}{
		"ID":     c.id,
		"Secret": c.secret,
	}
	returnData, err := c.AddResource(resource, resourceData)
	if err != nil {
		return err
	}

	c.sessionID = returnData["login"].([]interface{})[0].(map[string]interface{})["sessionid"].(string)
	return nil

}

// setNewToken returns a new access token
func (c *NitroClient) setNewToken() error {
	endpoint := fmt.Sprintf("https://api-%s.cloud.com/cctrustoauth2/root/tokens/clients", c.hostLocation)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.id)
	data.Set("client_secret", c.secret)

	client := &http.Client{}
	r, err := http.NewRequest("POST", endpoint, strings.NewReader(data.Encode())) // URL-encoded payload
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Accept", "application/json")

	res, err := client.Do(r)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res.Status)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("NewToken details", string(body))
	response := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		panic(err)
	}
	// log.Println(sec["access_token"].(string))
	c.accessToken = response["access_token"].(string)
	return nil
}

// MakeNitroRequest makes a API request to the NetScaler
func (c *NitroClient) MakeNitroRequest(n NitroRequestParams) ([]byte, error) {
	var buff []byte
	var err error

	if n.Method == "POST" || n.Method == "PUT" {
		// if n.Resource present in UrlResourceToBodyResource map, then use the bodyResource
		if bodyResource, ok := URLResourceToBodyResource[n.Resource]; ok {
			n.Resource = bodyResource
		}
		payload := map[string]interface{}{n.Resource: n.ResourceData}
		if n.ActionParams != "" {
			payload["params"] = map[string]interface{}{
				"action": n.ActionParams,
			}
		}
		buff, err = JSONMarshal(payload)
		if err != nil {
			return nil, err
		}
		if n.Resource != "login" {
			log.Println("MakeNitroRequest payload", toJSONIndent(payload)) // print json converted payload
		}
	} else if n.Method == "GET" || n.Method == "DELETE" {
		buff = []byte{}
	}

	urlstr := fmt.Sprintf("%s/%s", c.host, n.ResourcePath)

	req, err := http.NewRequest(n.Method, urlstr, bytes.NewBuffer(buff))
	if err != nil {
		return nil, err
	}

	// Authenticate
	// For stylebook APIs use Cookie in Header else use Authorization header
	if strings.Contains(urlstr, "stylebook") {
		req.Header.Set("Cookie", fmt.Sprintf("SESSID=%s", c.sessionID))
		req.Header.Set("sessionId", fmt.Sprintf("%s", c.sessionID))
		req.Header.Set("isCloud", "true")
	} else if strings.Contains(urlstr, "massvc") {
		req.Header.Set("Authorization", fmt.Sprintf("CwsAuth bearer=%s", c.accessToken))
		req.Header.Set("isCloud", "true")
	} else {
		// login API. No headers
	}

	// Standard headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// Request defined headers
	// log.Println("MakeNitroRequest method: Request defined headers")
	for k, v := range n.Headers {
		// log.Println(k, v)
		req.Header.Set(k, v)
	}

	if n.Resource != "login" {
		log.Printf("MakeNitroRequest request:%s, url:%s, headers:%v", req.Method, req.URL, toJSONIndent(req.Header))
	}
	// log.Println("MakeNitroRequest payload", string(buff)) // print json converted payload

	resp, err := c.client.Do(req)

	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	var body []byte

	if statusCodeSuccess(n.SuccessStatusCodes, resp.StatusCode) {
		body, _ = ioutil.ReadAll(resp.Body)
		log.Println("MakeNitroRequest resopnse", n.Method, "url:", urlstr, "status:", resp.StatusCode)
		return body, nil
	}
	log.Println("MakeNitroRequest resopnse", n.Method, "url:", urlstr, "status:", resp.StatusCode)
	body, _ = ioutil.ReadAll(resp.Body)
	return []byte{}, errors.New("failed: " + resp.Status + " (" + string(body) + ")")
}

func statusCodeSuccess(slice []int, val int) bool {
	for _, item := range slice {
		if val == item {
			return true
		}
	}
	return false
}

// JSONMarshal https://stackoverflow.com/questions/28595664/how-to-stop-json-marshal-from-escaping-and/28596225#28596225
func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

// WaitForActivityCompletion waits for the activity to complete
func (c *NitroClient) WaitForActivityCompletion(activityID string, timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return errors.New("TIMEOUT: WaitForActivityCompletion ActivityID: " + activityID)
		}
		time.Sleep(time.Second * 5)

		returnData, err := c.GetResource("activity_status", activityID)
		if err != nil {
			return err
		}
		activityStatus := returnData["activity_status"].([]interface{})
		log.Println("Activity Status", toJSONIndent(activityStatus))

		// check for "is_last" key in activityStatus array and if it is true, then check for "status" key. And if the value of "status" key is "Completed" or "Failed" then return the activityStatus
		for _, activity := range activityStatus {
			if activity.(map[string]interface{})["is_last"].(string) == "true" {
				if activity.(map[string]interface{})["status"].(string) == "Completed" {
					return nil
				}
				if activity.(map[string]interface{})["status"].(string) == "Failed" {
					return fmt.Errorf(activity.(map[string]interface{})["status"].(string) + ": " + activity.(map[string]interface{})["message"].(string))
				}
			}
		}
	}
}

// WaitForStylebookJobCompletion waits for a stylebook job to complete
func (c *NitroClient) WaitForStylebookJobCompletion(jobID string, timeout time.Duration) error {
	log.Printf("WaitForStylebookJobCompletion: jobID: %s", jobID)
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return errors.New("TIMEOUT: WaitForStylebookJobCompletion jobID: " + jobID)
		}
		time.Sleep(time.Second * 5)

		returnData, err := c.GetResource("jobs", jobID)
		if err != nil {
			return err
		}

		// find jobStatus
		jobStatus := returnData["job"].(map[string]interface{})["status"].(string)

		log.Println("Activity Status", toJSONIndent(jobStatus))

		if jobStatus == "completed" {
			log.Println("Job Status", toJSONIndent(returnData))
			return nil
		} else if jobStatus == "failed" {
			log.Println("Job Status", toJSONIndent(returnData))

			progressInfo := returnData["job"].(map[string]interface{})["progress_info"].([]interface{})
			// collect all the "message" from "progress_info" array when "status" is "failed"
			// and concatenate them into one string

			var failedMessage string
			for _, progress := range progressInfo {
				progressMap := progress.(map[string]interface{})
				if progressMap["status"].(string) == "failed" {
					failedMessage = failedMessage + "\n" + progressMap["message"].(string)
				}
			}
			return errors.New(failedMessage)
		}
	}
}

func contains(slice []string, val string) bool {
	for _, item := range slice {
		if strings.EqualFold(item, val) {
			return true
		}
	}
	return false
}

func toJSONIndent(v interface{}) string {
	b, _ := json.MarshalIndent(v, "", "    ")
	return string(b)
}

// GetResource returns a resource
func (c *NitroClient) GetResource(resource string, resourceID string) (map[string]interface{}, error) {
	log.Println("GetResource method:", resource, resourceID)
	var returnData map[string]interface{}

	var resourcePath string
	if contains(stylebookEndpoints, resource) {
		resourcePath = fmt.Sprintf("stylebook/nitro/v2/config/%s/%s", resource, resourceID)
	} else {
		resourcePath = fmt.Sprintf("massvc/%s/nitro/v2/config/%s/%s", c.customerID, resource, resourceID)
	}

	n := NitroRequestParams{
		Resource:           resource,
		ResourcePath:       resourcePath,
		Method:             "GET",
		SuccessStatusCodes: []int{200},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return returnData, err
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &returnData)
		if err != nil {
			return returnData, err
		}
		log.Printf("GetResource response %v", toJSONIndent(returnData))
	}
	return returnData, nil
}

// GetAllResource returns all resources
func (c *NitroClient) GetAllResource(resource string) (map[string]interface{}, error) {
	log.Println("GetAllResource method:", resource)
	var returnData map[string]interface{}

	var resourcePath string
	if contains(stylebookEndpoints, resource) {
		resourcePath = fmt.Sprintf("stylebook/nitro/v2/config/%s", resource)
	} else {
		resourcePath = fmt.Sprintf("massvc/%s/nitro/v2/config/%s", c.customerID, resource)
	}

	n := NitroRequestParams{
		Resource:           resource,
		ResourcePath:       resourcePath,
		Method:             "GET",
		SuccessStatusCodes: []int{200},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return returnData, err
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &returnData)
		if err != nil {
			return returnData, err
		}
		log.Printf("GetAllResource response %v", toJSONIndent(returnData))
	}
	return returnData, nil
}

// AddResource adds a resource
func (c *NitroClient) AddResource(resource string, resourceData interface{}) (map[string]interface{}, error) {
	if resource != "login" {
		log.Println("AddResource method:", resource, resourceData)
	}
	var returnData map[string]interface{}

	var resourcePath string
	if resource == "login" {
		resourcePath = fmt.Sprintf("nitro/v2/config/%s", resource)
	} else if contains(stylebookEndpoints, resource) {
		resourcePath = fmt.Sprintf("stylebook/nitro/v2/config/%s", resource)
	} else {
		resourcePath = fmt.Sprintf("massvc/%s/nitro/v2/config/%s", c.customerID, resource)
	}

	n := NitroRequestParams{
		Resource:           resource,
		ResourcePath:       resourcePath,
		ResourceData:       resourceData,
		Method:             "POST",
		SuccessStatusCodes: []int{200, 201, 202},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return returnData, err
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &returnData)
		if err != nil {
			return returnData, err
		}
		log.Printf("AddResource response %v", toJSONIndent(returnData))
	}
	return returnData, nil
}

// AddResourceWithActionParams adds a resource with action params
func (c *NitroClient) AddResourceWithActionParams(resource string, resourceData interface{}, actionParam string) (map[string]interface{}, error) {
	log.Println("AddResourceWithActionParams method:", resource, resourceData, actionParam)
	var returnData map[string]interface{}

	var resourcePath string
	if contains(stylebookEndpoints, resource) || strings.Contains(resource, "stylebook") {
		// https://adm.cloud.com/stylebook/nitro/v2/config/stylebooks/actions/import
		// https://adm.cloud.com/stylebook/nitro/v2/config/stylebooks/{{stylebook_namespace}}/{{stylebook_version}}/{{stylebook_name}}/actions/update
		resourcePath = fmt.Sprintf("stylebook/nitro/v2/config/%s/actions/%s", resource, actionParam)
		resource = actionParam
		actionParam = "" // reset actionParam so that it is not added to the requestPayload in MakeNitroRequest()
	} else {
		resourcePath = fmt.Sprintf("massvc/%s/nitro/v2/config/%s", c.customerID, resource)
	}

	n := NitroRequestParams{
		Resource:           resource,
		ResourcePath:       resourcePath,
		ResourceData:       resourceData,
		ActionParams:       actionParam,
		Method:             "POST",
		SuccessStatusCodes: []int{200, 201},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return returnData, err
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &returnData)
		if err != nil {
			return returnData, err
		}
		log.Printf("AddResourceWithActionParams response %v", toJSONIndent(returnData))
	}
	return returnData, nil
}

// UpdateResource updates a resource
func (c *NitroClient) UpdateResource(resource string, resourceData interface{}, resourceID string) (map[string]interface{}, error) {
	log.Println("UpdateResource method:", resource, resourceData, resourceID)
	var returnData map[string]interface{}

	var resourcePath string
	if contains(stylebookEndpoints, resource) {
		resourcePath = fmt.Sprintf("stylebook/nitro/v2/config/%s/%s", resource, resourceID)
	} else {
		resourcePath = fmt.Sprintf("massvc/%s/nitro/v2/config/%s/%s", c.customerID, resource, resourceID)
	}

	n := NitroRequestParams{
		Resource:           resource,
		ResourcePath:       resourcePath,
		ResourceData:       resourceData,
		Method:             "PUT",
		SuccessStatusCodes: []int{200, 201, 202},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return returnData, err
	}

	if len(body) == 0 {
		err = json.Unmarshal(body, &returnData)
		if err != nil {
			return returnData, err
		}
		log.Printf("UpdateResource response %v", toJSONIndent(returnData))
	}
	return returnData, nil
}

// DeleteResource deletes a resource
func (c *NitroClient) DeleteResource(resource string, resourceID string) (map[string]interface{}, error) {
	log.Println("DeleteResource method:", resource, resourceID)
	var returnData map[string]interface{}

	var resourcePath string
	if contains(stylebookEndpoints, resource) {
		resourcePath = fmt.Sprintf("stylebook/nitro/v2/config/%s/%s", resource, resourceID)
	} else {
		resourcePath = fmt.Sprintf("massvc/%s/nitro/v2/config/%s/%s", c.customerID, resource, resourceID)
	}

	n := NitroRequestParams{
		Resource:           resource,
		ResourcePath:       resourcePath,
		Method:             "DELETE",
		SuccessStatusCodes: []int{200, 202, 204},
	}

	body, err := c.MakeNitroRequest(n)
	if err != nil {
		return returnData, err
	}

	if len(body) > 0 {
		err = json.Unmarshal(body, &returnData)
		if err != nil {
			return returnData, err
		}
		// log.Printf("delete response %v", deleteResponse)
		log.Printf("DeleteResource response %v", toJSONIndent(returnData))
	}
	return returnData, nil
}
