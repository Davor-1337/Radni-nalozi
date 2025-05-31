import { Component } from '@angular/core';
import {
  MaterialServiceService,
  Material,
} from '../../../core/services/material-service.service';
import { MatTableModule } from '@angular/material/table';
import { CommonModule } from '@angular/common';
import { MatDialog } from '@angular/material/dialog';
import { AddComponent } from '../add/add.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ChartComponent } from '../chart/chart.component';
import { MatTableDataSource } from '@angular/material/table';
import { debounceTime } from 'rxjs/operators';
import { SearchService } from '../../../core/services/search.service';
@Component({
  selector: 'app-list',
  imports: [
    MatTableModule,
    CommonModule,
    ReactiveFormsModule,
    FormsModule,
    ChartComponent,
  ],
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
})
export class ListComponent {
  recentMaterials: any[] = [];
  dataSource = new MatTableDataSource<Material>([]);
  totalCount: number = 0;
  editMode: { row: any; field: string } | null = null;
  editValue: string = '';
  originalValue: string = '';
  displayedColumns: string[] = [
    'icon',
    'NazivMaterijala',
    'Cijena',
    'Kolicina',
  ];

  constructor(
    private materialService: MaterialServiceService,
    private searchService: SearchService,
    private dialog: MatDialog,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.fetchMaterials();

    this.searchService.searchTerm$
      .pipe(debounceTime(300))
      .subscribe((term) => this.applyFilter(term));
  }

  fetchMaterials(): void {
    this.materialService.getMaterials().subscribe({
      next: (data) => {
        this.dataSource.data = data;
      },
      error: (err) => console.error('Error fetching materials:', err),
    });
  }

  applyFilter(term: string) {
    if (!term) {
      return this.fetchMaterials();
    }
    this.materialService.filterMaterials({ search: term }).subscribe({
      next: (filtered) => {
        this.dataSource.data = filtered;
      },
      error: (err) => console.error('Error filtering materials:', err),
    });
  }

  openAddDialog() {
    const dialogRef = this.dialog.open(AddComponent, {
      width: '400px',
      height: '350px',
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
      Materijal_ID: row.Materijal_ID,
      NazivMaterijala:
        field === 'NazivMaterijala' ? this.editValue : row.NazivMaterijala,
      Cijena: field === 'Cijena' ? parseFloat(this.editValue) : row.Cijena,
      KolicinaUSkladistu:
        field === 'KolicinaUSkladistu'
          ? parseInt(this.editValue)
          : row.KolicinaUSkladistu,
    };

    this.materialService
      .updateMaterial(row.Materijal_ID, updatedData)
      .subscribe({
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
          this.snackBar.open('Greška pri ažuriranju materijala', 'Zatvori', {
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

  deleteMaterial(materialId: number, event: MouseEvent): void {
    event.stopPropagation();

    if (confirm('Da li ste sigurni da želite obrisati ovaj materijal?')) {
      this.materialService.deleteMaterial(materialId).subscribe({
        next: () => {
          this.recentMaterials = this.recentMaterials.filter(
            (material) => material.Materijal_ID !== materialId
          );
          this.snackBar.open('Materijal uspješno obrisan', 'Zatvori', {
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
}
