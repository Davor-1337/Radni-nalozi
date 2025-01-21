Ova aplikacija je razvijena kao deo diplomskog rada. Služi za praćenje i upravljanje radnim nalozima, praćenje potrošnje materijala, evidenciju vremena rada i generisanje faktura.
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

Aplikacija sadrzi sve handlere koje si mi zadao da uradim, s tim da sam ja jos dodao nekih stvari koje su mi se cinile zanimljive. Npr korisnik tabelu u kojoj ce se nalaziti
razliciti nalozi za korisnike, sa 3 razlicite role. Admin role koji ima pristup svim rutama, serviser role koji ima pristup rutama za mijenjanje radnog naloga za koji je zaduzen, unos materijala, sati itd. I klijent
koji ima pristup svojim radnim nalozima i fakturama. Napravio sam po middleware za svaki role koji vrsi provjeru preko tokena da li ulogovani nalog ima pristup odredjenoj ruti. 
Dodao sam automatsko arhiviranje naloga, gdje ce se nalog arhivirati automatski ako mu je Status: "Zavrsen" preko mjesec dana, a moguce je i manuelno arhiviranje.
Ja sam stavio da je admin zaduzen za kreiranje naloga serviserima i korisnicima, nisam siguran da li je to odgovarajuci pristup ali popravljacemo sta ne bude valjalo. 
Sto se tice notifikacije o promjeni statusa radnog naloga, ja sam napravio automatsko slanje emaila klijentu kojem pripada radni nalog. Znaci izvuce iz baze email klijenta i bice mu poslan email da je status
njegovog radnog naloga promijenjen. Mislim da je to to, ako se sjetim jos neceg pisacu na viber. 


