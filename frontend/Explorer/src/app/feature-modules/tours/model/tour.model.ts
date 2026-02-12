export enum Difficulty {
    EASY = 0,
    MEDIUM = 1,
    HARD = 2
}

export enum TransportType {
    WALK = 0,
    BIKE = 1,
    CAR = 2
}

export interface TourDuration {
    minutes: number;
    transportType: TransportType;
}

export interface Checkpoint {
    id?: number;
    name: string;
    description: string;
    latitude: number;
    longitude: number;
    image_url?: string;
}

export interface Tour {
    id?: number;
    name: string;
    description: string;
    difficulty: Difficulty;
    tags: string;
    price: number;
    author_id: number;
    status?: number;          // 0: Draft, 1: Published, 2: Archived
    distance?: number;        // Kilometraža sa mape
    durations?: TourDuration[]; // Niz vremena
    checkpoints?: Checkpoint[];
}