package mock

type MockStore struct {
	Expressions map[string]string
	Err         error
}

func (m *MockStore) GetExpressions() (map[string]string, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Expressions, nil
}
