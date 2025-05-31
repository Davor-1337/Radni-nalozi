import { Component, Inject, OnInit } from '@angular/core';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';
import { TechniciansService } from '../../../core/services/technicians.service';
import { MatTableModule } from '@angular/material/table';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-details',
  templateUrl: './details.component.html',
  styleUrls: ['./details.component.scss'],
  imports: [MatTableModule, CommonModule],
})
export class DetailsComponent implements OnInit {
  serviserId: number;
  workOrders: any[] = [];
  displayedColumns: string[] = [
    'icon',
    'Klijent',
    'OpisProblema',
    'Lokacija',
    'DatumDodjele',
    'DatumZavrsetka',
    'BrojRadnihSati',
  ];

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: { serviserId: number },
    private technicianService: TechniciansService
  ) {
    this.serviserId = data.serviserId;
  }

  ngOnInit(): void {
    this.fetchWorkOrders();
    this.workOrders.forEach((order) => {
      order.DatumZavrsetkaFormatted = this.formatDate(order.DatumZavrsetka);
    });
  }

  formatDate(datum: string | null | undefined): string {
    if (!datum || datum.includes('0001-01-01')) {
      return 'U toku';
    }

    const parsedDate = new Date(datum);
    return parsedDate.toLocaleDateString('en-GB'); // "dd/MM/yyyy" format
  }

  fetchWorkOrders(): void {
    this.technicianService
      .getWorkOrderDetailsByTechnicianId(this.serviserId)
      .subscribe({
        next: (data) => {
          this.workOrders = data;
          console.log(data);
        },
        error: (err) => {
          console.error('Greška prilikom učitavanja podataka:', err);
        },
      });
  }
}
