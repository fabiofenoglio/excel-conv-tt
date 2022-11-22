package excel

import (
	"fmt"
	"testing"
)

func Test_rewriteSchoolNameWithRules(t *testing.T) {

	tests := []struct {
		input string
		want  string
	}{
		{input: "", want: ""},
		{input: "some name", want: "some name"},
		{
			input: "Scuola: Secondaria I grado ISTITUTO COMPRENSIVO - ISTITUTO COMPRENSIVO DI MONTA'",
			want:  "SM I.C. DI MONTA'",
		},
		{
			input: "Scuola: Secondaria I grado IC 1 San Mauro",
			want:  "SM I.C. 1 San Mauro",
		},
		{
			input: "Scuola: Primaria ISTITUTO COMPRENSIVO - I.C. ALPIGNANO",
			want:  "EL I.C. ALPIGNANO",
		},
		{
			input: "Scuola: Secondaria II grado ISTRUZIONE SECONDARIA SUPERIORE - EINSTEIN",
			want:  "SUP I.S.S. - EINSTEIN",
		},
		{
			input: "Scuola: Primaria DIREZIONE DIDATTICA - D. D. 'P.P.LAMBERT'",
			want:  "EL DIREZIONE DIDATTICA - D. D. 'P.P.LAMBERT'",
		},
		{
			input: "Scuola: Primaria ISTITUTO COMPRENSIVO - IC SOMMARIVA DEL BOSCO'",
			want:  "EL I.C. SOMMARIVA DEL BOSCO'",
		},
		{
			input: "Scuola: Primaria ISTITUTO COMPRENSIVO - VERZUOLO - L. DA VINCI",
			want:  "EL I.C. VERZUOLO - L. DA VINCI",
		},
		{
			input: "Scuola: Primaria ISTITUTO COMPRENSIVO - CHERASCO - S. TARICCO",
			want:  "EL I.C. CHERASCO - S. TARICCO",
		},
		{
			input: "Scuola: Primaria I.C. San Giorgio Canavese",
			want:  "EL I.C. San Giorgio Canavese",
		},
		{
			input: "Scuola: Primaria CRESCEREINSIEME SOCIETA' COOPERATIVA SOCIALE ONLUS",
			want:  "EL CRESCEREINSIEME SOCIETA' COOPERATIVA SOCIALE ONLUS",
		},
		{
			input: "Scuola: Primaria IC  Moncalieri - Nasi",
			want:  "EL I.C. Moncalieri - Nasi",
		},
		{
			input: "Scuola: Secondaria I grado ISTRUZIONE SECONDARIA SUPERIORE - IS L. DES AMBROIS",
			want:  "SM I.S.S. - IS L. DES AMBROIS",
		},
		{
			input: "Scuola: Primaria Istituto    Comprensivo Marconi - Antonelli",
			want:  "EL I.C. Marconi - Antonelli",
		},
		{
			input: "Scuola: Primaria IC Sibilla Aleramo",
			want:  "EL I.C. Sibilla Aleramo",
		},
		{
			input: "Scuola: Primaria ISTITUTO COMPRENSIVO - I.C. - TORINO - CENA",
			want:  "EL I.C. TORINO - CENA",
		},
		{
			input: "Scuola: Primaria ISTITUTO COMPRENSIVO - I.C.S.MAURIZIO-MARIA MONTESSORI",
			want:  "EL I.C. I.C.S.MAURIZIO-MARIA MONTESSORI",
		},
		{
			input: "Altro ISTITUTO COMPRENSIVO - I.C. - TORINO - CENA",
			want:  "Altro I.C. TORINO - CENA",
		},
		{
			input: "Custom: Primaria SNIC Sibilla Aleramo",
			want:  "Custom: EL SNIC Sibilla Aleramo",
		},
		{
			input: "   Custom: Primaria ICT Sibilla Aleramo   ",
			want:  "Custom: EL ICT Sibilla Aleramo",
		},
		{
			input: "Custom: Primaria S.N.I.C. Sibilla Aleramo",
			want:  "Custom: EL S.N.I.C. Sibilla Aleramo",
		},
		{
			input: "Custom: Primaria S.N.I.C.T. Sibilla Aleramo",
			want:  "Custom: EL S.N.I.C.T. Sibilla Aleramo",
		},
		{
			input: "   Custom:    Primaria  \t\r\n  S.N.I.C.T.  Sibilla  Aleramo   ",
			want:  "Custom: EL S.N.I.C.T. Sibilla Aleramo",
		},
		{
			input: "Se scrivo ANPrimaria o Primariamente o XXXXPrimariaSSS deve restare cosi",
			want:  "Se scrivo ANPrimaria o Primariamente o XXXXPrimariaSSS deve restare cosi",
		},
		{
			input: "Primaria a inizio",
			want:  "EL a inizio",
		},
		{
			input: "a fine primaria",
			want:  "a fine EL",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			if got := rewriteSchoolNameWithRules(tt.input); got != tt.want {
				t.Errorf("rewriteSchoolNameWithRules() = %v, want %v", got, tt.want)
			}
		})
	}
}
