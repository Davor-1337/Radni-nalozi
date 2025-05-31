import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatDialogModule } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-user-dialog',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatSnackBarModule,
  ],
  templateUrl: './user-dialog.component.html',
  styleUrls: ['./user-dialog.component.scss'],
})
export class UserDialogComponent {
  currentPassword: string = '';
  newPassword: string = '';
  confirmPassword: string = '';

  constructor(
    public dialogRef: MatDialogRef<UserDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { username: string; role: string },
    private authService: AuthService,
    private snackBar: MatSnackBar
  ) {}

  changePassword() {
    if (this.newPassword !== this.confirmPassword) {
      this.snackBar.open('Nove lozinke se ne poklapaju.', 'Zatvori', {
        duration: 3000,
        panelClass: ['error-snackbar'],
      });
      return;
    }

    if (this.newPassword.length < 8) {
      this.snackBar.open('Lozinka mora imati barem 8 karaktera.', 'Zatvori', {
        duration: 3000,
        panelClass: ['error-snackbar'],
      });
      return;
    }

    this.authService
      .changePassword(this.currentPassword, this.newPassword)
      .subscribe({
        next: () => {
          this.snackBar.open('Lozinka uspješno promijenjena!', 'Zatvori', {
            duration: 3000,
            panelClass: ['success-snackbar'],
          });
          this.dialogRef.close();
        },
        error: (err) => {
          const msg = err.error?.message || 'Greška pri promjeni lozinke.';
          this.snackBar.open(msg, 'Zatvori', {
            duration: 3000,
            panelClass: ['error-snackbar'],
          });
        },
      });
  }

  close() {
    this.dialogRef.close();
  }

  showPasswordFields = false;

  togglePasswordFields() {
    this.showPasswordFields = !this.showPasswordFields;
  }
}
