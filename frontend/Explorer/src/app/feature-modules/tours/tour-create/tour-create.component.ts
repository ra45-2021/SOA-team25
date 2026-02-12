import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { TourService } from '../tour.service';
import { Tour, Difficulty } from '../model/tour.model';
import { AuthService } from 'src/app/infrastructure/auth/auth.service';

@Component({
  selector: 'xp-tour-create',
  templateUrl: './tour-create.component.html',
  styleUrls: ['./tour-create.component.css']
})
export class TourCreateComponent {
  tour: Tour = {
    name: '',
    description: '',
    difficulty: Difficulty.EASY,
    tags: '',
    price: 0,
    author_id: 0
  };

  constructor(
    private service: TourService, 
    private authService: AuthService,
    private router: Router
  ) {
    this.authService.user$.subscribe(user => {
      if (user) this.tour.author_id = user.id;
    });
  }

  createTour(): void {
  const cleanedTags = this.tour.tags
    .split(';')
    .map(t => t.trim())
    .filter(t => t !== '')
    .join(';');

  const tourToSend = { 
    ...this.tour, 
    tags: cleanedTags,
    difficulty: Number(this.tour.difficulty) 
  };

  console.log('Sending to backend:', tourToSend); 
  this.service.createTour(tourToSend).subscribe({
    next: (result) => {
      console.log('Tour created!', result);
      this.router.navigate(['/add-checkpoints', result.id]);
    },
    error: (err) => {
      console.error('Backend error:', err);
    }
  });
}
}