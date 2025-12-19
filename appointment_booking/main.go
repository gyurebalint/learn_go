package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Appointment struct {
	BookingId    string
	HospitalName string
	PatientName  string
	DoctorName   string
	Contact      string
	Date         string
	Time         string
}

func main() {
	// create and parse form
	formHtml := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Book Appointment</title>
<style>
        body { font-family: system-ui, -apple-system, sans-serif; background: #f4f4f9; display: flex; justify-content: center; align-items: center; min-height: 100vh; margin: 0; padding: 20px; box-sizing: border-box; }
        .container { background: white; padding: 1.5rem; border-radius: 8px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); width: 100%; max-width: 380px; }
        h2 { text-align: center; color: #333; margin-top: 0; margin-bottom: 1rem; font-size: 1.5rem; }
        .form-group { margin-bottom: 0.8rem; }
        label { display: block; margin-bottom: 0.3rem; color: #666; font-size: 0.85rem; }
        input, select { width: 100%; padding: 0.6rem; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; font-size: 0.95rem; background-color: white; }
        input:focus, select:focus { border-color: #007bff; outline: none; }
        button { width: 100%; padding: 0.75rem; background: #007bff; color: white; border: none; border-radius: 4px; font-size: 1rem; cursor: pointer; transition: background 0.2s; margin-top: 0.5rem; }
        button:hover { background: #0056b3; }
    </style>
</head>
<body>
    <div class="container">
        <h2>New Appointment</h2>
        <form action="/submit" method="POST" target="_blank">
            <div class="form-group">
                <label for="HospitalName">Hospital Name</label>
                <input type="text" id="HospitalName" name="HospitalName" placeholder="e.g. City General" required>
            </div>

            <div class="form-group">
                <label for="PatientName">Patient Name</label>
                <input type="text" id="PatientName" name="PatientName" placeholder="Full Name" required>
            </div>

			<div class="form-group">
                <label for="DoctorName">Doctor Name</label>
                <select id="DoctorName" name="DoctorName" required style="width: 100%; padding: 0.75rem; border: 1px solid #ddd; border-radius: 4px; font-size: 1rem; background-color: white;">
                    <option value="" disabled selected>Select a doctor</option>
                    <option value="Dr. Meredith Grey">Dr. Meredith Grey</option>
                    <option value="Dr. Gregory House">Dr. Gregory House</option>
                    <option value="Dr. Stephen Strange">Dr. Stephen Strange</option>
                </select>
            </div>

            <div class="form-group">
                <label for="Contact">Contact Number</label>
                <input type="tel" id="Contact" name="Contact" placeholder="+1 234 567 890" required>
            </div>

            <div class="form-group">
                <label for="Date">Preferred Date</label>
                <input type="date" id="Date" name="Date" required>
            </div>

            <div class="form-group">
                <label for="Time">Preferred Time</label>
                <input type="time" id="Time" name="Time" required>
            </div>

            <button type="submit">Book Appointment</button>
        </form>
    </div>
</body>
</html>
`
	readyAppointment := `
			<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>Appointment Confirmed</title>
				<style>
					body { font-family: system-ui, -apple-system, sans-serif; background: #eef2f5; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
					.container { background: white; padding: 2.5rem; border-radius: 8px; box-shadow: 0 4px 15px rgba(0,0,0,0.1); width: 100%; max-width: 400px; text-align: center; }
					h2 { color: #28a745; margin-bottom: 1.5rem; }
					.detail-row { display: flex; justify-content: space-between; padding: 0.8rem 0; border-bottom: 1px solid #eee; }
					.detail-row:last-child { border-bottom: none; }
					.label { color: #666; font-weight: 500; }
					.value { color: #333; font-weight: bold; }
					.footer { margin-top: 2rem; font-size: 0.85rem; color: #888; }
				</style>
			</head>
			<body>
				<div class="container">
					<h2>✅ Booking Confirmed</h2>
				<div/>
			
				<div class="detail-row">
					<span class="label">Hospital</span>
					<span class="value">{{.HospitalName}}</span>
				</div>
				<div class="detail-row">
						<span class="label">Patient:</span>
						<span class="value">{{.PatientName}}</span>
					</div>
					<div class="detail-row">
						<span class="label">Date:</span>
						<span class="value">{{.Date}}</span>
					</div>
					<div class="detail-row">
						<span class="label">Time:</span>
						<span class="value">{{.Time}}</span>
					</div>
			
					<div class="footer">
						Booking ID: {{.BookingId}}<br>
						Please arrive 15 minutes early.
					</div>
				</div>
			</body>
			</html>
			`

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		form, err := template.New("form").Parse(formHtml)
		if err != nil {
			fmt.Println("rendering of the form html is not working")
		}

		err = form.Execute(w, nil)
		if err != nil {
			fmt.Println("error loading form")
		}
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form html", http.StatusBadRequest)
			return
		}
		// prepare data
		appointment := Appointment{
			BookingId:    "Bk-001",
			HospitalName: r.FormValue("HospitalName"),
			PatientName:  r.FormValue("PatientName"),
			DoctorName:   r.FormValue("DoctorName"),
			Contact:      r.FormValue("Contact"),
			Date:         r.FormValue("Date"),
			Time:         r.FormValue("Time"),
		}

		receipt, err := template.New("receipt").Parse(readyAppointment)
		if err != nil {
			fmt.Println("could not parse receipt")
		}
		err = receipt.Execute(w, appointment)
		if err != nil {
			fmt.Println("error rendering receipt")
		}
	})

	fmt.Println("listening on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
