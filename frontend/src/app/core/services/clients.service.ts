import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

export interface Client {
  Klijent_ID: number;
  Naziv: string;
  KontaktOsoba: string;
  Email: string;
  Tel: string;
  Adresa: string;
  User_ID: number;
}

export interface ClientFilters {
  search: string;
}

@Injectable({
  providedIn: 'root',
})
export class ClientsService {
  private apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  getClients(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/klijenti`);
  }

  postClient(data: {
    Klijent_ID: number;
    Naziv: string;
    KontaktOsoba: string;
    Email: string;
    Tel: string;
    Adresa: string;
    User_ID: number;
  }): Observable<{ Klijent_ID: number }> {
    return this.http.post<{ Klijent_ID: number }>(
      `${this.apiUrl}/klijenti`,
      data
    );
  }

  deleteClient(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/klijenti/${id}`);
  }

  updateClient(id: number, updatedData: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/klijenti/${id}`, updatedData, {
      headers: { 'Content-Type': 'application/json' },
    });
  }

  getClientWithOrders(clientId: number): Observable<any> {
    return this.http.get(`${this.apiUrl}/klijenti/${clientId}/nalozi`);
  }

  filterClients(filters: ClientFilters) {
    return this.http
      .post<{ clients: Client[] }>(`${this.apiUrl}/klijenti/filter`, filters)
      .pipe(map((res) => res.clients));
  }
}
