package models

// API error codes
const (
	AcceptLanguageHeaderError          = "ErrorAcceptLanguageHeader"
	AreaDataIdGetError                 = "ErrorRetrievingAreaCode"
	AreaDataIdUpsertError              = "AreaDataIdUpsertError"
	AncestryDataGetError               = "ErrorRetrievingAncestryData"
	MarshallingAreaDataError           = "ErrorMarshallingAreaData"
	MarshallingAreaRelationshipsError  = "ErrorMarshallingAreaRelationshipData"
	InvalidAreaCodeError               = "InvalidAreaCode"
	AreaNameNotProvidedError           = "AreaNameNotProvidedError"
	InvalidAreaTypeError               = "InvalidAreaType"
	AreaNameActiveFromNotProvidedError = "AreaNameActiveFromNotProvidedError"
	AreaNameActiveToNotProvidedError   = "AreaNameActiveToNotProvidedError"
	AreaNameDetailsNotProvidedError    = "AreaNameDetailsNotProvidedError"
	BodyCloseError                     = "BodyCloseError"
	BodyReadError                      = "RequestBodyReadError"
	JSONUnmarshalError                 = "JSONUnmarshalError"
)

// API error descriptions
const (
	AcceptLanguageHeaderNotFoundDescription       = "accept language header not found"
	AcceptLanguageHeaderInvalidDescription        = "accept language header invalid"
	AreaDataGetErrorDescription                   = "area code not found"
	BodyClosedFailedDescription                   = "the request body failed to close"
	BodyReadFailedDescription                     = "endpoint returned an error reading the request body"
	ErrorUnmarshalFailedDescription               = "failed to unmarshal the request body"
	InvalidAreaCodeErrorDescription               = "the area code could not be validated"
	AreaNameDetailsNotProvidedErrorDescription    = "required field area_name not provided"
	AreaNameNotProvidedErrorDescription           = "required field area_name.name not provided"
	AreaNameActiveFromNotProvidedErrorDescription = "required field area_name.active_from not provided"
	AreaNameActiveToNotProvidedErrorDescription   = "required field area_name.active_to not provided"
	InvalidAreaTypeErrorDescription               = "failed to derive area type from area code"
)
