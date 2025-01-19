package utils

import (
    "fmt"
    "gopkg.in/gomail.v2"
)

func SendStatusChangeEmail(toEmail, clientName, workOrderID, newStatus string) error {
    // SMTP konfiguracija
    smtpHost := "smtp.gmail.com"
    smtpPort := 587             // Port za TLS
    smtpUser := "davor.markovic1337@gmail.com" // Tvoj e-mail
    smtpPass := "afva uqxt hqce dmpo" // Tvoja lozinka ili aplikacijski ključ

    // Sadržaj e-maila
    subject := fmt.Sprintf("Promena statusa radnog naloga #%s", workOrderID)
    body := fmt.Sprintf("Poštovani %s,\n\nStatus vašeg radnog naloga #%s je promenjen na: %s.\n\nHvala na poverenju!\n\nS poštovanjem,\nVaš servisni tim", clientName, workOrderID, newStatus)

    
    m := gomail.NewMessage()
    m.SetHeader("From", smtpUser)          // Pošiljalac
    m.SetHeader("To", toEmail)            // Primaoc
    m.SetHeader("Subject", subject)       // Naslov
    m.SetBody("text/plain", body)         // Tekst poruke

    // Slanje e-maila
    d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

    // Pokušaj slanja
    if err := d.DialAndSend(m); err != nil {
        return fmt.Errorf("failed to send email: %v", err)
    }

    return nil
}
