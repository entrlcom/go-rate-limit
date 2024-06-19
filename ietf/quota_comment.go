package ietf

type QuotaComment struct {
	token string
	value string
}

func (x QuotaComment) Token() string {
	return x.token
}

func (x QuotaComment) Value() string {
	return x.value
}

func NewQuotaComment(token, value string) QuotaComment {
	return QuotaComment{
		token: token,
		value: value,
	}
}
