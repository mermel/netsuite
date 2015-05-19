package netsuite

import (
  "net/url"
  "net/http"
  "fmt"
  "bytes"
  "io/ioutil"
  "os"

  "github.com/VioletGrey/error-handler"
)

type Response struct {
  Success bool
  Message, Id string
}

func GetNetsuiteUserRequest(id string, email string) (responseBody []byte) {
  var params map[string]string
  params = make(map[string]string)
  params["recordtype"] = "customer"

  if id != "" {
    params["id"] = id
  }
  if email != "" {
    params["email"] = email
  }

  Url := NetsuiteUrlWithParams("user", params)
  req, err := http.NewRequest("GET", Url.String(), nil)
  vgError.FailOnError(err, "Failure: Failed GET request to NetSuite", Url.String())

  resp := NetsuiteClient(req)
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  return body
}

func NetsuiteRequest(requestType string, scriptType string, requestBody []byte)(body []byte){
  Url := NetsuiteUrl(scriptType)
  req, err := http.NewRequest(requestType, Url.String(), bytes.NewBuffer(requestBody))

  resp := NetsuiteClient(req)
  reqBody, err := ioutil.ReadAll(req.Body)
  req.Body.Close()
  vgError.FailOnError(err, "Failure: Failed request", string(reqBody))
  defer resp.Body.Close()

  body, _ = ioutil.ReadAll(resp.Body)
  fmt.Println("response Body:", string(body))
  return body
}

func NetsuiteClient(req *http.Request) (resp *http.Response) {
  reqWithHeader := NetsuiteHeader(req)
  client := &http.Client{}
  resp, err := client.Do(reqWithHeader)
  vgError.FailOnError(err, "Unable to start NetSuite Client", "netsuiteClient")
  return resp
}

func NetsuiteHeader(req *http.Request) (request *http.Request) {
  authHeaderString := fmt.Sprintf("NLAuth nlauth_account=%s, nlauth_email=%s, nlauth_signature=%s, nlauth_role=%s", os.Getenv("NLAUTH_ACCOUNT"), os.Getenv("NLAUTH_EMAIL"), os.Getenv("NLAUTH_SIGNATURE"), os.Getenv("NLAUTH_ROLE"))
  req.Header.Set("Authorization", authHeaderString)
  req.Header.Set("Content-Type", "application/json")
  return req
}

func NetsuiteUrl(scriptType string) (Url *url.URL) {
  Urll := NetsuiteBaseUrl()
  parameters := url.Values{}
  switch scriptType {
    case "user":
      parameters.Add("script", "10")
    case "order":
      parameters.Add("script", "15")
  }
  parameters.Add("deploy", "1")
  Urll.RawQuery = parameters.Encode()

  return Urll
}

func NetsuiteUrlWithParams(scriptType string, params map[string]string) (Url *url.URL) {
  Urll := NetsuiteBaseUrl()
  parameters := url.Values{}
  switch scriptType {
    case "user":
      parameters.Add("script", "10")
    case "order":
      parameters.Add("script", "15")
  }
  parameters.Add("deploy", "1")
  for key, value := range params {
    parameters.Add(key, value)
  }
  Urll.RawQuery = parameters.Encode()

  fmt.Printf("Encoded URL is %q\n", Urll.String())
  return Urll
}

func NetsuiteBaseUrl() (Url *url.URL) {
  baseUrl := "https://rest.na1.netsuite.com/app/site/hosting/restlet.nl"
  Url, err := url.Parse(baseUrl)
  vgError.FailOnError(err, "Unable to parse NetSuite base URL", "invalidNetsuiteUrl")
  return Url
}

