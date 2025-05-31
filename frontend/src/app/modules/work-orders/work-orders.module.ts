import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { WorkOrdersRoutingModule } from './work-orders-routing.module';
import { ListComponent } from './list/list.component';
import { AddComponent } from './add/add.component';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    WorkOrdersRoutingModule,
    RouterModule,
    ListComponent,
    AddComponent,
  ],
})
export class WorkOrdersModule {}
