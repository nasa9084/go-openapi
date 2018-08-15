package openapi

type Validater = validater

var (
	ValidateOASVersion     = validateOASVersion
	ValidateComponentKeys  = validateComponentKeys
	ReduceComponentKeys    = reduceComponentKeys
	ReduceComponentObjects = reduceComponentObjects
	HasDuplicatedParameter = hasDuplicatedParameter
	ValidateStatusCode     = validateStatusCode
	MustURL                = mustURL
	ValidateAll            = validateAll
)
