package mail

import (
	"fmt"
	"github.com/lmurature/melist-api/src/api/config"
	"github.com/lmurature/melist-api/src/api/domain/apierrors"
	"net/smtp"
	"strings"
)

var auth smtp.Auth

func init() {
	mail, pass := config.EmailAddress, config.EmailPassword
	auth = smtp.PlainAuth("", mail, pass, config.SmtpHost)
}

func SendMail(emailAddress string,
	shareType string, inviterFirstName string,
	inviterLastName string, listTitle string, authUrl string) {
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: ¡Te invitaron a colaborar en una lista Melist!\r\n"+
		"\r\n"+
		"¡Hola!\n\n\nEl usuario %s %s te invitó a colaborar en su lista %s otorgandote acceso de %s.\n\n\nPara ingresar, debés registrarte en la plataforma utilizando el siguiente link: %s \n\n\n¡Feliz listado!\r\n",
		emailAddress, strings.ToTitle(inviterFirstName), strings.ToTitle(inviterLastName), listTitle, shareType, authUrl))
	err := smtp.SendMail(
		config.SmtpAddress,
		auth,
		"melistapplication@gmail.com",
		[]string{emailAddress},
		msg,
	)

	if err != nil {
		e := apierrors.NewInternalServerApiError(fmt.Sprintf("Error while trying to mail [%s, %s, %s, %s, %s]",
			emailAddress, inviterFirstName, inviterLastName, shareType, authUrl), err)
		fmt.Println(e)
		return
	}

	fmt.Println("Successfully mailed user ", emailAddress)
}
