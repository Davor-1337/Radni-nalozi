import { Component } from '@angular/core';
import {
  RequestsService,
  WorkOrder,
} from '../../../core/services/requests.service';
import { Zahtjev } from '../../../core/services/requests.service';
import { CommonModule } from '@angular/common';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialog } from '@angular/material/dialog';
import { AssignDialogComponent } from '../assign-dialog/assign-dialog.component';
import { WorkOrderService } from '../../../core/services/workOrder.service';
import { TechniciansService } from '../../../core/services/technicians.service';
import { AuthService } from '../../../core/services/auth.service';

import {
  Obavestenje,
  NotificationService,
} from '../../../core/services/notifications.service';

@Component({
  selector: 'app-list',
  imports: [CommonModule],
  templateUrl: './list.component.html',
  styleUrl: './list.component.scss',
})
export class ListComponent {
  zahtjevi: Zahtjev[] = [];
  radniNalozi: WorkOrder[] = [];
  obavestenja: Obavestenje[] = [];
  isAdmin = false;
  isTechnician = false;
  isClient = false;
  selectedWorkOrderDetails: WorkOrder | null = null;
  selectedNotificationId: number | null = null;
  constructor(
    private requestsService: RequestsService,
    private snackBar: MatSnackBar,
    private dialog: MatDialog,
    private workOrderService: WorkOrderService,
    private technicianService: TechniciansService,
    private notifService: NotificationService,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    const role = this.authService.getRole();
    this.isAdmin = role === 'admin';
    this.isTechnician = role === 'serviser';
    this.isClient = role === 'klijent';

    if (this.isAdmin) {
      this.loadRequests();
      this.loadPendingOrders();
    } else if (this.isTechnician || this.isClient) {
      this.loadNotifications();
    }
    console.log(this.obavestenja);
  }

  loadPendingOrders(): void {
    this.requestsService.getPendingWorkOrders().subscribe((data) => {
      this.radniNalozi = data;
    });
  }

  loadNotifications() {
    this.notifService.getObavestenja().subscribe((list) => {
      this.obavestenja = list;
    });
  }

  handleRequest(zahtjevId: number, akcija: string): void {
    this.requestsService.handleRequest(zahtjevId, akcija).subscribe({
      next: (res) => {
        const message =
          akcija === 'prihvati' ? 'Zahtjev je prihvaćen' : 'Zahtjev je odbijen';
        this.snackBar.open(message, 'Zatvori', {
          duration: 3000,
          horizontalPosition: 'right',
          verticalPosition: 'top',
          panelClass: ['success-snackbar'],
        });

        this.loadRequests();
      },
      error: (err) => {
        console.error('Greška prilikom obrade zahteva:', err);
        this.snackBar.open(
          'Došlo je do greške prilikom obrade zahtjeva',
          'Zatvori',
          {
            duration: 3000,
            horizontalPosition: 'right',
            verticalPosition: 'top',
            panelClass: ['error-snackbar'],
          }
        );
      },
    });
  }

  showWorkOrderDetails(message: string): void {
    const workOrderId = this.extractWorkOrderId(message);
    if (workOrderId === null) {
      this.snackBar.open(
        'Ne mogu pronaći ID radnog naloga u poruci',
        'Zatvori',
        {
          duration: 3000,
          panelClass: ['error-snackbar'],
        }
      );
      return;
    }

    this.selectedNotificationId = workOrderId;

    this.workOrderService.getWorkOrderDetails(workOrderId).subscribe({
      next: (details) => {
        this.selectedWorkOrderDetails = details;
      },
      error: (err) => {
        console.error('Greška pri dohvaćanju detalja:', err);
        this.snackBar.open(
          'Greška pri dohvaćanju detalja radnog naloga',
          'Zatvori',
          {
            duration: 3000,
            panelClass: ['error-snackbar'],
          }
        );
      },
    });
  }

  toggleDetails(message: string): void {
    const workOrderId = this.extractWorkOrderId(message);
    if (this.selectedNotificationId === workOrderId) {
      this.selectedNotificationId = null;
      this.selectedWorkOrderDetails = null;
    } else {
      this.showWorkOrderDetails(message);
    }
  }

  extractWorkOrderId(message: string): number | null {
    const match = message.match(/#(\d+)/);
    return match ? +match[1] : null;
  }

  loadRequests(): void {
    this.requestsService.getAll().subscribe((data) => {
      this.zahtjevi = data;
    });
  }

  obrisiObavestenje(poruka: string, id: number) {
    this.notifService.deleteNotification(id).subscribe({
      next: () => {
        this.obavestenja = this.obavestenja.filter(
          (o) => o.obavestenje_id !== id
        );

        if (this.selectedNotificationId === this.extractWorkOrderId(poruka)) {
          this.selectedNotificationId = null;
        }

        console.log('Obavještenje uspješno obrisano');
      },
      error: (err) => {
        console.error('Greška pri brisanju obavještenja:', err);
      },
    });
  }

  handleWorkOrder(nalogId: number, akcija: 'Otvoren' | 'Odbijen'): void {
    if (akcija === 'Otvoren') {
      this.technicianService.getTechnicians().subscribe({
        next: (serviseri) => {
          const dialogRef = this.dialog.open(AssignDialogComponent, {
            width: '400px',
            data: { users: serviseri },
          });

          dialogRef.afterClosed().subscribe((selectedServiserId: number) => {
            if (!selectedServiserId) return;

            this.workOrderService
              .assignWorkOrder({
                Nalog_ID: nalogId,
                Serviser_ID: selectedServiserId,
              })
              .subscribe({
                next: () => {
                  this.workOrderService
                    .updateStatus(nalogId, 'Otvoren')
                    .subscribe({
                      next: () => {
                        this.snackBar.open(
                          'Radni nalog dodijeljen servisera i otvoren',
                          'Zatvori',
                          { duration: 3000, panelClass: ['success-snackbar'] }
                        );
                        this.loadPendingOrders();
                      },
                      error: (err) => {
                        console.error('Greška pri ažuriranju statusa:', err);
                        this.snackBar.open(
                          'Greška pri postavljanju statusa',
                          'Zatvori',
                          { duration: 3000, panelClass: ['error-snackbar'] }
                        );
                      },
                    });
                },
                error: (err) => {
                  console.error('Greška pri dodjeli servisa:', err);
                  this.snackBar.open(
                    'Greška pri dodjeli naloga serviseru',
                    'Zatvori',
                    { duration: 3000, panelClass: ['error-snackbar'] }
                  );
                },
              });
          });
        },
        error: (err) => {
          console.error('Greška pri dohvaćanju servisera:', err);
          this.snackBar.open('Ne mogu dohvatiti listu servisera', 'Zatvori', {
            duration: 3000,
            panelClass: ['error-snackbar'],
          });
        },
      });
    } else {
      this.requestsService.handleWorkOrder(nalogId, akcija).subscribe({
        next: () => {
          this.workOrderService.updateStatus(nalogId, 'Odbijen').subscribe({
            next: () => {
              this.snackBar.open('Radni nalog odbijen', 'Zatvori', {
                duration: 3000,
                panelClass: ['error-snackbar'],
              });
              this.loadPendingOrders();
            },
            error: (err) => {
              console.error('Greška pri ažuriranju statusa:', err);
              this.snackBar.open('Greška pri postavljanju statusa', 'Zatvori', {
                duration: 3000,
                panelClass: ['error-snackbar'],
              });
            },
          });
        },
        error: (err) => {
          console.error('Greška pri odbijanju naloga:', err);
          this.snackBar.open('Greška pri odbijanju naloga', 'Zatvori', {
            duration: 3000,
            panelClass: ['error-snackbar'],
          });
        },
      });
    }
  }
}
