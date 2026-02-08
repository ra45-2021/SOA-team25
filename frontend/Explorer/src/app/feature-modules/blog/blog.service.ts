import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { environment } from 'src/env/environment';

export interface Blog {
  id: number;
  title: string;
  descriptionMarkdown: string;
  createdAt: string;
  images?: string[];
  authorUserId: string;
  authorUsername?: string;
}

export interface CreateBlogRequest {
  title: string;
  descriptionMarkdown: string;
  images: string[];
}

@Injectable({ providedIn: 'root' })
export class BlogService {
  constructor(private http: HttpClient) {}

  getAll(): Observable<Blog[]> {
  return this.http.get<Blog[]>(`${environment.apiHost}/blogs`);
  }

  create(req: CreateBlogRequest): Observable<Blog> {
    return this.http.post<Blog>(`${environment.apiHost}/blogs`, req);
  }

}
