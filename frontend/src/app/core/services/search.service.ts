import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class SearchService {
  private searchTerm$$ = new BehaviorSubject<string>('');
  public searchTerm$ = this.searchTerm$$.asObservable();

  update(term: string) {
    this.searchTerm$$.next(term);
  }
}
