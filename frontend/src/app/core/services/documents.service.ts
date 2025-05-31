import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

export interface InvoicePreview {
  work_order_id: number;
  total_hours: number;
  labor_cost: number;
  materials: { name: string; quantity: number; price: number }[];
  total_cost: number;
}

export interface InvoiceFinalized {
  message: string;
  workOrderID: number;
  Iznos: number;
  DatumFakture: string;
}

@Injectable({
  providedIn: 'root',
})
export class DocumentsService {
  private apiUrl = 'http://localhost:8080/api';
  constructor(private http: HttpClient) {}

  generateInvoice(nalogId: number): Observable<InvoicePreview> {
    return this.http.get<InvoicePreview>(
      `${this.apiUrl}/fakture/generisi/${nalogId}`
    );
  }

  finalizeInvoice(nalogId: number): Observable<InvoiceFinalized> {
    return this.http.post<InvoiceFinalized>(
      `${this.apiUrl}/fakture/generisi/${nalogId}`,
      {}
    );
  }

  getInvoices(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/fakture`);
  }

  updateInvoice(id: number, updatedData: any): Observable<any> {
    console.log('ID fakture:', id);
    return this.http.put(`${this.apiUrl}/fakture/${id}`, updatedData, {
      headers: { 'Content-Type': 'application/json' },
    });
  }

  deleteInvoice(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/fakture/${id}`);
  }

  postInvoice(data: {
    Nalog_ID: Number;
    Iznos: number;
  }): Observable<{ faktura_id: number }> {
    return this.http.post<{ faktura_id: number }>(
      `${this.apiUrl}/fakture`,
      data
    );
  }

  getShortReport(): Observable<any[]> {
    return this.http.get<any[]>(
      `${this.apiUrl}/izvjestaji/radni-nalog/skraceni`
    );
  }

  getInvoicesByClient(clientId: number): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/fakture/klijent/${clientId}`);
  }

  getReportsByClient(clientId: number): Observable<any[]> {
    return this.http.get<any[]>(
      `${this.apiUrl}/izvjestaji/klijent/${clientId}`
    );
  }
  // REPORTS ///////////////
  generateReport(orderId: number): Observable<Blob> {
    return this.http.get(
      `${this.apiUrl}/izvjestaji/radni-nalog/pdf/${orderId}`,
      {
        responseType: 'blob',
      }
    );
  }
}
