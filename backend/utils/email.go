package utils

import (
    "fmt"
    "gopkg.in/gomail.v2"
)

func SendStatusChangeEmail(toEmail, clientName, workOrderID, newStatus string) error {
    // SMTP konfiguracija
    smtpHost := "smtp.gmail.com"
    smtpPort := 587            
    smtpUser := "davor.markovic1337@gmail.com" 
    smtpPass := "afva uqxt hqce dmpo" 

    // Sadržaj e-maila
    subject := fmt.Sprintf("Promena statusa radnog naloga #%s", workOrderID)
    body := fmt.Sprintf("Poštovani %s,\n\nStatus vašeg radnog naloga #%s je promenjen na: %s.\n\nHvala na poverenju!\n\nS poštovanjem,\nVaš servisni tim", clientName, workOrderID, newStatus)

    
    m := gomail.NewMessage()
    m.SetHeader("From", smtpUser)         
    m.SetHeader("To", toEmail)           
    m.SetHeader("Subject", subject)       
    m.SetBody("text/plain", body)         

   
    d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPass)

   
    if err := d.DialAndSend(m); err != nil {
        return fmt.Errorf("failed to send email: %v", err)
    }

    return nil
}
