import { Component } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { BlogService } from '../blog.service';

@Component({
  selector: 'xp-blog-create',
  templateUrl: './blog-create.component.html',
  styleUrls: ['./blog-create.component.css'],
})
export class BlogCreateComponent {
  saving = false;

  form = new FormGroup({
    title: new FormControl('', [Validators.required]),
    descriptionMarkdown: new FormControl('', [Validators.required]),
  });

  constructor(private blogService: BlogService, private router: Router) {}

  submit(): void {
    if (!this.form.valid || this.saving) return;

    this.saving = true;
    this.blogService.create({
      title: this.form.value.title || '',
      descriptionMarkdown: this.form.value.descriptionMarkdown || '',
      images: [],
    }).subscribe({
      next: () => this.router.navigate(['/blogs']),
      error: () => (this.saving = false),
    });
  }
}
