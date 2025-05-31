import { Routes } from '@angular/router';
import { LayoutComponent } from './core/layout/layout.component';
import { HomeLayoutComponent } from './core/layout/home-layout/home-layout.component';
import { LoginComponent } from './modules/auth/login/login.component';
import { roleGuard } from './guards/role.guard';

export const routes: Routes = [
  { path: 'login', component: LoginComponent },

  { path: '', redirectTo: 'login', pathMatch: 'full' },

  {
    path: '',
    component: LayoutComponent,
    children: [
      {
        path: 'home',
        component: HomeLayoutComponent,
        canActivate: [roleGuard],
        data: { roles: ['admin', 'serviser'] },
      },
      {
        path: 'radni-nalozi',
        loadChildren: () =>
          import('./modules/work-orders/work-orders.module').then(
            (m) => m.WorkOrdersModule
          ),
        canActivate: [roleGuard],
        data: { roles: ['admin', 'serviser', 'klijent'] },
      },
      {
        path: 'materijali',
        loadChildren: () =>
          import('./modules/materials/materials.module').then(
            (m) => m.MaterialsModule
          ),
        canActivate: [roleGuard],
        data: { roles: ['admin', 'serviser'] },
      },
      {
        path: 'fakture',
        loadChildren: () =>
          import('./modules/documents/documents.module').then(
            (m) => m.DocumentsModule
          ),
        canActivate: [roleGuard],
        data: { roles: ['admin', 'klijent'] },
      },
      {
        path: 'klijenti',
        loadChildren: () =>
          import('./modules/clients/clients.module').then(
            (m) => m.ClientsModule
          ),
        canActivate: [roleGuard],
        data: { roles: ['admin'] },
      },
      {
        path: 'serviseri',
        loadChildren: () =>
          import('./modules/technicians/technicians.module').then(
            (m) => m.TechniciansModule
          ),
        canActivate: [roleGuard],
        data: { roles: ['admin'] },
      },
      {
        path: 'obavjestenja',
        loadChildren: () =>
          import('./modules/notifications/notifications.module').then(
            (m) => m.NotificationsModule
          ),
        canActivate: [roleGuard],
        data: { roles: ['admin', 'serviser', 'klijent'] },
      },
    ],
  },
];
