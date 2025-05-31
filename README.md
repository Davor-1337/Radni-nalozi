# Sistem za upravljanje radnim nalozima

## Funkcionalnosti aplikacije

###  Administratorske funkcije
- **Upravljanje radnim nalozima**
  - Kreiranje, pregled i ažuriranje radnih naloga
  - Dodjela naloga serviserima
  - Praćenje statusa naloga
  - Arhiviranje radnih naloga
- **Upravljanje korisnicima**
  - Odobravanje/odbijanje zahtjeva za registraciju (klijenti i serviseri)
  - Upravljanje korisničkim rolama i privilegijama
- **Upravljanje materijalima**
  - Evidencija potrošnje materijala
  - Praćenje zaliha
- **Financijski modul**
  - Generisanje faktura
  - Generisanje PDF izvjestaja
- **Sistem obavještenja**
  - Real-time notifikacije za zahtjeve za naloge
  - Obavještenja o zahtjevima za registraciju

### Klijentske funkcije
- **Radni nalozi**
  - Slanje zahtjeva za nove radne naloge
  - Pregled svih vlastitih naloga
  - Praćenje statusa (otvoren/prihvaćen/završen)
- **Dokumentacija**
  - Pregled generisanih faktura
  - Preuzimanje PDF izvještaja
- **Notifikacije**
  - Obavještenja o promjenama statusa naloga
  

###  Serviserske funkcije
- **Radni proces**
  - Pregled dodijeljenih radnih naloga
  - Unos vremena utrošenog na rad
  - Evidencija korištenih materijala
  - Završetak radnih naloga
  - Unos i pracenje materijala
- **Komunikacija**
  - Obavještenja o novim zadacima
  - Notifikacije o hitnim nalozima


## Tehnologije

**Backend:**
- Go (1.x+)
- SQL Server


**Frontend:**
- Angular CLI: 19.1.8
- TypeScript
- HTML/SCSS

**Baza podataka:**
- Microsoft SQL Server

## Preduvjeti

Prije instalacije, osigurajte da imate instalirano:
- [Go](https://golang.org/dl/) (verzija 1.x ili novija)
- [Node.js](https://nodejs.org/) (za Angular)
- [Angular CLI](https://angular.io/cli)
- [Microsoft SQL Server](https://www.microsoft.com/en-us/sql-server/sql-server-downloads)

## Instalacija i pokretanje

### Backend (Go)

1. Klonirajte repozitorij:
   ```bash
   git clone https://github.com/Davor-1337/Radni-nalozi
   cd radni-nalozi/backend

2. Instalirajte zavisnosti:
     go mod tidy
   
3. Konfigurisite bazu podataka (pogledajte Konfiguracija sekciju)
   
4. Pokrenite backend
  go run main.go


### Frontend (Angular)
1. Otvorite novi terminal i idite u frontend direktorij:
   
    cd radni-nalozi/frontend

2. Instalirajte zavisnosti:
    
    npm install

3. Pokrenite Angular aplikaciju:
   
    ng serve

4. Aplikacija će biti dostupna na: http://localhost:4200

##  Baza podataka - Postavljanje i konfiguracija

# 1. Priprema baze podataka
**Opcija A: Korištenje SQL Server Management Studija**
1. Pokrenite SQL Server Management Studio (SSMS)
2. Otvorite novi query prozor
3. Izvršite skriptu `backend/database/baza.sql`
4. Provjerite da je baza `RadniNaloziDB` kreirana

 #  2. Konfiguracija aplikacije
U backend direktoriju kreirajte .env fajl sa sledećim sadržajem:

"DB_CONNECTION=sqlserver://<username>:<password>@localhost:1433?database=RadniNaloziDB"
Zamijenite placeholdere:

<username> sa vašim SQL Server korisničkim imenom

<password> sa odgovarajućom lozinkom

3. Verifikacija
Pokrenite backend aplikaciju

Aplikacija bi trebala uspješno uspostaviti vezu sa bazom

Provjerite da se tabele i podaci učitavaju korektno






















   
