package model

type ExcelRow struct {
	Codice           string `column:"Codice" required:"true"`
	Data             string `column:"Data" required:"true"`
	Orario           string `column:"Orario" required:"true"`
	Educatore        string `column:"Educatore" required:"true"`
	Aula             string `column:"Aula" required:"true"`
	Attivita         string `column:"Attività" required:"true"`
	LinguaAttivita   string `column:"Lingua dell'attività"`
	NotaPrenotazione string `column:"Nota prenotazione"`
	NotaOperatore    string `column:"Nota operatore"`
	TipologiaScuola  string `column:"Tipologia Scuola"`
	NomeScuola       string `column:"Nome scuola"`
	Classe           string `column:"Classe"`
	Sezione          string `column:"Sezione"`
	Paganti          string `column:"Paganti"`
	Gratuiti         string `column:"Gratuiti"`
	Accompagnatori   string `column:"Accompagnatori"`
	Tipologia        string `column:"Tipologia"`
	Bus              string `column:"Bus"`
	Acconti          string `column:"Acconti"`
	StatoAcconti     string `column:"Stato acconti"`

	Other map[string]interface{}
}
