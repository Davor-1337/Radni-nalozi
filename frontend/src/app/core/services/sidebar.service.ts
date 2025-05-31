import { Injectable } from '@angular/core';
import { AuthService } from './auth.service';

@Injectable({ providedIn: 'root' })
export class SidebarService {
  constructor(private auth: AuthService) {}

  getMenuItems() {
    const role = this.auth.getCurrentUserRole();

    const allItems = [
      {
        name: 'Home',
        icon: 'fa-solid fa-house',
        href: '/home',
        roles: ['admin', 'serviser'],
      },
      {
        name: 'Radni nalozi',
        icon: 'fa-solid fa-newspaper',
        href: '/radni-nalozi',
        roles: ['admin', 'serviser', 'klijent'],
      },

      {
        name: 'Materijali',
        icon: 'fa-solid fa-toolbox',
        href: '/materijali',
        roles: ['admin', 'serviser'],
      },
      {
        name: 'Fakture',
        icon: 'fa-solid fa-print',
        href: '/fakture',
        roles: ['admin', 'klijent'],
      },
      // Admin-only stavke
      ...(role === 'admin'
        ? [
            {
              name: 'Klijenti',
              icon: 'fa-solid fa-user-tie',
              href: '/klijenti',
              roles: ['admin'],
            },
            {
              name: 'Serviseri',
              icon: 'fa-solid fa-user-group',
              href: '/serviseri',
              roles: ['admin'],
            },
          ]
        : []),
      // Opšte stavke
      {
        name: 'Obaveštenja',
        icon: 'fa-solid fa-bell',
        href: '/obavjestenja',
        roles: ['admin', 'serviser', 'klijent'],
      },
    ];

    return allItems.filter((item) => item.roles.includes(role));
  }
}
