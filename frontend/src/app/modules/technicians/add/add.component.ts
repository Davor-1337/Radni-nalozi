import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialogRef } from '@angular/material/dialog';
import { TechniciansService } from '../../../core/services/technicians.service';
import {
  FormBuilder,
  FormGroup,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';

export interface Technician {
  Serviser_ID: number;
  Ime: string;
  Prezime: string;
  Specijalnost: string;
  Telefon: string;
  User_ID: number;
}

@Component({
  selector: 'app-add',
  imports: [CommonModule, ReactiveFormsModule],
  templateUrl: './add.component.html',
  styleUrl: './add.component.scss',
})
export class AddComponent {
  techniciansForm: FormGroup;

  constructor(
    private fb: FormBuilder,
    private technicianService: TechniciansService,
    private snackBar: MatSnackBar,
    private dialogRef: MatDialogRef<AddComponent>
  ) {
    this.techniciansForm = this.fb.group({
      Ime: ['', Validators.required],
      Prezime: ['', Validators.required],
      Specijalnost: ['', Validators.required],
      Telefon: ['', [Validators.required, Validators.pattern('^[0-9]{6,15}$')]],
    });
  }

  maxId: number = 0;

  ngOnInit() {
    this.technicianService.getTechnicians().subscribe((technicians) => {
      const ids = technicians.map((t) => t.Serviser_ID);
      this.maxId = Math.max(...ids);
    });
  }

  closeForm() {
    this.dialogRef.close();
  }

  onSubmit() {
    console.log('Forma za servisera je submitovana');
    if (this.techniciansForm.valid) {
      const newId = this.maxId + 1;

      const technicianData: Technician = {
        Serviser_ID: newId,
        User_ID: newId,
        Ime: this.techniciansForm.value.Ime,
        Prezime: this.techniciansForm.value.Prezime,
        Specijalnost: this.techniciansForm.value.Specijalnost,
        Telefon: this.techniciansForm.value.Telefon,
      };
      console.log('Podaci koji se šalju na backend:', technicianData);

      this.technicianService.addTechnician(technicianData).subscribe({
        next: (response) => {
          this.showSuccess(`Serviser uspješno dodan!`);
          this.dialogRef.close(true);
          this.maxId++;
        },
        error: (err: any) => {
          this.showError(
            err.error?.message || 'Došlo je do greške pri dodavanju servisera'
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

  addTechnician(technician: Technician): void {
    this.technicianService.addTechnician(technician).subscribe({
      next: () => {
        // Osvežite listu servisera ili prikažite toast poruku
        console.log('Serviser uspešno dodat');
        // Dodajte kod za osvežavanje liste
      },
      error: (err) => {
        console.error('Greška pri dodavanju servisera:', err);
      },
    });
  }
}
