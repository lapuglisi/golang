package dbsol

/*
defNomeCampo("id_forma_comunicacao", 16, TD_STRING);
defNomeCampo("de_forma_comunicacao", 50, TD_STRING);
defNomeCampo("nm_biblioteca_comunicacao", 16, TD_STRING);
defNomeCampo("nm_arquivo_configuracao", 16, TD_STRING);
defNomeCampo("nm_arquivo_identificacao", 16, TD_STRING);
*/

type FormasComunicao struct {
	IdFormaComunicao          string `length:"16"`
	DescFormaComunicacao      string `length:"50"`
	NomeBibliotecaComunicacao string `length:"16"`
	NomeArquivoConfiguracao   string `length:"16"`
	NomeArquivoIdentificacao  string `length:"16"`
}
