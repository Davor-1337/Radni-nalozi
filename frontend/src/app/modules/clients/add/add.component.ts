import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialogRef } from '@angular/material/dialog';
import { ClientsService } from '../../../core/services/clients.service';
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
  clientsForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private clientService: ClientsService,
    private snackBar: MatSnackBar,
    private dialogRef: MatDialogRef<AddComponent>
  ) {
    this.clientsForm = this.fb.group({
      Klijent_ID: ['', Validators.required],
      Naziv: ['', Validators.required],
      KontaktOsoba: ['', Validators.required],
      Email: ['', Validators.required],
      Tel: ['', Validators.required],
      Adresa: ['', Validators.required],
      User_ID: [null],
    });
  }

  closeForm() {
    this.dialogRef.close();
  }

  onSubmit() {
    console.log('Forma je submitovana');
    if (this.clientsForm.valid) {
      const formData = {
        Klijent_ID: parseInt(this.clientsForm.value.Klijent_ID),
        Naziv: this.clientsForm.value.Naziv,
        KontaktOsoba: this.clientsForm.value.KontaktOsoba,
        Email: this.clientsForm.value.Email,
        Tel: this.clientsForm.value.Tel,
        Adresa: this.clientsForm.value.Adresa,
        User_ID: parseInt(this.clientsForm.value.User_ID),
      } satisfies {
        Klijent_ID: number;
        Naziv: string;
        KontaktOsoba: string;
        Email: string;
        Tel: string;
        Adresa: string;
        User_ID: number;
      };

      this.clientService.postClient(formData).subscribe({
        next: (response: { Klijent_ID: number }) => {
          this.showSuccess(`Klijent uspješno dodan!`);
          this.dialogRef.close(true);
        },
        error: (err: any) => {
          this.showError(
            err.error?.message || 'Došlo je do greške pri dodavanju klijenta'
          );
          console.error('Greška:', err);
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
