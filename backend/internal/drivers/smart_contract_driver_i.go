package drivers

type SmartContractDriverI interface {
	CreateTransaction(
		iFunctionName string,
		iArgs ...string,
	) ([]byte, error)

	Query(
		iFunctionName string,
		iArgs ...string,
	) ([]byte, error)
}
