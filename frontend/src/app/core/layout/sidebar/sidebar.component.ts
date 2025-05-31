import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router } from '@angular/router';
import { SidebarService } from '../../services/sidebar.service';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.component.html',
  imports: [CommonModule, RouterModule],
  styleUrls: ['./sidebar.component.scss'],
})
export class SidebarComponent {
  activeItem: any = null;
  menuItems = [
    { name: 'Home', icon: 'fa-solid fa-house', href: '/home' },
    {
      name: 'Work Orders',
      icon: 'fa-solid fa-newspaper',
      href: '/radni-nalozi',
    },
    { name: 'Materials', icon: 'fa-solid fa-toolbox', href: '/materijali' },

    { name: 'Finance documents', icon: 'fa-solid fa-print', href: '/fakture' },
    { name: 'Clients', icon: 'fa-solid fa-user-tie', href: '/klijenti' },
    { name: 'Tehnicians', icon: 'fa-solid fa-user-group', href: '/serviseri' },
    { name: 'Notifications', icon: 'fa-solid fa-bell', href: '/obavjestenja' },
  ];

  constructor(private router: Router, private sidebarService: SidebarService) {}

  ngOnInit() {
    this.loadMenuItems();
  }

  loadMenuItems() {
    this.menuItems = this.sidebarService.getMenuItems();
  }

  setActive(item: any) {
    this.activeItem = item;
    this.router.navigate([item.href]);
  }
}
