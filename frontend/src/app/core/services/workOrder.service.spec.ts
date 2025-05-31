import { TestBed } from '@angular/core/testing';

import { WorkOrderService } from './workOrder.service';

describe('ApiService', () => {
  let service: WorkOrderService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(WorkOrderService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
