import { Component } from '@angular/core';
import { WorkOrderService } from '../../../core/services/workOrder.service';
import { MatTableModule } from '@angular/material/table';
import { CommonModule } from '@angular/common';
import { MatDialogRef } from '@angular/material/dialog';
@Component({
  selector: 'app-archive',
  templateUrl: './archive.component.html',
  styleUrls: ['./archive.component.scss'],
  imports: [MatTableModule, CommonModule],
})
export class ArchiveComponent {
  recentWorkOrders: any[] = [];
  totalCount: number = 0;
  displayedColumns: string[] = ['icon', 'WorkOrder', 'title', 'date', 'client'];

  constructor(
    private workOrderService: WorkOrderService,
    private dialogRef: MatDialogRef<ArchiveComponent>
  ) {}

  closeDialog(): void {
    this.dialogRef.close(); // Ovo zatvara dijalog
  }

  ngOnInit(): void {
    this.fetchWorkOrders();
  }

  fetchWorkOrders(): void {
    this.workOrderService.getArchivedOrders().subscribe({
      next: (data) => {
        this.recentWorkOrders = data;
        this.totalCount = data.length;
      },
      error: (err) => {
        console.error('Error fetching work orders:', err);
      },
    });
  }
}
