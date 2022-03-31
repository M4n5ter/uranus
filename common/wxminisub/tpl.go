package wxminisub

//订单支付成功
const OrderPaySuccessTemplateID = "QIJPmfxaNqYzSjOlXGk1T6Xfw94JwbSPuOd3u_hi3WE"

type OrderPaySuccessDataParam struct {
	Sn             string
	PayTotal       string
	DepartPosition string
	DepartTime     string
	ArrivePosition string
	ArriveTime     string
}

func OrderPaySuccessData(params OrderPaySuccessDataParam) map[string]string {
	return map[string]string{
		"character_string6": params.Sn,
		"amount2":           params.PayTotal,
		"departPosition":    params.DepartPosition,
		"departTime":        params.DepartTime,
		"arrivePosition":    params.ArrivePosition,
		"arriveTime":        params.ArriveTime,
	}
}
