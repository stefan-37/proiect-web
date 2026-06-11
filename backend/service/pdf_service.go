package service

import (
	"bytes"
	"fmt"
	"time"

	"backend/models"

	"github.com/go-pdf/fpdf"
)

// GeneratePaymentReceipt renders a PDF receipt for a single subscription
// purchase and returns it as raw bytes, ready to be streamed from a handler
// with c.Data(http.StatusOK, "application/pdf", pdfBytes).
//
// Note on Romanian text: fpdf's built-in fonts are Latin-1, so characters like
// ă/â/î/ș/ț won't render correctly. Keep receipt strings ASCII, or register a
// UTF-8 TTF with pdf.AddUTF8Font and switch the SetFont calls to use it.
func GeneratePaymentReceipt(user *models.User, sub *models.Subscription, us *models.UserSubscription) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 22)
	pdf.CellFormat(0, 12, "Chitanta", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 6, fmt.Sprintf("Emisa la %s", time.Now().Format("02.01.2006 15:04")), "", 1, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(8)

	// Customer block
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "Client", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	line(pdf, "Nume:", user.Name)
	line(pdf, "Email:", user.Email)
	line(pdf, "Telefon:", user.Phone)
	pdf.Ln(6)

	// Subscription block
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "Abonament", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	line(pdf, "Tip:", sub.Type)
	freeClasses := "Nu"
	if sub.FreeClasses {
		freeClasses = "Da"
	}
	line(pdf, "Clase gratuite:", freeClasses)
	if !us.StartedAt.IsZero() {
		line(pdf, "Valabil de la:", us.StartedAt.Format("02.01.2006"))
	}
	if !us.ExpiresAt.IsZero() {
		line(pdf, "Expira la:", us.ExpiresAt.Format("02.01.2006"))
	}
	pdf.Ln(6)

	// Total
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(120, 10, "Total de plata", "", 0, "L", false, 0, "")
	pdf.CellFormat(50, 10, fmt.Sprintf("%.2f RON", sub.Price), "", 1, "R", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("generating receipt pdf: %w", err)
	}
	return buf.Bytes(), nil
}

// GenerateBookingReceipt renders a PDF receipt for a single class booking and
// returns it as raw bytes, ready to be streamed from a handler with
// c.Data(http.StatusOK, "application/pdf", pdfBytes).
//
// Same Latin-1 caveat as GeneratePaymentReceipt applies to Romanian diacritics.
func GenerateBookingReceipt(user *models.User, class *models.Class, bs *models.BookingSituation) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 20, 20)
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 22)
	pdf.CellFormat(0, 12, "Chitanta rezervare", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(120, 120, 120)
	pdf.CellFormat(0, 6, fmt.Sprintf("Emisa la %s", time.Now().Format("02.01.2006 15:04")), "", 1, "C", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(8)

	// Customer block
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "Client", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	line(pdf, "Nume:", user.Name)
	line(pdf, "Email:", user.Email)
	line(pdf, "Telefon:", user.Phone)
	pdf.Ln(6)

	// Class block
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 8, "Clasa", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	line(pdf, "Nume:", class.Name)
	if class.Description != "" {
		line(pdf, "Descriere:", class.Description)
	}
	if !class.ScheduledAt.IsZero() {
		line(pdf, "Programata la:", class.ScheduledAt.Format("02.01.2006 15:04"))
	}
	if !bs.CreatedAt.IsZero() {
		line(pdf, "Rezervat la:", bs.CreatedAt.Format("02.01.2006 15:04"))
	}
	pdf.Ln(6)

	// Total
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(20, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(4)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(120, 10, "Total de plata", "", 0, "L", false, 0, "")
	pdf.CellFormat(50, 10, fmt.Sprintf("%.2f RON", class.Price), "", 1, "R", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, fmt.Errorf("generating booking receipt pdf: %w", err)
	}
	return buf.Bytes(), nil
}

// line prints a "label    value" row, label in a fixed-width left column.
func line(pdf *fpdf.Fpdf, label, value string) {
	pdf.CellFormat(40, 7, label, "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 7, value, "", 1, "L", false, 0, "")
}
