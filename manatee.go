package manatee

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"

	"github.com/hooklift/gowsdl/soap"

	"github.com/hooklift/gowsdl"
)

// Manatee is a soap client.
type Manatee struct {
	wsdlURL string
	wsdl    *gowsdl.WSDL
	client  *soap.Client
}

// Call sends the requests data by calling SOAP action
func (c *Manatee) Call(soapAction string, request, response interface{}) error {
	err := c.init()
	if err != nil {
		return err
	}

	return c.client.Call(c.getActionEndpoint(soapAction), request, response)
}

func (c *Manatee) init() error {
	err := c.parseWSDL()
	if err != nil {
		return err
	}

	c.client = soap.NewClient(c.getServiceURL())

	return nil
}

func (c *Manatee) parseWSDL() error {
	if c.wsdl == nil {
		res, err := http.Get(c.wsdlURL)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		byteValue, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}

		return xml.Unmarshal(byteValue, &c.wsdl)
	}

	return nil
}

func (c *Manatee) getServiceURL() string {
	return c.wsdl.Service[0].Ports[0].SOAPAddress.Location
}

func (c *Manatee) getActionEndpoint(action string) string {
	for _, binding := range c.wsdl.Binding {
		for _, operation := range binding.Operations {
			if operation.Name == action {
				return operation.SOAPOperation.SOAPAction
			}
		}
	}

	return ""
}

// NewClient generates a new manatee client with given wsdlURL url.
func NewClient(wsdl string) *Manatee {
	return &Manatee{wsdlURL: wsdl}
}
