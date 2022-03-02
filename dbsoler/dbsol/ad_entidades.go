package dbsol

type Entidades struct {
	IdEntidade               string `length:"16"`
	NomeEntidade             string `length:"50"`
	TipoEntidade             int16  `length:"2"`
	TipoDestino              int16  `length:"2"`
	IdAplTrans               string `length:"16"`
	CdMetodoRoteamento       int16  `length:"2"`
	NomeBibliotecaRoteamento string `length:"16"`
	NuIniId                  int16  `length:"2"`
	IdEvento                 byte   `length:"2"`
}
