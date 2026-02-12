import { Component, OnInit } from '@angular/core';
import { TourService } from '../tour.service';
import { AuthService } from 'src/app/infrastructure/auth/auth.service';

@Component({
  selector: 'xp-my-tours',
  templateUrl: './my-tours.component.html',
  styleUrls: ['./my-tours.component.css']
})
export class MyToursComponent implements OnInit {
  purchasedTours: any[] = []; 

  constructor(
    private service: TourService,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    this.service.getPublishedTours().subscribe(tours => {
      this.purchasedTours = tours.map(t => ({
        tour: t,
        token: { is_executed: false, is_reviewed: false }
      }));
    });
  }

  leaveReview(item: any): void {
    console.log('Reviewing tour:', item.tour.id);
  }
}