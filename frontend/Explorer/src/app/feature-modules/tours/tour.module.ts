import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SharedModule } from 'src/app/shared/shared.module';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MaterialModule } from 'src/app/infrastructure/material/material.module'; 
import { RouterModule } from '@angular/router';

import { TourListComponent } from './tour-list/tour-list.component';
import { TourCreateComponent } from './tour-create/tour-create.component';
import { AddCheckpointsComponent } from './add-checkpoints/add-checkpoints.component';
import { MyToursComponent } from './my-tours/my-tours.component';
import { MatCardModule } from "@angular/material/card";
import { MapComponent } from 'src/app/shared/map/map.component';
import { MatDividerModule } from "@angular/material/divider";


@NgModule({
  declarations: [
    TourListComponent,
    TourCreateComponent,
    AddCheckpointsComponent,
    MyToursComponent
  ],
  imports: [
    CommonModule,
    SharedModule,
    FormsModule,
    ReactiveFormsModule,
    MaterialModule,
    RouterModule,
    MatCardModule,
    MapComponent,
    MatDividerModule
],
  exports: [
    TourListComponent,
    TourCreateComponent,
    AddCheckpointsComponent,
    MyToursComponent
  ]
})
export class TourModule { }