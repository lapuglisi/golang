package dbsol

/*
defNomeCampo("nm_maquina", 16, TD_STRING);
defNomeCampo("id_processo", 16, TD_STRING);
defNomeCampo("nm_ambiente", 16, TD_STRING);
defNomeCampo("tp_processo", 2, TD_SHORT);
defNomeCampo("id_aplicacao", 16, TD_STRING);
defNomeCampo("nm_programa", 16, TD_STRING);
defNomeCampo("de_processo", 50, TD_STRING);
defNomeCampo("ha_processo_com", 1, TD_CHAR);
defNomeCampo("_alinhamento", 3, TD_CHAR);
defNomeCampo("nu_tarefas_processamento", 4, TD_INT);
defNomeCampo("nu_itens_fila", 4, TD_INT);
defNomeCampo("cd_status", 2, TD_SHORT);
*/

type Processos struct {
	NomeMaquina            string `length:"16"`
	IdProcesso             string `length:"16"`
	NomeAmbiente           string `length:"16"`
	TipoProcesso           int16  `length:"2"`
	IdAplicacao            string `length:"16"`
	NomePrograma           string `length:"16"`
	DescProcesso           string `length:"50"`
	HaProcessoCom          int32  `length:"4"`
	NuTarefasProcessamento int32  `length:"4"`
	NuItensFila            int32  `length:"4"`
	CdStatus               int32  `length:"4"`
}
