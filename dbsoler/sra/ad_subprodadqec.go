package sra

import "time"

/*
	defNomeCampo("cd_pessoa_estabelecimento", 4, TD_INT);
	defNomeCampo("cd_produto_cartao_adquirente", 4, TD_INT);
	defNomeCampo("dt_ini_vigencia", 4, TD_DATA);
	defNomeCampo("dt_fim_vigencia", 4, TD_DATA);
	defNomeCampo("cd_status", 2, TD_SHORT);
*/

type SubProdAdqEc struct {
	CdPessoaEstabelecimento   int32     `length:"4"`
	CdProdutoCartaoAdquirente int32     `length:"4"`
	DataInicioVigencia        time.Time `length:"4"`
	DataFimVigencia           time.Time `length:"4"`
	CdStatus                  int16     `length:"2"`
}
