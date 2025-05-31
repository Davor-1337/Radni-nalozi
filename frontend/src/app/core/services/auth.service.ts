import { Injectable, inject } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { JwtHelperService } from '@auth0/angular-jwt';
import { HttpHeaders } from '@angular/common/http';
import { tap } from 'rxjs/operators';
@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private apiUrl = 'http://localhost:8080/api';
  private jwtHelper = inject(JwtHelperService);
  constructor(private http: HttpClient) {}

  login(credentials: { username: string; password: string }): Observable<any> {
    return this.http.post(`${this.apiUrl}/login`, credentials).pipe(
      tap((res: any) => {
        localStorage.setItem('token', res.token);
        localStorage.setItem('username', res.username);
      })
    );
  }

  logout(): void {
    localStorage.removeItem('token');
  }

  isAuthenticated(): boolean {
    return !!localStorage.getItem('token');
  }

  signUp(data: any): Observable<any> {
    return this.http.post(`${this.apiUrl}/zahtjevi`, data);
  }

  getCurrentUserRole(): string {
    const token = localStorage.getItem('token');
    if (!token) return '';

    const decoded = this.jwtHelper.decodeToken(token);
    return decoded?.Role || decoded?.role || '';
  }

  getDecodedToken() {
    const token = localStorage.getItem('token');
    if (!token) return null;

    const payload = token.split('.')[1];
    return JSON.parse(atob(payload));
  }

  getRole(): string {
    const decoded = this.getDecodedToken();
    return decoded?.Role;
  }

  getUserId(): number {
    const decoded = this.getDecodedToken();
    return decoded?.User_ID;
  }

  changePassword(oldPassword: string, newPassword: string): Observable<any> {
    const token = localStorage.getItem('token');
    const headers = new HttpHeaders({
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    });

    return this.http.put(
      `${this.apiUrl}/updatePassword`,
      {
        old_password: oldPassword,
        new_password: newPassword,
      },
      { headers }
    );
  }
}
