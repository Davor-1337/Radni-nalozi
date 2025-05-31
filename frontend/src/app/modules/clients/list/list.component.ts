import { Component } from '@angular/core';
import { ClientsService, Client } from '../../../core/services/clients.service';
import { MatTableModule } from '@angular/material/table';
import { CommonModule } from '@angular/common';
import { MatDialog } from '@angular/material/dialog';
import { AddComponent } from '../add/add.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { OrdersComponent } from '../orders/orders.component';
import { SearchService } from '../../../core/services/search.service';
import { MatTableDataSource } from '@angular/material/table';
import { debounceTime } from 'rxjs/operators';

@Component({
  selector: 'app-list',
  imports: [MatTableModule, CommonModule, ReactiveFormsModule, FormsModule],
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
})
export class ListComponent {
  clients: any[] = [];
  dataSource = new MatTableDataSource<Client>([]);
  totalCount: number = 0;
  editMode: { row: any; field: string } | null = null;
  editValue: string = '';
  originalValue: string = '';
  displayedColumns: string[] = [
    'icon',
    'Client',
    'Contact',
    'Email',
    'Tel',
    'Address',
    'AccountID',
  ];

  constructor(
    private clientService: ClientsService,
    private search: SearchService,
    private dialog: MatDialog,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.fetchClients();

    this.search.searchTerm$
      .pipe(debounceTime(300))
      .subscribe((term) => this.applyFilter(term));
  }

  fetchClients() {
    this.clientService.getClients().subscribe((data) => {
      this.dataSource.data = data;
    });
  }

  applyFilter(term: string) {
    if (!term) {
      return this.fetchClients();
    }
    this.clientService.filterClients({ search: term }).subscribe((data) => {
      this.dataSource.data = data;
    });
  }

  openAddDialog() {
    const dialogRef = this.dialog.open(AddComponent, {
      width: '500px',
      height: '650px',
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        console.log('Novi klijent dodan:', result);
      }
    });
  }

  openOrdersDialog(clientId: number, clientName: string): void {
    const dialogRef = this.dialog.open(OrdersComponent, {
      width: '800px',
      height: '600px',
      data: {
        clientId: clientId,
        clientName: clientName,
      },
      panelClass: 'archive-dialog',
    });
  }

  deleteClient(clientId: number, event: MouseEvent): void {
    event.stopPropagation();

    if (confirm('Da li ste sigurni da želite obrisati klijenta?')) {
      this.clientService.deleteClient(clientId).subscribe({
        next: () => {
          this.clients = this.clients.filter(
            (client) => client.Klijent_ID !== clientId
          );
          this.snackBar.open('Klijent obrisan', 'Zatvori', {
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
      Klijent_ID: row.Klijent_ID,
      Naziv: field === 'Client' ? this.editValue : row.Naziv,
      KontaktOsoba: field === 'Contact' ? this.editValue : row.KontaktOsoba,
      Email: field === 'Email' ? this.editValue : row.Email,
      Tel: field === 'Tel' ? this.editValue : row.Tel,
      Adresa: field === 'Address' ? this.editValue : row.Adresa,
      User_ID: row.User_ID,
    };

    this.clientService.updateClient(row.Klijent_ID, updatedData).subscribe({
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
        this.snackBar.open('Greška pri ažuriranju naloga', 'Zatvori', {
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
}
