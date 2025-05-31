import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Zahtjev {
  zahtjev_id: number;
  username: string;
  role: string;
  status: string;
  vrijemeKreiranja: string;
}

export interface WorkOrder {
  id: number;
  brojNaloga: string;
  prioritet: string;
  status: string;
  lokacija: string;
  datumOtvaranja: string;
  opisProblema: string;
  nazivKlijenta: string;
}

@Injectable({
  providedIn: 'root',
})
export class RequestsService {
  private apiUrl = 'http://localhost:8080/api';
  constructor(private http: HttpClient) {}

  getAll(): Observable<Zahtjev[]> {
    return this.http.get<Zahtjev[]>(`${this.apiUrl}/zahtjevi`);
  }

  handleRequest(zahtjevId: number, akcija: string): Observable<any> {
    const url = `${this.apiUrl}/zahtjevi/azuriraj`;
    const payload = {
      zahtjev_id: zahtjevId,
      akcija: akcija,
    };

    return this.http.post(url, payload);
  }

  getPendingWorkOrders(): Observable<WorkOrder[]> {
    return this.http.get<WorkOrder[]>(`${this.apiUrl}/radni-nalozi/na-cekanju`);
  }

  handleWorkOrder(
    nalogId: number,
    akcija: 'Otvoren' | 'Odbijen'
  ): Observable<any> {
    return this.http.post(`${this.apiUrl}/radni-nalozi/odaberi/${nalogId}`, {
      akcija,
    });
  }

  getTechnicians(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/serviseri`);
  }
}
