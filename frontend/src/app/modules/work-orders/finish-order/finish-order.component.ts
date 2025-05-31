import { Component, Inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { Material } from '../../../core/services/material-service.service';
import { MatSnackBar } from '@angular/material/snack-bar';
import {
  FormArray,
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import {
  MatDialogModule,
  MatDialogRef,
  MAT_DIALOG_DATA,
} from '@angular/material/dialog';
import { WorkOrderService } from '../../../core/services/workOrder.service';
import { MaterialServiceService } from '../../../core/services/material-service.service';
import { Observable } from 'rxjs';
import { forkJoin } from 'rxjs';

@Component({
  selector: 'app-finish-order',
  standalone: true,
  imports: [
    CommonModule,
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatIconModule,
    MatButtonModule,
  ],
  templateUrl: './finish-order.component.html',
  styleUrls: ['./finish-order.component.scss'],
})
export class FinishOrderComponent implements OnInit {
  form: FormGroup;
  loading = false;
  successMessage = '';
  errorMessage = '';
  materialsList$: Observable<Material[]>;

  constructor(
    @Inject(MAT_DIALOG_DATA)
    public data: { nalogId: number; serviserId: number },
    private dialogRef: MatDialogRef<FinishOrderComponent>,
    private fb: FormBuilder,
    private snackBar: MatSnackBar,
    private workOrderService: WorkOrderService,
    private materialService: MaterialServiceService
  ) {
    this.form = this.fb.group({
      materials: this.fb.array([]),
      hours: [null, [Validators.required, Validators.min(0.1)]],
    });

    this.materialsList$ = this.materialService.getMaterials();
  }

  ngOnInit() {
    if (this.materials.length === 0) {
      this.addMaterialRow();
    }
  }

  get materials(): FormArray {
    return this.form.get('materials') as FormArray;
  }

  addMaterialRow() {
    this.materials.push(
      this.fb.group({
        materialId: [null, Validators.required],
        quantity: [1, [Validators.required, Validators.min(1)]],
      })
    );
  }

  removeMaterialRow(index: number) {
    if (this.materials.length > 1) {
      this.materials.removeAt(index);
    }
  }

  finishWorkOrder() {
    if (this.form.invalid) {
    }

    this.loading = true;

    const nalogId = this.data.nalogId;
    const serviserId = this.data.serviserId;
    const hoursPayload = {
      serviser_ID: serviserId,
      brojRadnihSati: this.form.value.hours,
    };

    this.workOrderService.addHours(nalogId, hoursPayload).subscribe({
      next: () => {
        const materialCalls = this.materials.value.map((m: any) =>
          this.workOrderService.addMaterial(nalogId, {
            materijal_ID: m.materialId,
            kolicinaUtrosena: m.quantity,
          })
        );

        forkJoin(materialCalls).subscribe({
          next: () => {
            this.workOrderService.finishWorkOrder(nalogId).subscribe({
              next: () => {
                this.snackBar.open('Nalog je uspješno završen.', 'Zatvori', {
                  duration: 3000,
                  panelClass: ['success-snackbar'],
                });
                this.dialogRef.close(true);
              },
              error: () => this.handleError('Pri završavanju naloga'),
            });
          },
          error: () => this.handleError('Pri unosu materijala'),
        });
      },
      error: () => this.handleError('Pri unosu radnih sati'),
    });
  }

  private handleError(msg: string) {
    this.errorMessage = `Greška ${msg}.`;
    this.loading = false;
  }

  cancel() {
    this.dialogRef.close(false);
  }
}
