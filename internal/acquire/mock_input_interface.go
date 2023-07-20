package acquire

type MockInputInterface struct {
}

func (m MockInputInterface) GetInput(request InputRequest) (InputResponse, error) {
	return InputResponse{}, nil
}
