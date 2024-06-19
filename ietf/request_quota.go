package ietf

type RequestQuota struct {
	quotaUnits int64
}

func (x RequestQuota) QuotaUnits() int64 {
	return x.quotaUnits
}

func NewRequestQuota(quotaUnits int64) RequestQuota {
	return RequestQuota{quotaUnits: quotaUnits}
}
