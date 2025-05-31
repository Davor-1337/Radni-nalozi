package models

import (
	"Diplomski/database"
	"database/sql"
	"fmt"
	
	"time"
	
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/signature"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontfamily"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/orientation"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"

	"github.com/johnfercher/maroto/v2/pkg/props"
)

type RadniNalogIzvjestaj struct {
	NalogID int64 
	BrojNaloga string
	Klijent string 
	Adresa string
	OpisProblema string 
	UtroseniMaterijali string 
	CijenaMaterijala float64 
	BrojRadnihSati float64 
	DatumDodjele time.Time
	DatumZavrsetka time.Time
	Serviser string
	Iznos float64
}

type NalogIzvjestaj struct {
	BrojNaloga string
	Klijent string 
	OpisProblema string
	Nalog_ID string
}



func GetAllOrderReportsShort() ([]NalogIzvjestaj, error) {
    query := `SELECT rn.BrojNaloga, k.Naziv as Klijent, rn.OpisProblema, rn.Nalog_ID
              FROM RadniNalog rn
              JOIN Klijenti k on rn.Klijent_ID = k.Klijent_ID`

    rows, err := database.DB.Query(query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var reports []NalogIzvjestaj

    for rows.Next() {
        var report NalogIzvjestaj
        err := rows.Scan(&report.BrojNaloga, &report.Klijent, &report.OpisProblema, &report.Nalog_ID)
        if err != nil {
            return nil, err
        }
        reports = append(reports, report)
    }

    if err = rows.Err(); err != nil {
        return nil, err
    }

    return reports, nil
}

func GetOrderReportsShortForClient(clientID int64) ([]NalogIzvjestaj, error) {
	const query = `
			SELECT rn.BrojNaloga,
						 k.Naziv    AS Klijent,
						 rn.OpisProblema,
						 rn.Nalog_ID
				FROM RadniNalog rn
				JOIN Klijenti k ON rn.Klijent_ID = k.Klijent_ID
			 WHERE rn.Klijent_ID = @ClientID
	`
	rows, err := database.DB.Query(
			query,
			sql.Named("ClientID", clientID),
	)
	if err != nil {
			return nil, err
	}
	defer rows.Close()

	var reports []NalogIzvjestaj
	for rows.Next() {
			var r NalogIzvjestaj
			if err := rows.Scan(&r.BrojNaloga, &r.Klijent, &r.OpisProblema, &r.Nalog_ID); err != nil {
					return nil, err
			}
			reports = append(reports, r)
	}
	if err := rows.Err(); err != nil {
			return nil, err
	}
	return reports, nil
}

func GetWorkOrderReport(id int64) (*RadniNalogIzvjestaj, error) {
	query := `SELECT 
    rn.Nalog_ID, 
    rn.BrojNaloga,
	k.Naziv as Klijent, 
	k.Adresa,
    rn.OpisProblema, 
    (
        SELECT STUFF((
            SELECT ', ' + m.NazivMaterijala + ' (' + CAST(um.KolicinaUtrosena AS NVARCHAR) + ' kom)'
            FROM UtroseniMaterijal um
            JOIN Materijal m ON um.Materijal_ID = m.Materijal_ID
            WHERE um.Nalog_ID = rn.Nalog_ID
            FOR XML PATH(''), TYPE).value('.', 'NVARCHAR(MAX)'), 1, 2, '')
    ) AS UtroseniMaterijali,
    Sum(m.Cijena) AS CijenaMaterijala,  
    es.BrojRadnihSati,
    dz.DatumDodjele,
    dz.DatumZavrsetka,
	CONCAT(s.Ime, ' ', s.Prezime) AS Serviser,
	f.Iznos
FROM RadniNalog rn
JOIN EvidencijaSati es ON es.Nalog_ID = rn.Nalog_ID
JOIN DodjelaZadatka dz ON dz.Nalog_ID = rn.Nalog_ID
JOIN UtroseniMaterijal um ON um.Nalog_ID = rn.Nalog_ID
JOIN Materijal m ON um.Materijal_ID = m.Materijal_ID
JOIN Serviser s on s.Serviser_ID = dz.Serviser_ID
JOIN Klijenti k on k.Klijent_ID = rn.Klijent_ID
JOIN Faktura f on f.Nalog_ID = rn.Nalog_ID
WHERE rn.Nalog_ID = @Nalog_ID
GROUP BY rn.Nalog_ID, rn.BrojNaloga, k.Adresa, rn.OpisProblema, es.BrojRadnihSati, dz.DatumDodjele, dz.DatumZavrsetka, k.Naziv, s.Ime, s.Prezime, f.Iznos

`
row := database.DB.QueryRow(query, sql.Named("Nalog_ID", id))

var report RadniNalogIzvjestaj

err := row.Scan(&report.NalogID, &report.BrojNaloga, &report.Klijent, &report.Adresa, &report.OpisProblema, &report.UtroseniMaterijali, &report.CijenaMaterijala, &report.BrojRadnihSati, &report.DatumDodjele, &report.DatumZavrsetka, &report.Serviser, &report.Iznos)

if err != nil {
	return nil, err
}

return &report, nil
}

func GeneratePDF(report *RadniNalogIzvjestaj) ([]byte, error) {
	cfg := config.NewBuilder().
		WithOrientation(orientation.Vertical).
		WithPageSize(pagesize.A4).
		WithLeftMargin(15).
		WithBottomMargin(5).
		WithRightMargin(15).
		WithTopMargin(15).
		Build()

	m := maroto.New(cfg)

	addHeader(m)
	addReportDetails(m, report)
	addWorkOrderDetails(m, report)
	addFooter(m)

	doc, err := m.Generate()
	if err != nil {
		return nil, err
	}

	pdfBytes := doc.GetBytes()
	

	return pdfBytes, nil
}


func addHeader(m core.Maroto){
	m.AddRow(20,
			text.NewCol(12, "Izvjestaj", props.Text{
				Top: 5,
				Style: fontstyle.Bold,
				Align: align.Center,
				Size: 18,
			}))	
}


func addReportDetails(m core.Maroto, report *RadniNalogIzvjestaj) {
	m.AddRow(10,
		text.NewCol(6, "Datum: "+time.Now().Format("02 Jan 2006"), props.Text{
			Align: align.Left,
			Size:  10,
		}),
		text.NewCol(6, "Izvještaj #"+report.BrojNaloga, props.Text{
			Align: align.Right,
			Size:  10,
		}))
	m.AddRow(5, line.NewCol(12)) 
}


func addWorkOrderDetails(m core.Maroto, report *RadniNalogIzvjestaj) {
	
	addField := func(label, value string) {
		m.AddRow(15,
			text.NewCol(4, label+":", props.Text{
				Align: align.Left,
				Size:  11,
				Style: fontstyle.Bold,
			}),
			text.NewCol(8, value, props.Text{
				Align: align.Left,
				Size:  11,
			}),
		)
	}

	// Dodavanje polja u PDF
	addField("Klijent", report.Klijent)
	addField("Adresa", report.Adresa)
	addField("Opis Problema", report.OpisProblema)
	addField("Utrošeni Materijali", report.UtroseniMaterijali)
	addField("Cijena Materijala", fmt.Sprintf("%.2f KM", report.CijenaMaterijala))
	addField("Broj Radnih Sati", fmt.Sprintf("%.2f", report.BrojRadnihSati))
	addField("Datum Dodjele", report.DatumDodjele.Format("02.01.2006"))
	addField("Datum Završetka", report.DatumZavrsetka.Format("02.01.2006"))
	addField("Serviser", report.Serviser)
	addField("Ukupan Iznos", fmt.Sprintf("%.2f KM", report.Iznos))

	m.AddRow(5, line.NewCol(12)) 
}

func addFooter (m core.Maroto) {
	m.AddRow(30)
	m.AddRow(50,
	signature.NewCol(6, "Potpis", props.Signature{FontFamily: fontfamily.Courier}),
	code.NewQrCol(6, "https://codeheim.io", props.Rect{Percent: 75, Center: true}))
}

