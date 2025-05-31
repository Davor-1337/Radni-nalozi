import { Component } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { CommonModule } from '@angular/common';
import { MatDialog } from '@angular/material/dialog';
import { AddComponent } from '../add/add.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { DocumentsService } from '../../../core/services/documents.service';
import { PdfComponent } from '../pdf/pdf.component';
import { AuthService } from '../../../core/services/auth.service';
// import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-list',
  imports: [MatTableModule, CommonModule, ReactiveFormsModule, FormsModule],
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
})
export class ListComponent {
  recentInvoices: any[] = [];
  recentReports: any[] = [];
  isSelectingReport = false;
  selectedOrder: any = null;
  editMode: { row: any; field: string } | null = null;
  editValue: string = '';
  originalValue: string = '';
  displayedColumns: string[] = [
    'icon',
    'broj_naloga',
    'datum_fakture',
    'iznos',
  ];
  displayedColumnsReports: string[] = [
    'icon',
    'BrojNaloga',
    'Klijent',
    'OpisProblema',
  ];

  constructor(
    private documentService: DocumentsService,
    private dialog: MatDialog,
    private snackBar: MatSnackBar,
    public authService: AuthService
  ) {}

  ngOnInit(): void {
    const role = this.authService.getRole();
    const userId = this.authService.getUserId();

    if (role === 'admin') {
      this.fetchInvoices();
      this.fetchReports();
    } else if (role === 'klijent') {
      this.fetchInvoicesForClient(userId);
      this.fetchReportsForClient(userId);
    }
  }

  // Invoices /////////////////////////////
  fetchInvoices(): void {
    this.documentService.getInvoices().subscribe({
      next: (data) => {
        this.recentInvoices = data;
      },
      error: (err) => {
        console.error('Error fetching Invoices:', err);
      },
    });
  }
  fetchInvoicesForClient(clientId: number): void {
    this.documentService.getInvoicesByClient(clientId).subscribe({
      next: (data) => (this.recentInvoices = data),
      error: (err) => console.error('Error fetching client invoices', err),
    });
  }

  openAddDialog() {
    const dialogRef = this.dialog.open(AddComponent, {
      width: '400px',
      height: '210px',
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        console.log('Novi materijal dodat:', result);
      }
    });
  }

  isEditing(row: any, field: string): boolean {
    return this.editMode?.row === row && this.editMode?.field === field;
  }

  startEdit(row: any, field: string, event: MouseEvent): void {
    event.stopPropagation();
    this.editMode = { row, field };
    this.editValue = row[field];
    this.originalValue = row[field];

    setTimeout(() => {
      document.addEventListener('click', this.handleClickOutside.bind(this));
    });
  }

  saveEdit(): void {
    if (!this.editMode) return;

    const { row, field } = this.editMode;

    row[field] = this.editValue;

    if (this.editValue === this.originalValue) {
      this.cancelEdit();
      return;
    }

    const updatedData = {
      Faktura_ID: row.faktura_id,
      Iznos:
        field === 'Iznos' ? parseFloat(this.editValue) : parseFloat(row.iznos),
    };

    this.documentService.updateInvoice(row.faktura_id, updatedData).subscribe({
      next: () => {
        row[field] = this.editValue;
        this.snackBar.open('Uspešno ažurirano!', 'Zatvori', {
          duration: 2000,
          panelClass: ['success-snackbar'],
        });
        this.cancelEdit();
      },
      error: (err) => {
        console.error('Greška pri ažuriranju:', err);

        row[field] = this.originalValue;
        this.snackBar.open('Greška pri ažuriranju fakture', 'Zatvori', {
          duration: 3000,
          panelClass: ['error-snackbar'],
        });
        this.cancelEdit();
      },
    });
  }

  cancelEdit(): void {
    if (this.editMode) {
      document.removeEventListener('click', this.handleClickOutside.bind(this));
      this.editMode = null;
    }
  }

  handleClickOutside(event: MouseEvent): void {
    if (!this.editMode) return;

    const clickedInside = (event.target as HTMLElement).closest(
      '.editing-cell'
    );
    if (!clickedInside) {
      if (this.editValue !== this.originalValue) {
        this.saveEdit();
      } else {
        this.cancelEdit();
      }
    }
  }

  deleteInvoice(InvoicesId: number, event: MouseEvent): void {
    event.stopPropagation();

    if (confirm('Da li ste sigurni da želite obrisati ovu fakturu?')) {
      this.documentService.deleteInvoice(InvoicesId).subscribe({
        next: () => {
          this.recentInvoices = this.recentInvoices.filter(
            (Invoices) => Invoices.faktura_id !== InvoicesId
          );
          this.snackBar.open('Faktura uspješno obrisana', 'Zatvori', {
            duration: 3000,
            panelClass: ['delete-snackbar'],
          });
        },
        error: (err) => {
          this.snackBar.open('Greška pri brisanju', 'Zatvori', {
            duration: 3000,
            panelClass: ['delete-snackbar'],
          });
        },
      });
    }
  }

  // REPORTS /////////////////////////////////////////

  fetchReports(): void {
    this.documentService.getShortReport().subscribe({
      next: (data) => {
        console.log('Podaci sa bekenda:', data);
        this.recentReports = data;
      },
      error: (err) => {
        console.error('Error fetching reports:', err);
      },
    });
  }

  fetchReportsForClient(clientId: number): void {
    this.documentService.getReportsByClient(clientId).subscribe({
      next: (data) => (this.recentReports = data),
      error: (err) => console.error('Error fetching client reports', err),
    });
  }

  openReport(orderId: number): void {
    console.log('Clicked on order ID:', orderId);

    if (!orderId) {
      this.snackBar.open('Neispravan ID naloga.', 'Zatvori', {
        duration: 3000,
        panelClass: ['error-snackbar'],
      });
      return;
    }

    this.documentService.generateReport(orderId).subscribe({
      next: (pdfBlob) => {
        this.dialog.open(PdfComponent, {
          width: '80%',
          height: '160%',
          data: { pdfBlob },
          panelClass: 'pdf-dialog',
        });
      },
      error: (err) => {
        console.error('Greška pri generisanju izveštaja:', err);

        if (err.status === 400 || err.status === 409) {
          this.snackBar.open(
            'Ne možete generisati izveštaj jer radni nalog još nije završen.',
            'Zatvori',
            { duration: 4000, panelClass: ['error-snackbar'] }
          );
        } else {
          this.snackBar.open(
            'Greška pri generisanju izveštaja. Pokušajte ponovo.',
            'Zatvori',
            { duration: 3000, panelClass: ['error-snackbar'] }
          );
        }
      },
    });
  }

  startReportGeneration(): void {
    this.isSelectingReport = true;
    this.selectedOrder = null;

    document.querySelector('.main-container')?.classList.add('blur-active');
  }

  handleRowSelection(order: any): void {
    if (!this.isSelectingReport) return;

    this.selectedOrder = order;

    this.generateOrderReport(order.id);

    this.resetReportSelection();
  }

  resetReportSelection(): void {
    this.isSelectingReport = false;
    this.selectedOrder = null;
    document.querySelector('.main-container')?.classList.remove('blur-active');
  }

  generateOrderReport(orderId: number): void {
    console.log('Generating report for order ID:', orderId);
    this.documentService.generateReport(orderId).subscribe({
      next: (pdfBlob) => {
        console.log('PDF Blob received:', pdfBlob);

        const blobUrl = URL.createObjectURL(pdfBlob);
        window.open(blobUrl, '_blank');
      },
      error: (err) => {
        console.error('Greška pri generisanju izveštaja:', err);
      },
    });
  }
}
