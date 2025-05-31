import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface Technician {
  Serviser_ID: number;
  Ime: string;
  Prezime: string;
  Specijalnost: string;
  Telefon: string;
  User_ID: number;
}

export interface WorkOrderDetails {
  OpisProblema: string;
  NazivKlijenta: string;
  Lokacija: string;
  DatumDodjele: Date;
  DatumZavrsetka: Date;
  BrojRadnihSati: number;
}
@Injectable({
  providedIn: 'root',
})
export class TechniciansService {
  private apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  getTechnicians(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/serviseri`);
  }

  addTechnician(technician: Technician): Observable<any> {
    return this.http.post(`${this.apiUrl}/serviseri`, technician);
  }

  getWorkOrdersByTechnicianId(serviserId: number): Observable<any[]> {
    return this.http.get<any[]>(
      `${this.apiUrl}/serviseri/${serviserId}/radni-nalozi`
    );
  }

  getWorkOrderDetailsByTechnicianId(
    serviserId: number
  ): Observable<WorkOrderDetails[]> {
    return this.http.get<WorkOrderDetails[]>(
      `${this.apiUrl}/serviseri/${serviserId}/radni-nalozi/details`
    );
  }

  getTotalHours(serviserId: number) {
    return this.http.get<number>(`${this.apiUrl}/serviseri/${serviserId}/sati`);
  }
}
