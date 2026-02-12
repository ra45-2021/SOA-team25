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
}