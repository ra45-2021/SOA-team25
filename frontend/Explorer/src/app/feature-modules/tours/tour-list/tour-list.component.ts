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
  currentUserId: number = 0;

  constructor(
    private service: TourService,
    private authService: AuthService
  ) {}

  ngOnInit(): void {
    this.authService.user$.subscribe(user => {
      this.userRole = user.role ? user.role.toUpperCase() : '';
      this.currentUserId = user.id;
      
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
      }
    });
  }

  loadPublicTours(): void {
    this.service.getPublishedTours().subscribe({
      next: (result: Tour[]) => {
        this.publishedTours = result;
      }
    });
  }

  archiveTour(tourId: number): void {
    if (confirm('Are you sure you want to archive this tour? It will be hidden from tourists.')) {
      this.service.archiveTour(tourId).subscribe({
        next: () => {
          alert('Tour archived successfully.');
          this.loadGuideTours(this.currentUserId); 
        },
        error: (err) => {
          console.error('Archive error:', err);
          alert('Failed to archive tour: ' + (err.error?.error || 'Server error'));
        }
      });
    }
  }

  reactivateTour(tourId: number): void {
    if (confirm('Reactivate this tour? It will become visible to tourists again.')) {
      this.service.reactivateTour(tourId).subscribe({
        next: () => {
          alert('Tour is now active again!');
          this.loadGuideTours(this.currentUserId); 
        },
        error: (err) => {
          console.error('Reactivation error:', err);
          alert('Failed to reactivate tour: ' + (err.error?.error || 'Server error'));
        }
      });
    }
  }
}