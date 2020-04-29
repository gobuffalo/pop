package cempty

type mockTranslator struct{}

func (mockTranslator) Name() string {
	return "test"
}
