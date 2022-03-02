package dbsol

/*
defNomeCampo("nm_ambiente", 16, TD_STRING);
defNomeCampo("nm_tabela", 16, TD_STRING);
defNomeCampo("de_tabela", 60, TD_STRING);
defNomeCampo("nu_tam_chave", 4, TD_INT);
defNomeCampo("nu_tam_max_conteudo", 4, TD_INT);
defNomeCampo("cd_opcoes", 4, TD_INT);
*/

type Tabelas struct {
	Ambiente              string `length:"16"`
	Tabela                string `length:"16"`
	Descricao             string `length:"60"`
	TamanhoChave          int32  `length:"4"`
	TamanhoMaximoConteudo int32  `length:"4"`
	CdOpcoes              int32  `length:"4"`
}
