import { ComponentFixture, TestBed } from '@angular/core/testing';

import { WorkOrdersChartComponent } from './work-orders-chart.component';

describe('WorkOrdersChartComponent', () => {
  let component: WorkOrdersChartComponent;
  let fixture: ComponentFixture<WorkOrdersChartComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [WorkOrdersChartComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(WorkOrdersChartComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
