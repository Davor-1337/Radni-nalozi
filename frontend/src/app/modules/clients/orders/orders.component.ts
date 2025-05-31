import { Component, Inject, OnInit } from '@angular/core';
import { ClientsService } from '../../../core/services/clients.service';
import { MatTableModule } from '@angular/material/table';
import { CommonModule } from '@angular/common';
import { MAT_DIALOG_DATA } from '@angular/material/dialog';

interface WorkOrder {
  brojNaloga: string;
  datum: string;
  opisProblema: string;
  status: string;
  tipNaloga: string;
}

@Component({
  selector: 'app-orders',
  imports: [MatTableModule, CommonModule],
  templateUrl: './orders.component.html',
  styleUrl: './orders.component.scss',
})
export class OrdersComponent {
  workOrders: WorkOrder[] = [];
  displayedColumns: string[] = [
    'icon',
    'brojNaloga',
    'datum',
    'opisProblema',
    'status',
    'tipNaloga',
  ];

  constructor(
    private clientsService: ClientsService,
    @Inject(MAT_DIALOG_DATA)
    public data: { clientId: number; clientName: string }
  ) {
    console.log('Dialog komponenta kreirana sa clientId:', data.clientId);
  }

  ngOnInit(): void {
    console.log('ngOnInit pozvan');
    this.fetchWorkOrders(this.data.clientId);
  }

  fetchWorkOrders(clientId: number): void {
    console.log('Pozivam fetchWorkOrders sa:', clientId);
    this.clientsService.getClientWithOrders(clientId).subscribe({
      next: (response) => {
        this.workOrders = response;
        console.log('WORKORDERS:', this.workOrders);
      },
      error: (err) => {
        console.error('Greška pri dohvatu naloga:', err);
      },
    });
  }
}
