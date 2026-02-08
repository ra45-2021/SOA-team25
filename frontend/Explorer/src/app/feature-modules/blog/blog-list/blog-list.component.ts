import { Component, OnInit } from '@angular/core';
import { Blog, BlogService } from '../blog.service';

@Component({
  selector: 'xp-blog-list',
  templateUrl: './blog-list.component.html',
  styleUrls: ['./blog-list.component.css'],
})
export class BlogListComponent implements OnInit {
  blogs: Blog[] = [];
  loading = true;
  errorMsg = '';


  constructor(private blogService: BlogService) {}

  ngOnInit(): void {
  this.blogService.getAll().subscribe({
    next: (data) => {
  console.log(data);
  this.blogs = data ?? [];
  this.loading = false;

    },
    error: (err) => {
      this.blogs = [];
      this.errorMsg = err?.error?.error || 'Failed to load blogs';
      this.loading = false;
    }
  });
}
}
