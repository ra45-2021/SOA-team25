import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'src/env/environment';
import { Tour } from './model/tour.model';

@Injectable({
  providedIn: 'root'
})
export class TourService {
  private readonly baseUrl = 'http://localhost:8080/tours';

  constructor(private http: HttpClient) {}

  getPublishedTours(): Observable<Tour[]> {
    return this.http.get<Tour[]>(this.baseUrl);
  }

  getMyTours(authorId: number): Observable<Tour[]> {
    return this.http.get<Tour[]>(`${this.baseUrl}/my?authorId=${authorId}`);
  }

  createTour(tour: any): Observable<any> {
    return this.http.post<any>(this.baseUrl, tour);
  }

  addCheckpoint(tourId: number, checkpointData: FormData): Observable<any> {
    return this.http.post<any>(`${this.baseUrl}/${tourId}/checkpoints`, checkpointData);
  }

  publishTour(tourId: number, publishData: { distance: number, durations: any[] }): Observable<any> {
    return this.http.put<any>(`${this.baseUrl}/${tourId}/publish`, publishData);
  }

  archiveTour(tourId: number): Observable<any> {
    return this.http.put<any>(`${this.baseUrl}/${tourId}/archive`, {});
  }

  reactivateTour(tourId: number): Observable<any> {
    return this.http.put<any>(`${this.baseUrl}/${tourId}/reactivate`, {});
  }
}