package email

import (
	"fmt"
	"log"
)

// SendWarningEmail отправляет предупреждение на почту пользователя при изменении IP-адреса
func SendWarningEmail(userEmail, ipAddress string) error {
	// Здесь можно настроить отправку email через сторонний сервис (например, SMTP, SendGrid и т.д.)
	// Для упрощения, имитируем отправку email через лог

	// В реальной системе можно интегрировать здесь реальный email-сервис
	log.Printf("Sending warning email to %s: IP address changed to %s\n", userEmail, ipAddress)

	// Имитация успешной отправки
	fmt.Printf("Warning email sent to %s about IP address change to %s\n", userEmail, ipAddress)
	return nil
}
