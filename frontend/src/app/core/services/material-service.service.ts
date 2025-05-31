import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

export interface Material {
  Materijal_ID: number;
  NazivMaterijala: string;
  Kategorija: string;
  Cijena: number;
  KolicinaUSkladistu: number;
}

export interface MaterialFilters {
  search: string;
}

@Injectable({
  providedIn: 'root',
})
export class MaterialServiceService {
  private apiUrl = 'http://localhost:8080/api';
  constructor(private http: HttpClient) {}

  getMaterials(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/materijali`);
  }

  getMaterialsStats(): Observable<any[]> {
    return this.http.get<any[]>(`${this.apiUrl}/materijali/kategorija`);
  }

  updateMaterial(id: number, updatedData: any): Observable<any> {
    return this.http.put(`${this.apiUrl}/materijali/${id}`, updatedData, {
      headers: { 'Content-Type': 'application/json' },
    });
  }

  deleteMaterial(id: number): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/materijali/${id}`);
  }

  postMaterial(data: {
    NazivMaterijala: string;
    Cijena: number;
    KolicinaUSkladistu: number;
  }): Observable<{ materijal_id: number }> {
    return this.http.post<{ materijal_id: number }>(
      `${this.apiUrl}/materijali`,
      data
    );
  }

  filterMaterials(filters: MaterialFilters) {
    return this.http
      .post<{ materials: Material[] }>(
        `${this.apiUrl}/materijali/filter`,
        filters
      )
      .pipe(map((res) => res.materials));
  }
}
