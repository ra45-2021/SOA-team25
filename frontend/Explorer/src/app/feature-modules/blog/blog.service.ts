import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable, map } from 'rxjs';
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

type UploadResponse = string[] | { urls: string[] };

@Injectable({ providedIn: 'root' })
export class BlogService {
  private readonly baseUrl = environment.apiHost;

  private readonly blogsUrl = `${this.baseUrl}/blogs`;

  private readonly uploadUrl = `${this.baseUrl}/blogs/images`;

  constructor(private http: HttpClient) {}

  getAll(): Observable<Blog[]> {
    return this.http.get<Blog[]>(this.blogsUrl);
  }

  getById(id: number): Observable<Blog> {
    return this.http.get<Blog>(`${this.blogsUrl}/${id}`);
  }

  create(req: CreateBlogRequest): Observable<Blog> {
    return this.http.post<Blog>(this.blogsUrl, req);
  }

  uploadImages(files: File[]): Observable<string[]> {
    if (!files || files.length === 0) return new Observable<string[]>((sub) => {
      sub.next([]);
      sub.complete();
    });

    const fd = new FormData();

    for (const f of files) {
      fd.append('files', f, f.name);
    }

    return this.http.post<UploadResponse>(this.uploadUrl, fd).pipe(
      map((res) => {
        if (Array.isArray(res)) return res;
        if (res && Array.isArray((res as any).urls)) return (res as any).urls;
        return [];
      })
    );
  }

  delete(id: number): Observable<void> {
    return this.http.delete<void>(`${this.blogsUrl}/${id}`);
  }

  search(q: string): Observable<Blog[]> {
    const params = new HttpParams().set('q', q);
    return this.http.get<Blog[]>(this.blogsUrl, { params });
  }
}
