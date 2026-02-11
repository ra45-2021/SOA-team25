import { Component } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { BlogService } from '../blog.service';
import { finalize, switchMap, of } from 'rxjs';

@Component({
  selector: 'xp-blog-create',
  templateUrl: './blog-create.component.html',
  styleUrls: ['./blog-create.component.css'],
})
export class BlogCreateComponent {
  saving = false;
  errorMsg = '';
  selectedFiles: File[] = [];
  previews: string[] = [];

  form = new FormGroup({
    title: new FormControl('', [Validators.required]),
    descriptionMarkdown: new FormControl('', [Validators.required]),
  });

  constructor(private blogService: BlogService, private router: Router) {}

  onFilesSelected(evt: Event) {
    const input = evt.target as HTMLInputElement;
    if (!input.files?.length) return;

    const files = Array.from(input.files);
    this.selectedFiles = [...this.selectedFiles, ...files];

    for (const f of files) {
      const reader = new FileReader();
      reader.onload = (e: any) => {
        this.previews.push(e.target.result);
      };
      reader.readAsDataURL(f);
    }
    
    input.value = '';
  }

  removeImage(index: number) {
    this.selectedFiles.splice(index, 1);
    this.previews.splice(index, 1);
  }

  submit(): void {
    if (!this.form.valid || this.saving) return;

    this.saving = true;
    this.errorMsg = '';

    const rawValue = this.form.getRawValue();
    const title = rawValue.title || '';
    const descriptionMarkdown = rawValue.descriptionMarkdown || '';

    const upload$ = this.selectedFiles.length
      ? this.blogService.uploadImages(this.selectedFiles)
      : of([] as string[]);

    upload$
      .pipe(
        switchMap((urls) =>
          this.blogService.create({
            title,
            descriptionMarkdown,
            images: urls,
          })
        ),
        finalize(() => (this.saving = false))
      )
      .subscribe({
        next: () => this.router.navigate(['/blogs']),
        error: (err) => {
          this.errorMsg = err?.error?.error || 'Failed to create blog.';
        },
      });
  }
}