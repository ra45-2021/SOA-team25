export enum Difficulty {
    EASY = 0,
    MEDIUM = 1,
    HARD = 2
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
    checkpoints?: Checkpoint[];
}