import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialogRef } from '@angular/material/dialog';
import {
  DocumentsService,
  InvoiceFinalized,
  InvoicePreview,
} from '../../../core/services/documents.service';
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
  invoicesForm: FormGroup;
  preview: InvoicePreview | null = null;
  loadingGenerate = false;
  loadingFinalize = false;

  constructor(
    private fb: FormBuilder,
    private documentService: DocumentsService,
    private snackBar: MatSnackBar,
    private dialogRef: MatDialogRef<AddComponent>
  ) {
    this.invoicesForm = this.fb.group({
      Nalog_ID: ['', Validators.required],
    });
  }

  closeForm() {
    this.dialogRef.close();
  }

  onGenerate() {
    if (this.invoicesForm.invalid) {
      return;
    }

    this.loadingGenerate = true;
    this.preview = null;

    const nalogId = +this.invoicesForm.value.Nalog_ID;

    this.documentService.generateInvoice(nalogId).subscribe({
      next: (data) => {
        const isEmptyInvoice =
          data.total_hours === 0 &&
          data.labor_cost === 0 &&
          data.total_cost === 0 &&
          (!data.materials || data.materials.length === 0);

        if (isEmptyInvoice) {
          this.snackBar.open('Nalog još nije završen.', 'Zatvori', {
            duration: 4000,
            panelClass: ['error-snackbar'],
          });
          this.loadingGenerate = false;
          return;
        }

        this.preview = data;
        this.dialogRef.updateSize('400px', 'auto');
        this.loadingGenerate = false;
      },
      error: (err) => {
        this.loadingGenerate = false;
        this.snackBar.open(
          err.error?.message || 'Greška pri generisanju fakture',
          'Zatvori',
          {
            duration: 4000,
            panelClass: ['error-snackbar'],
          }
        );
      },
    });
  }

  onAdd() {
    if (!this.preview) {
      return;
    }
    this.loadingFinalize = true;

    const nalogId = this.preview.work_order_id;
    this.documentService.finalizeInvoice(nalogId).subscribe({
      next: (res: InvoiceFinalized) => {
        this.loadingFinalize = false;
        this.snackBar.open(res.message, 'Zatvori', {
          duration: 4000,
          panelClass: ['success-snackbar'],
        });
        this.dialogRef.close(true);
      },
      error: (err) => {
        this.loadingFinalize = false;
        this.snackBar.open(
          err.error?.message || 'Greška pri čuvanju fakture',
          'Zatvori',
          {
            duration: 4000,
            panelClass: ['error-snackbar'],
          }
        );
      },
    });
  }

  onSubmit() {
    console.log('Forma je submitovana');
    if (this.invoicesForm.valid) {
      const formData = {
        Nalog_ID: parseInt(this.invoicesForm.value.Nalog_ID),
        Iznos: parseFloat(this.invoicesForm.value.Iznos),
      } satisfies {
        Nalog_ID: number;
        Iznos: number;
      };

      this.documentService.postInvoice(formData).subscribe({
        next: (response: { faktura_id: number }) => {
          this.showSuccess(`Faktura uspješno dodana!`);
          this.dialogRef.close(true);
        },
        error: (err: any) => {
          this.showError(
            err.error?.message || 'Došlo je do greške pri dodavanju fakture.'
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
