import { Component } from '@angular/core';
import { TechniciansService } from '../../../core/services/technicians.service';
import { CommonModule } from '@angular/common';
import { MatDialog } from '@angular/material/dialog';
import { MatSnackBar } from '@angular/material/snack-bar';
import { AddComponent } from '../add/add.component';
import { MatTableModule } from '@angular/material/table';
import { DetailsComponent } from '../details/details.component';

@Component({
  selector: 'app-list',
  imports: [CommonModule, MatTableModule],
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
})
export class ListComponent {
  technicians: any[] = [];
  focusedTechnician: number | null = null;
  selectedServiserId: number | null = null;
  displayedColumnsOrders: string[] = [
    'icon',
    'BrojNaloga',
    'Klijent',
    'OpisProblema',
  ];
  technicianWorkOrders: any[] = [];
  totalHours: number = 0;

  constructor(
    private technicianService: TechniciansService,
    private dialog: MatDialog,
    private snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.fetchTechnicians();
    this.loadTechnicianWorkOrders(100);
    this.fetchTotalHours(100);
  }

  fetchTechnicians(): void {
    this.technicianService.getTechnicians().subscribe(
      (data) => {
        this.technicians = data;
        console.log('Podaci sa backenda:', data);
      },
      (error) => {
        console.error('Error loading technicians:', error);
      }
    );
  }
  setFocus(userId: number) {
    if (this.focusedTechnician === userId) {
      this.focusedTechnician = null;
    } else {
      this.focusedTechnician = userId;
    }
  }

  fetchTotalHours(serviserId: number): void {
    this.technicianService.getTotalHours(serviserId).subscribe({
      next: (response: any) => {
        this.totalHours = response.totalHours;
        console.log('Ukupni sati:', this.totalHours);
      },
      error: (error) => {
        console.error('Greška pri dohvaćanju sati:', error);
      },
    });
  }

  loadTechnicianWorkOrders(serviserId: number): void {
    this.technicianService.getWorkOrdersByTechnicianId(serviserId).subscribe({
      next: (workOrders) => {
        this.technicianWorkOrders = workOrders;
        console.log(workOrders);
      },
      error: (error) => {
        console.error('Greška prilikom učitavanja radnih naloga:', error);
      },
    });
  }

  openDialog() {
    if (this.selectedServiserId !== null) {
      this.dialog.open(DetailsComponent, {
        width: '1200px',
        height: '500px',
        maxWidth: 'none',
        data: { serviserId: this.selectedServiserId },
      });
    }
  }

  onServiserClick(serviserId: number): void {
    this.selectedServiserId = serviserId;
    this.loadTechnicianWorkOrders(serviserId);
    this.fetchTotalHours(serviserId); // <-- Dodaj ovo ovdje
  }

  openAddDialog() {
    const dialogRef = this.dialog.open(AddComponent, {
      width: '400px',
      height: '450px',
    });

    dialogRef.afterClosed().subscribe((result) => {
      if (result) {
        console.log('Novi materijal dodat:', result);
      }
    });
  }
}
