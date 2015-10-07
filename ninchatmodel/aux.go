package ninchatmodel

// Aux
type Aux struct {
	m map[interface{}]interface{}
}

func (aux *Aux) GetAux(key interface{}) interface{} {
	return aux.m[key]
}

func (aux *Aux) SetAux(key, value interface{}) {
	if aux.m == nil {
		aux.m = make(map[interface{}]interface{})
	}
	aux.m[key] = value
}
