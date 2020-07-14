package cmangos
import (
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

// SoapClient provides methods for interacting with a cmangos server
// via its soap API
type SoapClientOpts struct {
	// Username to authenticate
	Username  	string
	// Password to authenticate
	Password 	string
	// Address of the cmangos soap api server
	Address		string
}

type soapClient struct {
	authHeader	string
	address		string
}

type SoapClient interface {
	SendExecCmd(*ExecCmdRequest) (*ExecCmdResponse, error)
}

// Validate ensures that the client has all parameters filled for connecting
// to the server
func (c *SoapClientOpts) Validate() error {
	if c.Username == "" {
		return errors.New("client username must have a value")
	}

	if c.Password == "" {
		return errors.New("client password must have a value")
	}

	if c.Address == "" {
		return errors.New("client address must have a value")
	}

	return nil
}

func NewClient(opts *SoapClientOpts) (SoapClient, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &soapClient{
		authHeader: "Basic " + basicAuth(opts.Username, opts.Password),
		address: opts.Address,
	}, nil
}

func (c *soapClient) SendExecCmd(req *ExecCmdRequest) (*ExecCmdResponse, error) {

	template, err := template.New("InputRequest").Parse(execCmdTemplate)
	if err != nil {
		return nil, fmt.Errorf("error while marshling object. %s \n",err.Error())
	}

	doc := &bytes.Buffer{}
	err = template.Execute(doc, req)
	if err != nil {
		return nil, fmt.Errorf("template.Execute error %s \n",err.Error())
	}

	buffer := &bytes.Buffer{}
	encoder := xml.NewEncoder(buffer)
	err = encoder.Encode(doc.String())
	if err != nil {
		return nil, fmt.Errorf("encoder.Encode error. %s \n",err.Error())
	}

	r, err := http.NewRequest(http.MethodPost, c.address, bytes.NewBuffer([]byte(doc.String())))
	if err != nil {
		return nil, fmt.Errorf("error making a request. %s \n", err.Error())
	}

	r.Header.Add("Authorization", c.authHeader)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("error sending request. %s \n", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error getting request body")
	}

	defer resp.Body.Close()

	cmdResp := &ExecCmdResponse{}
	err = xml.Unmarshal(body, &cmdResp)
	if err != nil {
		return nil, fmt.Errorf("unable to parse command response\nbody: %s\nerror: %s\n", body, err)
	}

	return cmdResp, nil
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}