USE [master]
GO
/****** Object:  Database [RadniNaloziDB]    Script Date: 1/21/2025 5:27:17 PM ******/
CREATE DATABASE [RadniNaloziDB]
 CONTAINMENT = NONE
 ON  PRIMARY 
( NAME = N'RadniNaloziDB', FILENAME = N'C:\Program Files\Microsoft SQL Server\MSSQL11.MSSQLSERVER\MSSQL\DATA\RadniNaloziDB.mdf' , SIZE = 4096KB , MAXSIZE = UNLIMITED, FILEGROWTH = 1024KB )
 LOG ON 
( NAME = N'RadniNaloziDB_log', FILENAME = N'C:\Program Files\Microsoft SQL Server\MSSQL11.MSSQLSERVER\MSSQL\DATA\RadniNaloziDB_log.ldf' , SIZE = 1024KB , MAXSIZE = 2048GB , FILEGROWTH = 10%)
GO
ALTER DATABASE [RadniNaloziDB] SET COMPATIBILITY_LEVEL = 110
GO
IF (1 = FULLTEXTSERVICEPROPERTY('IsFullTextInstalled'))
begin
EXEC [RadniNaloziDB].[dbo].[sp_fulltext_database] @action = 'enable'
end
GO
ALTER DATABASE [RadniNaloziDB] SET ANSI_NULL_DEFAULT OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET ANSI_NULLS OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET ANSI_PADDING OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET ANSI_WARNINGS OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET ARITHABORT OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET AUTO_CLOSE OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET AUTO_CREATE_STATISTICS ON 
GO
ALTER DATABASE [RadniNaloziDB] SET AUTO_SHRINK OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET AUTO_UPDATE_STATISTICS ON 
GO
ALTER DATABASE [RadniNaloziDB] SET CURSOR_CLOSE_ON_COMMIT OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET CURSOR_DEFAULT  GLOBAL 
GO
ALTER DATABASE [RadniNaloziDB] SET CONCAT_NULL_YIELDS_NULL OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET NUMERIC_ROUNDABORT OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET QUOTED_IDENTIFIER OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET RECURSIVE_TRIGGERS OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET  DISABLE_BROKER 
GO
ALTER DATABASE [RadniNaloziDB] SET AUTO_UPDATE_STATISTICS_ASYNC OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET DATE_CORRELATION_OPTIMIZATION OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET TRUSTWORTHY OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET ALLOW_SNAPSHOT_ISOLATION OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET PARAMETERIZATION SIMPLE 
GO
ALTER DATABASE [RadniNaloziDB] SET READ_COMMITTED_SNAPSHOT OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET HONOR_BROKER_PRIORITY OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET RECOVERY SIMPLE 
GO
ALTER DATABASE [RadniNaloziDB] SET  MULTI_USER 
GO
ALTER DATABASE [RadniNaloziDB] SET PAGE_VERIFY CHECKSUM  
GO
ALTER DATABASE [RadniNaloziDB] SET DB_CHAINING OFF 
GO
ALTER DATABASE [RadniNaloziDB] SET FILESTREAM( NON_TRANSACTED_ACCESS = OFF ) 
GO
ALTER DATABASE [RadniNaloziDB] SET TARGET_RECOVERY_TIME = 0 SECONDS 
GO
USE [RadniNaloziDB]
GO
/****** Object:  StoredProcedure [dbo].[Arhiviraj]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[Arhiviraj] 
AS BEGIN
    -- Definišite trenutni datum
    DECLARE @TrenutniDatum DATETIME = GETDATE();

    -- Arhiviranje
    INSERT INTO ArhiviraniRadniNalozi
    SELECT rn.Nalog_ID, rn.Klijent_ID, rn.OpisProblema, rn.Prioritet, rn.DatumOtvaranja, rn.Status,
           rn.Lokacija, rn.BrojNaloga, @TrenutniDatum AS DatumArhiviranja
    FROM RadniNalog rn
    JOIN DodjelaZadatka dz ON dz.Nalog_ID = rn.Nalog_ID
    WHERE rn.Status = 'Zavrsen' 
      AND dz.DatumZavrsetka IS NOT NULL
      AND ISDATE(dz.DatumZavrsetka) = 1
      AND CAST(dz.DatumZavrsetka AS DATE) < DATEADD(MONTH, -2, @TrenutniDatum);

    -- Brisanje
    DELETE FROM RadniNalog
    WHERE Nalog_ID IN (
        SELECT rn.Nalog_ID
        FROM RadniNalog rn
        JOIN DodjelaZadatka dz ON dz.Nalog_ID = rn.Nalog_ID
        WHERE rn.Status = 'Zavrsen' 
          AND dz.DatumZavrsetka IS NOT NULL
          AND ISDATE(dz.DatumZavrsetka) = 1
          AND CAST(dz.DatumZavrsetka AS DATE) < DATEADD(MONTH, -2, @TrenutniDatum)
    );
END
GO
/****** Object:  StoredProcedure [dbo].[Arhiviranje]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE PROCEDURE [dbo].[Arhiviranje]
    @NalogID INT
AS
BEGIN
    SET NOCOUNT ON;

    -- Započni transakciju
    BEGIN TRANSACTION;

    BEGIN TRY
	DECLARE @TrenutniDatum DATETIME = GETDATE()
    -- Arhiviranje podataka za određeni radni nalog
    INSERT INTO ArhiviraniRadniNalozi (Nalog_ID, Klijent_ID, OpisProblema, Prioritet, DatumOtvaranja, Status, Lokacija, BrojNaloga, DatumArhiviranja)
    SELECT 
        rn.Nalog_ID, 
        rn.Klijent_ID, 
        rn.OpisProblema, 
        rn.Prioritet, 
        rn.DatumOtvaranja, 
        rn.Status,
        rn.Lokacija, 
        rn.BrojNaloga, 
        @TrenutniDatum AS DatumArhiviranja
    FROM 
        RadniNalog rn
    JOIN 
        DodjelaZadatka dz ON dz.Nalog_ID = rn.Nalog_ID
    WHERE 
        rn.Nalog_ID = @NalogID
        AND rn.Status = 'Zavrsen'

    -- Provjera da li je nešto ubačeno
    IF @@ROWCOUNT = 0
    BEGIN
        RAISERROR('Radni nalog sa ID %d nije pronađen ili nije završen.', 16, 1, @NalogID);
        RETURN;
    END

    -- Brisanje arhiviranih podataka iz glavne tabele
    DELETE FROM RadniNalog
    WHERE Nalog_ID = @NalogID
    AND Status = 'Zavrsen'

    COMMIT TRANSACTION;

    PRINT 'Arhiviranje uspješno izvršeno za Nalog_ID: ' + CAST(@NalogID AS NVARCHAR);
END TRY

    BEGIN CATCH
        -- Ako je došlo do greške, poništi transakciju
        ROLLBACK TRANSACTION;

        -- Vrati grešku korisniku
        DECLARE @ErrorMessage NVARCHAR(4000) = ERROR_MESSAGE();
        PRINT 'Došlo je do greške: ' + @ErrorMessage;
    END CATCH
END;

EXEC Arhiviranje @NalogID = 44
GO
/****** Object:  Table [dbo].[ArhiviraniRadniNalozi]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[ArhiviraniRadniNalozi](
	[Nalog_ID] [int] NOT NULL,
	[Klijent_ID] [int] NULL,
	[OpisProblema] [nvarchar](500) NULL,
	[Prioritet] [nvarchar](500) NULL,
	[DatumOtvaranja] [datetime] NULL,
	[Status] [nvarchar](50) NULL,
	[Lokacija] [nvarchar](100) NULL,
	[BrojNaloga] [nvarchar](20) NULL,
	[DatumArhiviranja] [datetime] NULL,
 CONSTRAINT [PK__Arhivira__F683A75916279883] PRIMARY KEY CLUSTERED 
(
	[Nalog_ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
/****** Object:  Table [dbo].[DodjelaZadatka]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[DodjelaZadatka](
	[Zadatak_ID] [int] IDENTITY(1,1) NOT NULL,
	[Nalog_ID] [int] NULL,
	[Serviser_ID] [int] NOT NULL,
	[DatumDodjele] [datetime] NOT NULL CONSTRAINT [DF__DodjelaZa__Datum__1DE57479]  DEFAULT (getdate()),
	[DatumZavrsetka] [datetime] NULL,
 CONSTRAINT [PK__DodjelaZ__330F844FFBF00FC6] PRIMARY KEY CLUSTERED 
(
	[Zadatak_ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
/****** Object:  Table [dbo].[EvidencijaSati]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[EvidencijaSati](
	[Evidencija_ID] [int] IDENTITY(1,1) NOT NULL,
	[Nalog_ID] [int] NULL,
	[Serviser_ID] [int] NOT NULL,
	[BrojRadnihSati] [decimal](5, 2) NOT NULL,
	[DatumUnosa] [datetime] NOT NULL CONSTRAINT [DF__Evidencij__Datum__29572725]  DEFAULT (getdate()),
 CONSTRAINT [PK__Evidenci__17B6E1D4C9CEC2E6] PRIMARY KEY CLUSTERED 
(
	[Evidencija_ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
/****** Object:  Table [dbo].[Faktura]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Faktura](
	[Faktura_ID] [int] IDENTITY(1,1) NOT NULL,
	[Nalog_ID] [int] NULL,
	[Iznos] [decimal](10, 2) NOT NULL,
	[DatumFakture] [datetime] NOT NULL,
 CONSTRAINT [PK__Faktura__70AA0049E13B43A8] PRIMARY KEY CLUSTERED 
(
	[Faktura_ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
/****** Object:  Table [dbo].[Klijenti]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
SET ANSI_PADDING ON
GO
CREATE TABLE [dbo].[Klijenti](
	[Klijent_ID] [int] NOT NULL,
	[Naziv] [varchar](255) NOT NULL,
	[KontaktOsoba] [varchar](255) NOT NULL,
	[Email] [varchar](255) NOT NULL,
	[Tel] [varchar](50) NULL,
	[Adresa] [varchar](255) NULL,
	[User_ID] [int] NULL,
 CONSTRAINT [PK__Klijenti__45BAB2C6068B4524] PRIMARY KEY CLUSTERED 
(
	[Klijent_ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
SET ANSI_PADDING OFF
GO
/****** Object:  Table [dbo].[Korisnici]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Korisnici](
	[User_ID] [int] NOT NULL,
	[Username] [nvarchar](50) NOT NULL,
	[Email] [nvarchar](100) NULL,
	[Password] [nvarchar](255) NOT NULL,
	[Role] [nvarchar](50) NULL,
 CONSTRAINT [PK__Korisnic__206D91904286829E] PRIMARY KEY CLUSTERED 
(
	[User_ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
/****** Object:  Table [dbo].[Materijal]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[Materijal](
	[Materijal_ID] [int] IDENTITY(1,1) NOT NULL,
	[NazivMaterijala] [nvarchar](100) NOT NULL,
	[Cijena] [decimal](10, 2) NOT NULL,
	[KolicinaUSkladistu] [int] NOT NULL,
PRIMARY KEY CLUSTERED 
(
	[Materijal_ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
/****** Object:  Table [dbo].[RadniNalog]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
SET ANSI_PADDING ON
GO
CREATE TABLE [dbo].[RadniNalog](
	[Nalog_ID] [int] IDENTITY(1,1) NOT NULL,
	[Klijent_ID] [int] NOT NULL,
	[OpisProblema] [nvarchar](500) NULL,
	[Prioritet] [nvarchar](50) NULL,
	[DatumOtvaranja] [datetime] NOT NULL DEFAULT (getdate()),
	[Status] [nvarchar](50) NULL DEFAULT ('Otvoren'),
	[Lokacija] [nvarchar](255) NULL,
	[BrojNaloga]  AS ('RN-'+right('0000'+CONVERT([varchar],[Nalog_ID]),(4))) PERSISTED,
PRIMARY KEY CLUSTERED 
(
	[Nalog_ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
SET ANSI_PADDING OFF
GO
/****** Object:  Table [dbo].[Serviser]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
SET ANSI_PADDING ON
GO
CREATE TABLE [dbo].[Serviser](
	[Serviser_ID] [int] NOT NULL,
	[Ime] [nvarchar](100) NOT NULL,
	[Prezime] [nvarchar](100) NOT NULL,
	[Specijalnost] [nvarchar](100) NULL,
	[Telefon] [varchar](50) NULL,
	[User_ID] [int] NULL,
 CONSTRAINT [PK__Serviser__62E9F5B7B26E8A88] PRIMARY KEY CLUSTERED 
(
	[Serviser_ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
SET ANSI_PADDING OFF
GO
/****** Object:  Table [dbo].[UtroseniMaterijal]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TABLE [dbo].[UtroseniMaterijal](
	[UtroseniMaterijal_ID] [int] IDENTITY(1,1) NOT NULL,
	[Nalog_ID] [int] NULL,
	[Materijal_ID] [int] NOT NULL,
	[KolicinaUtrosena] [int] NOT NULL,
	[DatumUnosa] [datetime] NOT NULL CONSTRAINT [DF__UtroseniM__Datum__24927208]  DEFAULT (getdate()),
 CONSTRAINT [PK__Utroseni__0E6351BCF8B132E7] PRIMARY KEY CLUSTERED 
(
	[UtroseniMaterijal_ID] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, IGNORE_DUP_KEY = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
) ON [PRIMARY]

GO
INSERT [dbo].[ArhiviraniRadniNalozi] ([Nalog_ID], [Klijent_ID], [OpisProblema], [Prioritet], [DatumOtvaranja], [Status], [Lokacija], [BrojNaloga], [DatumArhiviranja]) VALUES (17, 200, N'Klima uredjaj ne hladi ispravno.', N'Nizak', CAST(N'2024-12-20 15:30:00.000' AS DateTime), N'Zavrsen', N'Kneza Mihajla 12a', N'RN-0017', CAST(N'2025-01-17 13:07:57.390' AS DateTime))
INSERT [dbo].[ArhiviraniRadniNalozi] ([Nalog_ID], [Klijent_ID], [OpisProblema], [Prioritet], [DatumOtvaranja], [Status], [Lokacija], [BrojNaloga], [DatumArhiviranja]) VALUES (18, 201, N'Potrebna instalacija novog osvjetljenja u kancelariji', N'Visok', CAST(N'2024-12-22 12:30:00.000' AS DateTime), N'Zavrsen', N'Mihajla Pupina 122', N'RN-0018', CAST(N'2025-01-17 13:29:52.380' AS DateTime))
INSERT [dbo].[ArhiviraniRadniNalozi] ([Nalog_ID], [Klijent_ID], [OpisProblema], [Prioritet], [DatumOtvaranja], [Status], [Lokacija], [BrojNaloga], [DatumArhiviranja]) VALUES (19, 202, N'Problemi s radom hidraulicnog sistema na radnoj masini.', N'Visok', CAST(N'2024-12-23 13:00:00.000' AS DateTime), N'Zavrsen', N'Nikole Tesle 99a', N'RN-0019', CAST(N'2025-01-17 13:37:04.100' AS DateTime))
INSERT [dbo].[ArhiviraniRadniNalozi] ([Nalog_ID], [Klijent_ID], [OpisProblema], [Prioritet], [DatumOtvaranja], [Status], [Lokacija], [BrojNaloga], [DatumArhiviranja]) VALUES (20, 203, N'Popravka sistema za grijanje.', N'Visok', CAST(N'2024-12-20 13:00:00.000' AS DateTime), N'Zavrsen', N'Desanka Maksimovic 202', N'RN-0020', CAST(N'2025-01-17 18:36:19.100' AS DateTime))
INSERT [dbo].[ArhiviraniRadniNalozi] ([Nalog_ID], [Klijent_ID], [OpisProblema], [Prioritet], [DatumOtvaranja], [Status], [Lokacija], [BrojNaloga], [DatumArhiviranja]) VALUES (22, 201, N'Monitor treperi i povremeno gubi sliku.', N'Nizak', CAST(N'2024-12-20 09:12:00.000' AS DateTime), N'Zavrsen', N'Mihajla Pupina 122', N'RN-0022', CAST(N'2025-01-17 19:06:46.050' AS DateTime))
SET IDENTITY_INSERT [dbo].[DodjelaZadatka] ON 

INSERT [dbo].[DodjelaZadatka] ([Zadatak_ID], [Nalog_ID], [Serviser_ID], [DatumDodjele], [DatumZavrsetka]) VALUES (31, 23, 102, CAST(N'2024-12-30 10:58:17.923' AS DateTime), NULL)
INSERT [dbo].[DodjelaZadatka] ([Zadatak_ID], [Nalog_ID], [Serviser_ID], [DatumDodjele], [DatumZavrsetka]) VALUES (32, 24, 101, CAST(N'2024-12-15 00:00:00.000' AS DateTime), NULL)
INSERT [dbo].[DodjelaZadatka] ([Zadatak_ID], [Nalog_ID], [Serviser_ID], [DatumDodjele], [DatumZavrsetka]) VALUES (33, 25, 101, CAST(N'2024-12-26 00:00:00.000' AS DateTime), NULL)
INSERT [dbo].[DodjelaZadatka] ([Zadatak_ID], [Nalog_ID], [Serviser_ID], [DatumDodjele], [DatumZavrsetka]) VALUES (34, 26, 101, CAST(N'2024-12-22 00:00:00.000' AS DateTime), NULL)
INSERT [dbo].[DodjelaZadatka] ([Zadatak_ID], [Nalog_ID], [Serviser_ID], [DatumDodjele], [DatumZavrsetka]) VALUES (35, 27, 100, CAST(N'2024-12-27 00:00:00.000' AS DateTime), NULL)
SET IDENTITY_INSERT [dbo].[DodjelaZadatka] OFF
SET IDENTITY_INSERT [dbo].[EvidencijaSati] ON 

INSERT [dbo].[EvidencijaSati] ([Evidencija_ID], [Nalog_ID], [Serviser_ID], [BrojRadnihSati], [DatumUnosa]) VALUES (47, 23, 102, CAST(5.00 AS Decimal(5, 2)), CAST(N'2024-12-31 00:00:00.000' AS DateTime))
INSERT [dbo].[EvidencijaSati] ([Evidencija_ID], [Nalog_ID], [Serviser_ID], [BrojRadnihSati], [DatumUnosa]) VALUES (48, 24, 101, CAST(8.00 AS Decimal(5, 2)), CAST(N'2024-12-16 00:00:00.000' AS DateTime))
INSERT [dbo].[EvidencijaSati] ([Evidencija_ID], [Nalog_ID], [Serviser_ID], [BrojRadnihSati], [DatumUnosa]) VALUES (49, 25, 101, CAST(10.00 AS Decimal(5, 2)), CAST(N'2024-12-26 00:00:00.000' AS DateTime))
INSERT [dbo].[EvidencijaSati] ([Evidencija_ID], [Nalog_ID], [Serviser_ID], [BrojRadnihSati], [DatumUnosa]) VALUES (50, 26, 101, CAST(8.00 AS Decimal(5, 2)), CAST(N'2024-12-22 00:00:00.000' AS DateTime))
INSERT [dbo].[EvidencijaSati] ([Evidencija_ID], [Nalog_ID], [Serviser_ID], [BrojRadnihSati], [DatumUnosa]) VALUES (51, 27, 100, CAST(14.00 AS Decimal(5, 2)), CAST(N'2024-12-27 00:00:00.000' AS DateTime))
SET IDENTITY_INSERT [dbo].[EvidencijaSati] OFF
INSERT [dbo].[Klijenti] ([Klijent_ID], [Naziv], [KontaktOsoba], [Email], [Tel], [Adresa], [User_ID]) VALUES (200, N'Pinnacle ProTech', N'Jovana Knezevic', N'petarp435@gmail.com', N'+38766554202', N'Kneza Mihajla 12a', 200)
INSERT [dbo].[Klijenti] ([Klijent_ID], [Naziv], [KontaktOsoba], [Email], [Tel], [Adresa], [User_ID]) VALUES (201, N'Blue Vibe Innovations', N'Milana Dobrijevic', N'MilanaBVI@gmail.com', N'+38766020333', N'Mihajla Pupina 122', 201)
INSERT [dbo].[Klijenti] ([Klijent_ID], [Naziv], [KontaktOsoba], [Email], [Tel], [Adresa], [User_ID]) VALUES (202, N'Vertex Vision', N'Jelena Mirkovic', N'JelenaVV@gmail.com', N'+38766233095', N'Nikole Tesle 99a', 202)
INSERT [dbo].[Klijenti] ([Klijent_ID], [Naziv], [KontaktOsoba], [Email], [Tel], [Adresa], [User_ID]) VALUES (203, N'Lumina Core', N'Dijana Radulovic', N'DijanaLC@gmail.com', N'+38766603106', N'Desanka Maksimovic 202', 203)
INSERT [dbo].[Klijenti] ([Klijent_ID], [Naziv], [KontaktOsoba], [Email], [Tel], [Adresa], [User_ID]) VALUES (204, N'Swift Peak', N'Milica Aksentijevic', N'MilicaSP@gmail.com', N'+38766205205', N'Petar Kocic 44', 204)
INSERT [dbo].[Korisnici] ([User_ID], [Username], [Email], [Password], [Role]) VALUES (13, N'davor', N'davor4@gmail.com', N'$2a$14$oEpmy7BIqhHbsH80lR4aAesGFjyOZwnVS3l6hJJUTlVy9GHmFNyPG', N'admin')
INSERT [dbo].[Korisnici] ([User_ID], [Username], [Email], [Password], [Role]) VALUES (100, N'Milan212', N'milan212@gmail.com', N'$2a$14$Ei1hC1LijE2GRMW/PdL8fu4yH6VXbjTNKBuJ/HvlB/wvbBLkN01dG', N'serviser')
INSERT [dbo].[Korisnici] ([User_ID], [Username], [Email], [Password], [Role]) VALUES (101, N'MarkoJ23', N'markojankovic23@gmail.com', N'$2a$14$D8F0weQlWTLdlOeo.BcxrO.ftD/QhmCjtQFlEVixM6Hp6PNDQH5Li', N'serviser')
INSERT [dbo].[Korisnici] ([User_ID], [Username], [Email], [Password], [Role]) VALUES (102, N'Petar00', N'petard00@gmail.com', N'$2a$14$V2CZC4QpWDd4Yh9dOowTt.2D0CAmUhOq9uaI6dreBJaSgQblI/4Qm', N'serviser')
INSERT [dbo].[Korisnici] ([User_ID], [Username], [Email], [Password], [Role]) VALUES (103, N'LukaV231', N'lukavukmir@gmail.com', N'$2a$14$a6dy/WIBLP2R1seTZDkkOudAln6Tb6Ll3cyfbDH4oY2Lk3AJF.o2y', N'serviser')
INSERT [dbo].[Korisnici] ([User_ID], [Username], [Email], [Password], [Role]) VALUES (104, N'StefanF', N'stefanfilipovic@gmail.com', N'$2a$14$ePEGWTSMewZaCdWpVzs1..Gm3r6hEsygVI16BqIu9lmcvFBwPVZWa', N'serviser')
INSERT [dbo].[Korisnici] ([User_ID], [Username], [Email], [Password], [Role]) VALUES (200, N'PPT', N'pinnaclept@gmail.com', N'$2a$14$5CDppCn/24qh/dghbh8DXeR0EhdenS1.YFYOvX3B6oonaGiPQQqAu', N'klijent')
INSERT [dbo].[Korisnici] ([User_ID], [Username], [Email], [Password], [Role]) VALUES (201, N'BlueVibeI', N'bluevibeinovations@gmail.com', N'$2a$14$mXOg4BJfWMFS8JgE8to9peMcw.QEQeKrmJFqwXpJRuP7DxfEVG.Cq', N'klijent')
INSERT [dbo].[Korisnici] ([User_ID], [Username], [Email], [Password], [Role]) VALUES (202, N'Vertex', N'VertexVision@gmail.com', N'$2a$14$ovKuBh/PCvlKDo0OLV30yOnF9mCnRKveSR34JnrGwk7izgdU2CUBe', N'klijent')
INSERT [dbo].[Korisnici] ([User_ID], [Username], [Email], [Password], [Role]) VALUES (203, N'LuminaC', N'LuminaCore@gmail.com', N'$2a$14$V2/h78YRak7/q2qkYnugCuho.r8PPKcljUW5K292Cu1/0SVCRLaVS', N'klijent')
INSERT [dbo].[Korisnici] ([User_ID], [Username], [Email], [Password], [Role]) VALUES (204, N'SwiftP', N'SwiftPeak@gmail.com', N'$2a$14$d50CWEw3/lby0fA9da6L8uIyYGcwLsXN.TNKJhslUgVcgziuSZvA6', N'klijent')
SET IDENTITY_INSERT [dbo].[Materijal] ON 

INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (4, N'Elektricni kabl (10m)', CAST(14.99 AS Decimal(10, 2)), 45)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (5, N'Uticnica sa uzemljenjem', CAST(8.99 AS Decimal(10, 2)), 97)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (6, N'LED sijalica (10W)', CAST(9.99 AS Decimal(10, 2)), 200)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (7, N'Mazivo za zupcanike', CAST(29.99 AS Decimal(10, 2)), 30)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (8, N'Vijci M6x30', CAST(25.99 AS Decimal(10, 2)), 475)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (9, N'Zavarivacka elektroda', CAST(4.99 AS Decimal(10, 2)), 80)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (10, N'Ethernet kabl (5m)', CAST(11.99 AS Decimal(10, 2)), 148)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (11, N'SSD disk (256GB)', CAST(189.99 AS Decimal(10, 2)), 19)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (12, N'Maticna ploca (Intel)', CAST(239.99 AS Decimal(10, 2)), 8)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (13, N'Termalna pasta', CAST(6.99 AS Decimal(10, 2)), 49)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (14, N'Ljepak za drvo', CAST(8.99 AS Decimal(10, 2)), 40)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (15, N'Izolir traka', CAST(2.99 AS Decimal(10, 2)), 300)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (16, N'Zamjenski filter za klima uredaj', CAST(44.99 AS Decimal(10, 2)), 15)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (17, N'Vodootporna traka (3m)', CAST(7.99 AS Decimal(10, 2)), 100)
INSERT [dbo].[Materijal] ([Materijal_ID], [NazivMaterijala], [Cijena], [KolicinaUSkladistu]) VALUES (18, N'Bakarna zica (5kg)', CAST(59.99 AS Decimal(10, 2)), 25)
SET IDENTITY_INSERT [dbo].[Materijal] OFF
SET ANSI_PADDING ON
SET IDENTITY_INSERT [dbo].[RadniNalog] ON 

INSERT [dbo].[RadniNalog] ([Nalog_ID], [Klijent_ID], [OpisProblema], [Prioritet], [DatumOtvaranja], [Status], [Lokacija]) VALUES (23, 200, N'Ne radi ventilacija u kancelariji.', N'Nizak', CAST(N'2024-12-20 09:54:00.000' AS DateTime), N'Otvoren', N'Kneza Mihajla 12a')
INSERT [dbo].[RadniNalog] ([Nalog_ID], [Klijent_ID], [OpisProblema], [Prioritet], [DatumOtvaranja], [Status], [Lokacija]) VALUES (24, 201, N'Server u glavnom uredu se ne može pokrenuti. Korisnici ne mogu pristupiti mrežnim resursima. Potrebno je hitno dijagnosticirati i otkloniti problem', N'Visok', CAST(N'2024-12-15 00:00:00.000' AS DateTime), N'Otvoren', N'Mihajla Pupina 122')
INSERT [dbo].[RadniNalog] ([Nalog_ID], [Klijent_ID], [OpisProblema], [Prioritet], [DatumOtvaranja], [Status], [Lokacija]) VALUES (25, 202, N'U odjelu prodaje dolazi do čestih prekida mrežne veze. Svi uređaji povezani na mrežu imaju ograničen pristup internetu i internim aplikacijama', N'Nizak', CAST(N'2024-12-26 00:00:00.000' AS DateTime), N'Otvoren', N'Nikole Tesle 99a')
INSERT [dbo].[RadniNalog] ([Nalog_ID], [Klijent_ID], [OpisProblema], [Prioritet], [DatumOtvaranja], [Status], [Lokacija]) VALUES (26, 203, N'Korisnik prijavljuje da aplikacija za upravljanje zalihama prestaje raditi prilikom unosa novih proizvoda. Pojavljuje se poruka greške: ''Unexpected Error 404''.', N'Nizak', CAST(N'2024-12-22 00:00:00.000' AS DateTime), N'Otvoren', N'Desanka Maksimovic 202')
INSERT [dbo].[RadniNalog] ([Nalog_ID], [Klijent_ID], [OpisProblema], [Prioritet], [DatumOtvaranja], [Status], [Lokacija]) VALUES (27, 204, N'U skladištu su se pojavili problemi s električnim napajanjem. Utičnice na južnom zidu ne rade, što uzrokuje probleme s korištenjem opreme.', N'Visok', CAST(N'2024-12-27 00:00:00.000' AS DateTime), N'Otvoren', N'Petar Kocic 44')
SET IDENTITY_INSERT [dbo].[RadniNalog] OFF
SET ANSI_PADDING OFF
INSERT [dbo].[Serviser] ([Serviser_ID], [Ime], [Prezime], [Specijalnost], [Telefon], [User_ID]) VALUES (100, N'Milan', N'Petkovic', N'Elektricar', N'+38766533131', 100)
INSERT [dbo].[Serviser] ([Serviser_ID], [Ime], [Prezime], [Specijalnost], [Telefon], [User_ID]) VALUES (101, N'Marko', N'Jankovic', N'IT Tehnicar', N'+38766005496', 101)
INSERT [dbo].[Serviser] ([Serviser_ID], [Ime], [Prezime], [Specijalnost], [Telefon], [User_ID]) VALUES (102, N'Petar', N'Dakic', N'Generalni serviser', N'+38766105442', 102)
INSERT [dbo].[Serviser] ([Serviser_ID], [Ime], [Prezime], [Specijalnost], [Telefon], [User_ID]) VALUES (103, N'Luka', N'Vukmir', N'Mehanicar', N'+38766558917', 103)
INSERT [dbo].[Serviser] ([Serviser_ID], [Ime], [Prezime], [Specijalnost], [Telefon], [User_ID]) VALUES (104, N'Stefan', N'Filipovic', N'Elektricar', N'+38766820046', 104)
SET IDENTITY_INSERT [dbo].[UtroseniMaterijal] ON 

INSERT [dbo].[UtroseniMaterijal] ([UtroseniMaterijal_ID], [Nalog_ID], [Materijal_ID], [KolicinaUtrosena], [DatumUnosa]) VALUES (25, 23, 4, 5, CAST(N'2024-12-20 00:00:00.000' AS DateTime))
INSERT [dbo].[UtroseniMaterijal] ([UtroseniMaterijal_ID], [Nalog_ID], [Materijal_ID], [KolicinaUtrosena], [DatumUnosa]) VALUES (26, 23, 5, 3, CAST(N'2024-12-20 00:00:00.000' AS DateTime))
INSERT [dbo].[UtroseniMaterijal] ([UtroseniMaterijal_ID], [Nalog_ID], [Materijal_ID], [KolicinaUtrosena], [DatumUnosa]) VALUES (27, 24, 12, 1, CAST(N'2024-12-16 00:00:00.000' AS DateTime))
INSERT [dbo].[UtroseniMaterijal] ([UtroseniMaterijal_ID], [Nalog_ID], [Materijal_ID], [KolicinaUtrosena], [DatumUnosa]) VALUES (28, 24, 13, 1, CAST(N'2024-12-16 00:00:00.000' AS DateTime))
INSERT [dbo].[UtroseniMaterijal] ([UtroseniMaterijal_ID], [Nalog_ID], [Materijal_ID], [KolicinaUtrosena], [DatumUnosa]) VALUES (29, 25, 11, 1, CAST(N'2024-12-28 00:00:00.000' AS DateTime))
INSERT [dbo].[UtroseniMaterijal] ([UtroseniMaterijal_ID], [Nalog_ID], [Materijal_ID], [KolicinaUtrosena], [DatumUnosa]) VALUES (30, 26, 12, 1, CAST(N'2024-12-23 00:00:00.000' AS DateTime))
INSERT [dbo].[UtroseniMaterijal] ([UtroseniMaterijal_ID], [Nalog_ID], [Materijal_ID], [KolicinaUtrosena], [DatumUnosa]) VALUES (31, 26, 10, 2, CAST(N'2024-12-28 00:00:00.000' AS DateTime))
SET IDENTITY_INSERT [dbo].[UtroseniMaterijal] OFF
SET ANSI_PADDING ON

GO
/****** Object:  Index [UQ__Korisnic__536C85E4DE1A9B6C]    Script Date: 1/21/2025 5:27:17 PM ******/
ALTER TABLE [dbo].[Korisnici] ADD  CONSTRAINT [UQ__Korisnic__536C85E4DE1A9B6C] UNIQUE NONCLUSTERED 
(
	[Username] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, SORT_IN_TEMPDB = OFF, IGNORE_DUP_KEY = OFF, ONLINE = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
GO
SET ANSI_PADDING ON

GO
/****** Object:  Index [UQ__Korisnic__A9D10534CE94E7D0]    Script Date: 1/21/2025 5:27:17 PM ******/
ALTER TABLE [dbo].[Korisnici] ADD  CONSTRAINT [UQ__Korisnic__A9D10534CE94E7D0] UNIQUE NONCLUSTERED 
(
	[Email] ASC
)WITH (PAD_INDEX = OFF, STATISTICS_NORECOMPUTE = OFF, SORT_IN_TEMPDB = OFF, IGNORE_DUP_KEY = OFF, ONLINE = OFF, ALLOW_ROW_LOCKS = ON, ALLOW_PAGE_LOCKS = ON) ON [PRIMARY]
GO
ALTER TABLE [dbo].[Faktura] ADD  CONSTRAINT [DF__Faktura__DatumFa__2E1BDC42]  DEFAULT (getdate()) FOR [DatumFakture]
GO
ALTER TABLE [dbo].[ArhiviraniRadniNalozi]  WITH CHECK ADD  CONSTRAINT [FK__Arhiviran__Klije__3864608B] FOREIGN KEY([Klijent_ID])
REFERENCES [dbo].[Klijenti] ([Klijent_ID])
GO
ALTER TABLE [dbo].[ArhiviraniRadniNalozi] CHECK CONSTRAINT [FK__Arhiviran__Klije__3864608B]
GO
ALTER TABLE [dbo].[DodjelaZadatka]  WITH CHECK ADD  CONSTRAINT [FK_DodjelaZadatka_Nalog] FOREIGN KEY([Nalog_ID])
REFERENCES [dbo].[RadniNalog] ([Nalog_ID])
ON DELETE CASCADE
GO
ALTER TABLE [dbo].[DodjelaZadatka] CHECK CONSTRAINT [FK_DodjelaZadatka_Nalog]
GO
ALTER TABLE [dbo].[DodjelaZadatka]  WITH CHECK ADD  CONSTRAINT [FK_DodjelaZadatka_Serviser] FOREIGN KEY([Serviser_ID])
REFERENCES [dbo].[Serviser] ([Serviser_ID])
ON DELETE CASCADE
GO
ALTER TABLE [dbo].[DodjelaZadatka] CHECK CONSTRAINT [FK_DodjelaZadatka_Serviser]
GO
ALTER TABLE [dbo].[EvidencijaSati]  WITH CHECK ADD  CONSTRAINT [FK__Evidencij__Servi__2B3F6F97] FOREIGN KEY([Serviser_ID])
REFERENCES [dbo].[Serviser] ([Serviser_ID])
GO
ALTER TABLE [dbo].[EvidencijaSati] CHECK CONSTRAINT [FK__Evidencij__Servi__2B3F6F97]
GO
ALTER TABLE [dbo].[EvidencijaSati]  WITH CHECK ADD  CONSTRAINT [FK_EvidencijaSati_Nalog_ID] FOREIGN KEY([Nalog_ID])
REFERENCES [dbo].[RadniNalog] ([Nalog_ID])
ON DELETE CASCADE
GO
ALTER TABLE [dbo].[EvidencijaSati] CHECK CONSTRAINT [FK_EvidencijaSati_Nalog_ID]
GO
ALTER TABLE [dbo].[Faktura]  WITH CHECK ADD  CONSTRAINT [FK_Faktura_Nalog_ID] FOREIGN KEY([Nalog_ID])
REFERENCES [dbo].[RadniNalog] ([Nalog_ID])
ON DELETE CASCADE
GO
ALTER TABLE [dbo].[Faktura] CHECK CONSTRAINT [FK_Faktura_Nalog_ID]
GO
ALTER TABLE [dbo].[Klijenti]  WITH CHECK ADD  CONSTRAINT [FK__Klijenti__User_I__22751F6C] FOREIGN KEY([User_ID])
REFERENCES [dbo].[Korisnici] ([User_ID])
ON DELETE CASCADE
GO
ALTER TABLE [dbo].[Klijenti] CHECK CONSTRAINT [FK__Klijenti__User_I__22751F6C]
GO
ALTER TABLE [dbo].[RadniNalog]  WITH CHECK ADD  CONSTRAINT [FK__RadniNalo__Klije__145C0A3F] FOREIGN KEY([Klijent_ID])
REFERENCES [dbo].[Klijenti] ([Klijent_ID])
GO
ALTER TABLE [dbo].[RadniNalog] CHECK CONSTRAINT [FK__RadniNalo__Klije__145C0A3F]
GO
ALTER TABLE [dbo].[Serviser]  WITH CHECK ADD  CONSTRAINT [FK__Serviser__User_I__2180FB33] FOREIGN KEY([User_ID])
REFERENCES [dbo].[Korisnici] ([User_ID])
ON DELETE CASCADE
GO
ALTER TABLE [dbo].[Serviser] CHECK CONSTRAINT [FK__Serviser__User_I__2180FB33]
GO
ALTER TABLE [dbo].[UtroseniMaterijal]  WITH CHECK ADD  CONSTRAINT [FK_UtroseniMaterijal_Materijal] FOREIGN KEY([Materijal_ID])
REFERENCES [dbo].[Materijal] ([Materijal_ID])
ON DELETE CASCADE
GO
ALTER TABLE [dbo].[UtroseniMaterijal] CHECK CONSTRAINT [FK_UtroseniMaterijal_Materijal]
GO
ALTER TABLE [dbo].[UtroseniMaterijal]  WITH CHECK ADD  CONSTRAINT [FK_UtroseniMaterijal_RadniNalog1] FOREIGN KEY([Nalog_ID])
REFERENCES [dbo].[RadniNalog] ([Nalog_ID])
ON DELETE CASCADE
GO
ALTER TABLE [dbo].[UtroseniMaterijal] CHECK CONSTRAINT [FK_UtroseniMaterijal_RadniNalog1]
GO
/****** Object:  Trigger [dbo].[UpdateStockAfterMaterialInsert]    Script Date: 1/21/2025 5:27:17 PM ******/
SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO
CREATE TRIGGER [dbo].[UpdateStockAfterMaterialInsert]
ON [dbo].[UtroseniMaterijal]
AFTER INSERT
AS
BEGIN
    SET NOCOUNT ON;

    -- Ažuriraj zalihe materijala
    UPDATE Materijal
    SET KolicinaUSkladistu = KolicinaUSkladistu - inserted.KolicinaUtrosena
    FROM Materijal
    INNER JOIN inserted ON Materijal.Materijal_ID = inserted.Materijal_ID;

    -- Proveri negativne zalihe i vrati grešku ako je potrebno
    IF EXISTS (
        SELECT 1
        FROM Materijal
        WHERE KolicinaUSkladistu < 0
    )
    BEGIN
        RAISERROR ('Nema dovoljno zaliha za ovaj materijal.', 16, 1);
        ROLLBACK TRANSACTION;
    END
END;


GO
USE [master]
GO
ALTER DATABASE [RadniNaloziDB] SET  READ_WRITE 
GO
