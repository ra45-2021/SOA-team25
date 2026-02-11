import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { BlogService } from '../blog.service';
import { finalize } from 'rxjs';

@Component({
  selector: 'xp-blog-details',
  templateUrl: './blog-details.component.html',
  styleUrls: ['./blog-details.component.css']
})
export class BlogDetailsComponent implements OnInit {
  blog: any = null;
  loading = true;
  errorMsg = '';

  constructor(
    private route: ActivatedRoute,
    private blogService: BlogService,
    private router: Router
  ) {}

  ngOnInit(): void {
  const idParam = this.route.snapshot.paramMap.get('id');
  if (idParam) {
    this.fetchBlogDetails(+idParam); 
  } else {
    this.router.navigate(['/blogs']);
  }
}

fetchBlogDetails(id: number): void {
  this.loading = true;
  this.blogService.getById(id)
    .pipe(finalize(() => this.loading = false))
    .subscribe({
      next: (res) => {
        this.blog = res;
      },
      error: () => {
        this.errorMsg = 'Could not load blog details.';
      }
    });
}
}