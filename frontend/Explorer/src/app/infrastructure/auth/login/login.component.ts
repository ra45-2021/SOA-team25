import { Component } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { AuthService } from '../auth.service';
import { Router } from '@angular/router';
import { Login } from '../model/login.model';

@Component({
  selector: 'xp-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent {

  errorMessage = '';
  private errorTimeout: any;

  constructor(
    private authService: AuthService,
    private router: Router
  ) {}

  loginForm = new FormGroup({
    username: new FormControl('', [Validators.required]),
    password: new FormControl('', [Validators.required]),
  });

  login(): void {
  this.errorMessage = '';

  const login: Login = {
    username: this.loginForm.value.username || "",
    password: this.loginForm.value.password || "",
  };

  if (this.loginForm.valid) {
    this.authService.login(login).subscribe({
      next: () => {
        this.router.navigate(['/']);
      },
      error: (err) => {
        this.errorMessage =
          err?.error?.message ||
          err?.error?.error ||
          'Login failed. Please try again.';

        if (this.errorTimeout) clearTimeout(this.errorTimeout);

        this.errorTimeout = setTimeout(() => {
          this.errorMessage = '';
        }, 2000);
      }
    });
  }
}

}
