import { TestBed } from '@angular/core/testing';

import { TehniciansService } from './tehnicians.service';

describe('TehniciansService', () => {
  let service: TehniciansService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TehniciansService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
