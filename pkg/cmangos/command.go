package cmangos

import "encoding/xml"

type ExecCmdRequest struct {
	Command string
}

type ExecCmdResponse struct {
	XMLName  xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body *ExecCmdResponseBody
}

type ExecCmdResponseBody struct {
	XMLName      			xml.Name `xml:"Body"`
	ExecCmdResponseText		ExecCmdResponseText
	Fault		 			*Fault
}

type ExecCmdResponseText struct {
	XMLName		xml.Name 	`xml:"executeCommandResponse"`
	Result		string 		`xml:"result"`
}

type Fault struct {
	XMLName     xml.Name `xml:"Fault"`
	Faultcode   string   `xml:"faultcode"`
	Faultstring string   `xml:"faultstring"`
}

var execCmdTemplate = `
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope soap:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/" xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:soapenc="http://schemas.xmlsoap.org/soap/encoding/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
	<soap:Body>
		<executeCommand xmlns="urn:MaNGOS">
			<command xsi:type="xsd:string">{{ .Command }}</command>
		</executeCommand>
	</soap:Body>
</soap:Envelope>
`