package types

var (
	automationFeeCondition = uint64(50_000)
	automationFeeAction    = uint64(200_000)
)

func DefaultParams() Params {
	return Params{
		AutomationFeeCondition: automationFeeCondition,
		AutomationFeeAction:    automationFeeAction,
	}
}

func NewParams() Params {
	return DefaultParams()
}

func (p Params) Validate() error {
	return nil
}
