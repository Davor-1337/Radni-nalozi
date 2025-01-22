Ova aplikacija je razvijena kao dio diplomskog rada. Služi za praćenje i upravljanje radnim nalozima, praćenje potrošnje materijala, evidenciju vremena rada i generisanje faktura.
Da biste pokrenuli ovu aplikaciju, potrebno je da imate instalirane sledeće:

-Go (verzija 1.x ili novija) za backend deo aplikacije.
-SQL Server sa pristupom bazi podataka.

Potrebno je instalirati sve zavisnosti:
go mod tidy
Ova komanda će preuzeti sve potrebne biblioteke koje aplikacija koristi.

Da bi aplikacija mogla da se poveže sa bazom podataka, potrebno je posaviti odgovarajući konekcioni string u .env fajlu. Sledeći format treba da bude korišćen:
DB_CONNECTION=sqlserver://<korisnicko_ime>:<lozinka>@localhost:1433?database=RadniNaloziDB

Pokretanje backend-a
Aplikacija se pokrece pomocu sledece komande:

go run main.go

Podaci za logovanje usera sa role-om "Admin":
"Username": "davor",
"Password":  "davor" 

Podaci za logovanje usera sa role-om "Serviser":
"Username": "Milan212",
	"Password":  "Milan212" 

Podaci za logovanje usera sa role-om "Klijent":
{     
	"Username": "LuminaC",
	"Password":  "LuminaC"      
}
