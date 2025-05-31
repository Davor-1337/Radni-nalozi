import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../services/auth.service';
import { Router, RouterModule, NavigationEnd } from '@angular/router';
import { SearchService } from '../../services/search.service';
import { filter } from 'rxjs/operators';
import { MatDialog } from '@angular/material/dialog';
import { UserDialogComponent } from '../user-dialog/user-dialog.component';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.scss'],
  imports: [CommonModule, RouterModule],
})
export class HeaderComponent {
  currentDate: Date = new Date();
  username: string | null = null;
  currentUrl: string = '';

  private hideSearchOn: string[] = [
    '/fakture',
    '/serviseri',
    '/obavjestenja',
    '/home',
  ];
  ngOnInit(): void {
    this.username = localStorage.getItem('username');
    this.currentUrl = this.router.url;
  }

  constructor(
    private authService: AuthService,
    private router: Router,
    private searchService: SearchService,
    private dialog: MatDialog
  ) {
    this.router.events
      .pipe(filter((e) => e instanceof NavigationEnd))
      .subscribe((e: NavigationEnd) => {
        this.currentUrl = e.urlAfterRedirects;
      });
  }

  openUserDialog() {
    this.dialog.open(UserDialogComponent, {
      data: {
        username: this.username,
        role: this.authService.getRole(),
      },
    });
  }

  logout() {
    this.authService.logout();
    this.router.navigate(['/login']);
  }

  onSearchChange(term: string) {
    console.log('[Header] search term:', term);
    this.searchService.update(term);
  }

  isAdmin(): boolean {
    const role = this.authService.getRole();
    return role === 'admin';
  }

  shouldShowSearch(): boolean {
    if (!this.isAdmin()) {
      return false;
    }
    return !this.hideSearchOn.includes(this.currentUrl);
  }
}
