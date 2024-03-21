package helpers

import "time"

func GetGreeting(sender string) string {
	// Configurar o locale para pt-br
	// dayjs.locale('pt-br');
	// Obter a data e hora atual
	now := time.Now().UTC().Add(-3 * time.Hour)
	// Verificar o período do dia e retornar a saudação apropriada
	hour := now.Hour()
	if hour >= 4 && hour < 12 {
		return "Bom dia " + sender + ", dormiu bem?"
	} else if hour >= 12 && hour < 18 {
		if hour == 12 {
			return "Boa tarde " + sender + ", já almoçou?"
		}
		return "Boa tarde " + sender + ", como vai?"
	} else {
		return "Boa noite " + sender + ", tudo bem?"
	}
}
