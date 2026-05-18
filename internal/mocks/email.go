package mocks

type EmailTest struct{}

// SendEmail implements [email.EmailService].
func (e EmailTest) SendEmail(to string) error {
	return nil
}

func NewTestEmailService() EmailTest {
	return EmailTest{}
}
