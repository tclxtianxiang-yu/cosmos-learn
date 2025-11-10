package types

// NewParams 创建一份新的 Params 实例。
func NewParams() Params {
	return Params{}
}

// DefaultParams 返回默认参数集。
func DefaultParams() Params {
	return NewParams()
}

// Validate 校验参数集合的有效性。
func (p Params) Validate() error {

	return nil
}
