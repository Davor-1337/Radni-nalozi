import { Component, OnInit } from '@angular/core';
import { MatTableModule } from '@angular/material/table';
import { CommonModule } from '@angular/common';
import { CarouselComponent } from '../carousel/carousel.component';
import { WorkOrdersChartComponent } from '../work-orders-chart/work-orders-chart.component';
import { WorkOrderService } from '../../services/workOrder.service';

@Component({
  selector: 'app-home-layout',
  standalone: true,
  imports: [
    CommonModule,
    MatTableModule,
    CarouselComponent,
    WorkOrdersChartComponent,
  ],
  templateUrl: './home-layout.component.html',
  styleUrl: './home-layout.component.scss',
})
export class HomeLayoutComponent implements OnInit {
  completedOrders: number = 0;
  inProgressOrders: number = 0;
  recentWorkOrders: any[] = [];
  recentWorkOrders2: any[] = [];
  totalCount: number = 0;
  displayedColumns: string[] = [
    'icon',
    'WorkOrder',
    'title',
    'priority',
    'status',
    'date',
    'location',
    'client',
  ];

  constructor(private workOrderService: WorkOrderService) {}

  ngOnInit(): void {
    this.fetch4WorkOrders();
    this.fetchWorkOrders();
    this.fetchWorkOrderStatus();
    this.fetchTotalWorkOrderCount();
  }

  fetchWorkOrders(): void {
    this.workOrderService.getWorkOrders().subscribe({
      next: (response) => {
        this.recentWorkOrders2 = response;
      },
      error: (err) => {
        console.error('Error fetching work orders:', err);
      },
    });
  }

  fetch4WorkOrders(): void {
    this.workOrderService.get4WorkOrders().subscribe({
      next: (data) => {
        this.recentWorkOrders = data;
      },
      error: (err) => {
        console.error('Error fetching work orders:', err);
      },
    });
  }

  fetchTotalWorkOrderCount(): void {
    this.workOrderService.getTotalWorkOrderCount().subscribe({
      next: (data) => {
        console.log('Odgovor za total count:', data);
        this.totalCount = data.ukupno;
      },
      error: (err) => {
        console.error('Error fetching total work order count:', err);
      },
    });
  }

  fetchWorkOrderStatus(): void {
    this.workOrderService.getWorkOrderStatusCount().subscribe({
      next: (data) => {
        this.completedOrders = data.completed;
        this.inProgressOrders = data.inProgress;
      },
      error: (err) => {
        console.error('Error fetching work order status:', err);
      },
    });
  }
}
