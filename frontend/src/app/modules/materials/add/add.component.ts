import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialogRef } from '@angular/material/dialog';
import { MaterialServiceService } from '../../../core/services/material-service.service';
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
  materialsForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private materialService: MaterialServiceService,
    private snackBar: MatSnackBar,
    private dialogRef: MatDialogRef<AddComponent>
  ) {
    this.materialsForm = this.fb.group({
      NazivMaterijala: ['', Validators.required],
      Cijena: ['', Validators.required],
      KolicinaUSkladistu: ['', Validators.required],
    });
  }

  closeForm() {
    this.dialogRef.close();
  }

  onSubmit() {
    console.log('Forma je submitovana');
    if (this.materialsForm.valid) {
      const formData = {
        NazivMaterijala: this.materialsForm.value.NazivMaterijala,
        Cijena: parseFloat(this.materialsForm.value.Cijena),
        KolicinaUSkladistu: parseInt(
          this.materialsForm.value.KolicinaUSkladistu
        ),
      } satisfies {
        NazivMaterijala: string;
        Cijena: number;
        KolicinaUSkladistu: number;
      };

      this.materialService.postMaterial(formData).subscribe({
        next: (response: { materijal_id: number }) => {
          this.showSuccess(`Materijal uspješno dodan!`);
          this.dialogRef.close(true);
        },
        error: (err: any) => {
          this.showError(
            err.error?.message || 'Došlo je do greške pri dodavanju materijala'
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
