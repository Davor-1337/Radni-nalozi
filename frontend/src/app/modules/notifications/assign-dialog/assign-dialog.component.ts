import { Component, Inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import {
  MAT_DIALOG_DATA,
  MatDialogModule,
  MatDialogRef,
} from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatSelectModule } from '@angular/material/select';

interface Serviser {
  Serviser_ID: number;
  Ime: string;
  Prezime: string;
}

@Component({
  standalone: true,
  selector: 'app-assign-dialog',
  templateUrl: './assign-dialog.component.html',
  styleUrls: ['./assign-dialog.component.scss'],
  imports: [
    CommonModule,
    FormsModule,
    MatDialogModule,
    MatButtonModule,
    MatFormFieldModule,
    MatSelectModule,
  ],
})
export class AssignDialogComponent {
  selectedUserId: number | null = null;

  constructor(
    public dialogRef: MatDialogRef<AssignDialogComponent>,
    @Inject(MAT_DIALOG_DATA) public data: { users: Serviser[] }
  ) {}

  ngOnInit() {
    console.log('USERS passed in dialog data:', this.data.users);
  }

  assign() {
    if (this.selectedUserId !== null) {
      this.dialogRef.close(this.selectedUserId);
    }
  }

  cancel() {
    this.dialogRef.close();
  }
}
