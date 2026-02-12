import { Component, OnInit } from '@angular/core';
import { TourService } from '../tour.service';
import { Tour } from '../model/tour.model';
import { AuthService } from 'src/app/infrastructure/auth/auth.service';

@Component({
  selector: 'xp-tour-list',
  templateUrl: './tour-list.component.html',
  styleUrls: ['./tour-list.component.css']
})
export class TourListComponent implements OnInit {
  publishedTours: Tour[] = [];
  userRole: string = '';

  constructor(
    private service: TourService,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    this.authService.user$.subscribe(user => {
      console.log('Trenutni korisnik:', user);
      this.userRole = user.role ? user.role.toUpperCase() : '';
      
      if (this.userRole === 'GUIDE') {
        this.loadGuideTours(user.id);
      } else {
        this.loadPublicTours();
      }
    });
  }

  loadGuideTours(authorId: number): void {
    this.service.getMyTours(authorId).subscribe({
      next: (result: Tour[]) => {
        this.publishedTours = result;
      },
      error: (err) => {
        console.error('Greška pri dohvatanju tura vodiča:', err);
      }
    });
  }

  loadPublicTours(): void {
    this.service.getPublishedTours().subscribe({
      next: (result: Tour[]) => {
        this.publishedTours = result;
      },
      error: (err) => {
        console.error('Greška pri dohvatanju javnih tura:', err);
      }
    });
  }
}