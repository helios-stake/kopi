package types

type DexDenoms []*DexDenom

func (dd DexDenoms) Get(denom string) *DexDenom {
	for _, d := range dd {
		if d.Name == denom {
			return d
		}
	}

	return nil
}
