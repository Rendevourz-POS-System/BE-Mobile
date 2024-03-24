package mail

type Service interface {
	ExecuteSendEmail(
		subject, content string,
		to, cc, bcc []string,
		attach string,
	) error
}
