package dbsol

type Variaveis struct {
	CdEscopo  int32  `length:"4"`
	IdEscopo1 string `length:"16"`
	IdEscopo2 string `length:"16"`
	Nome      string `length:"64"`
	Tipo      int32  `length:"4"`
	Valor     string `length:"256"`
}
