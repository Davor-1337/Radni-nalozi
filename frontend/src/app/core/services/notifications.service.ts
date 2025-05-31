// core/services/notification.service.ts
import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Obavestenje {
  obavestenje_id: number;
  korisnik_id: number;
  tip: string;
  poruka: string;
  datum_kreiranja: string;
  procitano: boolean;
}

@Injectable({ providedIn: 'root' })
export class NotificationService {
  private apiUrl = 'http://localhost:8080/api';

  getObavestenja(): Observable<Obavestenje[]> {
    return this.http.get<Obavestenje[]>(`${this.apiUrl}/obavjestenja`);
  }

  deleteNotification(id: number): Observable<any> {
    return this.http.delete(`${this.apiUrl}/obavjestenja/${id}`);
  }

  constructor(private http: HttpClient) {}
}
