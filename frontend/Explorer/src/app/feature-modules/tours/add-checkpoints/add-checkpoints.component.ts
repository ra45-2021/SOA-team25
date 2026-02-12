import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TourService } from '../tour.service';
import { ActivatedRoute, Router } from '@angular/router';

export enum TransportType { Walk = 0, Bike = 1, Car = 2 }

@Component({
  selector: 'xp-add-checkpoints',
  templateUrl: './add-checkpoints.component.html',
  styleUrls: ['./add-checkpoints.component.css']
})
export class AddCheckpointsComponent implements OnInit {
  tourId: number = 0;
  checkpoints: any[] = [];
  durations: any[] = [];
  tourDistanceKm: number = 0;
  checkpointForm: FormGroup;
  
  selectedFiles: File[] = [];
  previews: string[] = [];
  
  TransportType = TransportType;
  selectedTransport: TransportType = TransportType.Walk;
  minutes: number = 0;
  editingCheckpointIndex: number | null = null;

  constructor(private fb: FormBuilder, private service: TourService, private route: ActivatedRoute, private router: Router) {
    this.checkpointForm = this.fb.group({
      name: ['', Validators.required],
      description: ['', Validators.required],
      latitude: [0, Validators.required],
      longitude: [0, Validators.required],
      image_url: [''] 
    });
  }

  ngOnInit(): void {
    this.tourId = Number(this.route.snapshot.paramMap.get('id'));
  }

  onLocationSelected(coords: { lat: number, lng: number }): void {
    this.checkpointForm.patchValue({ latitude: coords.lat, longitude: coords.lng });
  }

  onFilesSelected(event: any): void {
    const files = event.target.files;
    if (files) {
      for (let file of files) {
        this.selectedFiles.push(file);
        const reader = new FileReader();
        reader.onload = (e: any) => this.previews.push(e.target.result);
        reader.readAsDataURL(file);
      }
    }
  }

  removeImage(index: number): void {
    this.selectedFiles.splice(index, 1);
    this.previews.splice(index, 1);
  }

  addCheckpoint(): void {
    if (this.checkpointForm.invalid || this.selectedFiles.length === 0) {
      alert("Popuni sva polja i izaberi sliku!");
      return;
    }

    const formData = new FormData();
    formData.append('name', this.checkpointForm.value.name);
    formData.append('description', this.checkpointForm.value.description);
    formData.append('latitude', this.checkpointForm.value.latitude.toString());
    formData.append('longitude', this.checkpointForm.value.longitude.toString());
    
    formData.append('image', this.selectedFiles[0]);

    this.service.addCheckpoint(this.tourId, formData).subscribe({
      next: (res) => {
        this.checkpoints = [...this.checkpoints, res];
        this.checkpointForm.reset();
        this.previews = [];
        this.selectedFiles = [];
        alert("Checkpoint sačuvan!");
      },
      error: (err) => {
        console.error("Greška pri slanju:", err);
        alert("Server i dalje odbija zahtev. Proveri Network tab za detalje.");
      }
    });
  }

  onTransportChange(type: TransportType): void {
    this.minutes = Math.round(this.tourDistanceKm * (type === TransportType.Walk ? 12 : 3));
  }

  addDuration(): void {
    this.durations.push({ transport: TransportType[this.selectedTransport], minutes: this.minutes });
  }

  deleteCheckpoint(index: number): void {
    this.checkpoints.splice(index, 1);
    this.checkpoints = [...this.checkpoints];
  }

  onRouteDistanceUpdated(dist: number): void {
    this.tourDistanceKm = dist;
    this.onTransportChange(this.selectedTransport);
  }

  finalizeTour(): void {
    if (this.checkpoints.length < 2) {
      alert("Tura mora imati bar 2 tačke!");
      return;
    }
    this.router.navigate(['/tour-list']);
  }
}