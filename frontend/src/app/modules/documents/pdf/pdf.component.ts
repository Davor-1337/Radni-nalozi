import { Component, Inject } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';

@Component({
  selector: 'app-pdf',
  templateUrl: './pdf.component.html',
})
export class PdfComponent {
  pdfUrl: SafeResourceUrl;

  constructor(
    @Inject(MAT_DIALOG_DATA) public data: { pdfBlob: Blob },
    private sanitizer: DomSanitizer,
    public dialogRef: MatDialogRef<PdfComponent>
  ) {
    console.log('PDF Blob received:', data.pdfBlob);
    if (data.pdfBlob) {
      this.pdfUrl = this.sanitizer.bypassSecurityTrustResourceUrl(
        URL.createObjectURL(data.pdfBlob)
      );
    } else {
      console.error('Invalid PDF Blob');
    }
    this.pdfUrl = this.sanitizer.bypassSecurityTrustResourceUrl(
      URL.createObjectURL(data.pdfBlob)
    );

    console.log('Generated PDF URL:', this.pdfUrl);
  }
}
