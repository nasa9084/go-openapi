package openapi

type Validater = validater

var (
	ValidateOASVersion     = validateOASVersion
	ValidateComponentKeys  = validateComponentKeys
	ReduceComponentKeys    = reduceComponentKeys
	ReduceComponentObjects = reduceComponentObjects
	HasDuplicatedParameter = hasDuplicatedParameter
	MustURL                = mustURL
	ValidateAll            = validateAll
)
