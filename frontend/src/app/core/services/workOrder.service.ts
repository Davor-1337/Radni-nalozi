import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { WorkOrder } from './requests.service';

interface WorkOrderStat {
  MonthNumber: number;
  Count: number;
}

export interface Filters {
  date_from: string;
  date_to: string;
  search: string;
}

@Injectable({
  providedIn: 'root',
})
export class WorkOrderService {
  private apiUrl = 'http://localhost:8080/api';

  constructor(private http: HttpClient) {}

  get4WorkOrders(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/radni-nalozi/4`);
  }

  getWorkOrderDetails(id: number): Observable<WorkOrder> {
    return this.http.get<any>(`${this.apiUrl}/radni-nalozi/${id}`).pipe(
      map((response) => ({
        id: response.Nalog_ID,
        brojNaloga: response.BrojNaloga,
        prioritet: response.Prioritet,
        status: response.Status,
        lokacija: response.Lokacija,
        datumOtvaranja: response.DatumOtvaranja,
        opisProblema: response.OpisProblema,
        nazivKlijenta: '',
      }))
    );
  }

  getActiveWorkOrders(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/radni-nalozi/aktivni`);
  }

  getWorkOrders(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/radni-nalozi`);
  }

  postWorkOrder(data: {
    klijent_id: number;
    opisProblema: string;
    prioritet: string;
    status: string;
    lokacija: string;
  }): Observable<{ status: string; nalog_id: number }> {
    return this.http.post<{ status: string; nalog_id: number }>(
      `${this.apiUrl}/radni-nalozi`,
      data
    );
  }

  getWorkOrderStatusCount(): Observable<{
    completed: number;
    inProgress: number;
  }> {
    return this.http.get<{ completed: number; inProgress: number }>(
      `${this.apiUrl}/radni-nalozi/status-count`
    );
  }

  deleteWorkOrder(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/radni-nalozi/${id}`);
  }

  updateWorkOrder(id: number, updatedData: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/radni-nalozi/${id}`, updatedData, {
      headers: { 'Content-Type': 'application/json' },
    });
  }

  getTotalWorkOrderCount(): Observable<{ ukupno: number }> {
    return this.http.get<{ ukupno: number }>(
      `${this.apiUrl}/radni-nalozi/ukupno`
    );
  }

  getArchivedOrders(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/radni-nalozi/arhiva`);
  }

  getWorkOrderStats(): Observable<WorkOrderStat[]> {
    return this.http.get<WorkOrderStat[]>(`${this.apiUrl}/radni-nalozi/stats`);
  }
  filterWorkOrders(filters: Filters) {
    return this.http
      .post<{ work_orders: any[] }>(
        `${this.apiUrl}/radni-nalozi/filter`,
        filters
      )
      .pipe(map((res) => res.work_orders));
  }

  archiveWorkOrder(nalogId: number): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(
      `${this.apiUrl}/radni-nalozi/arhiva`,
      {
        nalog_id: nalogId,
      }
    );
  }

  updateStatus(
    nalogId: number,
    akcija: 'Otvoren' | 'Odbijen'
  ): Observable<any> {
    return this.http.put(`${this.apiUrl}/radni-nalozi/status/${nalogId}`, {
      work_order_id: nalogId,
      Status: akcija,
    });
  }

  assignWorkOrder(data: {
    Nalog_ID: number;
    Serviser_ID: number;
  }): Observable<any> {
    return this.http.post(`${this.apiUrl}/radni-nalozi/dodjela`, data);
  }

  //technician-orders
  getOrdersByTehnician(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/serviseri/radni-nalozi`);
  }

  addMaterial(
    workOrderId: number,
    material: { materijal_ID: number; kolicinaUtrosena: number }
  ): Observable<any> {
    return this.http.post(
      `${this.apiUrl}/radni-nalozi/${workOrderId}/materijal`,
      material
    );
  }

  addHours(
    workOrderId: number,
    hours: { serviser_ID: number; brojRadnihSati: number }
  ): Observable<any> {
    return this.http.post(
      `${this.apiUrl}/radni-nalozi/${workOrderId}/sati`,
      hours
    );
  }

  finishWorkOrder(workOrderId: number): Observable<any> {
    const url = `${this.apiUrl}/radni-nalozi/${workOrderId}/zavrsi`;
    return this.http.put(url, {});
  }

  //client-orders
  getOrdersByClient(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/klijenti/radni-nalozi`);
  }
}
