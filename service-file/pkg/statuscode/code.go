package statuscode

const (
	Success       = 200
	Error         = 500
	InvalidParams = 400

	EmailNotRegister   = 10002
	NeuronNotFound     = 10008
	UserHasNoNeuron    = 10009
	NameExisted        = 10010
	NameSameAsOriginal = 10011

	InvalidDocumentOwnership = 10012
	DocumentNotFound         = 10013
)
