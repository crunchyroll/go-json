// Package auto is an auto-generated schema for input JSON. Do not edit it by hand or check it in to
// source control. It is generated from
package auto

type Glossary_GlossDiv_GlossDiv struct {
	B string `json:"B"`
}

type GlossDef struct {
	Para         string   `json:"para"`
	GlossSeeAlso []string `json:"GlossSeeAlso"`
}

type GlossEntry struct {
	GlossSee  string    `json:"GlossSee"`
	ID        string    `json:"ID"`
	SortAs    string    `json:"SortAs"`
	GlossTerm string    `json:"GlossTerm"`
	Acronym   string    `json:"Acronym"`
	Abbrev    string    `json:"Abbrev"`
	GlossDef  *GlossDef `json:"GlossDef"`
}

func (x *GlossEntry) GetGlossDef() *GlossDef {
	if x.GlossDef == nil {
		return &GlossDef{}
	}

	return x.GlossDef
}

type Glossary_GlossDiv_GlossList struct {
	GlossEntry *GlossEntry `json:"GlossEntry"`
}

func (x *Glossary_GlossDiv_GlossList) GetGlossEntry() *GlossEntry {
	if x.GlossEntry == nil {
		return &GlossEntry{}
	}

	return x.GlossEntry
}

type Glossary_GlossDiv struct {
	Title     string                       `json:"title"`
	GlossDiv  *Glossary_GlossDiv_GlossDiv  `json:"GlossDiv"`
	GlossList *Glossary_GlossDiv_GlossList `json:"GlossList"`
}

func (x *Glossary_GlossDiv) GetGlossList() *Glossary_GlossDiv_GlossList {
	if x.GlossList == nil {
		return &Glossary_GlossDiv_GlossList{}
	}

	return x.GlossList
}

func (x *Glossary_GlossDiv) GetGlossDiv() *Glossary_GlossDiv_GlossDiv {
	if x.GlossDiv == nil {
		return &Glossary_GlossDiv_GlossDiv{}
	}

	return x.GlossDiv
}

type Glossary_GlossList struct {
	A float64 `json:"a"`
}

type Glossary struct {
	GlossDiv   *Glossary_GlossDiv  `json:"GlossDiv"`
	EmptyArray []interface{}       `json:"empty_array"`
	Number     float64             `json:"number"`
	GlossList  *Glossary_GlossList `json:"gloss_list"`
	Title      string              `json:"title"`
}

func (x *Glossary) GetGlossDiv() *Glossary_GlossDiv {
	if x.GlossDiv == nil {
		return &Glossary_GlossDiv{}
	}

	return x.GlossDiv
}

func (x *Glossary) GetGlossList() *Glossary_GlossList {
	if x.GlossList == nil {
		return &Glossary_GlossList{}
	}

	return x.GlossList
}

type Auto struct {
	Glossary *Glossary `json:"glossary"`
}

func (x *Auto) GetGlossary() *Glossary {
	if x.Glossary == nil {
		return &Glossary{}
	}

	return x.Glossary
}
