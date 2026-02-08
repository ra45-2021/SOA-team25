import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { LoginComponent } from './login/login.component';
import { MaterialModule } from '../material/material.module';
import { ReactiveFormsModule } from '@angular/forms';
import { RegistrationComponent } from './registration/registration.component';
import { MatSelectModule } from '@angular/material/select';
import { MatCardModule } from "@angular/material/card";


@NgModule({
  declarations: [
    LoginComponent,
    RegistrationComponent
  ],
  imports: [
    CommonModule,
    MaterialModule,
    ReactiveFormsModule,
    MatSelectModule,
    MatCardModule
],
  exports: [
    LoginComponent
  ]
})
export class AuthModule { }
