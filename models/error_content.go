package models

// API error codes
const (
	AcceptLanguageHeaderError = "ErrorAcceptLanguageHeader"
	AreaDataIdGetError = "ErrorRetrievingAreaCode"
	MarshallingAreaDataError = "ErrorMarshallingAreaData"
	MarshallingAreaRelationshipsError = "ErrorMarshallingAreaRelationshipData"
)

// API error descriptions
const (
	AcceptLanguageHeaderNotFoundDescription = "accept language header not found"
	AcceptLanguageHeaderInvalidDescription = "accept language header invalid"
	AreaDataGetErrorDescription = "area code not found"
)
