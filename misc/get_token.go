// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	// endpoint := fmt.Sprintf("%s%s", c.host, "/nitro/v2/config/login")
	endpoint := "https://api-us.cloud.com/cctrustoauth2/root/tokens/clients"

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", "987a9390-6a65-4f78-8587-790feb82d63a")
	data.Set("client_secret", "BZnvbtdZaWJ3jYJwYtnsCw==")

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
	log.Println(string(body))
	// convert body as map[string]string and return access token as string
	sec := map[string]interface{}{}
	if err := json.Unmarshal([]byte(body), &sec); err != nil {
		panic(err)
	}
	log.Println(sec["access_token"].(string))

}

// Sample Output

// ❯ go run ./misc/get_token.go
// 2022/03/21 15:45:18 200 OK
// 2022/03/21 15:45:18
//                         {
//                             "token_type": "bearer",
//                             "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiNjJiYjQ0OGItMjA3MS00ZWUxLTkwMDYtODY0YTc3YTA5ZDJhIiwicHJpbmNpcGFsIjoic3VtYW50aC5saW5nYXBwYUBjaXRyaXguY29tIiwiYWNjZXNzX3Rva2VuX3Njb3BlIjoiIiwicmVmcmVzaF90b2tlbiI6IiIsImFjY2Vzc190b2tlbiI6IiIsImRpc3BsYXlOYW1lIjoiU3VtYW50aCBMaW5nYXBwYSIsInJlZnJlc2hfZXhwaXJhdGlvbiI6IjE2NDc5MDA5MTgwMzMiLCJjdXN0b21lcnMiOiJbe1wiQ3VzdG9tZXJJZFwiOlwidmJkM25tMzJmbjV3XCIsXCJHZW9cIjpcIkFQLVNcIn1dIiwiZW1haWxfdmVyaWZpZWQiOiJUcnVlIiwiY3R4X2F1dGhfYWxpYXMiOiI5ODdhOTM5MC02YTY1LTRmNzgtODU4Ny03OTBmZWI4MmQ2M2EiLCJpZHAiOiJjaXRyaXhzdHMiLCJuYW1lIjoiU3VtYW50aCBMaW5nYXBwYSIsInN1YiI6IjYyYmI0NDhiLTIwNzEtNGVlMS05MDA2LTg2NGE3N2EwOWQyYSIsImVtYWlsIjoic3VtYW50aC5saW5nYXBwYUBjaXRyaXguY29tIiwiYW1yIjoiW1wiY2xpZW50XCJdIiwiZGlzY292ZXJ5Ijoie1wiSXNzdWVyXCI6XCJodHRwczovL3RydXN0LXVzLmNpdHJpeHdvcmtzcGFjZXNhcGkubmV0XCJ9IiwiY3R4X3VzZXIiOiJ7XCJPaWRcIjpcIk9JRDovY2l0cml4LzYyYmI0NDhiLTIwNzEtNGVlMS05MDA2LTg2NGE3N2EwOWQyYVwifSIsImN0eF9kaXJlY3RvcnlfY29udGV4dCI6IntcIklkZW50aXR5UHJvdmlkZXJcIjpcIkNpdHJpeFwifSIsIm5iZiI6MTY0Nzg1NzcxOCwiZXhwIjoxNjQ3ODYxMzE4LCJpYXQiOjE2NDc4NTc3MTgsImlzcyI6ImN3cyJ9.PhMPyF0TpLk2__UDKGq6jWQdlMupBhlK5bbf2JlikwjZQzSbxhxh7b9lvVmbBVJFQrCabGfjAj18pxJKW_IQwRxDB4rxEMdPr0Fr7bW7xQPdm5GYZf8ZkIQ9gNohSz1llRtD1GLH4IteL6YWaCN8y0wJ1MhNRG1cZGXQP68qNwim5GzqDwtbY2zGAJqpfdHOH6mQz2W28rr3Tsn4S6aQHkS1GaaYbhSHDmVfnpy0J9pJ2QYRdYqivc-V-1n0Exwgd0tw9w5mhnVVyG8AczT2asMpSC8_LPm3M7Yow-GZ9G_f6N0gK9UPbQPjovLrb7eHMv6iAVFCB6b56bS4AxViVQ",
//                             "expires_in": "3600"
//                         }
