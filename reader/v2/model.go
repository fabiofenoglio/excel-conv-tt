package reader

import "time"

type Row struct {
	// campi calcolati in fase di lettura
	ID        int
	rowNumber uint

	// campi che verranno esposti così come letti, senza conversione

	BookingCode          string `column:"codice" required:"true"`
	Operator             string `column:"educatore" required:"true"`
	Room                 string `column:"aula" required:"true"`
	Activity             string `column:"evento" required:"true"`
	ActivityLanguage     string `column:"lingua dell'attività"`
	BookingNote          string `column:"nota prenotazione"`
	OperatorNote         string `column:"nota operatore"`
	SchoolType           string `column:"tipologia scuola"`
	SchoolName           string `column:"nome scuola"`
	Class                string `column:"classe"`
	ClassSection         string `column:"sezione"`
	ClassTeacher         string `column:"insegnante"`
	ClassRefEmail        string `column:"email referente"`
	Type                 string `column:"tipologia"`
	Bus                  string `column:"bus"`
	PaymentAdvance       string `column:"acconti"`
	PaymentAdvanceStatus string `column:"stato acconti"`
	SpecialProjectName   string `column:"nome progetto speciale"`

	// campi che verranno convertiti prima di essere esposti

	DateRawString            string `column:"data" required:"true"`
	TimesRawString           string `column:"orario" required:"true"`
	NumPayingRawString       string `column:"paganti"`
	NumFreeRawString         string `column:"gratuiti"`
	NumAccompanyingRawString string `column:"accompagnatori"`
	ConfirmedRawString       string `column:"confermata"`

	// campi che saranno esposti dopo apposita conversione

	Date                        time.Time
	StartTime                   time.Time
	EndTime                     time.Time
	Duration                    time.Duration
	NumPaying                   int
	NumFree                     int
	NumAccompanying             int
	Confirmed                   *bool
	IsPlaceholderNumeroAttivita bool
}
