import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AuthService } from '../../../core/services/auth.service';
import { Router } from '@angular/router';
import { ReactiveFormsModule, FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { MatSnackBar } from '@angular/material/snack-bar';
import { jwtDecode } from 'jwt-decode';

@Component({
  selector: 'app-login',
  standalone: true,
  imports: [ReactiveFormsModule, FormsModule, CommonModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss'],
})
export class LoginComponent {
  errorMessage: string = '';
  loginForm: FormGroup;
  signupForm: FormGroup;
  activeForm: 'login' | 'signup' = 'login';

  constructor(
    private fb: FormBuilder,
    private authService: AuthService,
    private router: Router,
    private snackBar: MatSnackBar
  ) {
    this.loginForm = this.fb.group({
      username: ['', Validators.required],
      password: ['', Validators.required],
    });

    this.signupForm = this.fb.group(
      {
        username: ['', Validators.required],
        email: ['', [Validators.required, Validators.email]],
        password: ['', [Validators.required, Validators.minLength(6)]],
        confirmPassword: ['', Validators.required],
        role: ['', Validators.required],
      },
      {
        validator: this.passwordMatchValidator,
      }
    );

    this.signupForm.get('role')!.valueChanges.subscribe((role: string) => {
      this.updateAdditionalFields(role);
    });
  }

  updateAdditionalFields(role: string) {
    const additionalFields = [
      'name',
      'contactPerson',
      'address',
      'tel',
      'name',
      'surname',
      'specialty',
    ];

    additionalFields.forEach((field) => {
      if (this.signupForm.contains(field)) {
        this.signupForm.removeControl(field);
      }
    });

    if (role === 'klijent') {
      this.signupForm.addControl(
        'naziv',
        this.fb.control('', Validators.required)
      );
      this.signupForm.addControl(
        'contactPerson',
        this.fb.control('', Validators.required)
      );
      this.signupForm.addControl(
        'tel',
        this.fb.control('', Validators.required)
      );
      this.signupForm.addControl(
        'address',
        this.fb.control('', Validators.required)
      );
    }

    if (role === 'serviser') {
      this.signupForm.addControl(
        'name',
        this.fb.control('', Validators.required)
      );
      this.signupForm.addControl(
        'surname',
        this.fb.control('', Validators.required)
      );
      this.signupForm.addControl(
        'specialty',
        this.fb.control('', Validators.required)
      );
      this.signupForm.addControl(
        'tel',
        this.fb.control('', Validators.required)
      );
    }
  }

  signUp() {
    if (this.signupForm.valid) {
      const formData = this.signupForm.value;

      const payload: any = {
        username: formData.username,
        password: formData.password,
        email: formData.email,
        role: formData.role,
      };

      if (formData.role === 'serviser') {
        payload.name = formData.name;
        payload.surname = formData.surname;
        payload.specialty = formData.specialty;
        payload.tel = formData.tel;
      }

      if (formData.role === 'klijent') {
        payload.Naziv = formData.naziv;
        payload.contactPerson = formData.contactPerson;
        payload.tel = formData.tel;
        payload.address = formData.address;
      }

      console.log('Šaljem signup payload:', payload);

      this.authService.signUp(payload).subscribe({
        next: () => {
          this.snackBar.open('Zahtjev uspješno poslan!', 'Zatvori', {
            duration: 3000,
            panelClass: ['success-snackbar'],
          });
          this.router.navigate(['/login']);
        },
        error: () => {
          this.snackBar.open('Greška prilikom slanja zahtjeva.', 'Zatvori', {
            duration: 3000,
            panelClass: ['error-snackbar'],
          });
        },
      });
    } else {
      this.snackBar.open('Molimo ispunite sva polja ispravno.', 'Zatvori', {
        duration: 3000,
      });
    }
  }

  passwordMatchValidator(formGroup: FormGroup): null | undefined {
    const password = formGroup.get('password')?.value;
    const confirmPassword = formGroup.get('confirmPassword')?.value;
    if (password !== confirmPassword) {
      formGroup.get('confirmPassword')?.setErrors({ mismatch: true });
    } else {
      formGroup.get('confirmPassword')?.setErrors(null);
    }
    return null;
  }

  login() {
    if (this.loginForm.valid) {
      this.authService.login(this.loginForm.value).subscribe({
        next: (res) => {
          console.log('Login response:', res);

          if (res.token) {
            localStorage.setItem('token', res.token);
            console.log('Token saved:', res.token);

            let role = '';

            try {
              const decodedToken: any = jwtDecode(res.token);
              role = decodedToken.Role;
              console.log('Decoded role:', role);
            } catch (e) {
              console.error('Failed to decode token:', e);
              this.errorMessage = 'Greška prilikom čitanja korisničke uloge.';
              return;
            }

            this.errorMessage = '';

            console.log('Role iz tokena:', role);
            let targetRoute = '/home';
            if (role === 'klijent') {
              targetRoute = '/radni-nalozi';
            }

            this.router.navigate([targetRoute]).then(
              () => console.log('Navigacija uspješna:', targetRoute),
              (err) => console.error('Greška pri navigaciji:', err)
            );
          } else {
            console.error('Token missing in response');
            this.errorMessage =
              'Autentifikacija nije uspela. Pokušajte ponovo.';
          }
        },
        error: (err) => {
          console.error('Login error:', err);
          if (err.status === 401) {
            this.errorMessage = 'Pogrešan korisnički naziv ili lozinka.';
          } else {
            this.errorMessage = 'Došlo je do greške. Pokušajte kasnije.';
          }
        },
      });
    }
  }
  setFormToSignup() {
    this.activeForm = 'signup';
  }
}
