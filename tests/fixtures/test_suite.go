package fixtures

type TestSuite struct {
	Solana *SolanaTestSuite
}

func NewTestSuite() *TestSuite {
	t := &TestSuite{
		Solana: NewSolanaTestSuite(),
	}
	return t
}
