import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialogRef } from '@angular/material/dialog';
import { AuthService } from '../../../core/services/auth.service';
import { WorkOrderService } from '../../../core/services/workOrder.service';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';

@Component({
  selector: 'app-add',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './add.component.html',
  styleUrl: './add.component.scss',
})
export class AddComponent {
  workOrderForm: FormGroup;

  public role: string = '';
  constructor(
    private fb: FormBuilder,
    private workOrderService: WorkOrderService,
    private snackBar: MatSnackBar,
    private dialogRef: MatDialogRef<AddComponent>,
    private authService: AuthService
  ) {
    this.role = this.authService.getRole();
    this.workOrderForm = this.fb.group({
      opisProblema: ['', Validators.required],
      prioritet: ['', Validators.required],
      status: [
        this.role === 'klijent' ? 'Na cekanju' : '',
        Validators.required,
      ],
      lokacija: ['', Validators.required],
      client: ['', Validators.required],
    });

    if (this.role === 'klijent') {
      const statusCtrl = this.workOrderForm.get('status')!;
      statusCtrl.setValue('Na cekanju');
      statusCtrl.clearValidators();
      statusCtrl.updateValueAndValidity();

      const clientCtrl = this.workOrderForm.get('client')!;

      const userId = this.authService.getUserId();
      console.log('User ID:', userId);

      clientCtrl.setValue(userId);
      console.log('Client value after setValue:', clientCtrl.value);
      clientCtrl.clearValidators();
      clientCtrl.updateValueAndValidity();
    }
  }

  closeForm() {
    this.dialogRef.close();
  }

  onSubmit() {
    if (this.workOrderForm.valid) {
      const formData = {
        klijent_id: Number(this.workOrderForm.value.client),
        opisProblema: this.workOrderForm.value.opisProblema,
        prioritet: this.workOrderForm.value.prioritet,
        status:
          this.role === 'klijent'
            ? 'Na cekanju'
            : this.workOrderForm.value.status,
        lokacija: this.workOrderForm.value.lokacija,
      };
      console.log('Client ID:', this.workOrderForm.value.client);

      this.workOrderService.postWorkOrder(formData).subscribe({
        next: (response: { status: string; nalog_id: number }) => {
          const message =
            this.role === 'klijent'
              ? 'Zahtjev uspješno poslan na odobrenje.'
              : 'Radni nalog uspješno kreiran!';
          this.showSuccess(message);
          this.dialogRef.close(true);
        },
        error: (err: any) => {
          this.showError(
            err.error?.message || 'Došlo je do greške pri kreiranju naloga'
          );
        },
      });
    }
  }

  private showSuccess(message: string): void {
    this.snackBar.open(message, 'Zatvori', {
      duration: 5000,
      panelClass: ['success-snackbar'],
      horizontalPosition: 'right',
      verticalPosition: 'top',
    });
  }

  private showError(message: string): void {
    this.snackBar.open(message, 'Zatvori', {
      duration: 7000,
      panelClass: ['error-snackbar'],
      horizontalPosition: 'right',
      verticalPosition: 'top',
    });
  }
}
