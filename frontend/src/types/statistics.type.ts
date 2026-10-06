export interface Distance {
    value: number;
    unit: "m" | "km";
}

export interface UserStatistics {
    finished_games: number;
    accuracy: number;
    best_score: number;
    closest_guess: Distance | null;
}