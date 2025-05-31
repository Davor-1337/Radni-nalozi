import { Component, ViewChild, ChangeDetectorRef } from '@angular/core';
import { WorkOrderService } from '../../../core/services/workOrder.service';
import { MatTableModule } from '@angular/material/table';
import { CommonModule } from '@angular/common';
import { MatDialog, MatDialogRef } from '@angular/material/dialog';
import { AddComponent } from '../add/add.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ArchiveComponent } from '../archive/archive.component';
import { debounceTime } from 'rxjs/operators';
import { SearchService } from '../../../core/services/search.service';
import { MatTable } from '@angular/material/table';
import { MatTableDataSource } from '@angular/material/table';
import { AuthService } from '../../../core/services/auth.service';
import { FinishOrderComponent } from '../finish-order/finish-order.component';

@Component({
  selector: 'app-list',
  standalone: true,
  imports: [MatTableModule, CommonModule, ReactiveFormsModule, FormsModule],
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
})
export class ListComponent {
  @ViewChild('table', { static: true }) table!: MatTable<any>;
  recentWorkOrders: any[] = [];
  totalCount: number = 0;
  dataSource = new MatTableDataSource<any>([]);
  editMode: { row: any; field: string } | null = null;
  editValue: string = '';
  originalValue: string = '';
  displayedColumns: string[] = [];
  public role!: string;

  constructor(
    private workOrderService: WorkOrderService,
    private searchService: SearchService,
    private authService: AuthService,
    private dialog: MatDialog,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.role = this.authService.getRole();

    if (this.role === 'admin' || this.role === 'serviser') {
      this.displayedColumns = [
        'icon',
        'WorkOrder',
        'id',
        'title',
        'priority',
        'status',
        'date',
        'location',
        'archiveIcon',
      ];
    } else {
      // klijent
      this.displayedColumns = [
        'icon',
        'tipNaloga',
        'WorkOrder',
        'title',
        'status',
        'date',
      ];
    }

    if (this.role === 'admin') {
      this.fetchWorkOrders();
      this.searchService.searchTerm$
        .pipe(debounceTime(300))
        .subscribe((term) => this.applyFilter(term));
    } else if (this.role === 'serviser') {
      this.fetchTehnicianWorkOrders();
    } else {
      this.fetchClientWorkOrders();
    }
  }

  openFinishDialog(nalogId: number) {
    const serviserId = this.authService.getUserId();

    const dialogRef = this.dialog.open(FinishOrderComponent, {
      width: '500px',
      data: { nalogId, serviserId },
    });

    dialogRef.afterClosed().subscribe((finished) => {
      if (finished) {
        this.fetchTehnicianWorkOrders();
      }
    });
  }

  archiveOrder(element: any, event: MouseEvent): void {
    event.stopPropagation();

    if (element.Status !== 'Zavrsen') {
      this.snackBar.open(
        'Ne može se arhivirati — nalog nije završen.',
        'Zatvori',
        {
          duration: 3000,
          panelClass: ['error-snackbar'],
        }
      );
      return;
    }

    this.workOrderService.archiveWorkOrder(element.Nalog_ID).subscribe({
      next: () => {
        if (this.role === 'admin') {
          this.fetchWorkOrders();
        } else if (this.role === 'serviser') {
          this.fetchTehnicianWorkOrders();
        } else {
          this.fetchClientWorkOrders();
        }
        this.snackBar.open('Nalog je arhiviran', 'Zatvori', {
          duration: 2000,
          panelClass: ['success-snackbar'],
        });
      },
      error: () => {
        this.snackBar.open('Greška pri arhiviranju', 'Zatvori', {
          duration: 3000,
          panelClass: ['error-snackbar'],
        });
      },
    });
  }

  fetchWorkOrders() {
    this.workOrderService.getWorkOrders().subscribe((data) => {
      console.log('[List] dataSource.data =', this.dataSource.data);
      this.recentWorkOrders = data;
      this.dataSource.data = data;
      this.totalCount = data.length;
    });
  }
  applyFilter(term: string) {
    if (!term) {
      return this.fetchWorkOrders();
    }

    this.workOrderService
      .filterWorkOrders({
        date_from: '',
        date_to: '',
        search: term,
      })
      .subscribe((data) => {
        this.recentWorkOrders = data;
        this.dataSource.data = data;
        this.totalCount = data.length;
      });
  }
  fetchTehnicianWorkOrders(): void {
    this.workOrderService.getOrdersByTehnician().subscribe({
      next: (data) => {
        this.dataSource.data = data;
      },
      error: (error) => {
        console.error('Error fetching work orders for serviser:', error);
      },
    });
  }

  fetchClientWorkOrders(): void {
    this.workOrderService.getOrdersByClient().subscribe({
      next: (data) => {
        this.dataSource.data = data;
      },
      error: (error) => {
        console.error('Error fetching work orders for serviser:', error);
      },
    });
  }

  openAddDialog() {
    let dialogRef: MatDialogRef<AddComponent>;

    if (this.role === 'admin') {
      dialogRef = this.dialog.open(AddComponent, {
        width: '500px',
        height: '500px',
        data: { role: this.role },
      });
    } else if (this.role === 'klijent') {
      dialogRef = this.dialog.open(AddComponent, {
        width: '400px',
        height: '400px',
        data: { role: this.role },
      });
    } else {
      return;
    }

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        console.log('Novi nalog dodan:', result);
      }
    });
  }

  deleteOrder(orderId: number, event: MouseEvent): void {
    event.stopPropagation();

    if (confirm('Da li ste sigurni da želite obrisati ovaj nalog?')) {
      this.workOrderService.deleteWorkOrder(orderId).subscribe({
        next: () => {
          this.recentWorkOrders = this.recentWorkOrders.filter(
            (order) => order.Nalog_ID !== orderId
          );
          this.snackBar.open('Nalog obrisan', 'Zatvori', {
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
    if (this.role === 'klijent') return;
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
      Nalog_ID: row.Nalog_ID,
      Klijent_ID: row.Klijent_ID,
      OpisProblema:
        field === 'OpisProblema' ? this.editValue : row.OpisProblema,
      Prioritet: field === 'Prioritet' ? this.editValue : row.Prioritet,
      DatumOtvaranja: row.DatumOtvaranja,
      Status: field === 'Status' ? this.editValue : row.Status,
      Lokacija: field === 'Lokacija' ? this.editValue : row.Lokacija,
      BrojNaloga: row.BrojNaloga,
    };

    this.workOrderService.updateWorkOrder(row.Nalog_ID, updatedData).subscribe({
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

  openArchive(): void {
    const dialogRef = this.dialog.open(ArchiveComponent, {
      width: '800px',
      height: '50%',
      panelClass: 'archive-dialog',
      data: {},
    });
  }
}
