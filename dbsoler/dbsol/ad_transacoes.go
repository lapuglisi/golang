package dbsol

/*
defNomeCampo("nm_transacao", 16, TD_STRING);
defNomeCampo("de_transacao", 50, TD_STRING);
defNomeCampo("_alinhamento", 2, TD_SHORT);
defNomeCampo("nu_tam_max_msg", 4, TD_INT);
defNomeCampo("nu_itens_fila", 4, TD_INT);
*/

type Transacoes struct {
	NomeTransacao string `length:"16"`
	DescTransacao string `length:"50"`
	AL            int16  `length:"2"`
	NuTamMaxMsg   int32  `length:"4"`
	NuItensFila   int32  `length:"4"`
}
