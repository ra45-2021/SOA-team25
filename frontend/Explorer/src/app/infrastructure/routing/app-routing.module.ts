import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { HomeComponent } from 'src/app/feature-modules/layout/home/home.component';
import { LoginComponent } from '../auth/login/login.component';
import { AuthGuard } from '../auth/auth.guard';
import { RegistrationComponent } from '../auth/registration/registration.component';
import { BlogListComponent } from 'src/app/feature-modules/blog/blog-list/blog-list.component';
import { BlogCreateComponent } from 'src/app/feature-modules/blog/blog-create/blog-create.component';
import { BlogDetailsComponent } from 'src/app/feature-modules/blog/blog-details/blog-details.component';
import { TourListComponent } from 'src/app/feature-modules/tours/tour-list/tour-list.component';
import { TourCreateComponent } from 'src/app/feature-modules/tours/tour-create/tour-create.component';
import { MyToursComponent } from 'src/app/feature-modules/tours/my-tours/my-tours.component';
import { AddCheckpointsComponent } from 'src/app/feature-modules/tours/add-checkpoints/add-checkpoints.component';

const routes: Routes = [
  {path: 'home', component: HomeComponent},
  {path: 'login', component: LoginComponent},
  {path: 'register', component: RegistrationComponent},

  { path: 'blogs', component: BlogListComponent, canActivate: [AuthGuard] },
  { path: 'blogs/new', component: BlogCreateComponent, canActivate: [AuthGuard] },
  { path: 'blogs/:id', component: BlogDetailsComponent },
  { path: 'tour-list', component: TourListComponent, canActivate: [AuthGuard] },
  { path: 'create-tour', component: TourCreateComponent, canActivate: [AuthGuard] },
  { path: 'my-tours', component: MyToursComponent, canActivate: [AuthGuard] },
  { path: 'add-checkpoints/:id', component: AddCheckpointsComponent, canActivate: [AuthGuard] },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
