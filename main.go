package main

import (
	"net/http"
	"os"

	"github.com/email-service/data"
	"github.com/email-service/util"
	"github.com/gorilla/handlers"
	mailgun "github.com/mailgun/mailgun-go"
)

var (
	// domain        = os.Getenv("DOMAIN_MAILGUN_SANDBOX")
	// apiKeyPrivate = os.Getenv("PRIVATE_KEY_MAILGUN")
	httpPort          = os.Getenv("PORT")
	listenIP          = "localhost"
	pathEmail         = "/sendemail"
	pathTemplatesHTML = "templates/inlined/alert.html"
)

const (
	messageMethodNotAllowed      = "Method Not Allowed"
	messageMustHaveAPIKeyPublic  = "API Key Public invalid"
	messageMustHaveAPIKeyPrivate = "API Key Private invalid"
	messageMustHaveDomain        = "Domain mailgun invalid"
	messageMustHaveFrom          = "Your must have key from for email onwer"
	messageMustHaveTo            = "Your must have key to for email recipient"
	messageMustHaveSubject       = "Your must have subject"
	messageMustHaveBodyOrHTML    = "Your must have body or html"
)

const (
	methodPost       = "POST"
	keyDomain        = "domain"
	keyAPIKeyPublic  = "api_key_public"
	keyAPIKeyPrivate = "api_key_private"
	keyFrom          = "from"
	keyTo            = "to"
	keySubject       = "subject"
	keyBody          = "body"
	keyHTML          = "html"
)

func main() {
	http.HandleFunc(pathEmail, handlerEmail)
	http.ListenAndServe(":"+httpPort, handlers.LoggingHandler(os.Stdout, http.DefaultServeMux))
}

func handlerEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != methodPost {
		util.PrintErrorMessage(w, http.StatusMethodNotAllowed, messageMethodNotAllowed)
		return
	}
	domain := r.FormValue(keyDomain)
	apiKeyPublic := r.FormValue(keyAPIKeyPublic)
	apiKeyPrivate := r.FormValue(keyAPIKeyPrivate)
	from := r.FormValue(keyFrom)
	to := r.FormValue(keyTo)
	subject := r.FormValue(keySubject)
	bodyMessage := r.FormValue(keyBody)
	html := r.FormValue(keyHTML)

	var emailModel data.Email
	emailModel.Domain = domain
	emailModel.APIKeyPublic = apiKeyPublic
	emailModel.APIKeyPrivate = apiKeyPrivate
	emailModel.Subject = subject
	emailModel.From = from
	emailModel.To = to
	emailModel.Body = bodyMessage
	emailModel.HTML = html

	errorMessage := validateParams(emailModel)
	if len(errorMessage) != 0 {
		util.PrintErrorMessage(w, http.StatusBadRequest, errorMessage)
	} else {
		if len(bodyMessage) > 0 {
			emailModel.HTML = ""
		} else if len(html) > 0 {
			emailModel.Body = ""
		}
		err := sendEmail(emailModel)
		if err != nil {
			util.PrintErrorMessage(w, http.StatusBadRequest, err.Error())
		} else {
			var model data.Success
			model.Messgae = getMessageSuccess(emailModel.To)
			util.PrintSuccessMessage(w, model)
		}
	}
}

func validateParams(model data.Email) (errorMessage string) {
	if len(model.Domain) == 0 {
		return messageMustHaveDomain
	} else if len(model.APIKeyPrivate) == 0 {
		return messageMustHaveAPIKeyPrivate
	} else if len(model.APIKeyPublic) == 0 {
		return messageMustHaveAPIKeyPublic
	} else if len(model.Subject) == 0 {
		return messageMustHaveSubject
	} else if len(model.From) == 0 {
		return messageMustHaveFrom
	} else if len(model.To) == 0 {
		return messageMustHaveTo
	} else if len(model.Body) == 0 && len(model.HTML) == 0 {
		return messageMustHaveBodyOrHTML
	}
	return ""
}

func sendEmail(model data.Email) (err error) {
	mg := mailgun.NewMailgun(model.Domain, model.APIKeyPrivate, model.APIKeyPublic)
	message := mg.NewMessage(
		model.From,
		model.Subject,
		model.Body,
		model.To)
	message.SetHtml(model.HTML)
	_, _, err = mg.Send(message)
	if err != nil {
		return err
	}
	return nil
}

func getMessageSuccess(email string) (message string) {
	return "send email to " + email + " successfully"
}

// TODO: read templates html with inlin style only from file to string
// html, err := ioutil.ReadFile(pathTemplatesHTML)
// if err != nil {
// 	log.Fatal(err)
// 	return
// }
// fmt.Printf("%s", string(html))
