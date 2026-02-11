import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from 'src/app/feature-modules/layout/home/home.component';
import { LoginComponent } from '../auth/login/login.component';
import { AuthGuard } from '../auth/auth.guard';
import { RegistrationComponent } from '../auth/registration/registration.component';
import { BlogListComponent } from 'src/app/feature-modules/blog/blog-list/blog-list.component';
import { BlogCreateComponent } from 'src/app/feature-modules/blog/blog-create/blog-create.component';
import { BlogDetailsComponent } from 'src/app/feature-modules/blog/blog-details/blog-details.component';

const routes: Routes = [
  {path: 'home', component: HomeComponent},
  {path: 'login', component: LoginComponent},
  {path: 'register', component: RegistrationComponent},

  { path: 'blogs', component: BlogListComponent, canActivate: [AuthGuard] },
  { path: 'blogs/new', component: BlogCreateComponent, canActivate: [AuthGuard] },
  { path: 'blogs/:id', component: BlogDetailsComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
