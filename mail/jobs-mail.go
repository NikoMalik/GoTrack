package mail

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
)

//go:embed templates/*
var templatesFS embed.FS

// MailData holds info for sending an email
type MailData struct {
	ToName       string
	ToAddress    string
	FromName     string
	FromAddress  string
	AdditionalTo []string
	Subject      string
	Content      template.HTML
	Template     string
	CC           []string
	UseHermes    bool
	Attachments  []string
	StringMap    map[string]string
	IntMap       map[string]int
	FloatMap     map[string]float32
	RowSets      map[string]interface{}
}

// MailJob is the unit of work to be performed when sending an email to chan
type MailJob struct {
	MailMessage MailData
}

type Worker struct {
	ID         int
	JobQueue   chan MailJob
	WorkerPool chan chan MailJob
	QuitChan   chan bool
}

func NewWorker(id int, workerPool chan chan MailJob) *Worker {
	return &Worker{
		ID:         id,
		JobQueue:   make(chan MailJob),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobQueue
			select {
			case job := <-w.JobQueue:
				err := w.processMailQueueJob(job.MailMessage)
				if err != nil {
					fmt.Printf("Error processing job: %v\n", err)
				}
			case <-w.QuitChan:
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

// Dispatcher holds info for a dispatcher
type Dispatcher struct {
	workerPool chan chan MailJob
	maxWorkers int
	jobQueue   chan MailJob
}

// NewDispatcher creates, and returns a new Dispatcher object.
func NewDispatcher(jobQueue chan MailJob, maxWorkers int) *Dispatcher {
	workerPool := make(chan chan MailJob, maxWorkers)
	return &Dispatcher{
		jobQueue:   jobQueue,
		maxWorkers: maxWorkers,
		workerPool: workerPool,
	}
}

// run runs the workers
func (d *Dispatcher) run() {
	for i := 0; i < d.maxWorkers; i++ {
		worker := NewWorker(i+1, d.workerPool)
		worker.Start()
	}

	go d.dispatch()
}

// dispatch dispatches worker
func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.jobQueue:
			go func() {
				workerJobQueue := <-d.workerPool
				workerJobQueue <- job
			}()
		}
	}
}

func (w *Worker) processMailQueueJob(mail MailData) error {
	var preferenceMap map[string]string

	data := struct {
		Content       template.HTML
		From          string
		FromName      string
		PreferenceMap map[string]string
		IntMap        map[string]int
		StringMap     map[string]string
		FloatMap      map[string]float32
		RowSets       map[string]interface{}
	}{
		Content:       mail.Content,
		FromName:      mail.FromName,
		From:          mail.FromAddress,
		PreferenceMap: preferenceMap,
		IntMap:        mail.IntMap,
		StringMap:     mail.StringMap,
		FloatMap:      mail.FloatMap,
		RowSets:       mail.RowSets,
	}

	tmpl := "mail.tmpl"
	if mail.Template != "" {
		tmpl = mail.Template
	}

	t := template.Must(template.ParseFS(templatesFS, "templates/"+tmpl))

	var body bytes.Buffer
	if err := t.ExecuteTemplate(&body, tmpl, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	sender := os.Getenv("SMTP_SENDER")
	if sender == "" {
		return errors.New("smtp sender not provided")
	}

	smtpServer := os.Getenv("SMTP_SERVER")
	if smtpServer == "" {
		return errors.New("smtp server not provided")
	}

	smtpUsername := os.Getenv("SMTP_USERNAME")
	if smtpUsername == "" {
		return errors.New("smtp username not provided")
	}

	smtpPassword := os.Getenv("SMTP_PASSWORD")
	if smtpPassword == "" {
		return errors.New("smtp password not provided")
	}

	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "587" // default
	}

	auth := smtp.PlainAuth("", smtpUsername, smtpPassword, smtpServer)

	to := append(mail.AdditionalTo, mail.ToAddress)
	msg := []byte("To: " + mail.ToAddress + "\r\n" +
		"Subject: " + mail.Subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n\r\n" +
		body.String())

	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, sender, to, msg)
	if err != nil {
		return fmt.Errorf("error sending mail: %v", err)
	}

	return nil
}
