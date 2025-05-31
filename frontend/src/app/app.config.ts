import { ApplicationConfig, provideZoneChangeDetection } from '@angular/core';
import { provideRouter } from '@angular/router';
import { routes } from './app.routes';
import { SidebarComponent } from './core/layout/sidebar/sidebar.component';
import { HeaderComponent } from './core/layout/header/header.component';
import { importProvidersFrom } from '@angular/core';
import { CommonModule } from '@angular/common';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { WorkOrderService } from './core/services/workOrder.service';
import { MaterialServiceService } from './core/services/material-service.service';
import { DocumentsService } from './core/services/documents.service';
import { authInterceptor } from './core/services/auth.interceptor';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { MatSnackBarModule } from '@angular/material/snack-bar';
import { provideAnimations } from '@angular/platform-browser/animations';
import { JwtHelperService } from '@auth0/angular-jwt';
import { registerLocaleData } from '@angular/common';
import localeBs from '@angular/common/locales/bs';
import { LOCALE_ID } from '@angular/core';

registerLocaleData(localeBs, 'bs-BA');

export const appConfig: ApplicationConfig = {
  providers: [
    { provide: JwtHelperService, useFactory: () => new JwtHelperService() },
    provideHttpClient(withInterceptors([authInterceptor])),
    importProvidersFrom(CommonModule),
    provideZoneChangeDetection({ eventCoalescing: true }),
    provideRouter(routes),
    SidebarComponent,
    { provide: LOCALE_ID, useValue: 'bs-BA' },
    HeaderComponent,
    WorkOrderService,
    DocumentsService,
    MaterialServiceService,
    importProvidersFrom(FormsModule, ReactiveFormsModule),
    importProvidersFrom(MatSnackBarModule),
    provideAnimations(),
  ],
};
