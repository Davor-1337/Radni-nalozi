import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { WorkOrderService } from '../../services/workOrder.service';

@Component({
  selector: 'app-carousel',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './carousel.component.html',
  styleUrl: './carousel.component.scss',
})
export class CarouselComponent {
  recentWorkOrders: any[] = [];

  constructor(private workOrderService: WorkOrderService) {}

  ngOnInit(): void {
    this.fetchWorkOrders();
  }

  fetchWorkOrders(): void {
    this.workOrderService.getActiveWorkOrders().subscribe({
      next: (data) => {
        this.recentWorkOrders = data;
      },
      error: (err) => {
        console.error('Error fetching work orders:', err);
      },
    });
  }

  activeIndex = 0;

  next() {
    if (this.recentWorkOrders.length > 0) {
      this.activeIndex = (this.activeIndex + 1) % this.recentWorkOrders.length;
    }
  }
}
