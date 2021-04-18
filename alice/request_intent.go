package alice

type Intent struct {
	Slots map[string]*IntentSlot `json:"slots"`
}

type IntentSlot struct {
	Type EntityType `json:"type"`
	Tokens SlotTokens `json:"tokens"`
	Value interface{} `json:"value"`
}

type SlotTokens struct {
	Start int `json:"start"`
	End int `json:"end"`
}